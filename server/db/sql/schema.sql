USE CACOPHONY_DB;

CREATE TABLE IF NOT EXISTS users (
     user_id VARCHAR(36) PRIMARY KEY,
     username VARCHAR(256) NOT NULL,
     password VARCHAR(256) NOT NULL
);

CREATE TABLE IF NOT EXISTS conversation (
    conversation_id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    conversation_name VARCHAR(256)
);

CREATE TABLE IF NOT EXISTS user_subscriptions (
    user_subscriptions_id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    user_id VARCHAR(36),
    conversation_id INT UNSIGNED,
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (conversation_id) REFERENCES conversation(conversation_id)
);

CREATE TABLE IF NOT EXISTS messages (
    message_id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    message_timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    message TEXT NOT NULL,
    edited BOOL DEFAULT FALSE,
    user_id VARCHAR(36) NOT NULL,
    conversation_id INT UNSIGNED NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (conversation_id) REFERENCES conversation(conversation_id)
);
