import { createParamDecorator, ExecutionContext } from '@nestjs/common';
import { Request } from 'express';

/**
 * Decorator to get the current user id from the request Only work on protected routes
 * @param jwtService has to be explicitly passed in because it is not injectable
 * @param ctx gets set by NestJS
 * @returns the user id of the user that made the request
 */
export const CurrentUserId = createParamDecorator((data: string, ctx: ExecutionContext) => {
  const req: Request = ctx.switchToHttp().getRequest();
  return req['userId'];
});
