DROP TABLE IF EXISTS users;
CREATE TABLE users (
  id INT AUTO_INCREMENT NOT NULL,
  email VARCHAR(128) NOT NULL,
  username VARCHAR(64) NOT NULL,
  password_hash BINARY(64) NOT NULL,
  PRIMARY KEY (id)
);

