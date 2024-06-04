/* eslint-disable no-console */

/**
 * migrate.ts
 *
 * This file is an entrypoint for running prisma database migrations.
 * It handles concurrent migration attempts via distributed locks.
 */
import { NestFactory } from '@nestjs/core';
import { execSync } from 'child_process';
import { access } from 'fs/promises';
import { join } from 'path';

import { ConfigService } from '@nestjs/config';
import { AppModule } from './app.module';
import { PrismaModule } from './common/prisma/prisma.module';
import { PrismaService } from './common/prisma/prisma.service';

/**
 * Probes the access to the specified path.
 *
 * @param path  The path
 * @returns     True if the path exists and can be accessed
 */
const probe = async (path: string): Promise<boolean> => {
  try {
    console.log(`Probing ${path}`);
    await access(path);
    return true;
  } catch (err) {
    console.warn(`Could not find ${path}`);
    return false;
  }
};

/**
 * Searches for the `prisma/schema.prisma` file in the current directory and the
 * parent directory.
 *
 * @returns  The absolute path to the schema.prisma file or NULL if it wasn't found
 */
const findPrismaSchema = async (): Promise<string | null> => {
  console.log('Looking for prisma schema..');

  const fromCurrentDir = join(__dirname, 'prisma', 'schema.prisma');
  if (await probe(fromCurrentDir)) return fromCurrentDir;

  const fromParentDir = join(__dirname, '..', 'prisma', 'schema.prisma');
  if (await probe(fromParentDir)) return fromParentDir;

  console.error(`Could not find prisma directory`);
  return null;
};

/**
 * Searches for the `node_modules/.bin/prisma` CLI in the current directory and
 * the parent directory.
 *
 * @returns  The absolute path to the prisma CLI or NULL if it wasn't found
 */
const findPrismaCli = async (): Promise<string | null> => {
  console.log('Looking for prisma CLI..');

  const fromCurrentDir = join(__dirname, 'node_modules', '.bin', 'prisma');
  if (await probe(fromCurrentDir)) return fromCurrentDir;

  const fromParentDir = join(__dirname, '..', 'node_modules', '.bin', 'prisma');
  if (await probe(fromParentDir)) return fromParentDir;

  console.error(`Could not find prisma CLI`);
  return null;
};

/**
 * The applications entry point
 */
async function bootstrap(): Promise<void> {
  const prismaSchema = await findPrismaSchema();
  if (!prismaSchema) process.exit(1);

  const prismaCli = await findPrismaCli();
  if (!prismaCli) process.exit(1);

  // The config service provides us with the necessary environment variables
  const app = await NestFactory.create(AppModule);
  const configService = app.get(ConfigService);

  // Avoid running mirations in parallel
  const prismaModule = await NestFactory.create(PrismaModule);
  const prismaService = prismaModule.get<PrismaService>(PrismaService);
  const releaser = await prismaService.acquireDistributedLock('PRISMA_MIGRATION', 15000);

  try {
    // Run the migration, will exit the current process with the same status code
    // as received from the child process.
    execSync(`${prismaCli} migrate deploy --schema ${prismaSchema}`, {
      stdio: 'inherit',
      env: {
        DATABASE_URL: configService.get('DATABASE_URL'),
        PATH: process.env.PATH,
      },
    });

    await releaser();
  } catch (err) {
    console.error(err);
    throw err;
  } finally {
    await releaser();
  }
}

void bootstrap();
