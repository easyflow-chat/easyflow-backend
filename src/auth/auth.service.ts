import {
  Injectable,
  InternalServerErrorException,
  Logger,
} from '@nestjs/common';
import { JwtService } from '@nestjs/jwt';
import * as bcrypt from 'bcrypt';
import { ErrorCodes } from 'enums/error-codes.enum';
import { UsersService } from '../user/users.service';

@Injectable()
export class AuthService {
  constructor(
    private usersService: UsersService,
    private jwtservice: JwtService,
  ) {}

  private readonly logger = new Logger(AuthService.name);

  async login(email: string, pass: string): Promise<{ accessToken: string }> {
    this.logger.log(`Attempting login for user with email: ${email}`);
    const user = await this.usersService.findUserByEmail(email);
    if (!user) {
      this.logger.error(
        `Login for user with email: ${email} failed, not found in database`,
      );
      throw new InternalServerErrorException({
        error: ErrorCodes.WRONG_CREDENTIALS,
      });
    }
    if (!(await bcrypt.compare(pass, user.password))) {
      this.logger.error(
        `Login for user with email: ${email} failed, invalid password`,
      );
      throw new InternalServerErrorException({
        error: ErrorCodes.WRONG_CREDENTIALS,
      });
    }
    const payload = { sub: user.id, email: user.email };
    const accessToken = await this.jwtservice.signAsync(payload);
    return { accessToken };
  }
}
