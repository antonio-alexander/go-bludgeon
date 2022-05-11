-- DROP DATABASE IF EXISTS bludgeon;
CREATE DATABASE IF NOT EXISTS bludgeon;

USE bludgeon;

-- DROP VIEW IF EXISTS employees_v1;
CREATE VIEW employees_v1 AS
SELECT
    id AS employee_id,
    first_name,
    last_name,
    email_address,
    version,
    UNIX_TIMESTAMP(last_updated) AS last_updated,
    last_updated_by
FROM
    employees;

-- DROP VIEW IF EXISTS timers_v1;
CREATE VIEW timers_v1 AS
SELECT
    id as timer_id,
    (SELECT MIN(UNIX_TIMESTAMP(start)) FROM time_slices WHERE timer_id = timers.id) AS start,
    IF(timers.completed, (SELECT MAX(UNIX_TIMESTAMP(finish)) FROM time_slices WHERE timer_id = timers.id), NULL ) AS finish,
    (SELECT SUM(TIME_TO_SEC(elapsed_time)) FROM time_slices WHERE timer_id = timers.id) AS elapsed_time,
    comment,
    archived,
    completed,
    timers.employee_id AS employee_id,
    (SELECT id FROM time_slices WHERE finish IS NULL AND timer_id = timers.id) AS active_time_slice_id,
    version,
    UNIX_TIMESTAMP(last_updated) AS last_updated,
    last_updated_by
FROM 
    timers;

-- DROP VIEW IF EXISTS time_slices_v1;
CREATE VIEW time_slices_v1 AS
SELECT
    id AS time_slice_id,
    UNIX_TIMESTAMP(start) AS start,
    UNIX_TIMESTAMP(finish) AS finish,
    completed,
    (SELECT SUM(TIME_TO_SEC(elapsed_time))) AS elapsed_time,
    (SELECT id FROM timers WHERE id = timer_id) AS timer_id,
    version,
    UNIX_TIMESTAMP(last_updated) AS last_updated,
    last_updated_by
FROM
    time_slices;