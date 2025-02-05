#!/bin/bash

# Load .env variables and export them
if [ -f .env ]; then
    echo "Loading environment variables from .env"
    export $(grep -v '^#' .env | xargs)
else
    echo "Error: .env file not found!"
    exit 1
fi

# Verify that variables are loaded
echo  "Database User: $DB_USER"
echo "Database Name: $DB_NAME"

# Update package list
echo "######### Updating package list ###########"
sudo apt update -y

# Upgrade installed packages
echo "########## Upgrading installed packages #############"
sudo apt upgrade -y

# Install PostgreSQL and additional utilities
echo "######### Installing PostgreSQL and additional packages ######"
sudo apt install -y postgresql postgresql-contrib unzip curl wget git vim

# Install Go
echo "############ Installing Golang ###################"
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
go version


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

# Create Linux group and user
echo "###### Creating Linux group '$LINUX_GROUP' ######"
sudo groupadd $LINUX_GROUP || echo "Group already exists."

echo "###### Creating Linux user '$LINUX_USER' #######"
sudo useradd -m -g $LINUX_GROUP -s /bin/bash $LINUX_USER || echo "User already exists."

# Unzip Application
echo "##### Creating /csye6225/ ######"
mkdir /opt/csye6225

echo "##### Unzipping #####"
unzip -o *.zip -d /opt/csye6225/

# Setting Permissions and Users

echo "####### Setting users and user groups ######"
sudo chown -R $LINUX_USER:$LINUX_GROUP /opt/csye6225/

echo "####### Setting permissions ######"
sudo chmod -R 755 /opt/csye6225/

echo "Setup complete!"

# Starting Application

source ~/.bashrc
go version

echo "####### Downloading dependencies #####"
APP_PATH=$(find /opt/csye6225/ -type d -name webapp)

cd "$APP_PATH"
go mod tidy

echo "####### Starting Application #########"
go run main.go

