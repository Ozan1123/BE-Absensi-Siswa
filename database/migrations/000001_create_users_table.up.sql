CREATE TABLE users (
    id BIGINT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    nisn VARCHAR(10) UNIQUE,
    full_name VARCHAR(100),
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(225) NOT NULL,
    role ENUM ("siswa", "guru", "admin", "superadmin") DEFAULT "siswa",
    class_group VARCHAR(20),
    parent_phone VARCHAR(20)
);