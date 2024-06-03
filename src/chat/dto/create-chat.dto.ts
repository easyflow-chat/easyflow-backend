import { IsArray, IsObject, IsOptional, IsString } from 'class-validator';

export class CreatChatDTO {
  @IsString()
  name: string;

  @IsOptional()
  @IsString()
  picture: string | undefined;

  @IsOptional()
  @IsString()
  description: string | undefined;

  @IsArray()
  @IsString({ each: true })
  users: string[];

  @IsArray()
  @IsObject({ each: true })
  userKeys: { key: string; userId: string }[];
}
