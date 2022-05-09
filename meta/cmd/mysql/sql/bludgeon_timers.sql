-- DROP DATABASE IF EXISTS bludgeon;
CREATE DATABASE IF NOT EXISTS bludgeon;

USE bludgeon;

-- DROP TABLE IF EXISTS timers;
CREATE TABLE IF NOT EXISTS timers (
    id BIGINT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    uuid VARCHAR(36) DEFAULT (UUID()),
    comment TEXT NOT NULL DEFAULT "",
    archived BOOLEAN NOT NULL DEFAULT FALSE,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    employee_id BIGINT NOT NULL,
    version INT NOT NULL DEFAULT 1,
    last_updated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    last_updated_by TEXT NOT NULL DEFAULT CURRENT_USER,
    FOREIGN KEY (employee_id) REFERENCES employees(id)
) ENGINE = InnoDB;

-- DROP TRIGGER IF EXISTS timers_audit_info_update;
CREATE TRIGGER timers_audit_info_update
BEFORE UPDATE ON timers FOR EACH ROW
    SET new.id = old.id, new.uuid = old.uuid, new.version = old.version+1, new.last_updated = CURRENT_TIMESTAMP(6), new.last_updated_by = CURRENT_USER;

-- DROP TABLE IF EXISTS timers_audit;
CREATE TABLE IF NOT EXISTS timers_audit (
    timer_id BIGINT NOT NULL,
    timer_uuid VARCHAR(36) NOT NULL,
    comment TEXT ,
    archived BOOLEAN,
    completed BOOLEAN,
    employee_id BIGINT,
    version INT NOT NULL,
    last_updated DATETIME NOT NULL,
    last_updated_by TEXT NOT NULL,
    PRIMARY KEY (timer_id, version),
    FOREIGN KEY (timer_id) REFERENCES timers(id) ON DELETE CASCADE
) ENGINE = InnoDB;

-- DROP TRIGGER IF EXISTS timers_audit_insert;
CREATE TRIGGER timers_audit_insert
AFTER INSERT ON timers FOR EACH ROW
    INSERT INTO timers_audit(timer_id, timer_uuid, comment, archived, completed, employee_id, version, last_updated, last_updated_by)
     VALUES(new.id, new.uuid, new.comment, new.archived, new.completed, new.employee_id, new.version, new.last_updated, new.last_updated_by);

-- DROP TRIGGER IF EXISTS timers_audit_update;
CREATE TRIGGER timers_audit_update
AFTER UPDATE ON timers FOR EACH ROW
    INSERT INTO timers_audit(timer_id, timer_uuid, comment, archived, completed, employee_id, version, last_updated, last_updated_by)
    VALUES(new.id, new.uuid, new.comment, new.archived, new.completed, new.employee_id, new.version, new.last_updated, new.last_updated_by);
