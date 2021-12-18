# start from the latest mariadb (latest or alpine)
FROM yobasystems/alpine-mariadb:latest

# copy bludgeon sql files to entrypoint to execute on first run 
COPY ./bludgeon_mysql.sql /docker-entrypoint-initdb.d/

# MariaDB port
EXPOSE 3306
