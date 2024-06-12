import { Logger } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { JwtService } from '@nestjs/jwt';
import {
  ConnectedSocket,
  MessageBody,
  OnGatewayConnection,
  OnGatewayDisconnect,
  OnGatewayInit,
  SubscribeMessage,
  WebSocketGateway,
  WebSocketServer,
} from '@nestjs/websockets';
import { Message } from '@prisma/client';
import { parse } from 'cookie';
import { signedCookie } from 'cookie-parser';
import { Server, Socket } from 'socket.io';
import { ChatService } from 'src/chat/chat.service';
import { SendMessageDTO } from 'src/chat/dto/send-message.dto';

@WebSocketGateway({
  cors: {
    origin: new ConfigService().get('FRONTEND_URL'),
    allowedHeaders: ['Content-Type', 'cookie'],
    credentials: true,
  },
})
export class WebsocketGateway implements OnGatewayInit, OnGatewayConnection, OnGatewayDisconnect {
  constructor(
    private readonly configService: ConfigService,
    private readonly chatService: ChatService,
    private readonly jwtService: JwtService,
  ) {}

  private readonly logger = new Logger(WebsocketGateway.name);

  @WebSocketServer() server: Server;

  afterInit(): void {
    this.logger.log('Websocket initialized');
  }

  async handleConnection(client: Socket): Promise<void> {
    const { sockets } = this.server.sockets;

    const cookie = parse(client.handshake.headers.cookie).access_token;

    // Extract and verify JWT from cookies
    const token = signedCookie(cookie, this.configService.get('COOKIE_SECRET'));

    if (token) {
      try {
        const decoded = await this.jwtService.verifyAsync(token, {
          secret: this.configService.get('JWT_SECRET'),
        });
        client.handshake.auth.userId = decoded.id;
        this.logger.log(`User ${decoded.id} connected`);
        this.logger.debug(`Number of connected clients: ${sockets.size}`);
        const chats = await this.chatService.getChatPreviews(decoded.id);
        chats.forEach(chat => client.join(chat.id));
      } catch (err) {
        this.logger.error('JWT verification failed', err.message);
        client.disconnect(true);
      }
    } else {
      this.logger.warn('No token found in cookies');
      client.disconnect(true);
    }
  }

  handleDisconnect(client: Socket): void {
    this.logger.log(`Client id: ${client.id} disconnected`);
  }

  @SubscribeMessage('send_message')
  async handleMessage(@MessageBody() data: SendMessageDTO, @ConnectedSocket() client: Socket): Promise<Message> {
    const message = await this.chatService.sendMessage(
      client.handshake.auth.userId,
      data.chatId,
      data.content,
      data.iv,
    );
    this.server.to(data.chatId).emit('receive_message', message);
    return message;
  }
}
