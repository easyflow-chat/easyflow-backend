import { Module } from '@nestjs/common';
import { ChatModule } from 'src/chat/chat.module';
import { ChatService } from 'src/chat/chat.service';
import { PrismaModule } from 'src/common/prisma/prisma.module';
import { WebsocketGateway } from './websocket.gateway';

@Module({
  imports: [ChatModule, PrismaModule],
  providers: [WebsocketGateway, ChatService],
})
export class WebsocketModule {}
