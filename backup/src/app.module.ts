import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import * as Joi from 'joi';
import { ChatModule } from './chat/chat.module';
import { AuthModule } from './common/auth/auth.module';
import { UserModule } from './user/user.module';

@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      envFilePath: '.env',
      validationSchema: Joi.object({
        NODE_ENV: Joi.string().valid('development', 'production', 'test').default('development'),
        DATABASE_URL: Joi.string().required(),
        JWT_SECRET: Joi.string().required(),
        JWT_EXPIRATION_TIME: Joi.string().required(),
        COOKIE_SECRET: Joi.string().required(),
        PORT: Joi.number().required(),
        DATABASE_QUERY_LOGGING: Joi.boolean().default(false),
        DATABASE_DEBUG_MODE: Joi.boolean().default(false),
      }),
      validationOptions: {
        allowUnknown: true,
      },
      cache: true,
    }),
    UserModule,
    AuthModule,
    ChatModule,
  ],
})
export class AppModule {}
