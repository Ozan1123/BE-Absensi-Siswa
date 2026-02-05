CREATE TABLE attedance_tokens (
    id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    token_code VARCHAR(10) UNIQUE NOT NULL,
    created_by BIGINT NOT NULL,
    is_active BOOLEAN,
    late_after DATETIME,
    valid_until DATETIME,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (created_by) REFERENCES users(id) 
);