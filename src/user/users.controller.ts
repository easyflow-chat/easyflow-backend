import { Body, Controller, Delete, Get, Param, Post, Put } from '@nestjs/common';
import { Prisma } from '@prisma/client';
import { CurrentUserId } from 'src/common/auth/current-user-id.decrator';
import { Public } from 'src/common/auth/public.decorator';
import { CreateUserDto } from './dto/create-user.dto';
import { UpdateUserDto } from './dto/update-user.dto';
import { UsersService } from './users.service';

//TODO: change to user
@Controller('users')
export class UsersController {
  constructor(private usersService: UsersService) {}

  @Public()
  @Post('signup')
  create(@Body() createUserDto: CreateUserDto): Promise<void> {
    return this.usersService.createUser(createUserDto);
  }

  @Get()
  findOne(@CurrentUserId() id: string): Promise<
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
    return this.usersService.findUserById(id);
  }

  @Put(':id')
  update(
    @Param('id') id: string,
    @Body() updateUserDto: UpdateUserDto,
  ): Promise<
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
    return this.usersService.updateUser(id, updateUserDto);
  }

  @Delete(':id')
  remove(@Param('id') id: string): Promise<
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
    return this.usersService.deleteUser(id);
  }
}
