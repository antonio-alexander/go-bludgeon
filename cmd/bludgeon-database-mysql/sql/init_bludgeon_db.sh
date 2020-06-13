#!/bin/ash
# I legit just combined some scripts I found online and got this to work
# I'm not even a little practiced with shell-scrits, enjoy.

# Start the MySQL daemon in the background.
mysqld --user=root &
mysql_pid=$!

#ping mysql until it's up and running after starting
until mysqladmin ping >/dev/null 2>&1; do
  echo -n "."; sleep 0.2
done

#load the configuration into the sql database
mysql -u root < /bludgeon/bludgeon_mysql.sql

# Tell the MySQL daemon to shutdown.
mysqladmin shutdown

# Wait for the MySQL daemon to exit.
wait $mysql_pid