-- DROP DATABASE IF EXISTS bludgeon;
CREATE DATABASE IF NOT EXISTS bludgeon;

USE bludgeon;

-- DROP TABLE IF EXISTS time_slices;
CREATE TABLE IF NOT EXISTS time_slices (
    id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT (UUID()),
    start DATETIME(6) DEFAULT CURRENT_TIMESTAMP(6),
    finish DATETIME(6),
    completed BOOLEAN DEFAULT FALSE,
    elapsed_time DATETIME(6) AS (
        if(
            finish IS NOT NULL,
            TIMEDIFF(finish, start),
            TIMEDIFF(CURRENT_TIMESTAMP(6), start)
        )
    ),
    timer_id VARCHAR(36) NOT NULL,
    aux_id BIGINT AUTO_INCREMENT,
    version INT NOT NULL DEFAULT 1,
    last_updated DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    last_updated_by TEXT NOT NULL DEFAULT CURRENT_USER,
    FOREIGN KEY (timer_id) REFERENCES timers(id) ON DELETE CASCADE,
    CONSTRAINT check_start_finish CHECK (finish > start OR finish IS NULL),
    INDEX(aux_id)
) ENGINE = InnoDB;

-- DROP TRIGGER IF EXISTS time_slices_audit_info_update;
CREATE TRIGGER time_slices_audit_info_update
BEFORE UPDATE ON time_slices FOR EACH ROW
    SET new.id = old.id, new.aux_id = old.aux_id, new.version = old.version+1, new.last_updated = CURRENT_TIMESTAMP(6), new.last_updated_by = CURRENT_USER;

-- DROP TRIGGER validate_time_slice_start_insert;
DELIMITER $$
CREATE TRIGGER validate_time_slice_start_insert
BEFORE INSERT
    ON time_slices FOR EACH ROW BEGIN
        IF (SELECT COUNT(*) FROM (SELECT id, start, finish FROM (SELECT id, start, finish FROM time_slices WHERE timer_id = new.timer_id ) AS timer_time_slices WHERE new.start BETWEEN timer_time_slices.start AND timer_time_slices.finish) AS conflict_time_slices) > 0
        THEN
            SIGNAL SQLSTATE '45000'
                SET MESSAGE_TEXT = 'Cannot insert time slice, start conflicts with existing time slices';
        END IF;
END$$
DELIMITER ;

-- DROP TRIGGER validate_time_slice_start_update;
DELIMITER $$
CREATE TRIGGER validate_time_slice_start_update
BEFORE INSERT
    ON time_slices FOR EACH ROW BEGIN
        IF (SELECT COUNT(*) FROM (SELECT id, start, finish FROM (SELECT id, start, finish FROM time_slices WHERE timer_id = new.timer_id ) AS timer_time_slices WHERE new.start BETWEEN timer_time_slices.start AND timer_time_slices.finish AND id <> new.id) AS conflict_time_slices) > 0
        THEN
            SIGNAL SQLSTATE '45000'
                SET MESSAGE_TEXT = 'Cannot update time slice, start conflicts with existing time slices';
        END IF;
END$$
DELIMITER ;

-- DROP TRIGGER validate_active_time_slice_insert;
DELIMITER $$
CREATE TRIGGER validate_active_time_slice_insert
BEFORE INSERT
    ON time_slices FOR EACH ROW BEGIN
        IF (SELECT COUNT(*) FROM (SELECT id FROM time_slices WHERE timer_id = new.timer_id AND finish IS NULL) AS validate_active_time_slice) > 0
        THEN
            SIGNAL SQLSTATE '45000'
                SET MESSAGE_TEXT = 'Cannot insert time slice, active time slice already exists for timer';
        END IF;
END$$
DELIMITER ;

-- DROP TRIGGER validate_active_time_slice_update;
DELIMITER $$
CREATE TRIGGER validate_active_time_slice_update
BEFORE UPDATE
    ON time_slices FOR EACH ROW BEGIN
        IF (SELECT COUNT(*) FROM (SELECT id FROM time_slices WHERE timer_id = new.timer_id AND finish IS NULL AND id <> new.id) AS validate_active_time_slice) > 0
        THEN
            SIGNAL SQLSTATE '45000'
                SET MESSAGE_TEXT = 'Cannot update time slice, active time slice already exists for timer';
        END IF;
END$$
DELIMITER ;

-- DROP TABLE IF EXISTS time_slices_audit;
CREATE TABLE IF NOT EXISTS time_slices_audit (
    time_slice_id VARCHAR(36) NOT NULL,
    start DATETIME(6),
    finish DATETIME(6),
    completed BOOLEAN DEFAULT false,
    elapsed_time DATETIME(6) AS (
        IF (
            finish IS NOT NULL,
            TIMEDIFF(finish, start),
            TIMEDIFF(CURRENT_TIMESTAMP(6), start)
        )
    ),
    timer_id VARCHAR(36),
    version INT NOT NULL,
    last_updated DATETIME(6) NOT NULL,
    last_updated_by TEXT NOT NULL,
    PRIMARY KEY (time_slice_id, version),
    FOREIGN KEY (time_slice_id) REFERENCES time_slices(id) ON DELETE CASCADE
) ENGINE = InnoDB;

-- DROP TRIGGER IF EXISTS time_slices_audit_insert;
CREATE TRIGGER time_slices_audit_insert
AFTER INSERT ON time_slices FOR EACH ROW
    INSERT INTO time_slices_audit(time_slice_id, start, finish, completed, timer_id, version, last_updated, last_updated_by)
    VALUES(new.id, new.start, new.finish, new.completed, new.timer_id, new.version, new.last_updated, new.last_updated_by);

-- DROP TRIGGER IF EXISTS time_slices_audit_update;
CREATE TRIGGER time_slices_audit_update
AFTER UPDATE ON time_slices FOR EACH ROW
    INSERT INTO time_slices_audit(time_slice_id, start, finish, completed, timer_id, version, last_updated, last_updated_by)
    VALUES(new.id, new.start, new.finish, new.completed, new.timer_id, new.version, new.last_updated, new.last_updated_by);
