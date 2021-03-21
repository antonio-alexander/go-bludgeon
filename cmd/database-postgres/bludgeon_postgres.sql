-- Database generated with pgModeler (PostgreSQL Database Modeler).
-- pgModeler  version: 0.9.3
-- PostgreSQL version: 13.0
-- Project Site: pgmodeler.io
-- Model Author: ---
-- object: bludgeon | type: ROLE --
-- DROP ROLE IF EXISTS bludgeon;
CREATE ROLE bludgeon WITH ;
-- ddl-end --


-- Database creation must be performed outside a multi lined SQL file. 
-- These commands were put in this file only as a convenience.
-- 
-- object: bludgeon | type: DATABASE --
-- DROP DATABASE IF EXISTS bludgeon;
CREATE DATABASE bludgeon
	OWNER = postgres;
-- ddl-end --


-- object: bludgeon | type: SCHEMA --
-- DROP SCHEMA IF EXISTS bludgeon CASCADE;
CREATE SCHEMA bludgeon;
-- ddl-end --
ALTER SCHEMA bludgeon OWNER TO bludgeon;
-- ddl-end --

SET search_path TO pg_catalog,public,bludgeon;
-- ddl-end --

-- object: bludgeon.timer | type: TABLE --
-- DROP TABLE IF EXISTS bludgeon.timer CASCADE;
CREATE TABLE bludgeon.timer (
	timer_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT BY 1 MINVALUE 0 MAXVALUE 9223372036854775807 START WITH 1 CACHE 1 ),
	timer_uuid char(36) NOT NULL,
	timer_start bigint NOT NULL,
	timer_finish bigint,
	timer_comment text,
	timer_archived bool NOT NULL DEFAULT FALSE,
	timer_billed bool NOT NULL DEFAULT FALSE,
	timer_completed bool NOT NULL DEFAULT FALSE,
	task_id bigint,
	employee_id bigint NOT NULL,
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
	timer_id bigint NOT NULL,
	slice_start bigint NOT NULL,
	slice_finish bigint,
	slice_archived bool NOT NULL DEFAULT FALSE,
	slice_elapsed_time bigint GENERATED ALWAYS AS ((slice_finish-slice_start)) STORED,
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
	client_name text,
	CONSTRAINT client_pk PRIMARY KEY (client_id),
	CONSTRAINT client_uuid_unique UNIQUE (client_uuid)

);
-- ddl-end --
ALTER TABLE bludgeon.client OWNER TO postgres;
-- ddl-end --

-- object: bludgeon.employee | type: TABLE --
-- DROP TABLE IF EXISTS bludgeon.employee CASCADE;
CREATE TABLE bludgeon.employee (
	employee_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT BY 1 MINVALUE 0 MAXVALUE 9223372036854775807 START WITH 1 CACHE 1 ),
	employee_uuid char(36) NOT NULL,
	employee_first_name text NOT NULL,
	employee_last_name text NOT NULL,
	CONSTRAINT employee_pk PRIMARY KEY (employee_id),
	CONSTRAINT employee_uuid_unique UNIQUE (employee_uuid)

);
-- ddl-end --
ALTER TABLE bludgeon.employee OWNER TO postgres;
-- ddl-end --

-- object: bludgeon.project | type: TABLE --
-- DROP TABLE IF EXISTS bludgeon.project CASCADE;
CREATE TABLE bludgeon.project (
	project_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT BY 1 MINVALUE 0 MAXVALUE 9223372036854775807 START WITH 1 CACHE 1 ),
	project_uuid char(36) NOT NULL,
	project_description text,
	CONSTRAINT project_pk PRIMARY KEY (project_id),
	CONSTRAINT project_uuid_unique UNIQUE (project_uuid)

);
-- ddl-end --
ALTER TABLE bludgeon.project OWNER TO postgres;
-- ddl-end --

-- object: bludgeon.project_client | type: TABLE --
-- DROP TABLE IF EXISTS bludgeon.project_client CASCADE;
CREATE TABLE bludgeon.project_client (
	project_id bigint NOT NULL,
	client_id bigint NOT NULL,
	CONSTRAINT project_client_pk PRIMARY KEY (project_id,client_id)

);
-- ddl-end --
ALTER TABLE bludgeon.project_client OWNER TO postgres;
-- ddl-end --

-- object: bludgeon.timer_slice_active | type: TABLE --
-- DROP TABLE IF EXISTS bludgeon.timer_slice_active CASCADE;
CREATE TABLE bludgeon.timer_slice_active (
	slice_id bigint NOT NULL,
	timer_id bigint NOT NULL,
	CONSTRAINT time_slice_active_pk PRIMARY KEY (slice_id,timer_id),
	CONSTRAINT timer_id_unique UNIQUE (timer_id)

);
-- ddl-end --
ALTER TABLE bludgeon.timer_slice_active OWNER TO postgres;
-- ddl-end --

-- object: client_project_idx | type: INDEX --
-- DROP INDEX IF EXISTS bludgeon.client_project_idx CASCADE;
CREATE INDEX client_project_idx ON bludgeon.project_client
	USING btree
	(
	  client_id,
	  project_id
	);
-- ddl-end --

-- object: timer_slice_idx | type: INDEX --
-- DROP INDEX IF EXISTS bludgeon.timer_slice_idx CASCADE;
CREATE INDEX timer_slice_idx ON bludgeon.timer_slice_active
	USING btree
	(
	  timer_id,
	  slice_id
	);
-- ddl-end --

-- object: bludgeon.task | type: TABLE --
-- DROP TABLE IF EXISTS bludgeon.task CASCADE;
CREATE TABLE bludgeon.task (
	task_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ,
	task_uuid char(36) NOT NULL,
	project_id bigint NOT NULL,
	task_description text,
	CONSTRAINT task_pk PRIMARY KEY (task_id),
	CONSTRAINT task_uuid_unique UNIQUE (task_uuid)

);
-- ddl-end --
ALTER TABLE bludgeon.task OWNER TO postgres;
-- ddl-end --

-- object: task_id_fk | type: CONSTRAINT --
-- ALTER TABLE bludgeon.timer DROP CONSTRAINT IF EXISTS task_id_fk CASCADE;
ALTER TABLE bludgeon.timer ADD CONSTRAINT task_id_fk FOREIGN KEY (task_id)
REFERENCES bludgeon.task (task_id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: employee_id_fk | type: CONSTRAINT --
-- ALTER TABLE bludgeon.timer DROP CONSTRAINT IF EXISTS employee_id_fk CASCADE;
ALTER TABLE bludgeon.timer ADD CONSTRAINT employee_id_fk FOREIGN KEY (employee_id)
REFERENCES bludgeon.employee (employee_id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: timer_id_fk | type: CONSTRAINT --
-- ALTER TABLE bludgeon.slice DROP CONSTRAINT IF EXISTS timer_id_fk CASCADE;
ALTER TABLE bludgeon.slice ADD CONSTRAINT timer_id_fk FOREIGN KEY (timer_id)
REFERENCES bludgeon.timer (timer_id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: project_id_fk | type: CONSTRAINT --
-- ALTER TABLE bludgeon.project_client DROP CONSTRAINT IF EXISTS project_id_fk CASCADE;
ALTER TABLE bludgeon.project_client ADD CONSTRAINT project_id_fk FOREIGN KEY (project_id)
REFERENCES bludgeon.project (project_id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: client_id_fk | type: CONSTRAINT --
-- ALTER TABLE bludgeon.project_client DROP CONSTRAINT IF EXISTS client_id_fk CASCADE;
ALTER TABLE bludgeon.project_client ADD CONSTRAINT client_id_fk FOREIGN KEY (client_id)
REFERENCES bludgeon.client (client_id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: timer_id_fk | type: CONSTRAINT --
-- ALTER TABLE bludgeon.timer_slice_active DROP CONSTRAINT IF EXISTS timer_id_fk CASCADE;
ALTER TABLE bludgeon.timer_slice_active ADD CONSTRAINT timer_id_fk FOREIGN KEY (timer_id)
REFERENCES bludgeon.timer (timer_id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: slice_id_fk | type: CONSTRAINT --
-- ALTER TABLE bludgeon.timer_slice_active DROP CONSTRAINT IF EXISTS slice_id_fk CASCADE;
ALTER TABLE bludgeon.timer_slice_active ADD CONSTRAINT slice_id_fk FOREIGN KEY (slice_id)
REFERENCES bludgeon.slice (slice_id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --

-- object: project_id_fk | type: CONSTRAINT --
-- ALTER TABLE bludgeon.task DROP CONSTRAINT IF EXISTS project_id_fk CASCADE;
ALTER TABLE bludgeon.task ADD CONSTRAINT project_id_fk FOREIGN KEY (project_id)
REFERENCES bludgeon.project (project_id) MATCH FULL
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- ddl-end --


