#!/bin/bash

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
echo "###### Creating Linux group csye6225 ######"
sudo groupadd csye6225 || echo "Group already exists."

echo "###### Creating Linux user csye6225 #######"
sudo useradd -m -g csye6225 -s /usr/sbin/nologin csye6225 || echo "User already exists."
