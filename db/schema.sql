USE server;

CREATE TABLE
    `user` (
                `id` VARCHAR(10) NOT NULL COMMENT 'ユーザー ID',
                `name` VARCHAR(10) NOT NULL COMMENT 'ユーザー名',
                `address` VARCHAR(50) NOT NULL COMMENT 'メールアドレス',
                `status` VARCHAR(10) NOT NULL COMMENT 'プレイ回数',
                `password` VARCHAR(10) NOT NULL COMMENT 'パスワード',
                `chat_number` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT 'チャット回数',
                `token` VARCHAR(100) NOT NULL COMMENT 'ユーザートークン',
                `created_at` DATETIME NOT NULL COMMENT '作成日時',
                `updated_at` DATETIME NOT NULL COMMENT '更新日時',
                PRIMARY KEY (`id`),
                INDEX `user_updated_at` (`updated_at`)
) COMMENT = 'ユーザー';

CREATE TABLE
    `admin` (
                         `id` VARCHAR(10) NOT NULL COMMENT '管理者アカウント ID',
                         `token` VARCHAR(100) COMMENT 'アドミントークン',
                         `password` VARCHAR(10) NOT NULL COMMENT 'パスワード',
                         `created_at` DATETIME NOT NULL COMMENT '作成日時',
                         `updated_at` DATETIME NOT NULL COMMENT '更新日時',
                         PRIMARY KEY (`id`)
) COMMENT = '管理者アカウント';

CREATE TABLE
    `User_Profile` (
                       `id` VARCHAR(10) NOT NULL COMMENT 'ユーザー ID',
                       `Comment` VARCHAR(100) COMMENT 'メッセージ',
                       `Friend_ID` VARCHAR(10) COMMENT 'フレンド一覧',
                       `created_at` DATETIME NOT NULL COMMENT '作成日時',
                       `updated_at` DATETIME NOT NULL COMMENT '更新日時',
                        PRIMARY KEY (`id`),
                        FOREIGN KEY (`id`) REFERENCES `user` (`id`),
                        FOREIGN KEY (`Friend_ID`) REFERENCES `user` (`id`)

) COMMENT = 'ユーザープロフィール';

INSERT INTO `admin`(id, password,created_at,updated_at)
VALUES ('admin', 'admin',NOW(),NOW());

INSERT INTO `user`(id,name,address,status,password,chat_number,token,created_at,updated_at)
VALUES ('test_1234','テストユーザー','test1@test.com','online','test',10,'test_1234_12345',NOW(),NOW());

INSERT INTO `user`(id,name,address,status,password,chat_number,token,created_at,updated_at)
VALUES ('test_5678','テストユーザー','test2@test.com','online','test',10,'test_5678_1234',NOW(),NOW());

INSERT INTO `User_Profile`(id,Comment,Friend_ID,created_at ,updated_at)
VALUES ('test_1234','テストユーザーです。','test_5678',NOW(), NOW());

INSERT INTO `User_Profile`(id,Comment,Friend_ID,created_at ,updated_at)
VALUES ('test_5678','テストユーザーです。','test_1234',NOW(), NOW());
