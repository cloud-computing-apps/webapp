#!/bin/bash

echo "############ Installing Cloudwatch Agent ###################"
wget https://s3.amazonaws.com/amazoncloudwatch-agent/ubuntu/amd64/latest/amazon-cloudwatch-agent.deb
sudo dpkg -i amazon-cloudwatch-agent.deb
sudo apt-get install -f
  
sudo mv /tmp/config.json /opt/aws/amazon-cloudwatch-agent/bin/config.json

sudo touch /var/log/webapp.log
sudo chown csye6225:csye6225 /var/log/webapp.log