#!/bin/ash

# Start mysql with runtime configuration
mysqld_safe --user=root --bind-address=0.0.0.0 --skip-networking=0 --port=3306 --verbose=1 &

# Store mysql pid
mysql_pid=$!

# Ping mysql until it's up and running after starting
until mysqladmin ping >/dev/null 2>&1; do
    echo -n "."; sleep 0.2
done

# Configure data consistency
if [ "$BLUDGEON_MICROSERVICE" = "true" ]; then
    echo "Data consistency configured for microservice"
    mysql -uroot < /bludgeon/bludgeon_microservice.sql
else 
    echo "Data consistency configured for ACID"
    mysql -uroot < /bludgeon/bludgeon_acid.sql
fi

# Wait for the MySQL daemon to exit.
wait $mysql_pid