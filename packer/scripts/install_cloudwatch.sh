#!/bin/bash

echo "############ Installing Cloudwatch Agent ###################"
wget https://s3.amazonaws.com/amazoncloudwatch-agent/ubuntu/amd64/latest/amazon-cloudwatch-agent.deb
sudo dpkg -i amazon-cloudwatch-agent.deb
sudo apt-get install -f
  
sudo mv /tmp/config.json /opt/aws/amazon-cloudwatch-agent/bin/config.json

sudo touch /var/log/webapp.log
sudo chown csye6225:csye6225 /var/log/webapp.log

echo "############ Installing Packages ###################"
sudo apt install -y unzip curl wget vim
sudo curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
sudo unzip awscliv2.zip
sudo ./aws/install