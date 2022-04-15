-- DROP DATABASE IF EXISTS bludgeon;
CREATE DATABASE IF NOT EXISTS bludgeon;

USE bludgeon;

-- DROP TABLE IF EXISTS employees;
CREATE TABLE IF NOT EXISTS employees (
    id BIGINT PRIMARY KEY NOT NULL AUTO_INCREMENT,
    uuid VARCHAR(36) DEFAULT (UUID()),
    first_name TEXT DEFAULT '',
    last_name TEXT DEFAULT '',
    email_address TEXT NOT NULL,
    version INT NOT NULL DEFAULT 1,
    last_updated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    last_updated_by TEXT NOT NULL DEFAULT CURRENT_USER,
    UNIQUE(uuid),
    UNIQUE(email_address)
) ENGINE = InnoDB;

-- DROP TRIGGER IF EXISTS employees_audit_info_update;
CREATE TRIGGER employees_audit_info_update
BEFORE UPDATE ON employees FOR EACH ROW
    SET new.id = old.id, new.uuid = old.uuid, new.version = old.version+1, new.last_updated = CURRENT_TIMESTAMP(6), new.last_updated_by = CURRENT_USER;

-- DROP TABLE IF EXISTS employees_audit;
CREATE TABLE IF NOT EXISTS employees_audit (
    employee_id BIGINT NOT NULL,
    employee_uuid VARCHAR(36) NOT NULL,
    first_name TEXT,
    last_name TEXT,
    email_address TEXT,
    version INT NOT NULL,
    last_updated DATETIME NOT NULL,
    last_updated_by TEXT NOT NULL,
    PRIMARY KEY (employee_id, version),
    FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
) ENGINE = InnoDB;

-- DROP TRIGGER IF EXISTS employees_audit_insert;
CREATE TRIGGER employees_audit_insert
AFTER INSERT ON employees FOR EACH ROW
    INSERT INTO employees_audit(employee_id, employee_uuid, first_name, last_name, email_address, version, last_updated, last_updated_by)
     VALUES (new.id, new.uuid, new.first_name,  new.last_name, new.email_address, new.version, new.last_updated, new.last_updated_by);

-- DROP TRIGGER IF EXISTS employees_audit_update;
CREATE TRIGGER employees_audit_update
AFTER UPDATE ON employees FOR EACH ROW
    INSERT INTO employees_audit(employee_id, employee_uuid, first_name, last_name, email_address, version, last_updated, last_updated_by)
     VALUES(new.id, new.uuid, new.first_name,  new.last_name, new.email_address, new.version, new.last_updated, new.last_updated_by);
