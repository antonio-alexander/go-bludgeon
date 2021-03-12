-- Database generated with pgModeler (PostgreSQL Database Modeler).
-- pgModeler  version: 0.9.3
-- PostgreSQL version: 12.0
-- Project Site: pgmodeler.io
-- Model Author: ---

-- Database creation must be performed outside a multi lined SQL file. 
-- These commands were put in this file only as a convenience.
-- 
-- object: new_database | type: DATABASE --
-- DROP DATABASE IF EXISTS new_database;
CREATE DATABASE new_database;
-- ddl-end --


-- object: bludgeon | type: SCHEMA --
-- DROP SCHEMA IF EXISTS bludgeon CASCADE;
CREATE SCHEMA bludgeon;
-- ddl-end --
ALTER SCHEMA bludgeon OWNER TO postgres;
-- ddl-end --

SET search_path TO pg_catalog,public,bludgeon;
-- ddl-end --

-- object: bludgeon.timer | type: TABLE --
-- DROP TABLE IF EXISTS bludgeon.timer CASCADE;
CREATE TABLE bludgeon.timer (
	timer_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT BY 1 MINVALUE 0 MAXVALUE 9223372036854775807 START WITH 1 CACHE 1 ),
	timer_uuid char(36) NOT NULL,
	timestamp_start bigint NOT NULL,
	timestamp_finish smallint,
	comment text,
	CONSTRAINT timer_pk PRIMARY KEY (timer_id),
	CONSTRAINT timer_uuid_unique UNIQUE (timer_uuid)

);
-- ddl-end --
ALTER TABLE bludgeon.timer OWNER TO postgres;
-- ddl-end --

-- object: bludgeon.slice | type: TABLE --
-- DROP TABLE IF EXISTS bludgeon.slice CASCADE;
CREATE TABLE bludgeon.slice (
	slice_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ,
	slice_uuid char(36) NOT NULL,
	timestamp_start bigint NOT NULL,
	timestamp_finish bigint,
	CONSTRAINT slice_pk PRIMARY KEY (slice_id),
	CONSTRAINT slice_uuid_unique UNIQUE (slice_uuid)

);
-- ddl-end --
ALTER TABLE bludgeon.slice OWNER TO postgres;
-- ddl-end --

-- object: bludgeon.client | type: TABLE --
-- DROP TABLE IF EXISTS bludgeon.client CASCADE;
CREATE TABLE bludgeon.client (
	client_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT BY 1 MINVALUE 0 MAXVALUE 9223372036854775807 START WITH 1 CACHE 1 ),
	client_uuid char(36) NOT NULL,
	name text,
	rate float,
	CONSTRAINT client_pk PRIMARY KEY (client_id)

);
-- ddl-end --
ALTER TABLE bludgeon.client OWNER TO postgres;
-- ddl-end --

-- object: bludgeon.employee | type: TABLE --
-- DROP TABLE IF EXISTS bludgeon.employee CASCADE;
CREATE TABLE bludgeon.employee (
	employee_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT BY 1 MINVALUE 0 MAXVALUE 9223372036854775807 START WITH 1 CACHE 1 ),
	first_name text NOT NULL,
	last_name text NOT NULL,
	CONSTRAINT employee_pk PRIMARY KEY (employee_id)

);
-- ddl-end --
ALTER TABLE bludgeon.employee OWNER TO postgres;
-- ddl-end --

-- object: bludgeon.project | type: TABLE --
-- DROP TABLE IF EXISTS bludgeon.project CASCADE;
CREATE TABLE bludgeon.project (
	project_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT BY 1 MINVALUE 0 MAXVALUE 9223372036854775807 START WITH 1 CACHE 1 ),
	project_uuid char(36) NOT NULL,
	description text,
	CONSTRAINT project_pk PRIMARY KEY (project_id)

);
-- ddl-end --
ALTER TABLE bludgeon.project OWNER TO postgres;
-- ddl-end --

-- object: bludgeon.time_slice | type: TABLE --
-- DROP TABLE IF EXISTS bludgeon.time_slice CASCADE;
CREATE TABLE bludgeon.time_slice (
	timer_id bigint NOT NULL,
	slice_id bigint NOT NULL
);
-- ddl-end --
ALTER TABLE bludgeon.time_slice OWNER TO postgres;
-- ddl-end --

-- object: bludgeon.timer_elapsed | type: TABLE --
-- DROP TABLE IF EXISTS bludgeon.timer_elapsed CASCADE;
CREATE TABLE bludgeon.timer_elapsed (
	timer_id bigint NOT NULL,
	elapsed_time bigint,
	CONSTRAINT timer_elapsed_pk PRIMARY KEY (timer_id)

);
-- ddl-end --
ALTER TABLE bludgeon.timer_elapsed OWNER TO postgres;
-- ddl-end --

-- object: bludgeon.timer_client | type: TABLE --
-- DROP TABLE IF EXISTS bludgeon.timer_client CASCADE;
CREATE TABLE bludgeon.timer_client (
	timer_id bigint NOT NULL,
	client_id bigint NOT NULL,
	CONSTRAINT timer_client_pk PRIMARY KEY (timer_id,client_id)

);
-- ddl-end --
ALTER TABLE bludgeon.timer_client OWNER TO postgres;
-- ddl-end --

-- object: bludgeon.timer_employee | type: TABLE --
-- DROP TABLE IF EXISTS bludgeon.timer_employee CASCADE;
CREATE TABLE bludgeon.timer_employee (
	timer_id bigint NOT NULL,
	employee_id bigint NOT NULL,
	CONSTRAINT timer_employee_pk PRIMARY KEY (timer_id,employee_id)

);
-- ddl-end --
ALTER TABLE bludgeon.timer_employee OWNER TO postgres;
-- ddl-end --

-- object: bludgeon.project_timer | type: TABLE --
-- DROP TABLE IF EXISTS bludgeon.project_timer CASCADE;
CREATE TABLE bludgeon.project_timer (
	timer_id bigint NOT NULL,
	project_id bigint NOT NULL,
	CONSTRAINT project_timer_pk PRIMARY KEY (timer_id,project_id)

);
-- ddl-end --
ALTER TABLE bludgeon.project_timer OWNER TO postgres;
-- ddl-end --

-- object: timer_id | type: CONSTRAINT --
-- ALTER TABLE bludgeon.time_slice DROP CONSTRAINT IF EXISTS timer_id CASCADE;
ALTER TABLE bludgeon.time_slice ADD CONSTRAINT timer_id FOREIGN KEY (timer_id)
REFERENCES bludgeon.timer (timer_id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: slice_id | type: CONSTRAINT --
-- ALTER TABLE bludgeon.time_slice DROP CONSTRAINT IF EXISTS slice_id CASCADE;
ALTER TABLE bludgeon.time_slice ADD CONSTRAINT slice_id FOREIGN KEY (slice_id)
REFERENCES bludgeon.slice (slice_id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: timer_id | type: CONSTRAINT --
-- ALTER TABLE bludgeon.timer_elapsed DROP CONSTRAINT IF EXISTS timer_id CASCADE;
ALTER TABLE bludgeon.timer_elapsed ADD CONSTRAINT timer_id FOREIGN KEY (timer_id)
REFERENCES bludgeon.timer (timer_id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: timer_id | type: CONSTRAINT --
-- ALTER TABLE bludgeon.timer_client DROP CONSTRAINT IF EXISTS timer_id CASCADE;
ALTER TABLE bludgeon.timer_client ADD CONSTRAINT timer_id FOREIGN KEY (timer_id)
REFERENCES bludgeon.timer (timer_id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: client_id | type: CONSTRAINT --
-- ALTER TABLE bludgeon.timer_client DROP CONSTRAINT IF EXISTS client_id CASCADE;
ALTER TABLE bludgeon.timer_client ADD CONSTRAINT client_id FOREIGN KEY (client_id)
REFERENCES bludgeon.client (client_id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: timer_id | type: CONSTRAINT --
-- ALTER TABLE bludgeon.timer_employee DROP CONSTRAINT IF EXISTS timer_id CASCADE;
ALTER TABLE bludgeon.timer_employee ADD CONSTRAINT timer_id FOREIGN KEY (timer_id)
REFERENCES bludgeon.timer (timer_id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: employee_id | type: CONSTRAINT --
-- ALTER TABLE bludgeon.timer_employee DROP CONSTRAINT IF EXISTS employee_id CASCADE;
ALTER TABLE bludgeon.timer_employee ADD CONSTRAINT employee_id FOREIGN KEY (employee_id)
REFERENCES bludgeon.employee (employee_id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: project_id | type: CONSTRAINT --
-- ALTER TABLE bludgeon.project_timer DROP CONSTRAINT IF EXISTS project_id CASCADE;
ALTER TABLE bludgeon.project_timer ADD CONSTRAINT project_id FOREIGN KEY (project_id)
REFERENCES bludgeon.project (project_id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: timer_id | type: CONSTRAINT --
-- ALTER TABLE bludgeon.project_timer DROP CONSTRAINT IF EXISTS timer_id CASCADE;
ALTER TABLE bludgeon.project_timer ADD CONSTRAINT timer_id FOREIGN KEY (project_id)
REFERENCES bludgeon.timer (timer_id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --


