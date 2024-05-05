import { Injectable, InternalServerErrorException, Logger } from '@nestjs/common';
import { JwtService } from '@nestjs/jwt';
import { User } from '@prisma/client';
import * as bcrypt from 'bcrypt';
import { ErrorCodes } from 'enums/error-codes.enum';
import { UserService } from 'src/user/user.service';

@Injectable()
export class AuthService {
  constructor(
    private userService: UserService,
    private jwtService: JwtService,
  ) {}

  private readonly logger = new Logger(AuthService.name);

  async login(email: string, pass: string): Promise<string> {
    this.logger.log(`Attempting login for user with email: ${email}`);
    let user: Omit<User, 'profilePicture'> | undefined;
    try {
      user = await this.userService.findUserByEmail(email);
    } catch (err) {
      this.logger.error(`Login for user with email: ${email} failed, not found in database`);
      throw new InternalServerErrorException({
        error: ErrorCodes.WRONG_CREDENTIALS,
      });
    }
    if (!user) {
      this.logger.error(`Login for user with email: ${email} failed, no user returned from database`);
      throw new InternalServerErrorException({
        error: ErrorCodes.API_ERROR,
      });
    }
    if (!(await bcrypt.compare(pass, user.password))) {
      this.logger.error(`Login for user with email: ${email} failed, invalid password`);
      throw new InternalServerErrorException({
        error: ErrorCodes.WRONG_CREDENTIALS,
      });
    }
    const payload = { id: user.id, email: user.email };
    return this.jwtService.signAsync(payload);
  }
}
