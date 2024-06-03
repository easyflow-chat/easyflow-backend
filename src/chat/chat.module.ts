import { Module } from '@nestjs/common';
import { PrismaClient } from '@prisma/client';
import { ChatController } from './chat.controller';
import { ChatService } from './chat.service';

@Module({
  imports: [PrismaClient],
  providers: [ChatService, PrismaClient],
  controllers: [ChatController],
})
export class ChatModule {}
