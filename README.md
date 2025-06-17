## Bender service

This is a system-wide systemd service which has root privileges, It is because this service is a part of the system which responsible for the creation and deletion of the playground container and associating bridge networks. This service accepts http requests from the backend.

This service can be installed on a system using the install script provided:

```bash
chmod +x ./install.sh
./install.sh
```

You can check the status by using the following command:
```bash
sudo systemctl status bender
```

Stop service:
```bash
sudo systemctl stop bender
```

Start and restart service:
```bash
sudo systemctl start bender
sudo systemctl restart bender
```

The service exposes a port 8080 which accepts requests from the backend