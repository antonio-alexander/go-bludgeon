
-- DROP DATABASE IF EXISTS bludgeon;
CREATE DATABASE IF NOT EXISTS bludgeon;

USE bludgeon;

-- DROP TABLE IF EXISTS timer
CREATE TABLE IF NOT EXISTS timer (     
    timer_id BIGINT NOT NULL AUTO_INCREMENT,
    -- task_id BIGINT,
    -- employee_id BIGINT NOT NULL,
    timer_uuid TEXT(36) NOT NULL,
    timer_start BIGINT NOT NULL,
    timer_finish BIGINT,
    timer_comment TEXT NOT NULL DEFAULT "",
    timer_archived BOOLEAN NOT NULL DEFAULT FALSE,
    timer_billed BOOLEAN NOT NULL DEFAULT FALSE,
    timer_completed BOOLEAN NOT NULL DEFAULT FALSE,

    PRIMARY KEY (timer_id),
    -- FOREIGN KEY (employee_id)
    --     REFERENCES employee(employee_id),
    -- FOREIGN KEY (task_id)
    --     REFERENCES task(task_id),
    UNIQUE(timer_uuid(36)),
    INDEX(timer_id)

)ENGINE=InnoDB;

-- DROP TABLE IF EXISTS slice
CREATE TABLE IF NOT EXISTS slice (     
    slice_id BIGINT NOT NULL AUTO_INCREMENT,
    timer_id BIGINT NOT NULL,
    slice_uuid TEXT(36) NOT NULL,
    slice_start BIGINT NOT NULL,
    slice_finish BIGINT,
    slice_archived BOOLEAN,
    slice_elapsed_time BIGINT AS (slice_finish-slice_start),

    PRIMARY KEY (slice_id),
    FOREIGN KEY (timer_id)
        REFERENCES timer(timer_id),
    UNIQUE(slice_uuid(36))

)ENGINE=InnoDB;

-- DROP TABLE IF EXISTS timer_slice_active
CREATE TABLE IF NOT EXISTS timer_slice_active (
    timer_id BIGINT NOT NULL,
    slice_id BIGINT NOT NULL,

    PRIMARY KEY (timer_id, slice_id),
    FOREIGN KEY (timer_id) 
        REFERENCES timer(timer_id),
    FOREIGN KEY (slice_id) 
        REFERENCES slice(slice_id),
    UNIQUE(timer_id),
    INDEX (slice_id, timer_id)

)Engine=InnoDB;

-- DROP TABLE IF EXISTS client
CREATE TABLE IF NOT EXISTS client (
    client_id BIGINT NOT NULL AUTO_INCREMENT,
    client_uuid TEXT(36) NOT NULL,
    client_name TEXT,

    PRIMARY KEY (client_id),
    UNIQUE(client_uuid(36))

)ENGINE=InnoDB;

-- DROP TABLE IF EXISTS employee
CREATE TABLE IF NOT EXISTS employee (
    employee_id BIGINT NOT NULL AUTO_INCREMENT,
    employee_uuid TEXT,
    employee_first_name TEXT,
    employee_last_name TEXT,

    PRIMARY KEY (employee_id),
    UNIQUE(employee_uuid)

)ENGINE=InnoDB;

-- DROP TABLE IF EXISTS project_client
CREATE TABLE IF NOT EXISTS project_client (
    project_id BIGINT NOT NULL AUTO_INCREMENT,
    client_id BIGINT NOT NULL AUTO_INCREMENT,

    PRIMARY KEY (project_id, client_id),
    FOREIGN KEY (project_id) 
        REFERENCES project(project_id),
    FOREIGN KEY (client_id) 
        REFERENCES client(client_id),
    UNIQUE(project_uuid),
    INDEX(client_id, project_id)
    
)ENGINE=InnoDB;

-- DROP TABLE IF EXISTS project
CREATE TABLE IF NOT EXISTS project (
    project_id BIGINT NOT NULL AUTO_INCREMENT,
    project_uuid TEXT(36) NOT NULL,
    project_description TEXT,

    PRIMARY KEY (project_id),
    UNIQUE(project_uuid),
    INDEX(project_id)
    
)ENGINE=InnoDB;

-- DROP TABLE IF EXISTS task
CREATE TABLE IF NOT EXISTS task (
    task_id BIGINT NOT NULL AUTO_INCREMENT,
    project_id BIGINT NOT NULL,
    task_uuid TEXT(36) NOT NULL,
    task_description TEXT,

    PRIMARY KEY (project_id),
    FOREIGN KEY (project_id) 
        REFERENCES project(project_id),
    UNIQUE(project_uuid),
    INDEX(project_id)
    
)ENGINE=InnoDB;