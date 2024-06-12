import { IsArray, IsDefined, IsOptional, IsString } from 'class-validator';

export class CreatChatDTO {
  @IsDefined()
  @IsString()
  name: string;

  @IsOptional()
  @IsString()
  picture: string | undefined;

  @IsOptional()
  @IsString()
  description: string | undefined;

  @IsDefined()
  @IsArray()
  @IsString({ each: true })
  users: string[];

  @IsDefined()
  @IsArray()
  userKeys: { key: string; userId: string }[];
}
