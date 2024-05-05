/*
  Warnings:

  - Added the required column `name` to the `User` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE `User` ADD COLUMN `bio` VARCHAR(1000) NULL,
    ADD COLUMN `name` VARCHAR(255) NOT NULL;
