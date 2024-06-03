import { Injectable, InternalServerErrorException, Logger } from '@nestjs/common';
import { Chat, Message, Prisma, PrismaClient, User } from '@prisma/client';
import { withAccelerate } from '@prisma/extension-accelerate';
import { ErrorCodes } from 'enums/error-codes.enum';
import { CreatChatDTO } from './dto/create-chat.dto';

@Injectable()
export class ChatService {
  private readonly prisma = new PrismaClient().$extends(withAccelerate());
  private readonly logger = new Logger(ChatService.name);

  createChat(userId: User['id'], createChatDto: CreatChatDTO): Promise<Chat> {
    this.logger.log(`Creating chat with name: ${createChatDto.name} for user with id: ${userId}`);
    return this.prisma.$transaction(tx => {
      try {
        return tx.chat.create({
          data: {
            name: createChatDto.name,
            picture: createChatDto.picture,
            description: createChatDto.description,
            users: {
              createMany: { data: createChatDto.users.map(id => ({ userId: id })) },
            },
            userKeys: {
              createMany: {
                data: createChatDto.userKeys.map(userKey => ({
                  key: userKey.key,
                  userId: userKey.userId,
                })),
              },
            },
          },
        });
      } catch (err) {
        this.logger.error(`Failed to create chat. Error: ${err}`);
        throw new InternalServerErrorException({ error: ErrorCodes.API_ERROR });
      }
    });
  }

  getChatPreviews(userId: User['id']): Promise<
    Prisma.ChatGetPayload<{
      include: {
        messages: {
          select: {
            content: true;
          };
        };
      };
    }>[]
  > {
    this.logger.log(`Finding chats for user with id: ${userId}`);
    return this.prisma.chat.findMany({
      where: {
        users: {
          some: {
            userId,
          },
        },
      },
      include: {
        messages: {
          select: {
            content: true,
          },
          orderBy: {
            createdAt: 'desc',
          },
          take: 1,
        },
      },
    });
  }

  getChatById(
    userId: User['id'],
    chatId: Chat['id'],
  ): Promise<
    Prisma.ChatGetPayload<{
      include: {
        messages: {
          select: {
            id: true;
            createdAt: true;
            updatedAt: true;
            content: true;
            iv: true;
            sender: {
              select: {
                id: true;
                name: true;
              };
            };
          };
        };
        users: {
          select: {
            user: {
              select: {
                id: true;
                name: true;
                profilePicture: true;
              };
            };
          };
        };
        userKeys: {
          select: {
            key: true;
            userId: true;
          };
        };
      };
    }>
  > {
    this.logger.log(`Finding chat with id: ${chatId} for user with id: ${userId}`);
    let chat;
    try {
      chat = this.prisma.chat.findUnique({
        where: {
          id: chatId,
        },
        include: {
          messages: {
            select: {
              id: true,
              createdAt: true,
              updatedAt: true,
              content: true,
              iv: true,
              sender: {
                select: {
                  id: true,
                  name: true,
                },
              },
            },
            orderBy: {
              createdAt: 'desc',
            },
          },
          users: {
            select: {
              user: {
                select: {
                  id: true,
                  name: true,
                  profilePicture: true,
                },
              },
            },
          },
          userKeys: {
            select: {
              key: true,
              userId: true,
            },
          },
        },
      });
    } catch (err) {
      this.logger.error(`Failed to find chat with id: ${chatId} for user with id: ${userId}. Error: ${err}`);
      throw new InternalServerErrorException({ error: ErrorCodes.API_ERROR });
    }
    if (!chat) {
      this.logger.error(`Failed to find chat with id: ${chatId} for user with id: ${userId}`);
      throw new InternalServerErrorException({ error: ErrorCodes.NOT_FOUND });
    }
    return chat;
  }

  sendMessage(
    senderId: User['id'],
    chatId: Chat['id'],
    content: Message['content'],
    iv: Message['iv'],
  ): Promise<
    Prisma.MessageGetPayload<{
      include: {
        sender: {
          select: {
            id: true;
            name: true;
          };
        };
      };
    }>
  > {
    this.logger.log(`Sending message for user with id: ${senderId} in chat with id: ${chatId}`);
    return this.prisma.$transaction(async tx => {
      try {
        const chat = await tx.chat.findUnique({
          where: {
            id: chatId,
          },
        });

        if (!chat) {
          this.logger.error(`Failed to send message for user with id: ${senderId} in chat with id: ${chatId}`);
          throw new InternalServerErrorException({ error: ErrorCodes.NOT_FOUND });
        }

        return tx.message.create({
          data: {
            content,
            senderId,
            chatId,
            iv,
          },
          include: {
            sender: {
              select: {
                id: true,
                name: true,
              },
            },
          },
        });
      } catch (err) {
        this.logger.error(
          `Failed to send message for user with id: ${senderId} in chat with id: ${chatId}. Error: ${err}`,
        );
        throw new InternalServerErrorException({ error: ErrorCodes.API_ERROR });
      }
    });
  }
}
