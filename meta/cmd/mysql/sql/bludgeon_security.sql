-- DROP DATABASE IF EXISTS bludgeon;
CREATE DATABASE IF NOT EXISTS bludgeon;

USE bludgeon;

-- DROP USER 'bludgeon'@'%';
CREATE USER 'bludgeon'@'%' IDENTIFIED BY 'bludgeon';

-- DROP USER 'bludgeon'@'localhost';
CREATE USER 'bludgeon'@'localhost' IDENTIFIED BY 'bludgeon';

GRANT ALL PRIVILEGES ON bludgeon.* TO 'bludgeon'@'%';
GRANT ALL PRIVILEGES ON bludgeon.* TO 'bludgeon'@'localhost';

FLUSH PRIVILEGES;