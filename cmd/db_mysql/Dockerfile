#build from mariadb container
FROM mariadb:latest 
#set work directory
WORKDIR /sql
#copy current directory into the gopath
COPY ./cmd/db_mysql/bludgeon_mysql.sql /sql
#install delve and build the hell-http app
