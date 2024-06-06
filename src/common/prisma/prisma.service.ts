import { Injectable, Logger, OnModuleDestroy, OnModuleInit } from '@nestjs/common';
import { Prisma, PrismaClient } from '@prisma/client';
import { withAccelerate } from '@prisma/extension-accelerate';
import { Mutex, MutexInterface, withTimeout } from 'async-mutex';

import { ConfigService } from '@nestjs/config';
import { Timer } from '../timer/timer';

@Injectable()
export class PrismaService extends PrismaClient implements OnModuleInit, OnModuleDestroy {
  /**
   * The internal client is used for session dependent queries
   *
   * Its connection pool size is limited to a single connection.
   * This is required for the LOCK statements to work since they
   * depend on the connection, e.g. you can only release a lock
   * if it's been obtained by the same session.
   *
   * For more information on this behaviour visit:
   * [1] https://www.prisma.io/docs/guides/performance-and-optimization/connection-management
   * [2] https://mariadb.com/kb/en/get_lock/#description
   */
  #internalClient;

  /**
   * The service logger
   */
  readonly #logger: Logger = new Logger(PrismaService.name);

  /**
   * An exclusive lock to block concurrent access to modify
   * entries of the #mutexes
   */
  readonly #modifyMutexesMutex = new Mutex();

  /**
   * The map of named mutex of the current process
   */
  readonly #mutexes: Map<string, Mutex> = new Map();

  public constructor(configService: ConfigService) {
    // The client used by the app itself
    super({
      datasources: {
        db: {
          url: configService.get('DATABASE_URL'),
        },
      },
      log: [
        {
          emit: 'event',
          level: 'query',
        },
        {
          emit: 'stdout',
          level: 'error',
        },
        {
          emit: 'stdout',
          level: 'info',
        },
        {
          emit: 'stdout',
          level: 'warn',
        },
      ],
    });

    // eslint-disable-next-line
    // @ts-ignore
    this.$on('query', (e: { query: string; params: string; duration: number }) => {
      if (configService.get('DATABASE_QUERY_LOGGING')) {
        this.#logger.debug(
          `Query: ${e.query} / Params: ${configService.get('DATABASE_DEBUG_MODE') ? e.params : '[REDACTED]'}  / Duration: ${
            e.duration
          }ms`,
        );
      }
    });

    // Set up the internal client with a connection pool size of 1
    const internalClientURL = new URL(configService.get('DATABASE_URL'));
    internalClientURL.searchParams.append('connection_limit', '1');

    this.#internalClient = new PrismaClient({
      datasources: {
        db: {
          url: internalClientURL.toString(),
        },
      },
    }).$extends(withAccelerate());

    // eslint-disable-next-line
    // @ts-ignore
    this.$use((params, next) => {
      // this middleware will be executed for every request and will return only non-deleted items
      // if you want to include deleted items use the following filter in your instance call: { deletedAt: 'any' }
      // if you want to include ONLY deleted items use the following filter in your instance call: { deletedAt: { not: null } }

      const entitiesWithSoftDelete = [
        'User',
        'Job',
        'Order',
        'OrderState',
        'Screening',
        'Check',
        'Customer',
        'Article',
        'Contact',
      ];

      // eslint-disable-next-line
      //@ts-ignore
      if (entitiesWithSoftDelete.includes(params.model)) {
        if (!params.args?.where?.deletedAt) {
          if (
            ['findMany', 'findFirst', 'findUnique', 'findFirstOrThrow', 'findUniqueOrThrow', 'count'].includes(
              params.action,
            )
          ) {
            params = { ...params, args: { ...params.args, where: { ...params.args?.where, deletedAt: null } } };
          }
        } else if (params.args?.where?.deletedAt === 'any') {
          if (
            ['findMany', 'findFirst', 'findUnique', 'findFirstOrThrow', 'findUniqueOrThrow', 'count'].includes(
              params.action,
            )
          ) {
            params = { ...params, args: { ...params.args, where: { ...params.args?.where, deletedAt: undefined } } };
          }
        }
      }
      return next(params);
    });
  }

  /**
   * Connect to the database when the application starts
   */
  async onModuleInit(): Promise<void> {
    this.#logger.debug('Connecting to internal database');
    await this.#internalClient.$connect();
    this.#logger.debug('Connecting to application database');
    await this.$connect();
    this.#logger.debug('Database connected');
  }

  /**
   * Disconnect the database on shutdown
   */
  async onModuleDestroy(): Promise<void> {
    this.#logger.debug('Disconnectiong application database connection');
    await this.$disconnect();
    this.#logger.debug('Disconnecting internal database connection');
    await this.#internalClient.$disconnect();
    this.#logger.debug('Database disconnected');
  }

  /**
   * Runs a database transaction
   *
   * @param func  The transaction
   * @param opts  The transaction options
   * @returns     The transaction result
   */
  public transaction<T>(
    func: (th: Prisma.TransactionClient) => Promise<T>,
    opts: {
      timeout?: number;
    } = {
      timeout: 30000,
    },
  ): Promise<T> {
    return this.$transaction(func, opts);
  }

  /**
   * Creates a locked transaction. Other locked transactions with the
   * specified lock name will not be able to perform any actions.
   *
   * The lock is set both on process as well as database level.
   * Therefore this will lock out both the current process as well as
   * other instances trying to acquire the lock with the specified name.
   *
   * Note: Database disconnects / closing sessions will automatically
   *       release the database lock.
   *       See https://dev.mysql.com/doc/refman/8.0/en/locking-functions.html
   *       for more information on named locks
   *
   * @param innerFunc           The transaction to run
   * @param lockName            The lock name
   * @param acquireTimeout      The timeout for acquiring the lock in milliseconds.
   *                            Once the timeout is reached the transaction won't be
   *                            executed and an exception is thrown. Defaults to 5 seconds
   * @param transactionTimeout  The transaction timeout in milliseconds
   * @returns                   The transaction result
   */
  public async lockedTransaction<T>(
    innerFunc: (th: Prisma.TransactionClient) => Promise<T>,
    lockName = 'DB_LOCK',
    acquireTimeout = 30000,
    transactionTimeout = 5000,
  ): Promise<T> {
    this.#logger.debug(`Initiating locked transaction '${lockName}'`);

    const releaser = await this.acquireDistributedLock(lockName, acquireTimeout);

    // Run the transaction
    try {
      this.#logger.debug(`Running transaction "${lockName}"`);
      const res = await this.transaction(
        async th => {
          const ret = await innerFunc(th);
          return ret;
        },
        {
          timeout: transactionTimeout,
        },
      );

      this.#logger.debug(`Finished transaction "${lockName}"`);

      // This should not throw to avoid double lock release
      try {
        await releaser();
      } catch (e) {
        this.#logger.error(
          `Could not release transaction lock "${lockName}": ${e instanceof Error ? e.message : 'N/A'}`,
        );
      }

      return res;
    } catch (e) {
      await releaser();
      throw e;
    }
  }

  /**
   * Checks whether the lock with the specified name is held
   * by any session (not necessarily the current one)
   *
   * @param lockName  The lock name
   * @returns         True if the lock is held by any session
   */
  public isLockUsedByAnySession(lockName: string): Promise<boolean> {
    return this.#isLockUsedByAnySession(lockName);
  }

  /**
   * Acquires a distributed lock
   *
   * @param name     The lock nmae
   * @param timeout  The acquire timeout
   * @returns        The releaser
   */
  public async acquireDistributedLock(name: string, timeout: number): Promise<() => Promise<void>> {
    const acquireTimer = Timer.start();

    // Acquire process lock
    const releaseMutex = await this.#acquireMutex(name, timeout);

    const remainingTime = timeout - acquireTimer.elapsed;
    if (remainingTime < 0) {
      this.#logger.error(`Could not acquire mutex for "${name}"`);
      releaseMutex();
      throw new Error(`Reached transaction mutex timeout`);
    }

    // Acquire database lock
    try {
      this.#logger.debug(`Acquiring database lock`);
      const hasDatabaseLock = await this.#acquireNamedLock(
        name,
        remainingTime > 1000 ? Math.ceil(remainingTime / 1000) : 1,
      );

      if (!hasDatabaseLock) {
        this.#logger.error(`Could not acquire lock for "${name}"`);
        throw new Error(`Failed to acquire transaction lock`);
      }
    } catch (err) {
      this.#logger.debug(`Releasing mutex "${name}"`);
      releaseMutex();
      throw err;
    }

    // Return releaser
    return async () => {
      this.#logger.debug(`Releasing database lock "${name}"`);

      try {
        await this.#releaseLock(name);
        this.#logger.debug(`Released database lock "${name}"`);
      } catch {
        this.#logger.error(`Could not release database lock ${name}`);
      }

      this.#logger.debug(`Releasing mutex "${name}"`);
      releaseMutex();
      this.#logger.debug(`Released mutex "${name}"`);
    };
  }

  /**
   * Acquires the mutex with the specified name
   *
   * @param name       The mutex name
   * @param timeout    The acquire timeout in milliseconds
   * @param traceUuid  The debug trace UUID
   * @returns          The releaser
   */
  async #acquireMutex(name: string, timeout: number): Promise<MutexInterface.Releaser> {
    const timer = Timer.start();
    this.#logger.debug(`Loading named mutex ${name}, (MMM: ${this.#modifyMutexesMutex.isLocked()})`);

    const mutex = await this.#getOrCreateNamedMutex(name, timeout);
    if (timeout - timer.elapsed < 0) throw new Error();
    this.#logger.debug(`Acquiring named mutex ${name}, (NM: ${mutex.isLocked()})`);

    return withTimeout(mutex, timeout).acquire();
  }

  /**
   * Gets or creates a named mutex
   *
   * @param name     The mutex name
   * @param timeout  The timeout to load the mutex in milliseconds
   * @returns        The mutex
   */
  async #getOrCreateNamedMutex(name: string, timeout: number): Promise<Mutex> {
    const releaser = await withTimeout(this.#modifyMutexesMutex, timeout).acquire();

    try {
      const existingMutex = this.#mutexes.get(name);
      if (existingMutex) {
        releaser();
        return existingMutex;
      }

      const newMutex = new Mutex();
      this.#mutexes.set(name, newMutex);
      releaser();
      return newMutex;
    } catch (e) {
      this.#logger.error(`Could not upsert mutex "${name}": ${e instanceof Error ? e.message : 'N/A'}`);
      releaser();
      throw e;
    }
  }

  /**
   * Checks whether a lock is used by any session
   *
   * @param lockName  The lock name
   * @returns         True if the lock is in use
   */
  async #isLockUsedByAnySession(lockName: string): Promise<boolean> {
    const [{ is_used }] = await this.#executeLockQuery<{
      is_used: null | number | bigint;
    }>("SELECT IS_USED_LOCK('?') as is_used", lockName);

    return (typeof is_used === 'number' || typeof is_used === 'bigint') && is_used > 0;
  }

  /**
   * Attempts to acquire the specified lock
   *
   * @param lockName        The lock name
   * @param acquireTimeout  The timeout for acquiring the lock in seconds
   * @returns               True if the lock was acquired
   */
  async #acquireNamedLock(lockName: string, acquireTimeout: number): Promise<boolean> {
    const [{ acquired }] = await this.#executeLockQuery<{
      acquired: null | number | bigint;
    }>(`SELECT GET_LOCK('?', ${acquireTimeout} ) AS acquired`, lockName);

    return acquired != null && [1, BigInt(1)].includes(acquired);
  }

  /**
   * Releases the specified lock
   *
   * @param lockName  The lock name
   */
  async #releaseLock(lockName: string): Promise<void> {
    await this.#executeLockQuery(`SELECT RELEASE_LOCK('?')`, lockName);
  }

  /**
   * Executes a lock query. The query must include exactly one placeholder
   * for where the lock name should be injected. The placeholder must have the
   * following format: '?'
   *
   * @param query     The lock query
   * @param lockName  The lock name
   * @returns         The query result
   */
  async #executeLockQuery<T>(query: string, lockName: string): Promise<T[]> {
    if (query.split("'?'").length !== 2) throw new Error('Must have exactly one lock name placeholder');

    const encodedLock = this.#safeLockName(lockName);
    this.#logger.debug(`Using lock name ${encodedLock} for ${query}`);
    this.#logger.debug(new Error().stack?.replace('Error', 'Lock Query Stack Trace'));

    const res = await this.#internalClient.$queryRawUnsafe(query.replace("'?'", `'${encodedLock}'`));

    this.#logger.debug(JSON.stringify(res));
    return res as T[];
  }

  /**
   * Encodes the lock name as base64 URL to get rid of any
   * potentially conflicting characters
   *
   * @param lockName  The lock name
   * @returns         The encoded lock name
   */
  #safeLockName(lockName: string): string {
    return Buffer.from(lockName, 'utf-8').toString('base64url').replace(/=/, '');
  }
}
