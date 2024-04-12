import { Body, Controller, Post, Res } from '@nestjs/common';
import { Response } from 'express';
import { AuthService } from './auth.service';
import { Public } from './public.decorator';

@Controller('auth')
export class AuthController {
  constructor(private authService: AuthService) {}

  @Public()
  @Post('login')
  async login(
    @Res({ passthrough: true }) response: Response,
    @Body() request: { email: string; password: string },
  ): Promise<void> {
    response.cookie('access_token', await this.authService.login(request.email, request.password), {
      sameSite: 'lax',
      path: '/',
      signed: true,
    });
  }
}
