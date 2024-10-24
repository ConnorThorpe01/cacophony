USE CACOPHONY_DB;

CREATE TABLE IF NOT EXISTS users (
     user_id VARCHAR(36) PRIMARY KEY,
     username VARCHAR(256) NOT NULL UNIQUE,
     email varchar(256) NOT NULL UNIQUE,
     password VARCHAR(256) NOT NULL
);

CREATE TABLE IF NOT EXISTS friends (
    friends_id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    user_id_1 VARCHAR(36),
    user_id_2 VARCHAR(36),
    FOREIGN KEY (user_id_1) REFERENCES users(user_id),
    FOREIGN KEY (user_id_2) REFERENCES users(user_id),
    UNIQUE(user_id_1, user_id_2)
);

CREATE TABLE IF NOT EXISTS `group` (
    group_id VARCHAR(36) PRIMARY KEY,
    group_name VARCHAR(256)
);

CREATE TABLE IF NOT EXISTS user_subscriptions (
    user_subscriptions_id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    user_id VARCHAR(36),
    group_id VARCHAR(36),
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (group_id) REFERENCES `group`(group_id)
);

CREATE TABLE IF NOT EXISTS messages (
    message_id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    message_timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    message TEXT NOT NULL,
    edited BOOL DEFAULT FALSE,
    user_id VARCHAR(36) NOT NULL,
    to_id VARCHAR(36),
    group_id VARCHAR(36),
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (group_id) REFERENCES `group`(group_id)
);
