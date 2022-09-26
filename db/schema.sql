CREATE TABLE
    `users` (
                `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ユーザー ID',
                `name` VARCHAR(10) NOT NULL COMMENT 'ユーザー名',
                `status` VARCHAR(10) NOT NULL COMMENT 'プレイ回数',
                `chat_number` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'チャット回数',
                `token` VARCHAR(100) NOT NULL COMMENT 'ユーザートークン',
                `created_at` DATETIME NOT NULL COMMENT '作成日時',
                `updated_at` DATETIME NOT NULL COMMENT '更新日時',
                PRIMARY KEY (`id`),
                INDEX `user_updated_at` (`updated_at`)
) COMMENT = 'ユーザー';

CREATE TABLE
    `administrators` (
                         `id` VARCHAR(10) NOT NULL COMMENT '管理者アカウント ID',
                         `token` VARCHAR(100) NOT NULL COMMENT 'アドミントークン',
                         `last_logged_in_at` DATETIME COMMENT '最終ログイン日時',
                         `created_at` DATETIME NOT NULL COMMENT '作成日時',
                         `updated_at` DATETIME NOT NULL COMMENT '更新日時',
                         PRIMARY KEY (`id`)
) COMMENT = '管理者アカウント';
