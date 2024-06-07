import { IsNotEmpty, IsString } from 'class-validator';

export class JoinChatDTO {
  @IsString()
  @IsNotEmpty()
  userId: string;

  @IsString()
  @IsNotEmpty()
  chatId: string;
}
