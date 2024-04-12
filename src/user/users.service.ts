import { Injectable, InternalServerErrorException, Logger } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { Prisma, PrismaClient, User } from '@prisma/client';
import * as bcrypt from 'bcrypt';
import { ErrorCodes } from 'enums/error-codes.enum';
import { CreateUserDto } from './dto/create-user.dto';
import { UpdateUserDto } from './dto/update-user.dto';

@Injectable()
export class UsersService {
  constructor(
    private prisma: PrismaClient,
    private configService: ConfigService,
  ) {}

  private readonly logger = new Logger(UsersService.name);

  createUser(createUserDto: CreateUserDto): Promise<void> {
    return this.prisma.$transaction(async (tx: PrismaClient) => {
      try {
        this.logger.log(`Attempting to create user with email: ${createUserDto.email}`);
        const userInDatabase = await tx.user.findUnique({
          where: {
            email: createUserDto.email,
          },
        });
        if (userInDatabase) {
          this.logger.error(`User with email: ${createUserDto.email} already exists`);
          throw new InternalServerErrorException({
            error: ErrorCodes.ALREADY_EXISTS,
          });
        }
        const salt = await bcrypt.genSalt(this.configService.get('SALT_ROUNDS'));
        const hash = await bcrypt.hash(createUserDto.password, salt);
        await tx.user.create({
          data: {
            email: createUserDto.email,
            password: hash,
          },
          select: {
            id: true,
            createdAt: true,
            updatedAt: true,
            email: true,
            password: false,
          },
        });
      } catch (err) {
        this.logger.error(`Failed to create user with email: ${createUserDto.email}`);
        throw new InternalServerErrorException({ error: ErrorCodes.API_ERROR });
      }
    });
  }

  async findUserById(id: string): Promise<
    Prisma.UserGetPayload<{
      select: {
        id: true;
        createdAt: true;
        updatedAt: true;
        email: true;
        password: false;
      };
    }>
  > {
    this.logger.log(`Attempting to find user with id: ${id}`);
    const user = await this.prisma.user.findUnique({
      where: { id },
      select: {
        id: true,
        createdAt: true,
        updatedAt: true,
        email: true,
        password: false,
      },
    });
    if (!user) {
      this.logger.error(`User with id: ${id} not found`);
      throw new InternalServerErrorException({ error: ErrorCodes.NOT_FOUND });
    }
    return user;
  }

  async findUserByEmail(email: string): Promise<User> {
    this.logger.log(`Attempting to find user with email: ${email}`);
    const user = await this.prisma.user.findUnique({
      where: { email },
    });
    if (!user) {
      this.logger.error(`User with email: ${email} not found`);
      throw new InternalServerErrorException({ error: ErrorCodes.NOT_FOUND });
    }
    return user;
  }

  updateUser(
    id: string,
    updateUserDto: UpdateUserDto,
  ): Promise<
    Prisma.UserGetPayload<{
      select: {
        id: true;
        createdAt: true;
        updatedAt: true;
        email: true;
        password: false;
      };
    }>
  > {
    return this.prisma.$transaction(async (tx: PrismaClient) => {
      try {
        this.logger.log(`Attempting to update user with id: ${id}`);
        const UserInDatabase = await tx.user.findUnique({
          where: {
            id,
          },
        });
        if (!UserInDatabase) {
          this.logger.error(`User with id: ${id} not found`);
          throw new InternalServerErrorException({
            error: ErrorCodes.NOT_FOUND,
          });
        }
        const user = await tx.user.update({
          where: {
            id,
          },
          data: updateUserDto,
          select: {
            id: true,
            createdAt: true,
            updatedAt: true,
            email: true,
            password: false,
          },
        });
        return user;
      } catch (err) {
        this.logger.error(`Failed to update user with id: ${id}`);
        throw new InternalServerErrorException({ error: ErrorCodes.API_ERROR });
      }
    });
  }

  deleteUser(id: string): Promise<
    Prisma.UserGetPayload<{
      select: {
        id: true;
        createdAt: true;
        updatedAt: true;
        email: true;
        password: false;
      };
    }>
  > {
    return this.prisma.$transaction(async (tx: PrismaClient) => {
      try {
        this.logger.log(`Attempting to delete user with id: ${id}`);
        const UserInDatabase = await tx.user.findUnique({
          where: {
            id,
          },
        });
        if (!UserInDatabase) {
          this.logger.error(`User with id: ${id} not found`);
          throw new InternalServerErrorException({
            error: ErrorCodes.NOT_FOUND,
          });
        }
        const user = await tx.user.delete({
          where: { id },
          select: {
            id: true,
            createdAt: true,
            updatedAt: true,
            email: true,
            password: false,
          },
        });
        return user;
      } catch (err) {
        this.logger.error(`Failed to delete user with id: ${id}`);
        throw new InternalServerErrorException({ error: ErrorCodes.API_ERROR });
      }
    });
  }
}
