# MariaDB on Alpine Linux
FROM alpine:latest
#set working directory
WORKDIR /var/lib/mysql
#copy sql file and database init script
COPY ./sql/bludgeon_mysql.sql /bludgeon/bludgeon_mysql.sql
COPY ./sql/init_bludgeon_db.sh /bludgeon/init_bludgeon_db.sh
# MariaDB server and client
RUN \
    apk update \
    && apk add mariadb mariadb-client \
    && mkdir -p /run/mysqld \
    && chown -R root:root /run/mysqld \
    && mkdir -p /var/lib/mysql \
    && chown -R root:root /var/lib/mysql \
    && mysql_install_db --user=root --ldata=/var/lib/mysql \
    && sh /bludgeon/init_bludgeon_db.sh
#start the entrypoint
CMD  mysqld_safe --user=root --verbose=0
# MariaDB port
EXPOSE 3306
