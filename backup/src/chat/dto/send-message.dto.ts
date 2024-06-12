import { IsNotEmpty, IsString } from 'class-validator';

export class SendMessageDTO {
  @IsString()
  @IsNotEmpty()
  content: string;

  @IsString()
  @IsNotEmpty()
  chatId: string;

  @IsString()
  @IsNotEmpty()
  iv: string;
}
