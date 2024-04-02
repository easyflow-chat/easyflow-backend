import { Body, Controller, Post } from '@nestjs/common';
import { AuthService } from './auth.service';

@Controller('auth')
export class AuthController {
  constructor(private authService: AuthService) {}

  @Post('login')
  async login(
    @Body() request: { email: string; password: string },
  ): Promise<{ accessToken: string }> {
    const accessToken = await this.authService.login(
      request.email,
      request.password,
    );
    return accessToken;
  }
}
