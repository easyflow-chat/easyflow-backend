import { Injectable, InternalServerErrorException } from '@nestjs/common';
import { Prisma, PrismaClient, User } from '@prisma/client';
import * as bcrypt from 'bcrypt';
import { SALT_OR_ROUNDS } from '../constants/constants';
import { CreateUserDto } from './dto/create-user.dto';
import { UpdateUserDto } from './dto/update-user.dto';

@Injectable()
export class UsersService {
  constructor(private prisma: PrismaClient) {}
  async createUser(createUserDto: CreateUserDto): Promise<
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
    try {
      const salt = await bcrypt.genSalt(SALT_OR_ROUNDS);
      const hash = await bcrypt.hash(createUserDto.password, salt);
      const user = await this.prisma.user.create({
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
      return user;
    } catch (err) {
      throw new InternalServerErrorException({ error: err });
    }
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
    return user;
  }

  async findUserByEmail(email: string): Promise<User> {
    const user = await this.prisma.user.findUnique({
      where: { email },
    });
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
    const user = this.prisma.user.update({
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
  }

  async deleteUser(id: string): Promise<
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
    const user = await this.prisma.user.delete({
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
  }
}
