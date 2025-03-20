#!/bin/bash

# Create folder
sudo mkdir -p /opt/csye6225/

# Move binary from tmp
sudo mv /tmp/webapp /opt/csye6225/webapp

# Setting Permissions and Users
echo "####### Setting users and user groups ######"
sudo chown -R csye6225:csye6225 /opt/csye6225/

echo "####### Setting permissions ######"
sudo chmod -R 755 /opt/csye6225/

