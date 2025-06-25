#!/bin/bash

go build . 

sudo mkdir -p /usr/local/bin

sudo cp ./bender /usr/local/bin/bender

sudo chmod +x /usr/local/bin/bender

sudo  chown root:root /usr/local/bin/bender

service="[Unit]
Description=Bender service
After=network.target redis.service

[Service]
Type=simple
ExecStart=/usr/local/bin/bender
Restart=on-failure
User=$(whoami)
WorkingDirectory=/usr/local/bin
SupplementaryGroups=docker

[Install]
WantedBy=multi-user.target"

echo "$service" | sudo tee /etc/systemd/system/bender.service > /dev/null

sudo chown root:docker /var/run/docker.sock
sudo chmod 0660 /var/run/docker.sock
sudo systemctl stop bender || true 
sudo systemctl daemon-reexec

sudo systemctl daemon-reload

sudo systemctl enable bender

sudo systemctl start bender
sudo systemctl restart bender
