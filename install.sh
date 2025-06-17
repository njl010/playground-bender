#!/bin/bash

go build .

sudo mkdir -p /usr/local/bin

sudo cp ./bender /usr/local/bin/bender

sudo chmod +x /usr/local/bin/bender

service="[Unit]
Description=Bender service
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/bender
Restart=on-failure
User=$(whoami)
WorkingDirectory=/home/$(whoami)

[Install]
WantedBy=multi-user.target"

echo "$service" | sudo tee /etc/systemd/system/bender.service > /dev/null

sudo systemctl enable bender

sudo systemctl start bender