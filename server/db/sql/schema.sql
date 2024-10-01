CREATE TABLE IF NOT EXISTS users(
    user_id VARCHAR(36) PRIMARY KEY,
    password varchar(256) NOT NULL,
    username varchar(256) NOT NULL
);

CREATE TABLE IF NOT EXISTS conversation(
    conversation_id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    conversation_name varchar(256)
);

CREATE TABLE IF NOT EXISTS user_subscriptions(
    user_subscriptions_id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    FOREIGN KEY user_id VARCHAR(36) REFERENCES users (user_id),
    FOREIGN KEY conversation_id INT UNSIGNED REFERENCES conversation (conversation_id)
);

CREATE TABLE IF NOT EXISTS messages(
    message_id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    message_timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    message VARCHAR(32768) NOT NULL,
    FOREIGN KEY user_id VARCHAR(36) NOT NULL REFERENCES users (user_id),
    FOREIGN KEY conversation_id INT UNSIGNED NOT NULL REFERENCES conversation (conversation_id)
    edited BOOL DEFAULT FALSE
);