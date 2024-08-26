CREATE TABLE notes (         
    id CHAR(36) PRIMARY KEY NOT NULL UNIQUE ,           
    text TEXT NOT NULL,
    author VARCHAR(255) NOT NULL,                   
    created_at DATETIME NOT NULL,  
    updated_at DATETIME NOT NULL
);

CREATE TABLE users (
    username VARCHAR(255) PRIMARY KEY NOT NULL UNIQUE, 
    password VARCHAR(255) NOT NULL           
);

CREATE TABLE sessions (
    username CHAR(36) PRIMARY KEY NOT NULL UNIQUE,
    token VARCHAR(255) NOT NULL,    
    FOREIGN KEY (username) REFERENCES users(username)
);

INSERT INTO users(username, password) VALUES 
("admin", "admin"),
("user123", "user321");