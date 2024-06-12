import { Body, Controller, Get, Param, Post } from '@nestjs/common';
import { Chat, Message, User } from '@prisma/client';
import { CurrentUserId } from 'src/common/auth/current-user-id.decrator';
import { ChatService } from './chat.service';
import { CreatChatDTO } from './dto/create-chat.dto';
import { SendMessageDTO } from './dto/send-message.dto';

@Controller('chat')
export class ChatController {
  constructor(private readonly chatService: ChatService) {}

  @Get('/preview')
  etChatsForUserById(@CurrentUserId() userId: User['id']): Promise<Chat[]> {
    return this.chatService.getChatPreviews(userId);
  }

  @Get('/:id')
  getChatById(@CurrentUserId() userId: User['id'], @Param('id') chatId: string): Promise<Chat> {
    return this.chatService.getChatById(userId, chatId);
  }

  @Post()
  createChat(@CurrentUserId() userId: User['id'], @Body() createChatDto: CreatChatDTO): Promise<Chat> {
    return this.chatService.createChat(userId, createChatDto);
  }

  @Post('/send-message')
  sendMessage(@CurrentUserId() userId: User['id'], @Body() sendMessageDTO: SendMessageDTO): Promise<Message> {
    return this.chatService.sendMessage(userId, sendMessageDTO.chatId, sendMessageDTO.content, sendMessageDTO.iv);
  }
}
