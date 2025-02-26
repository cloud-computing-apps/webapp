#!/bin/bash

source /tmp/load_env.sh

# Install PostgreSQL and additional utilities
echo "######### Installing PostgreSQL and additional packages ######"
sudo apt install -y postgresql postgresql-contrib unzip curl wget vim

# Enable and start PostgreSQL service
echo "########### Enabling PostgreSQL service #########"
sudo systemctl enable --now postgresql

# Set PostgreSQL password
echo "##### Setting password ######"
export PGPASSWORD="$DB_PASSWORD"

sudo -u postgres psql -c "ALTER USER $DB_USER WITH PASSWORD '$DB_PASSWORD';"

# Modify pg_hba.conf
PG_HBA="/etc/postgresql/16/main/pg_hba.conf"
sudo sed -i "s/peer/md5/g" $PG_HBA
sudo systemctl restart postgresql

# Create DB and user
psql -U "$DB_USER" -h localhost -d postgres -tc "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME';" | grep -q 1 || psql -U "$DB_USER" -h localhost -d postgres -c "CREATE DATABASE $DB_NAME;"

unset PGPASSWORD
