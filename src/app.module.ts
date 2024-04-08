import { Module } from '@nestjs/common';
import { AuthModule } from './auth/auth.module';
import { UsersModule } from './user/users.module';

@Module({
  imports: [UsersModule, AuthModule],
})
export class AppModule {}
