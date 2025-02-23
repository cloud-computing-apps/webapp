#!/bin/bash

source /tmp/load_env.sh

# Create folder
sudo mkdir -p /opt/csye6225/

# Move binary from tmp
sudo mv /tmp/webapp /opt/csye6225/webapp
sudo mv /tmp/.env /opt/csye6225/.env

# Setting Permissions and Users
echo "####### Setting users and user groups ######"
sudo chown -R "$LINUX_USER":"$LINUX_GROUP" /opt/csye6225/

echo "####### Setting permissions ######"
sudo chmod -R 755 /opt/csye6225/

