import { Module } from '@nestjs/common';
import { PrismaClient } from '@prisma/client';
import { UserController } from './user.controller';
import { UserService } from './user.service';
import { ConfigModule } from '@nestjs/config';

@Module({
  imports: [ConfigModule],
  controllers: [UserController],
  providers: [UserService, PrismaClient],
  exports: [UserService],
})
export class UserModule {}
