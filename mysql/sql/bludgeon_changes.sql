-- DROP DATABASE IF EXISTS bludgeon;
CREATE DATABASE IF NOT EXISTS bludgeon;

USE bludgeon;

-- DROP TABLE IF EXISTS changes;
CREATE TABLE IF NOT EXISTS changes (
    id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (UUID()),
    aux_id BIGINT AUTO_INCREMENT,
    data_id VARCHAR(36) NOT NULL,
    version INT NOT NULL DEFAULT 1,
    type TEXT,
    service TEXT,
    action TEXT,
    when_changed DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    changed_by TEXT NOT NULL DEFAULT CURRENT_USER,
    INDEX(aux_id),
    UNIQUE(data_id, version, type, service, action)
) ENGINE = InnoDB;

-- DROP TABLE IF EXISTS registrations;
CREATE TABLE IF NOT EXISTS registrations (
    id VARCHAR(36) PRIMARY KEY NOT NULL,
    aux_id BIGINT AUTO_INCREMENT,
    INDEX(aux_id)
) ENGINE = InnoDB;

-- DROP TABLE IF EXISTS registration_changes;
CREATE TABLE IF NOT EXISTS registration_changes (
    registration_id VARCHAR(36) NOT NULL,
    change_id VARCHAR(36) NOT NULL,
    FOREIGN KEY (registration_id)
        REFERENCES registrations(id)
        ON DELETE CASCADE,
    FOREIGN KEY (change_id)
        REFERENCES changes(id),
    PRIMARY KEY(registration_id, change_id)
) ENGINE = InnoDB;