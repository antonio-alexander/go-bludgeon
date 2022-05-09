#!/bin/ash

# Start the MySQL daemon in the background.
cd '/usr' ; mysqld --user=root --datadir='/var/lib/mysql/' &
mysql_pid=$!

#ping mysql until it's up and running after starting
until mysqladmin ping >/dev/null 2>&1; do
    echo -n "."; sleep 0.2
done

#load the configuration into the sql database
mysql -uroot < /bludgeon/bludgeon_security.sql
mysql -uroot < /bludgeon/bludgeon_employees.sql
mysql -uroot < /bludgeon/bludgeon_timers.sql
mysql -uroot < /bludgeon/bludgeon_time_slices.sql
mysql -uroot < /bludgeon/bludgeon_views.sql

# Tell the MySQL daemon to shutdown.
mysqladmin shutdown

# Wait for the MySQL daemon to exit.
wait $mysql_pid