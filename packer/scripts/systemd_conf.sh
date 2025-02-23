#!/bin/bash
sudo mv /tmp/webapp.service /etc/systemd/system/webapp.service

sudo systemctl daemon-reload
sudo systemctl enable webapp.service
sudo systemctl start webapp.service