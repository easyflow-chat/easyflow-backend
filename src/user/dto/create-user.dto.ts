import { IsString } from 'class-validator';

export class CreateUserDto {
  @IsString()
  email: string;

  @IsString()
  name: string;

  @IsString()
  password: string;

  @IsString()
  publicKey: string;

  @IsString()
  privateKey: string;

  @IsString()
  iv: string;
}
