CREATE TABLE IF NOT EXISTS `book`
(
    `id`             INT(11)      NOT NULL AUTO_INCREMENT,
    `title`          VARCHAR(255) NOT NULL,
    `subtitle`       VARCHAR(255) NULL,
    `author`         VARCHAR(255) NOT NULL,
    `category`       VARCHAR(255) NOT NULL,
    `notes`          BLOB         NULL,
    `slug`           VARCHAR(255) NOT NULL,
    `is_recommended` BOOL         NOT NULL DEFAULT FALSE,
    `finished_at`    DATE         NOT NULL,
    `created_at`     DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`     DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at`     DATETIME     NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_slug` (`slug`),
    KEY `idx_title` (`title`) USING BTREE,
    KEY `idx_author` (`author`) USING BTREE,
    KEY `idx_category` (`category`) USING BTREE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `book_log`
(
    `book_id`     INT(11) NOT NULL,
    `finished_at` DATE    NOT NULL,
    CONSTRAINT `book_log_ibfk_1` FOREIGN KEY (`book_id`) REFERENCES `book` (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `book_highlight`
(
    `id`         INT(11)      NOT NULL AUTO_INCREMENT,
    `book_id`    INT(11)      NOT NULL,
    `content`    BLOB         NOT NULL,
    `comment`    BLOB         NULL,
    `chapter`    VARCHAR(255) NULL,
    `page`       INT(11)      NULL,
    `created_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_book_id` (`book_id`) USING BTREE,
    CONSTRAINT `book_highlight_ibfk_1` FOREIGN KEY (`book_id`) REFERENCES `book` (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;