package models

import "testing"

func TestGetClusterPropertyFromXml(t *testing.T) {
	_ = `<?xml version="1.0"?><!DOCTYPE resource-agent SYSTEM "ra-api-1.dtd">
	<resource-agent name="pacemaker-schedulerd">
	  <version>1.0</version>
	  <longdesc lang="en">Cluster properties used by Pacemaker's scheduler, formerly known as pengine</longdesc>
	  <shortdesc lang="en">scheduler properties</shortdesc>
	  <parameters>
		<parameter name="no-quorum-policy" unique="0">
		  <shortdesc lang="en">What to do when the cluster does not have quorum</shortdesc>
		  <content type="enum" default="stop"/>
		  <longdesc lang="en">What to do when the cluster does not have quorum  Allowed values: stop, freeze, ignore, suicide</longdesc>
		</parameter>
		<parameter name="symmetric-cluster" unique="0">
		  <shortdesc lang="en">All resources can run anywhere by default</shortdesc>
		  <content type="boolean" default="true"/>
		  <longdesc lang="en">All resources can run anywhere by default</longdesc>
		</parameter>
		<parameter name="maintenance-mode" unique="0">
		  <shortdesc lang="en">Should the cluster monitor resources and start/stop them as required</shortdesc>
		  <content type="boolean" default="false"/>
		  <longdesc lang="en">Should the cluster monitor resources and start/stop them as required</longdesc>
		</parameter>
		<parameter name="start-failure-is-fatal" unique="0">
		  <shortdesc lang="en">Always treat start failures as fatal</shortdesc>
		  <content type="boolean" default="true"/>
		  <longdesc lang="en">When set to TRUE, the cluster will immediately ban a resource from a node if it fails to start there. When FALSE, the cluster will instead check the resource's fail count against its migration-threshold.</longdesc>
		</parameter>
		<parameter name="enable-startup-probes" unique="0">
		  <shortdesc lang="en">Should the cluster check for active resources during startup</shortdesc>
		  <content type="boolean" default="true"/>
		  <longdesc lang="en">Should the cluster check for active resources during startup</longdesc>
		</parameter>
		<parameter name="stonith-enabled" unique="0">
		  <shortdesc lang="en">Failed nodes are STONITH'd</shortdesc>
		  <content type="boolean" default="true"/>
		  <longdesc lang="en">Failed nodes are STONITH'd</longdesc>
		</parameter>
		<parameter name="stonith-action" unique="0">
		  <shortdesc lang="en">Action to send to STONITH device ('poweroff' is a deprecated alias for 'off')</shortdesc>
		  <content type="enum" default="reboot"/>
		  <longdesc lang="en">Action to send to STONITH device ('poweroff' is a deprecated alias for 'off')  Allowed values: reboot, off, poweroff</longdesc>
		</parameter>
		<parameter name="stonith-timeout" unique="0">
		  <shortdesc lang="en">How long to wait for the STONITH action (reboot,on,off) to complete</shortdesc>
		  <content type="time" default="60s"/>
		  <longdesc lang="en">How long to wait for the STONITH action (reboot,on,off) to complete</longdesc>
		</parameter>
		<parameter name="have-watchdog" unique="0">
		  <shortdesc lang="en">Enable watchdog integration</shortdesc>
		  <content type="boolean" default="false"/>
		  <longdesc lang="en">Set automatically by the cluster if SBD is detected.  User configured values are ignored.</longdesc>
		</parameter>
		<parameter name="concurrent-fencing" unique="0">
		  <shortdesc lang="en">Allow performing fencing operations in parallel</shortdesc>
		  <content type="boolean" default="false"/>
		  <longdesc lang="en">Allow performing fencing operations in parallel</longdesc>
		</parameter>
		<parameter name="startup-fencing" unique="0">
		  <shortdesc lang="en">STONITH unseen nodes</shortdesc>
		  <content type="boolean" default="true"/>
		  <longdesc lang="en">Advanced Use Only!  Not using the default is very unsafe!</longdesc>
		</parameter>
		<parameter name="cluster-delay" unique="0">
		  <shortdesc lang="en">Round trip delay over the network (excluding action execution)</shortdesc>
		  <content type="time" default="60s"/>
		  <longdesc lang="en">The "correct" value will depend on the speed and load of your network and cluster nodes.</longdesc>
		</parameter>
		<parameter name="batch-limit" unique="0">
		  <shortdesc lang="en">The number of jobs that the TE is allowed to execute in parallel</shortdesc>
		  <content type="integer" default="0"/>
		  <longdesc lang="en">The "correct" value will depend on the speed and load of your network and cluster nodes.</longdesc>
		</parameter>
		<parameter name="migration-limit" unique="0">
		  <shortdesc lang="en">The number of migration jobs that the TE is allowed to execute in parallel on a node</shortdesc>
		  <content type="integer" default="-1"/>
		  <longdesc lang="en">The number of migration jobs that the TE is allowed to execute in parallel on a node</longdesc>
		</parameter>
		<parameter name="stop-all-resources" unique="0">
		  <shortdesc lang="en">Should the cluster stop all active resources</shortdesc>
		  <content type="boolean" default="false"/>
		  <longdesc lang="en">Should the cluster stop all active resources</longdesc>
		</parameter>
		<parameter name="stop-orphan-resources" unique="0">
		  <shortdesc lang="en">Should deleted resources be stopped</shortdesc>
		  <content type="boolean" default="true"/>
		  <longdesc lang="en">Should deleted resources be stopped</longdesc>
		</parameter>
		<parameter name="stop-orphan-actions" unique="0">
		  <shortdesc lang="en">Should deleted actions be cancelled</shortdesc>
		  <content type="boolean" default="true"/>
		  <longdesc lang="en">Should deleted actions be cancelled</longdesc>
		</parameter>
		<parameter name="remove-after-stop" unique="0">
		  <shortdesc lang="en">Remove resources from the executor after they are stopped</shortdesc>
		  <content type="boolean" default="false"/>
		  <longdesc lang="en">Always set this to false.  Other values are, at best, poorly tested and potentially dangerous.</longdesc>
		</parameter>
		<parameter name="pe-error-series-max" unique="0">
		  <shortdesc lang="en">The number of scheduler inputs resulting in ERRORs to save</shortdesc>
		  <content type="integer" default="-1"/>
		  <longdesc lang="en">Zero to disable, -1 to store unlimited</longdesc>
		</parameter>
		<parameter name="pe-warn-series-max" unique="0">
		  <shortdesc lang="en">The number of scheduler inputs resulting in WARNINGs to save</shortdesc>
		  <content type="integer" default="5000"/>
		  <longdesc lang="en">Zero to disable, -1 to store unlimited</longdesc>
		</parameter>
		<parameter name="pe-input-series-max" unique="0">
		  <shortdesc lang="en">The number of other scheduler inputs to save</shortdesc>
		  <content type="integer" default="4000"/>
		  <longdesc lang="en">Zero to disable, -1 to store unlimited</longdesc>
		</parameter>
		<parameter name="node-health-strategy" unique="0">
		  <shortdesc lang="en">The strategy combining node attributes to determine overall node health.</shortdesc>
		  <content type="enum" default="none"/>
		  <longdesc lang="en">Requires external entities to create node attributes (named with the prefix '#health') with values: 'red', 'yellow' or 'green'.  Allowed values: none, migrate-on-red, only-green, progressive, custom</longdesc>
		</parameter>
		<parameter name="node-health-base" unique="0">
		  <shortdesc lang="en">The base score assigned to a node</shortdesc>
		  <content type="integer" default="0"/>
		  <longdesc lang="en">Only used when node-health-strategy is set to progressive.</longdesc>
		</parameter>
		<parameter name="node-health-green" unique="0">
		  <shortdesc lang="en">The score 'green' translates to in rsc_location constraints</shortdesc>
		  <content type="integer" default="0"/>
		  <longdesc lang="en">Only used when node-health-strategy is set to custom or progressive.</longdesc>
		</parameter>
		<parameter name="node-health-yellow" unique="0">
		  <shortdesc lang="en">The score 'yellow' translates to in rsc_location constraints</shortdesc>
		  <content type="integer" default="0"/>
		  <longdesc lang="en">Only used when node-health-strategy is set to custom or progressive.</longdesc>
		</parameter>
		<parameter name="node-health-red" unique="0">
		  <shortdesc lang="en">The score 'red' translates to in rsc_location constraints</shortdesc>
		  <content type="integer" default="-INFINITY"/>
		  <longdesc lang="en">Only used when node-health-strategy is set to custom or progressive.</longdesc>
		</parameter>
		<parameter name="placement-strategy" unique="0">
		  <shortdesc lang="en">The strategy to determine resource placement</shortdesc>
		  <content type="enum" default="default"/>
		  <longdesc lang="en">The strategy to determine resource placement  Allowed values: default, utilization, minimal, balanced</longdesc>
		</parameter>
	  </parameters>
	</resource-agent>`
	// TODO: add test here
}
