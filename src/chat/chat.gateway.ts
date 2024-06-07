import { Logger } from '@nestjs/common';
import {
  MessageBody,
  OnGatewayConnection,
  OnGatewayDisconnect,
  OnGatewayInit,
  SubscribeMessage,
  WebSocketGateway,
  WebSocketServer,
} from '@nestjs/websockets';
import { Message } from '@prisma/client';
import { Server, Socket } from 'socket.io';
import { CurrentUserId } from '../common/auth/current-user-id.decrator';
import { ChatService } from './chat.service';
import { SendMessageDTO } from './dto/send-message.dto';

@WebSocketGateway()
export class ChatGateway implements OnGatewayInit, OnGatewayConnection, OnGatewayDisconnect {
  constructor(private readonly chatService: ChatService) {}

  private readonly logger = new Logger(ChatGateway.name);

  @WebSocketServer() server: Server;

  afterInit(): void {
    this.logger.log('Websocket initialized');
  }

  handleConnection(client: Socket): void {
    const { sockets } = this.server.sockets;
    this.logger.log(`Client id: ${client.id} connected`);
    this.logger.debug(`Number of connected clients: ${sockets.size}`);
  }

  handleDisconnect(client: Socket): void {
    this.logger.log(`Client id: ${client.id} disconnected`);
  }

  @SubscribeMessage('send_message')
  handleMessage(@MessageBody() data: SendMessageDTO, @CurrentUserId() userId: string): Promise<Message> {
    return this.chatService.sendMessage(userId, data.chatId, data.content, data.iv);
  }
}
