# start from the latest postgres
FROM postgres:alpine

# copy bludgeon sql files to entrypoint to execute on first run 
COPY ./bludgeon_postgres.sql /docker-entrypoint-initdb.d

# postgres port
EXPOSE 5432
