/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: yangzhao_kl <yangzhao1@kylinos.cn>
 * Date: Fri Jan 8 20:56:40 2021 +0800
 */
package models

import (
	"fmt"
	"strings"
	"testing"

	"github.com/beevik/etree"
)

func TestGetResourceInfoByrscID(t *testing.T) {
	out := `xml:
	<primitive class="ocf" id="dummy" provider="pacemaker" type="Dummy">
	  <operations>
		<op id="dummy-migrate_from-interval-0s" interval="0s" name="migrate_from" timeout="20s"/>
		<op id="dummy-migrate_to-interval-0s" interval="0s" name="migrate_to" timeout="20s"/>
		<op id="dummy-monitor-interval-3s" interval="3s" name="monitor"/>
		<op id="dummy-reload-interval-0s" interval="0s" name="reload" timeout="20s"/>
		<op id="dummy-start-interval-0s" interval="0s" name="start" timeout="20s"/>
		<op id="dummy-stop-interval-0s" interval="0s" name="stop" timeout="20s"/>
	  </operations>
	</primitive>
	`
	xml := strings.Split(string(out), ":\n")[1]
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xml); err != nil {
		// return ""
	}

	// fmt.Println(doc.Space)
	// fmt.Println(doc.Tag)
	// fmt.Println(doc.Attr)
	// fmt.Println(doc.Child)
	fmt.Println(doc.Root().Tag)
	// var value map[string]interface{}
	// var doc *etree.Document
}

func TestGetGroupRscs(t *testing.T) {
	out := `Resource Group: group1
	dummy2     (ocf::heartbeat:Dummy): Started ha1
	dummy6     (ocf::heartbeat:Dummy): Started ha1
	dummy7     (ocf::heartbeat:Dummy): Started ha1
xml:
<group id="group1">
 <primitive class="ocf" id="dummy2" provider="heartbeat" type="Dummy">
   <meta_attributes id="dummy2-meta_attributes"/>
   <operations>
	 <op id="dummy2-migrate_from-interval-0s" interval="0s" name="migrate_from" timeout="20s"/>
	 <op id="dummy2-migrate_to-interval-0s" interval="0s" name="migrate_to" timeout="20s"/>
	 <op id="dummy2-monitor-interval-10s" interval="10s" name="monitor" timeout="20s"/>
	 <op id="dummy2-reload-interval-0s" interval="0s" name="reload" timeout="20s"/>
	 <op id="dummy2-start-interval-0s" interval="0s" name="start" timeout="20s"/>
	 <op id="dummy2-stop-interval-0s" interval="0s" name="stop" timeout="20s"/>
   </operations>
 </primitive>
 <primitive class="ocf" id="dummy6" provider="heartbeat" type="Dummy">
   <meta_attributes id="dummy6-meta_attributes"/>
 </primitive>
 <primitive class="ocf" id="dummy7" provider="heartbeat" type="Dummy">
   <meta_attributes id="dummy7-meta_attributes"/>
 </primitive>
</group>
	`

	xml := strings.Split(string(out), ":\n")[1]
	fmt.Println(xml)
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xml); err != nil {
		fmt.Println(err)
	}
	et := doc.FindElements("//primitive")
	rscs := []string{}

	for _, pri := range et {
		rscs = append(rscs, pri.SelectAttrValue("id", ""))
	}
	fmt.Println(rscs)
}

// [root@ha1 ~]# cibadmin --query --scope resources
// <resources>
//   <clone id="sysinfo-clone">
//     <primitive class="ocf" id="sysinfo" provider="pacemaker" type="SysInfo">
//       <operations>
//         <op id="sysinfo-monitor-interval-60s" interval="60s" name="monitor" timeout="20s"/>
//         <op id="sysinfo-start-interval-0s" interval="0s" name="start" timeout="20s"/>
//         <op id="sysinfo-stop-interval-0s" interval="0s" name="stop" timeout="20s"/>
//       </operations>
//     </primitive>
//   </clone>
//   <primitive class="ocf" id="dummy" provider="pacemaker" type="Dummy">
//     <operations>
//       <op id="dummy-migrate_from-interval-0s" interval="0s" name="migrate_from" timeout="20s"/>
//       <op id="dummy-migrate_to-interval-0s" interval="0s" name="migrate_to" timeout="20s"/>
//       <op id="dummy-monitor-interval-3s" interval="3s" name="monitor"/>
//       <op id="dummy-reload-interval-0s" interval="0s" name="reload" timeout="20s"/>
//       <op id="dummy-start-interval-0s" interval="0s" name="start" timeout="20s"/>
//       <op id="dummy-stop-interval-0s" interval="0s" name="stop" timeout="20s"/>
//     </operations>
//   </primitive>
// </resources>

// [root@ha1 ~]# crm_resource --resource dummy --query-xml
//  dummy  (ocf::pacemaker:Dummy): Started ha1
// xml:
// <primitive class="ocf" id="dummy" provider="pacemaker" type="Dummy">
//   <operations>
//     <op id="dummy-migrate_from-interval-0s" interval="0s" name="migrate_from" timeout="20s"/>
//     <op id="dummy-migrate_to-interval-0s" interval="0s" name="migrate_to" timeout="20s"/>
//     <op id="dummy-monitor-interval-3s" interval="3s" name="monitor"/>
//     <op id="dummy-reload-interval-0s" interval="0s" name="reload" timeout="20s"/>
//     <op id="dummy-start-interval-0s" interval="0s" name="start" timeout="20s"/>
//     <op id="dummy-stop-interval-0s" interval="0s" name="stop" timeout="20s"/>
//   </operations>
// </primitive>

// [root@ha1 ~]# crm_resource --resource group1 --query-xml
//  Resource Group: group1
//      dummy2     (ocf::heartbeat:Dummy): Started ha2
// xml:
// <group id="group1">
//   <primitive class="ocf" id="dummy2" provider="heartbeat" type="Dummy">
//     <operations>
//       <op id="dummy2-migrate_from-interval-0s" interval="0s" name="migrate_from" timeout="20s"/>
//       <op id="dummy2-migrate_to-interval-0s" interval="0s" name="migrate_to" timeout="20s"/>
//       <op id="dummy2-monitor-interval-10s" interval="10s" name="monitor" timeout="20s"/>
//       <op id="dummy2-reload-interval-0s" interval="0s" name="reload" timeout="20s"/>
//       <op id="dummy2-start-interval-0s" interval="0s" name="start" timeout="20s"/>
//       <op id="dummy2-stop-interval-0s" interval="0s" name="stop" timeout="20s"/>
//     </operations>
//   </primitive>
// </group>

// [root@ha1 ~]# crm_resource --resource sysinfo-clone --query-xml
//  Clone Set: sysinfo-clone [sysinfo]
//      Started: [ ha1 ha2 ]
// xml:
// <clone id="sysinfo-clone">
//   <primitive class="ocf" id="sysinfo" provider="pacemaker" type="SysInfo">
//     <operations>
//       <op id="sysinfo-monitor-interval-60s" interval="60s" name="monitor" timeout="20s"/>
//       <op id="sysinfo-start-interval-0s" interval="0s" name="start" timeout="20s"/>
//       <op id="sysinfo-stop-interval-0s" interval="0s" name="stop" timeout="20s"/>
//     </operations>
//   </primitive>
// </clone>

// [root@ha1 ~]# cibadmin -Q
// <cib crm_feature_set="3.2.0" validate-with="pacemaker-3.2" epoch="56" num_updates="14818" admin_epoch="0" cib-last-written="Tue Dec 29 14:41:29 2020" update-origin="ha1" update-client="cibadmin" update-user="hacluster" have-quorum="1" dc-uuid="2">
//   <configuration>
//     <crm_config>
//       <cluster_property_set id="cib-bootstrap-options">
//         <nvpair id="cib-bootstrap-options-have-watchdog" name="have-watchdog" value="false"/>
//         <nvpair id="cib-bootstrap-options-dc-version" name="dc-version" value="2.0.3-1.oe1-4b1f869f0f"/>
//         <nvpair id="cib-bootstrap-options-cluster-infrastructure" name="cluster-infrastructure" value="corosync"/>
//         <nvpair id="cib-bootstrap-options-cluster-name" name="cluster-name" value="hacluster"/>
//         <nvpair id="cib-bootstrap-options-stonith-enabled" name="stonith-enabled" value="false"/>
//         <nvpair id="cib-bootstrap-options-no-quorum-policy" name="no-quorum-policy" value="ignore"/>
//       </cluster_property_set>
//     </crm_config>
//     <nodes>
//       <node id="1" uname="ha1">
//         <instance_attributes id="nodes-1"/>
//       </node>
//       <node id="2" uname="ha2"/>
//     </nodes>
//     <resources>
//       <clone id="sysinfo-clone">
//         <primitive class="ocf" id="sysinfo" provider="pacemaker" type="SysInfo">
//           <operations>
//             <op id="sysinfo-monitor-interval-60s" interval="60s" name="monitor" timeout="20s"/>
//             <op id="sysinfo-start-interval-0s" interval="0s" name="start" timeout="20s"/>
//             <op id="sysinfo-stop-interval-0s" interval="0s" name="stop" timeout="20s"/>
//           </operations>
//         </primitive>
//       </clone>
//       <primitive class="ocf" id="dummy" provider="pacemaker" type="Dummy">
//         <operations>
//           <op id="dummy-migrate_from-interval-0s" interval="0s" name="migrate_from" timeout="20s"/>
//           <op id="dummy-migrate_to-interval-0s" interval="0s" name="migrate_to" timeout="20s"/>
//           <op id="dummy-monitor-interval-3s" interval="3s" name="monitor"/>
//           <op id="dummy-reload-interval-0s" interval="0s" name="reload" timeout="20s"/>
//           <op id="dummy-start-interval-0s" interval="0s" name="start" timeout="20s"/>
//           <op id="dummy-stop-interval-0s" interval="0s" name="stop" timeout="20s"/>
//         </operations>
//       </primitive>
//     </resources>
//     <constraints/>
//   </configuration>
//   <status>
//     <node_state id="2" uname="ha2" in_ccm="true" crmd="online" crm-debug-origin="do_update_resource" join="member" expected="member">
//       <transient_attributes id="2">
//         <instance_attributes id="status-2">
//           <nvpair id="status-2-arch" name="arch" value="x86_64"/>
//           <nvpair id="status-2-os" name="os" value="Linux-4.19.90-2012.1.0.0050.oe1.x86_64"/>
//           <nvpair id="status-2-free_swap" name="free_swap" value="3750"/>
//           <nvpair id="status-2-cpu_info" name="cpu_info" value="Intel(R) Core(TM) i7-8700 CPU @ 3.20GHz"/>
//           <nvpair id="status-2-cpu_speed" name="cpu_speed" value="6384.00"/>
//           <nvpair id="status-2-cpu_cores" name="cpu_cores" value="4"/>
//           <nvpair id="status-2-cpu_load" name="cpu_load" value="0.02,"/>
//           <nvpair id="status-2-ram_total" name="ram_total" value="3450"/>
//           <nvpair id="status-2-ram_free" name="ram_free" value="650"/>
//           <nvpair id="status-2-root_free" name="root_free" value="30"/>
//           <nvpair id="status-2-.health_disk" name="#health_disk" value="green"/>
//         </instance_attributes>
//       </transient_attributes>
//       <lrm id="2">
//         <lrm_resources>
//           <lrm_resource id="sysinfo" type="SysInfo" class="ocf" provider="pacemaker">
//             <lrm_rsc_op id="sysinfo_last_0" operation_key="sysinfo_start_0" operation="start" crm-debug-origin="build_active_RAs" crm_feature_set="3.2.0" transition-key="3:4:0:a056ee7d-3b46-4f8d-8a46-bc2dfc75fc75" transition-magic="0:0;3:4:0:a056ee7d-3b46-4f8d-8a46-bc2dfc75fc75" exit-reason="" on_node="ha2" call-id="7" rc-code="0" op-status="0" interval="0" last-rc-change="10137" last-run="10137" exec-time="119243" queue-time="29" op-digest="f2317cad3d54cec5d7d7aa7d0bf35cf8"/>
//             <lrm_rsc_op id="sysinfo_monitor_60000" operation_key="sysinfo_monitor_60000" operation="monitor" crm-debug-origin="build_active_RAs" crm_feature_set="3.2.0" transition-key="3:6:0:a056ee7d-3b46-4f8d-8a46-bc2dfc75fc75" transition-magic="0:0;3:6:0:a056ee7d-3b46-4f8d-8a46-bc2dfc75fc75" exit-reason="" on_node="ha2" call-id="8" rc-code="0" op-status="0" interval="60000" last-rc-change="10137" exec-time="106981" queue-time="17" op-digest="4811cef7f7f94e3a35a70be7916cb2fd"/>
//           </lrm_resource>
//           <lrm_resource id="dummy" type="Dummy" class="ocf" provider="pacemaker">
//             <lrm_rsc_op id="dummy_last_0" operation_key="dummy_stop_0" operation="stop" crm-debug-origin="do_update_resource" crm_feature_set="3.2.0" transition-key="12:17729:0:a056ee7d-3b46-4f8d-8a46-bc2dfc75fc75" transition-magic="0:0;12:17729:0:a056ee7d-3b46-4f8d-8a46-bc2dfc75fc75" exit-reason="" on_node="ha2" call-id="57" rc-code="0" op-status="0" interval="0" last-rc-change="1029368" last-run="1029368" exec-time="13514" queue-time="15" op-digest="f2317cad3d54cec5d7d7aa7d0bf35cf8" op-force-restart=" envfile  op_sleep  passwd  state " op-restart-digest="f2317cad3d54cec5d7d7aa7d0bf35cf8" op-secure-params=" passwd " op-secure-digest="f2317cad3d54cec5d7d7aa7d0bf35cf8"/>
//             <lrm_rsc_op id="dummy_monitor_3000" operation_key="dummy_monitor_3000" operation="monitor" crm-debug-origin="do_update_resource" crm_feature_set="3.2.0" transition-key="14:17725:0:a056ee7d-3b46-4f8d-8a46-bc2dfc75fc75" transition-magic="0:0;14:17725:0:a056ee7d-3b46-4f8d-8a46-bc2dfc75fc75" exit-reason="" on_node="ha2" call-id="55" rc-code="0" op-status="0" interval="3000" last-rc-change="1029323" exec-time="12324" queue-time="68" op-digest="4811cef7f7f94e3a35a70be7916cb2fd" op-secure-params=" passwd " op-secure-digest="f2317cad3d54cec5d7d7aa7d0bf35cf8"/>
//           </lrm_resource>
//         </lrm_resources>
//       </lrm>
//     </node_state>
//     <node_state id="1" uname="ha1" in_ccm="true" crmd="online" crm-debug-origin="do_update_resource" join="member" expected="member">
//       <lrm id="1">
//         <lrm_resources>
//           <lrm_resource id="sysinfo" type="SysInfo" class="ocf" provider="pacemaker">
//             <lrm_rsc_op id="sysinfo_last_0" operation_key="sysinfo_start_0" operation="start" crm-debug-origin="build_active_RAs" crm_feature_set="3.2.0" transition-key="5:41:0:a056ee7d-3b46-4f8d-8a46-bc2dfc75fc75" transition-magic="0:0;5:41:0:a056ee7d-3b46-4f8d-8a46-bc2dfc75fc75" exit-reason="" on_node="ha1" call-id="45" rc-code="0" op-status="0" interval="0" last-rc-change="29007" last-run="29007" exec-time="99736" queue-time="32" op-digest="f2317cad3d54cec5d7d7aa7d0bf35cf8"/>
//             <lrm_rsc_op id="sysinfo_monitor_60000" operation_key="sysinfo_monitor_60000" operation="monitor" crm-debug-origin="build_active_RAs" crm_feature_set="3.2.0" transition-key="6:41:0:a056ee7d-3b46-4f8d-8a46-bc2dfc75fc75" transition-magic="0:0;6:41:0:a056ee7d-3b46-4f8d-8a46-bc2dfc75fc75" exit-reason="" on_node="ha1" call-id="46" rc-code="0" op-status="0" interval="60000" last-rc-change="29007" exec-time="101711" queue-time="17" op-digest="4811cef7f7f94e3a35a70be7916cb2fd"/>
//           </lrm_resource>
//           <lrm_resource id="dummy" type="Dummy" class="ocf" provider="pacemaker">
//             <lrm_rsc_op id="dummy_last_0" operation_key="dummy_start_0" operation="start" crm-debug-origin="do_update_resource" crm_feature_set="3.2.0" transition-key="13:17729:0:a056ee7d-3b46-4f8d-8a46-bc2dfc75fc75" transition-magic="0:0;13:17729:0:a056ee7d-3b46-4f8d-8a46-bc2dfc75fc75" exit-reason="" on_node="ha1" call-id="78" rc-code="0" op-status="0" interval="0" last-rc-change="1047822" last-run="1047822" exec-time="12155" queue-time="21" op-digest="f2317cad3d54cec5d7d7aa7d0bf35cf8" op-force-restart=" envfile  op_sleep  passwd  state " op-restart-digest="f2317cad3d54cec5d7d7aa7d0bf35cf8" op-secure-params=" passwd " op-secure-digest="f2317cad3d54cec5d7d7aa7d0bf35cf8"/>
//             <lrm_rsc_op id="dummy_monitor_3000" operation_key="dummy_monitor_3000" operation="monitor" crm-debug-origin="do_update_resource" crm_feature_set="3.2.0" transition-key="14:17729:0:a056ee7d-3b46-4f8d-8a46-bc2dfc75fc75" transition-magic="0:0;14:17729:0:a056ee7d-3b46-4f8d-8a46-bc2dfc75fc75" exit-reason="" on_node="ha1" call-id="79" rc-code="0" op-status="0" interval="3000" last-rc-change="1047822" exec-time="8472" queue-time="19" op-digest="4811cef7f7f94e3a35a70be7916cb2fd" op-secure-params=" passwd " op-secure-digest="f2317cad3d54cec5d7d7aa7d0bf35cf8"/>
//           </lrm_resource>
//         </lrm_resources>
//       </lrm>
//       <transient_attributes id="1">
//         <instance_attributes id="status-1">
//           <nvpair id="status-1-arch" name="arch" value="x86_64"/>
//           <nvpair id="status-1-free_swap" name="free_swap" value="3850"/>
//           <nvpair id="status-1-cpu_speed" name="cpu_speed" value="6384.00"/>
//           <nvpair id="status-1-cpu_load" name="cpu_load" value="0.04,"/>
//           <nvpair id="status-1-cpu_info" name="cpu_info" value="Intel(R) Core(TM) i7-8700 CPU @ 3.20GHz"/>
//           <nvpair id="status-1-ram_total" name="ram_total" value="3450"/>
//           <nvpair id="status-1-cpu_cores" name="cpu_cores" value="4"/>
//           <nvpair id="status-1-root_free" name="root_free" value="27"/>
//           <nvpair id="status-1-ram_free" name="ram_free" value="650"/>
//           <nvpair id="status-1-.health_disk" name="#health_disk" value="green"/>
//           <nvpair id="status-1-os" name="os" value="Linux-4.19.90-2012.1.0.0050.oe1.x86_64"/>
//         </instance_attributes>
//       </transient_attributes>
//     </node_state>
//   </status>
// </cib>

// [root@ha1 ~]# crm_mon -1 --as-xml
// <crm_mon version="2.0.3">
//   <summary>
//     <stack type="corosync"/>
//     <current_dc present="true" version="2.0.3-1.oe1-4b1f869f0f" name="ha2" id="2" with_quorum="true"/>
//     <last_update time="Fri Jan  8 13:48:44 2021"/>
//     <last_change time="Tue Dec 29 14:41:29 2020" user="hacluster" client="cibadmin" origin="ha1"/>
//     <nodes_configured number="2"/>
//     <resources_configured number="3" disabled="0" blocked="0"/>
//     <cluster_options stonith-enabled="false" symmetric-cluster="true" no-quorum-policy="ignore" maintenance-mode="false"/>
//   </summary>
//   <nodes>
//     <node name="ha1" id="1" online="true" standby="false" standby_onfail="false" maintenance="false" pending="false" unclean="false" shutdown="false" expected_up="true" is_dc="false" resources_running="2" type="member"/>
//     <node name="ha2" id="2" online="true" standby="false" standby_onfail="false" maintenance="false" pending="false" unclean="false" shutdown="false" expected_up="true" is_dc="true" resources_running="1" type="member"/>
//   </nodes>
//   <resources>
//     <clone id="sysinfo-clone" multi_state="false" unique="false" managed="true" failed="false" failure_ignored="false">
//       <resource id="sysinfo" resource_agent="ocf::pacemaker:SysInfo" role="Started" active="true" orphaned="false" blocked="false" managed="true" failed="false" failure_ignored="false" nodes_running_on="1">
//         <node name="ha2" id="2" cached="true"/>
//       </resource>
//       <resource id="sysinfo" resource_agent="ocf::pacemaker:SysInfo" role="Started" active="true" orphaned="false" blocked="false" managed="true" failed="false" failure_ignored="false" nodes_running_on="1">
//         <node name="ha1" id="1" cached="true"/>
//       </resource>
//     </clone>
//     <resource id="dummy" resource_agent="ocf::pacemaker:Dummy" role="Started" active="true" orphaned="false" blocked="false" managed="true" failed="false" failure_ignored="false" nodes_running_on="1">
//       <node name="ha1" id="1" cached="true"/>
//     </resource>
//   </resources>
//   <node_attributes>
//     <node name="ha1">
//       <attribute name="arch" value="x86_64"/>
//       <attribute name="cpu_cores" value="4"/>
//       <attribute name="cpu_info" value="Intel(R) Core(TM) i7-8700 CPU @ 3.20GHz"/>
//       <attribute name="cpu_load" value="0.06,"/>
//       <attribute name="cpu_speed" value="6384.00"/>
//       <attribute name="free_swap" value="3850"/>
//       <attribute name="os" value="Linux-4.19.90-2012.1.0.0050.oe1.x86_64"/>
//       <attribute name="ram_free" value="700"/>
//       <attribute name="ram_total" value="3450"/>
//       <attribute name="root_free" value="27"/>
//     </node>
//     <node name="ha2">
//       <attribute name="arch" value="x86_64"/>
//       <attribute name="cpu_cores" value="4"/>
//       <attribute name="cpu_info" value="Intel(R) Core(TM) i7-8700 CPU @ 3.20GHz"/>
//       <attribute name="cpu_load" value="0.16,"/>
//       <attribute name="cpu_speed" value="6384.00"/>
//       <attribute name="free_swap" value="3750"/>
//       <attribute name="os" value="Linux-4.19.90-2012.1.0.0050.oe1.x86_64"/>
//       <attribute name="ram_free" value="650"/>
//       <attribute name="ram_total" value="3450"/>
//       <attribute name="root_free" value="30"/>
//     </node>
//   </node_attributes>
//   <node_history>
//     <node name="ha2">
//       <resource_history id="sysinfo" orphan="false" migration-threshold="1000000">
//         <operation_history call="7" task="start" last-rc-change="Thu Jan  1 10:48:57 1970" last-run="Thu Jan  1 10:48:57 1970" exec-time="119243ms" queue-time="29ms" rc="0" rc_text="ok"/>
//         <operation_history call="8" task="monitor" interval="60000ms" last-rc-change="Thu Jan  1 10:48:57 1970" exec-time="106981ms" queue-time="17ms" rc="0" rc_text="ok"/>
//       </resource_history>
//       <resource_history id="dummy" orphan="false" migration-threshold="1000000">
//         <operation_history call="55" task="monitor" interval="3000ms" last-rc-change="Tue Jan 13 05:55:23 1970" exec-time="12324ms" queue-time="68ms" rc="0" rc_text="ok"/>
//         <operation_history call="57" task="stop" last-rc-change="Tue Jan 13 05:56:08 1970" last-run="Tue Jan 13 05:56:08 1970" exec-time="13514ms" queue-time="15ms" rc="0" rc_text="ok"/>
//       </resource_history>
//     </node>
//     <node name="ha1">
//       <resource_history id="sysinfo" orphan="false" migration-threshold="1000000">
//         <operation_history call="45" task="start" last-rc-change="Thu Jan  1 16:03:27 1970" last-run="Thu Jan  1 16:03:27 1970" exec-time="99736ms" queue-time="32ms" rc="0" rc_text="ok"/>
//         <operation_history call="46" task="monitor" interval="60000ms" last-rc-change="Thu Jan  1 16:03:27 1970" exec-time="101711ms" queue-time="17ms" rc="0" rc_text="ok"/>
//       </resource_history>
//       <resource_history id="dummy" orphan="false" migration-threshold="1000000">
//         <operation_history call="78" task="start" last-rc-change="Tue Jan 13 11:03:42 1970" last-run="Tue Jan 13 11:03:42 1970" exec-time="12155ms" queue-time="21ms" rc="0" rc_text="ok"/>
//         <operation_history call="79" task="monitor" interval="3000ms" last-rc-change="Tue Jan 13 11:03:42 1970" exec-time="8472ms" queue-time="19ms" rc="0" rc_text="ok"/>
//       </resource_history>
//     </node>
//   </node_history>
// </crm_mon>

// [root@ha1 ~]# cibadmin --query --scope resources
// <resources>
//   <clone id="sysinfo-clone">
//     <primitive class="ocf" id="sysinfo" provider="pacemaker" type="SysInfo">
//       <operations>
//         <op id="sysinfo-monitor-interval-60s" interval="60s" name="monitor" timeout="20s"/>
//         <op id="sysinfo-start-interval-0s" interval="0s" name="start" timeout="20s"/>
//         <op id="sysinfo-stop-interval-0s" interval="0s" name="stop" timeout="20s"/>
//       </operations>
//     </primitive>
//   </clone>
//   <primitive class="ocf" id="dummy" provider="pacemaker" type="Dummy">
//     <operations>
//       <op id="dummy-migrate_from-interval-0s" interval="0s" name="migrate_from" timeout="20s"/>
//       <op id="dummy-migrate_to-interval-0s" interval="0s" name="migrate_to" timeout="20s"/>
//       <op id="dummy-monitor-interval-3s" interval="3s" name="monitor"/>
//       <op id="dummy-reload-interval-0s" interval="0s" name="reload" timeout="20s"/>
//       <op id="dummy-start-interval-0s" interval="0s" name="start" timeout="20s"/>
//       <op id="dummy-stop-interval-0s" interval="0s" name="stop" timeout="20s"/>
//     </operations>
//   </primitive>
// </resources>
