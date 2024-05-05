import { Body, Controller, Delete, Get, Post, Put } from '@nestjs/common';
import { Prisma, User } from '@prisma/client';
import { CurrentUserId } from 'src/common/auth/current-user-id.decrator';
import { Public } from 'src/common/auth/public.decorator';
import { CreateUserDto } from './dto/create-user.dto';
import { UpdateUserDto } from './dto/update-user.dto';
import { UserService } from './user.service';

@Controller('user')
export class UserController {
  constructor(private userService: UserService) {}

  @Public()
  @Post('signup')
  create(@Body() createUserDto: CreateUserDto): Promise<void> {
    return this.userService.createUser(createUserDto);
  }

  @Get()
  findOne(@CurrentUserId() id: User['id']): Promise<
    Prisma.UserGetPayload<{
      select: {
        id: true;
        createdAt: true;
        updatedAt: true;
        email: true;
        password: false;
      };
    }>
  > {
    return this.userService.findUserById(id);
  }

  @Put()
  update(
    @CurrentUserId() id: User['id'],
    @Body() updateUserDto: UpdateUserDto,
  ): Promise<
    Prisma.UserGetPayload<{
      select: {
        id: true;
        createdAt: true;
        updatedAt: true;
        email: true;
        password: false;
        profilePicture: false;
      };
    }>
  > {
    return this.userService.updateUser(id, updateUserDto);
  }

  @Delete()
  remove(@CurrentUserId() id: User['id']): Promise<
    Prisma.UserGetPayload<{
      select: {
        id: true;
        createdAt: true;
        updatedAt: true;
        email: true;
        password: false;
      };
    }>
  > {
    return this.userService.deleteUser(id);
  }

  @Get('profile-picture')
  getProfilePicture(@CurrentUserId() id: User['id']): Promise<User['profilePicture']> {
    return this.userService.getProfilePicture(id);
  }
}
