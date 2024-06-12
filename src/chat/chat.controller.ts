import { Body, Controller, Get, Param, Post } from '@nestjs/common';
import { Chat, ChatUserKeys, User } from '@prisma/client';
import { CurrentUserId } from 'src/common/auth/current-user-id.decrator';
import { ChatService } from './chat.service';
import { CreatChatDTO } from './dto/create-chat.dto';

@Controller('chat')
export class ChatController {
  constructor(private readonly chatService: ChatService) {}

  @Get('/preview')
  etChatsForUserById(@CurrentUserId() userId: User['id']): Promise<Chat[]> {
    return this.chatService.getChatPreviews(userId);
  }

  @Get('/data/:id')
  getChatById(@CurrentUserId() userId: User['id'], @Param('id') chatId: string): Promise<Chat> {
    return this.chatService.getChatById(userId, chatId);
  }

  @Get('/keys')
  getChatKeys(@CurrentUserId() userId: User['id']): Promise<ChatUserKeys[]> {
    return this.chatService.getChatKeys(userId);
  }

  @Post()
  createChat(@CurrentUserId() userId: User['id'], @Body() createChatDto: CreatChatDTO): Promise<Chat> {
    return this.chatService.createChat(userId, createChatDto);
  }
}
