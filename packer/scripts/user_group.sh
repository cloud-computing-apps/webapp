#!/bin/bash

sudo chmod +x /tmp/load_env.sh
source /tmp/load_env.sh

# Update package list
echo "######### Updating package list ###########"
sudo apt update -y

# Upgrade installed packages
echo "########## Upgrading installed packages #############"
sudo apt upgrade -y

echo "########## Uninstalling git #############"
sudo apt-get purge -y git
sudo apt-get autoremove -y

# Create Linux group and user
echo "###### Creating Linux group '$LINUX_GROUP' ######"
sudo groupadd "$LINUX_GROUP" || echo "Group already exists."

echo "###### Creating Linux user '$LINUX_USER' #######"
sudo useradd -m -g "$LINUX_GROUP" -s /usr/sbin/nologin "$LINUX_USER" || echo "User already exists."
