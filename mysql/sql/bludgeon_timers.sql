-- DROP DATABASE IF EXISTS bludgeon;
CREATE DATABASE IF NOT EXISTS bludgeon;

USE bludgeon;

-- DROP TABLE IF EXISTS timers;
CREATE TABLE IF NOT EXISTS timers (
    id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (UUID()),
    comment TEXT NOT NULL DEFAULT "",
    archived BOOLEAN NOT NULL DEFAULT FALSE,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    employee_id VARCHAR(36),
    aux_id BIGINT AUTO_INCREMENT,
    version INT NOT NULL DEFAULT 1,
    last_updated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    last_updated_by TEXT NOT NULL DEFAULT CURRENT_USER,
    INDEX(aux_id)
) ENGINE = InnoDB;

-- DROP TRIGGER IF EXISTS timers_audit_info_update;
CREATE TRIGGER timers_audit_info_update
BEFORE UPDATE ON timers FOR EACH ROW
    SET new.id = old.id, new.aux_id = old.aux_id, new.version = old.version+1, new.last_updated = CURRENT_TIMESTAMP(6), new.last_updated_by = CURRENT_USER;

-- DROP TABLE IF EXISTS timers_audit;
CREATE TABLE IF NOT EXISTS timers_audit (
    timer_id VARCHAR(36) NOT NULL,
    comment TEXT ,
    archived BOOLEAN,
    completed BOOLEAN,
    employee_id VARCHAR(36),
    version INT NOT NULL,
    last_updated DATETIME NOT NULL,
    last_updated_by TEXT NOT NULL,
    PRIMARY KEY (timer_id, version),
    FOREIGN KEY (timer_id) REFERENCES timers(id) ON DELETE CASCADE
) ENGINE = InnoDB;

-- DROP TRIGGER IF EXISTS timers_audit_insert;
CREATE TRIGGER timers_audit_insert
AFTER INSERT ON timers FOR EACH ROW
    INSERT INTO timers_audit(timer_id, comment, archived, completed, employee_id, version, last_updated, last_updated_by)
     VALUES(new.id, new.comment, new.archived, new.completed, new.employee_id, new.version, new.last_updated, new.last_updated_by);

-- DROP TRIGGER IF EXISTS timers_audit_update;
CREATE TRIGGER timers_audit_update
AFTER UPDATE ON timers FOR EACH ROW
    INSERT INTO timers_audit(timer_id, comment, archived, completed, employee_id, version, last_updated, last_updated_by)
    VALUES(new.id, new.comment, new.archived, new.completed, new.employee_id, new.version, new.last_updated, new.last_updated_by);
