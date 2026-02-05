CREATE TABLE attedance_logs(
    id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT,
    token_id BIGINT,
    status ENUM("hadir", "telat"),
    captured_ip VARCHAR(45) NULL,
    clock_in_time DATETIME,

    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (token_id) REFERENCES attedance_tokens(id)
);