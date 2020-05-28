CREATE DATABASE IF NOT EXISTS bludgeon;
CREATE USER IF NOT EXISTS 'bludgeon'@'%' identified by 'bludgeon';
USE bludgeon;

CREATE TABLE IF NOT EXISTS timer (     
    id BIGINT NOT NULL AUTO_INCREMENT,
    uuid TEXT,
    activesliceid BIGINT,
    activesliceuuid TEXT,    
    start BIGINT,
    finish BIGINT,
    elapsedtime BIGINT,
    INDEX(id, uuid(36)),

    PRIMARY KEY (id)
    -- FOREIGN KEY (employeeid)
        -- REFERENCES employee(id)
        -- ON UPDATE CASCADE ON DELETE RESTRICT
)ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS timeslice (     
    id BIGINT NOT NULL AUTO_INCREMENT,
    uuid TEXT,
    timerid BIGINT,
    timeruuid TEXT,    
    start BIGINT,
    finish BIGINT,
    elapsedtime BIGINT,
    INDEX(id, uuid(36)),

    PRIMARY KEY (id)
    -- FOREIGN KEY (employeeid)
        -- REFERENCES employee(id)
        -- ON UPDATE CASCADE ON DELETE RESTRICT
)ENGINE=InnoDB;

-- CREATE TABLE IF NOT EXISTS client (
--     id BIGINT NOT NULL AUTO_INCREMENT,
--     name TEXT,
--     rate FLOAT

--     PRIMARY KEY (id)
-- )ENGINE=InnoDB;

-- CREATE TABLE IF NOT EXISTS employee (
--     id BIGINT NOT NULL AUTO_INCREMENT,
--     first_name TEXT,
--     last_name TEXT,
--     INDEX(id, unitid),

--     PRIMARY KEY (id),
-- )ENGINE=InnoDB;

-- CREATE TABLE IF NOT EXISTS project (
--     id BIGINT NOT NULL AUTO_INCREMENT,
--     client_id BIGINT NOT NULL,
--     description TEXT,
--     INDEX(id, unitid),

--     PRIMARY KEY (id),
--     FOREIGN KEY (clientid)
--         REFERENCES client(id)
-- )ENGINE=InnoDB;

GRANT ALL ON bludgeon.* to 'bludgeon'@'%'