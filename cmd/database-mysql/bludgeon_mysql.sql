
-- DROP DATABASE IF EXISTS bludgeon;
CREATE DATABASE IF NOT EXISTS bludgeon;

USE bludgeon;

-- DROP TABLE IF EXISTS timer
CREATE TABLE IF NOT EXISTS timer (     
    timer_id BIGINT NOT NULL AUTO_INCREMENT,
    timer_uuid TEXT(36) NOT NULL,
    timestamp_start BIGINT,
    timestamp_finish BIGINT,
    comment TEXT,

    PRIMARY KEY (timer_id),
    INDEX(timer_id),
    UNIQUE(timer_uuid(36))

)ENGINE=InnoDB;

-- DROP TABLE IF EXISTS slice
CREATE TABLE IF NOT EXISTS slice (     
    slice_id BIGINT NOT NULL AUTO_INCREMENT,
    slice_uuid TEXT(36) NOT NULL,
    timestamp_start BIGINT NOT NULL,
    timestamp_finish BIGINT,

    PRIMARY KEY (slice_id),
    UNIQUE(slice_uuid(36))

)ENGINE=InnoDB;

-- DROP TABLE IF EXISTS timer_slice
CREATE TABLE IF NOT EXISTS timer_slice (
    timer_id BIGINT NOT NULL,
    slice_id BIGINT NOT NULL,

    PRIMARY KEY (timer_id, slice_id),
    INDEX (slice_id, timer_id),
    FOREIGN KEY (timer_id) REFERENCES timer(timer_id),
    FOREIGN KEY (slice_id) REFERENCES slice(slice_id)

)Engine=InnoDB;

-- DROP TABLE IF EXISTS timer_slice_active
CREATE TABLE IF NOT EXISTS timer_slice_active (
    timer_id BIGINT NOT NULL,
    slice_id BIGINT NOT NULL,

    PRIMARY KEY (timer_id, slice_id),
    FOREIGN KEY (timer_id) REFERENCES timer(timer_id),
    FOREIGN KEY (slice_id) REFERENCES slice(slice_id)

)Engine=InnoDB;

-- DROP TABLE IF EXISTS client
CREATE TABLE IF NOT EXISTS client (
    client_id BIGINT NOT NULL AUTO_INCREMENT,
    client_uuid TEXT(36) NOT NULL,
    first_name TEXT,
    last_name TEXT,

    PRIMARY KEY (client_id),
    UNIQUE(client_uuid(36))

)ENGINE=InnoDB;

-- DROP TABLE IF EXISTS timer_client
CREATE TABLE IF NOT EXISTS timer_client (
    timer_id BIGINT NOT NULL,
    client_id BIGINT NOT NULL,

    PRIMARY KEY (timer_id, client_id),
    FOREIGN KEY (timer_id) REFERENCES timer(timer_id),
    FOREIGN KEY (client_id) REFERENCES client(client_id),
    INDEX (client_id, timer_id)

)ENGINE=InnoDB;

-- DROP TABLE IF EXISTS employee
CREATE TABLE IF NOT EXISTS employee (
    employee_id BIGINT NOT NULL AUTO_INCREMENT,
    employee_uuid TEXT,
    name TEXT,
    rate FLOAT,

    PRIMARY KEY (employee_id),
    UNIQUE(employee_id)

)ENGINE=InnoDB;

-- DROP TABLE IF EXISTS timer_employee
CREATE TABLE IF NOT EXISTS timer_employee (
    timer_id BIGINT NOT NULL,
    employee_id BIGINT NOT NULL,

    PRIMARY KEY (timer_id, employee_id),
    FOREIGN KEY (timer_id) REFERENCES timer(timer_id),
    FOREIGN KEY (employee_id) REFERENCES employee(employee_id),
    INDEX (employee_id, timer_id)

)ENGINE=InnoDB;

-- DROP TABLE IF EXISTS project
CREATE TABLE IF NOT EXISTS project (
    project_id BIGINT NOT NULL AUTO_INCREMENT,
    project_uuid TEXT(36) NOT NULL,
    description TEXT,

    PRIMARY KEY (project_id),
    INDEX(project_id)
    
)ENGINE=InnoDB;

-- DROP TABLE IF EXISTS timer_project
CREATE TABLE IF NOT EXISTS timer_project (
    timer_id BIGINT NOT NULL,
    project_id BIGINT NOT NULL,

    PRIMARY KEY (timer_id, project_id),
    FOREIGN KEY (timer_id) REFERENCES timer(timer_id),
    FOREIGN KEY (project_id) REFERENCES project(project_id),
    INDEX (project_id, timer_id)

)ENGINE=InnoDB;
