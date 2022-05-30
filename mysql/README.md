# MySQL (github.com/antonio-alexander/go-bludgeon/mysql)

MySQL provides a common database for bludgeon. It contains tables, views, and relationships for timers and employees, as well as security providing access to the user bludgeon. Although the expectation is that you'd deploy a "single" database with multiple tables. See the section [Data Consistency](#data-consistency). for more information regarding relationships between the different tables.

## Getting Started

To build the mysql image, you can use one of the following commands from the root of /mysql:

```sh
docker compose build --no-cache
```

```sh
 docker build --no-cache -f ./cmd/Dockerfile . -t ghcr.io/antonio-alexander/go-bludgeon-mysql:amd64_latest
  --build-arg GIT_COMMIT=$GITHUB_SHA --build-arg GIT_BRANCH=$GITHUB_REF --build-arg PLATFORM=linux/amd64
 docker build --no-cache -f ./cmd/Dockerfile . -t ghcr.io/antonio-alexander/go-bludgeon-mysql:armv7_latest
  --build-arg GIT_COMMIT=$GITHUB_SHA --build-arg GIT_BRANCH=$GITHUB_REF --build-arg PLATFORM=linux/amd64
```

Keep in mind, that if you don't build without cache, the way in which the Dockerfile is configured may cause you to have two images that are on the same architecture rather than two different architectures. Alternatively, you could pull the most recent image with the following command:

```sh
docker pull ghcr.io/antonio-alexander/go-bludgeon-mysql:latest
```

Once the image is available locally, you can bring the mysql database up with:

```sh
docker compose up -d
```

Once the database is up, you can interact with it using the following commands:

```sh
docker exec -it mysql /bin/ash
```

```sh
docker exec -it mysql mysql -u root bludgeon
```

Once up and running you can interact with the tables directly or use the mysql meta api for employees or timers.

## Data Consistency

With any kidn of microservice architecture and databases, there are generally two ways to enforce data consistency between different services. Often the microservice paradigm is that every microservice has a completely separate database, both logically and sometimes instance wise (if you want to take it there).

So far, there's only one logical connection between the different microservices: there is a one-to-many (1:N) relationship between employees and timers. In terms of data consistency, consider the following situations:

- What if an employee is deleted and there are timers associated with it?
- How can you verify if an associated employee is valid?

Data consistency can be enforced in an ACID way using a foreign key constraint between the employee_id column of timers and the id column of the employees tables. If BLUDGEON_MICROSERVICE is set to true, this foreign key will be removed/not present, while if its set to false or not present, the foreign key constraint will be added.

This can be verified with the following:

```sh
docker exec -it mysql mysql -u root bludgeon
MariaDB [bludgeon]> SELECT COLUMN_NAME, CONSTRAINT_NAME, REFERENCED_COLUMN_NAME, REFERENCED_TABLE_NAME
    ->     FROM information_schema.KEY_COLUMN_USAGE WHERE TABLE_NAME = 'timers'; 
```

If BLUDGEON_MICROSERVICE is false:

```log
+-------------+-----------------+------------------------+-----------------------+
| COLUMN_NAME | CONSTRAINT_NAME | REFERENCED_COLUMN_NAME | REFERENCED_TABLE_NAME |
+-------------+-----------------+------------------------+-----------------------+
| id          | PRIMARY         | NULL                   | NULL                  |
| employee_id | fk_employee_id  | id                     | employees             |
+-------------+-----------------+------------------------+-----------------------+
2 rows in set (0.001 sec)
```

If BLUDGEON_MICROSERVICE is true:

```log
+-------------+-----------------+------------------------+-----------------------+
| COLUMN_NAME | CONSTRAINT_NAME | REFERENCED_COLUMN_NAME | REFERENCED_TABLE_NAME |
+-------------+-----------------+------------------------+-----------------------+
| id          | PRIMARY         | NULL                   | NULL                  |
+-------------+-----------------+------------------------+-----------------------+
1 row in set (0.001 sec)
```
