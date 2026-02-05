CREATE TABLE users (
    id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    nisn VARCHAR(10) UNIQUE,
    full_name VARCHAR(100),
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(225) NOT NULL,
    role ENUM ("guru", "siswa") DEFAULT "siswa",
    class_group VARCHAR(20)
);