import { CanActivate, ExecutionContext, Injectable, UnauthorizedException } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { Reflector } from '@nestjs/core';
import { JwtService } from '@nestjs/jwt';
import { ErrorCodes } from 'enums/error-codes.enum';
import { Request } from 'express';

@Injectable()
export class AuthGuard implements CanActivate {
  constructor(
    private readonly jwtService: JwtService,
    private readonly reflector: Reflector,
    private readonly configService: ConfigService,
  ) {}

  async canActivate(context: ExecutionContext): Promise<boolean> {
    const isPublic = this.reflector.getAllAndOverride<boolean>('isPublic', [context.getHandler(), context.getClass()]);
    if (isPublic) {
      return true;
    }
    const request = context.switchToHttp().getRequest();
    const token = this.extractTokenFromHeader(request);
    if (!token) {
      throw new UnauthorizedException({ error: ErrorCodes.UNAUTHORIZED });
    }
    try {
      const payload = await this.jwtService.verifyAsync<{ id: string; email: string }>(token, {
        secret: this.configService.get('JWT_SECRET'),
      });
      request['userId'] = payload.id;
    } catch {
      throw new UnauthorizedException({ error: ErrorCodes.UNAUTHORIZED });
    }
    return true;
  }

  private extractTokenFromHeader(request: Request): string | undefined {
    return request.signedCookies.access_token;
  }
}
