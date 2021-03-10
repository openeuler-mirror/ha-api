#!/usr/bin/bash

rm -rf /tmp/kylinha-log

hostname=`hostname`
gentime=` date +'%Y%m%d%H%M%S'`
file_path=/usr/share/heartbeat-gui/ha-api/static/kylinha-log-$hostname-$gentime.tar
file_name=kylinha-log-$hostname-$gentime.tar

#expect -c '
#spawn sosreport -a --tmp-dir /tmp/kylinha-log
#expect "Press ENTER to continue, or CTRL-C to quit."
#send "\n"
#expect "Please enter your first initial and last name*"
#send "\n"
#expect "Please enter the case id that you are generating this report for []*"
#send "\n"
#expect eof
#'
temp_dir=/tmp/kylinha-log-$hostname-$gentime

mkdir $temp_dir

echo "#crm_mon -1" >> $temp_dir/commands.log
crm_mon -1 >> $temp_dir/commands.log
echo "" >> $temp_dir/commands.log

echo "#mount"  >> $temp_dir/commands.log
mount  >> $temp_dir/commands.log
echo "" >> $temp_dir/commands.log

echo "#iptables -L" >> $temp_dir/commands.log
iptables -L >> $temp_dir/commands.log
echo "" >> $temp_dir/commands.log

echo "#getenforce" >> $temp_dir/commands.log
getenforce >> $temp_dir/commands.log
echo "" >> $temp_dir/commands.log

echo "#ifconfig" >> $temp_dir/commands.log
ifconfig >> $temp_dir/commands.log
echo "" >> $temp_dir/commands.log

echo "#ps -aux|grep -i -e pacemaker -e corosync -e pcsd" >> $temp_dir/commands.log
ps -aux|grep -i -e pacemaker -e corosync -e pcsd >> $temp_dir/commands.log
echo "" >> $temp_dir/commands.log

echo "#systemctl status ha-api" >> $temp_dir/commands.log
systemctl status ha-api >> $temp_dir/commands.log
echo "" >> $temp_dir/commands.log

cp -rf /var/log/messages*  /var/log/pcsd/* /var/log/cluster/* /var/log/pacemaker/pacemaker.log* /var/lib/pacemaker /etc/hosts /etc/corosync $temp_dir

cd  /tmp/
tar -cf  $file_path kylinha-log-$hostname-$gentime
echo $file_name
