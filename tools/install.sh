# install bin and static files
# check path
if [ ! -d "/usr/share/heartbeat-gui/ha-api" ]; then
  mkdir -p /usr/share/heartbeat-gui/ha-api
fi
cp ../ha-api /usr/share/heartbeat-gui/ha-api/
chmod +x /usr/share/heartbeat-gui/ha-api/ha-api
cp -r ../views /usr/share/heartbeat-gui/ha-api/

# install loggen.sh script
cp ./loggen.sh /usr/share/heartbeat-gui/ha-api/loggen.sh
chmod 755 /usr/share/heartbeat-gui/ha-api/loggen.sh

# create log storage dir
mkdir -p /usr/share/heartbeat-gui/ha-api/static

# install ls /usr/lib/systemd/system/ha-api.service
# currently not used as a service
cp ./ha-api.service /usr/lib/systemd/system/ha-api.service
