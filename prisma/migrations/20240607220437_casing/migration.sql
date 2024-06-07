/*
  Warnings:

  - You are about to drop the column `chatId` on the `chat_user_keys` table. All the data in the column will be lost.
  - You are about to drop the column `userId` on the `chat_user_keys` table. All the data in the column will be lost.
  - You are about to drop the column `chatId` on the `message` table. All the data in the column will be lost.
  - You are about to drop the column `senderId` on the `message` table. All the data in the column will be lost.
  - You are about to drop the column `chatId` on the `user_chat` table. All the data in the column will be lost.
  - You are about to drop the column `userId` on the `user_chat` table. All the data in the column will be lost.
  - Added the required column `chat_id` to the `chat_user_keys` table without a default value. This is not possible if the table is not empty.
  - Added the required column `user_id` to the `chat_user_keys` table without a default value. This is not possible if the table is not empty.
  - Added the required column `chat_id` to the `message` table without a default value. This is not possible if the table is not empty.
  - Added the required column `sender_id` to the `message` table without a default value. This is not possible if the table is not empty.
  - Added the required column `chat_id` to the `user_chat` table without a default value. This is not possible if the table is not empty.
  - Added the required column `user_id` to the `user_chat` table without a default value. This is not possible if the table is not empty.

*/
-- DropForeignKey
ALTER TABLE `chat_user_keys` DROP FOREIGN KEY `chat_user_keys_chatId_fkey`;

-- DropForeignKey
ALTER TABLE `chat_user_keys` DROP FOREIGN KEY `chat_user_keys_userId_fkey`;

-- DropForeignKey
ALTER TABLE `message` DROP FOREIGN KEY `message_chatId_fkey`;

-- DropForeignKey
ALTER TABLE `message` DROP FOREIGN KEY `message_senderId_fkey`;

-- DropForeignKey
ALTER TABLE `user_chat` DROP FOREIGN KEY `user_chat_chatId_fkey`;

-- DropForeignKey
ALTER TABLE `user_chat` DROP FOREIGN KEY `user_chat_userId_fkey`;

-- AlterTable
ALTER TABLE `chat_user_keys` DROP COLUMN `chatId`,
    DROP COLUMN `userId`,
    ADD COLUMN `chat_id` VARCHAR(191) NOT NULL,
    ADD COLUMN `user_id` VARCHAR(191) NOT NULL;

-- AlterTable
ALTER TABLE `message` DROP COLUMN `chatId`,
    DROP COLUMN `senderId`,
    ADD COLUMN `chat_id` VARCHAR(191) NOT NULL,
    ADD COLUMN `sender_id` VARCHAR(191) NOT NULL;

-- AlterTable
ALTER TABLE `user_chat` DROP COLUMN `chatId`,
    DROP COLUMN `userId`,
    ADD COLUMN `chat_id` VARCHAR(191) NOT NULL,
    ADD COLUMN `user_id` VARCHAR(191) NOT NULL;

-- AddForeignKey
ALTER TABLE `chat_user_keys` ADD CONSTRAINT `chat_user_keys_chat_id_fkey` FOREIGN KEY (`chat_id`) REFERENCES `chat`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `chat_user_keys` ADD CONSTRAINT `chat_user_keys_user_id_fkey` FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `user_chat` ADD CONSTRAINT `user_chat_user_id_fkey` FOREIGN KEY (`user_id`) REFERENCES `user`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `user_chat` ADD CONSTRAINT `user_chat_chat_id_fkey` FOREIGN KEY (`chat_id`) REFERENCES `chat`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `message` ADD CONSTRAINT `message_chat_id_fkey` FOREIGN KEY (`chat_id`) REFERENCES `chat`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE `message` ADD CONSTRAINT `message_sender_id_fkey` FOREIGN KEY (`sender_id`) REFERENCES `user`(`id`) ON DELETE RESTRICT ON UPDATE CASCADE;
