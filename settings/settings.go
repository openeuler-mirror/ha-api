/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * ha-api is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: liqiuyu
 * Date: 2022-04-19 16:49:51
 * LastEditTime: 2022-04-20 11:14:14
 * Description: 集群中的一些配置
 ******************************************************************************/
package settings

import (
	"path/filepath"
)

const (
	SystemctlBinary   = "/bin/systemctl"
	ChkconfigBinary   = "/sbin/chkconfig"
	ServiceBinary     = "/sbin/service"
	PacemakerBinaries = "/usr/sbin/"

	corosync_binaries                        = "/usr/sbin/"
	corosync_qnet_binaries                   = "/usr/bin/"
	corosync_conf_dir                        = "/etc/corosync/"
	corosync_qdevice_net_client_ca_file_name = "qnetd-cacert.crt"
	// Must be set to 256 for corosync to work in FIPS environment.
	corosync_authkey_bytes = 256
	corosync_log_file      = "/var/log/cluster/corosync.log"

	pacemaker_authkey_file = "/etc/pacemaker/authkey"
	// Using the same value as for corosync. Higher values MAY work in FIPS.
	pacemaker_authkey_bytes = 256

	booth_authkey_file_mode = 0o600
	// # Booth does not support keys longer than 64 bytes.
	booth_authkey_bytes  = 64
	cluster_conf_file    = "/etc/cluster/cluster.conf"
	fence_agent_binaries = "/usr/sbin/"
	PacemakerSchedulerd  = "/usr/libexec/pacemaker/pacemaker-schedulerd"
	PacemakerControld    = "/usr/libexec/pacemaker/pacemaker-controld"
	PacemakerBased       = "/usr/libexec/pacemaker/pacemaker-based"
	pacemaker_fenced     = "/usr/libexec/pacemaker/pacemaker-fenced"
	pcs_version          = "0.10.2"

	crm_mon_schema              = "/usr/share/pacemaker/crm_mon.rng"
	agent_metadata_schema       = "/usr/share/resource-agents/ra-api-1.dtd"
	pcsd_cert_location          = "/var/lib/pcsd/pcsd.crt"
	pcsd_key_location           = "/var/lib/pcsd/pcsd.key"
	pcsd_users_conf_location    = "/var/lib/pcsd/pcs_users.conf"
	pcsd_settings_conf_location = "/var/lib/pcsd/pcs_settings.conf"
	pcsd_exec_location          = "/usr/lib/pcsd/"
	pcsd_log_location           = "/var/log/pcsd/pcsd.log"
	pcsd_default_port           = 2224
	pcsd_config                 = "/etc/sysconfig/pcsd"
	cib_dir                     = "/var/lib/pacemaker/cib/"
	pacemaker_uname             = "hacluster"
	pacemaker_gname             = "haclient"
	sbd_binary                  = "/usr/sbin/sbd"
	sbd_watchdog_default        = "/dev/watchdog"
	sbd_config                  = "/etc/sysconfig/sbd"
	// # this limit is also mentioned in docs, change there as well
	sbd_max_device_num = 3

	pacemaker_wait_timeout_status = 124
	booth_config_dir              = "/etc/booth"
	booth_binary                  = "/usr/sbin/booth"
	default_request_timeout       = 60
	pcs_bundled_dir               = "/usr/lib/pcs/bundled/"

	default_ssl_ciphers = "DEFAULT:!RC4:!3DES:@STRENGTH"

	// # Ssl options are based on default options in python (maybe with some extra
	// # options). Format here is the same as the PCSD_SSL_OPTIONS environment
	// # variable format (string with coma as a delimiter).
	// default_ssl_options = ",".join([
	//     "OP_NO_COMPRESSION",
	//     "OP_CIPHER_SERVER_PREFERENCE",
	//     "OP_SINGLE_DH_USE",
	//     "OP_SINGLE_ECDH_USE",
	//     "OP_NO_SSLv2",
	//     "OP_NO_SSLv3",
	//     "OP_NO_TLSv1",
	//     "OP_NO_TLSv1_1",
	//     "OP_NO_RENEGOTIATION",
	// ])
	// # Set pcsd_gem_path to None if there are no bundled ruby gems and the path does
	// # not exists.
	pcsd_gem_path   = "vendor/bundle/ruby"
	ruby_executable = "/usr/bin/ruby"

	gui_session_lifetime_seconds = 60 * 60
)

var CrmResourceBinary = filepath.Join(PacemakerBinaries, "crm_resource")
var corosync_conf_file = filepath.Join(corosync_conf_dir, "corosync.conf")
var corosync_uidgid_dir = filepath.Join(corosync_conf_dir, "uidgid.d/")
var corosync_qdevice_net_server_certs_dir = filepath.Join(corosync_conf_dir, "qnetd/nssdb")
var corosync_qdevice_net_client_certs_dir = filepath.Join(corosync_conf_dir, "qdevice/net/nssdb")
var corosync_authkey_file = filepath.Join(corosync_conf_dir, "authkey")

var crm_report = filepath.Join(PacemakerBinaries, "crm_report")
var crm_verify = filepath.Join(PacemakerBinaries, "crm_verify")
var cibadmin = filepath.Join(PacemakerBinaries, "cibadmin")

var pcs_bundled_pacakges_dir = filepath.Join(pcs_bundled_dir, "packages")

// # message types are also mentioned in docs, change there as well
var sbd_message_types = []string{"test", "reset", "off", "crashdump", "exit", "clear"}
