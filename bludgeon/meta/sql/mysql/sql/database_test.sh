#!/bin/bash

# black box docker automation - a script to allow starting and stopping of docker containers remotely or locally
               
if [ "$#" -eq 0 ]; then
    echo "Use case: Use this to start/stop the blackbox container
    Syntax: [start|stop] [mariadb|postgres] [Port] [Password]"
    else
        if [ "$1" = "start" ]; then
            echo "...starting container"     
            if [ "$2" = "mariadb" ]; then
                docker run --name blackbox-mariadb -p $3:3306 -e MYSQL_ROOT_PASSWORD=$4 -d mariadb:latest
            elif [ "$2" = "postgres" ]; then
                docker run --name blackbox-postgres -p $3:5432 -e POSTGRES_PASSWORD=$4 -d postgres:latest
            fi
        elif [ "$1" = "status" ]; then
            if [ "$2" = "mariadb" ]; then
                docker container ls --filter name="blackbox-mariadbb"
            elif [ "$2" = "postgres" ]; then
                docker container ls --filter name="blackbox-postgres"
            fi         
        elif [ "$1" = "stop" ]; then
            echo "...stopping container"
            if [ "$2" = "mariadb" ]; then
                docker stop blackbox-mariadb
                docker wait blackbox-mariadb
                docker rm blackbox-mariadb
            elif [ "$2" = "postgres" ]; then
                docker stop blackbox-postgres
                docker wait blackbox-postgres
                docker rm blackbox-postgres
        fi
    fi
fi

## if start
# check to see if name is running
# if running, stop and clean up
# start container with configured items
## if stop
# check to see if running
# if running, stop configured items

## MariaDB
# docker run --name blackbox-mariadb -p 3306:3306 -e MYSQL_ROOT_PASSWORD=Password -d mariadb:latest
# docker rm blabox-mariadb

## PostgreSQL
# docker run --name blackbox-postgres -p 5432:5432 -e POSTGRES_PASSWORD=Password -d postgres:latest
# docker rm blackbox-postgres