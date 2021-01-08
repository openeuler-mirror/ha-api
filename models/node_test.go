package models

// [root@ha1 ~]# crm_mon --as-xml
// <crm_mon version="2.0.3">
//   <summary>
//     <stack type="corosync"/>
//     <current_dc present="true" version="2.0.3-1.oe1-4b1f869f0f" name="ha2" id="2" with_quorum="true"/>
//     <last_update time="Fri Jan  8 17:22:32 2021"/>
//     <last_change time="Fri Jan  8 14:54:31 2021" user="hacluster" client="cibadmin" origin="ha1"/>
//     <nodes_configured number="2"/>
//     <resources_configured number="4" disabled="0" blocked="0"/>
//     <cluster_options stonith-enabled="false" symmetric-cluster="true" no-quorum-policy="ignore" maintenance-mode="false"/>
//   </summary>
//   <nodes>
//     <node name="ha1" id="1" online="true" standby="false" standby_onfail="false" maintenance="false" pending="false" unclean="false" shutdown="false" expected_up="true" is_dc="false" resources_running="2" type="member"/>
//     <node name="ha2" id="2" online="true" standby="false" standby_onfail="false" maintenance="false" pending="false" unclean="false" shutdown="false" expected_up="true" is_dc="true" resources_running="2" type="member"/>
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
//     <resource id="dummy2" resource_agent="ocf::heartbeat:Dummy" role="Started" active="true" orphaned="false" blocked="false" managed="true" failed="false" failure_ignored="false" nodes_running_on="1">
//       <node name="ha2" id="2" cached="true"/>
//     </resource>
//   </resources>
//   <node_attributes>
//     <node name="ha1">
//       <attribute name="arch" value="x86_64"/>
//       <attribute name="cpu_cores" value="4"/>
//       <attribute name="cpu_info" value="Intel(R) Core(TM) i7-8700 CPU @ 3.20GHz"/>
//       <attribute name="cpu_load" value="0.09,"/>
//       <attribute name="cpu_speed" value="6384.00"/>
//       <attribute name="free_swap" value="3850"/>
//       <attribute name="os" value="Linux-4.19.90-2012.1.0.0050.oe1.x86_64"/>
//       <attribute name="ram_free" value="800"/>
//       <attribute name="ram_total" value="3450"/>
//       <attribute name="root_free" value="27"/>
//     </node>
//     <node name="ha2">
//       <attribute name="arch" value="x86_64"/>
//       <attribute name="cpu_cores" value="4"/>
//       <attribute name="cpu_info" value="Intel(R) Core(TM) i7-8700 CPU @ 3.20GHz"/>
//       <attribute name="cpu_load" value="average:"/>
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
//       <resource_history id="dummy2" orphan="false" migration-threshold="1000000">
//         <operation_history call="62" task="start" last-rc-change="Fri Jan 23 06:09:20 1970" last-run="Fri Jan 23 06:09:20 1970" exec-time="8447ms" queue-time="31ms" rc="0" rc_text="ok"/>
//         <operation_history call="63" task="monitor" interval="10000ms" last-rc-change="Fri Jan 23 06:09:20 1970" exec-time="7228ms" queue-time="31ms" rc="0" rc_text="ok"/>
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
