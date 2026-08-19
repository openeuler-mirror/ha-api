/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: yangzhao_kl <yangzhao1@kylinos.cn>
 * Date: Fri Jan 8 20:56:40 2021 +0800
 */
package models

import (
	"testing"

	"github.com/beevik/etree"
)

// ==================== validateResourceID ====================

func TestValidateResourceID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{"valid simple id", "dummy", false},
		{"valid with underscore", "my_resource", false},
		{"valid with hyphen", "my-resource", false},
		{"valid with colon", "clone:0", false},
		{"valid with dot", "my.resource", false},
		{"valid with digits", "rsc123", false},
		{"valid complex id", "sysinfo-clone:1", false},
		{"empty id", "", true},
		{"invalid with space", "my resource", true},
		{"invalid with slash", "my/resource", true},
		{"invalid with special char", "rsc@123", true},
		{"invalid with brackets", "rsc[0]", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateResourceID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateResourceID(%q) error = %v, wantErr %v", tt.id, err, tt.wantErr)
			}
		})
	}
}

// ==================== validateIdentifier ====================

func TestValidateIdentifier(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		fieldName string
		wantErr   bool
	}{
		{"valid class", "ocf", "class", false},
		{"valid provider", "pacemaker", "provider", false},
		{"valid with hyphen", "heart-beat", "provider", false},
		{"valid with underscore", "my_type", "type", false},
		{"empty value", "", "class", true},
		{"invalid with colon", "ocf:heartbeat", "class", true},
		{"invalid with dot", "my.type", "type", true},
		{"invalid with space", "my type", "provider", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateIdentifier(tt.value, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateIdentifier(%q, %q) error = %v, wantErr %v", tt.value, tt.fieldName, err, tt.wantErr)
			}
		})
	}
}

// ==================== splitXMLOutput ====================

func TestSplitXMLOutput(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			"valid output with xml prefix",
			"xml:\n<primitive id=\"dummy\"/>",
			"<primitive id=\"dummy\"/>",
			false,
		},
		{
			"valid output with text prefix",
			"dummy (ocf::pacemaker:Dummy): Started ha1\nxml:\n<primitive id=\"dummy\"/>",
			"<primitive id=\"dummy\"/>",
			false,
		},
		{
			"no separator",
			"<primitive id=\"dummy\"/>",
			"",
			true,
		},
		{
			"empty input",
			"",
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := splitXMLOutput([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("splitXMLOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("splitXMLOutput() = %q, want %q", got, tt.want)
			}
		})
	}
}

// ==================== GetResourceSvcFromXml ====================

func TestGetResourceSvcFromXml(t *testing.T) {
	tests := []struct {
		name    string
		xmlStr  string
		wantSvc string
	}{
		{
			"standard ocf resource agent",
			`<resource resource_agent="ocf::pacemaker:Dummy"/>`,
			"Dummy",
		},
		{
			"heartbeat resource agent",
			`<resource resource_agent="ocf::heartbeat:IPaddr2"/>`,
			"IPaddr2",
		},
		{
			"systemd resource agent",
			`<resource resource_agent="systemd:my-service"/>`,
			"my-service",
		},
		{
			"no resource_agent attribute",
			`<resource id="dummy"/>`,
			"",
		},
		{
			"single segment resource agent",
			`<resource resource_agent="Dummy"/>`,
			"Dummy",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := etree.NewDocument()
			if err := doc.ReadFromString(tt.xmlStr); err != nil {
				t.Fatalf("failed to parse XML: %v", err)
			}
			got := GetResourceSvcFromXml(doc.Root())
			if got != tt.wantSvc {
				t.Errorf("GetResourceSvcFromXml() = %q, want %q", got, tt.wantSvc)
			}
		})
	}
}

// ==================== hasAttribute ====================

func TestHasAttribute(t *testing.T) {
	doc := etree.NewDocument()
	if err := doc.ReadFromString(`<rsc_colocation rsc="dummy" with-rsc="ip" score="INFINITY" with-rsc-role="Master"/>`); err != nil {
		t.Fatalf("failed to parse XML: %v", err)
	}
	root := doc.Root()

	tests := []struct {
		name     string
		attrName string
		want     bool
	}{
		{"existing attribute", "rsc", true},
		{"existing with-rsc-role", "with-rsc-role", true},
		{"existing score", "score", true},
		{"non-existing attribute", "rsc-role", false},
		{"non-existing attribute 2", "nonexistent", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasAttribute(root, tt.attrName)
			if got != tt.want {
				t.Errorf("hasAttribute() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ==================== getOtherRsc ====================

func TestGetOtherRsc(t *testing.T) {
	tests := []struct {
		name    string
		rsc     string
		rscWith string
		want    string
	}{
		{"rsc is empty", "", "ip", "ip"},
		{"rsc is not empty", "dummy", "ip", "dummy"},
		{"both empty", "", "", ""},
		{"both non-empty returns rsc", "dummy", "ip2", "dummy"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getOtherRsc(tt.rsc, tt.rscWith)
			if got != tt.want {
				t.Errorf("getOtherRsc() = %q, want %q", got, tt.want)
			}
		})
	}
}

// ==================== getLevelFromScore ====================

func TestGetLevelFromScore(t *testing.T) {
	tests := []struct {
		name  string
		score string
		want  string
	}{
		{"master node", "20000", "Master Node"},
		{"slave 1", "16000", "Slave 1"},
		{"slave 2", "15000", "Slave 2"},
		{"slave 3", "14000", "Slave 3"},
		{"slave 4", "13000", "Slave 4"},
		{"unknown score", "10000", ""},
		{"empty score", "", ""},
		{"infinity", "INFINITY", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getLevelFromScore(tt.score)
			if got != tt.want {
				t.Errorf("getLevelFromScore(%q) = %q, want %q", tt.score, got, tt.want)
			}
		})
	}
}

// ==================== getScoreForLevel ====================

func TestGetScoreForLevel(t *testing.T) {
	tests := []struct {
		name  string
		level string
		want  int
	}{
		{"master node", "Master Node", 20000},
		{"slave 1", "Slave 1", 16000},
		{"slave 2", "Slave 2", 15000},
		{"slave 3", "Slave 3", 14000},
		{"slave 4", "Slave 4", 13000},
		{"unknown level", "Unknown", 0},
		{"empty level", "", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getScoreForLevel(tt.level)
			if got != tt.want {
				t.Errorf("getScoreForLevel(%q) = %d, want %d", tt.level, got, tt.want)
			}
		})
	}
}

// ==================== extractRscID ====================

func TestExtractRscID(t *testing.T) {
	tests := []struct {
		name  string
		opKey string
		want  string
	}{
		{"stop operation", "dummy_stop_0", "dummy"},
		{"start operation", "dummy_start_0", "dummy"},
		{"monitor operation", "dummy_monitor_3000", "dummy"},
		{"demote operation", "ms_demote_0", "ms"},
		{"promote operation", "ms_promote_0", "ms"},
		{"empty key", "", ""},
		{"no known operation", "dummy_validate_0", "dummy_validate_0"},
		{"complex id with stop", "my-rsc_stop_0", "my-rsc"},
		{"complex id with monitor", "ipaddr2_monitor_10000", "ipaddr2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractRscID(tt.opKey)
			if got != tt.want {
				t.Errorf("extractRscID(%q) = %q, want %q", tt.opKey, got, tt.want)
			}
		})
	}
}

// ==================== getRscColocation ====================

func TestGetRscColocation(t *testing.T) {
	xml := `<constraints>
		<rsc_colocation id="coloc-1" rsc="dummy" with-rsc="ip" score="INFINITY"/>
		<rsc_colocation id="coloc-2" rsc="dummy" with-rsc="fs" score="-INFINITY"/>
		<rsc_colocation id="coloc-3" rsc="web" with-rsc="dummy" score="INFINITY" with-rsc-role="Master"/>
		<rsc_colocation id="coloc-4" rsc="db" with-rsc="dummy" score="-INFINITY" rsc-role="Slave"/>
		<rsc_colocation id="coloc-5" rsc="other" with-rsc="other2" score="INFINITY"/>
	</constraints>`

	doc := etree.NewDocument()
	if err := doc.ReadFromString(xml); err != nil {
		t.Fatalf("failed to parse XML: %v", err)
	}
	elements := doc.FindElements("//rsc_colocation")

	t.Run("colocation for dummy", func(t *testing.T) {
		sameNode, diffNode := getRscColocation(elements, "dummy")
		// coloc-1: rsc="dummy" score=INFINITY -> sameNode: "ip"
		// coloc-2: rsc="dummy" score=-INFINITY -> diffNode: "fs"
		// coloc-3: with-rsc="dummy" but with-rsc-role modifies rscWith to "dummy/Master" != "dummy", no match
		// coloc-4: with-rsc="dummy" matches, rsc-role="Slave" -> diffNode: "db/Slave"

		if len(sameNode) != 1 {
			t.Errorf("expected 1 sameNode entry, got %d: %v", len(sameNode), sameNode)
		} else if sameNode[0] != "ip" {
			t.Errorf("sameNode[0] = %q, want %q", sameNode[0], "ip")
		}

		if len(diffNode) != 2 {
			t.Errorf("expected 2 diffNode entries, got %d: %v", len(diffNode), diffNode)
		} else {
			if diffNode[0] != "fs" {
				t.Errorf("diffNode[0] = %q, want %q", diffNode[0], "fs")
			}
			if diffNode[1] != "db/Slave" {
				t.Errorf("diffNode[1] = %q, want %q", diffNode[1], "db/Slave")
			}
		}
	})

	t.Run("colocation for non-existent resource", func(t *testing.T) {
		sameNode, diffNode := getRscColocation(elements, "nonexistent")
		if len(sameNode) != 0 || len(diffNode) != 0 {
			t.Errorf("expected empty results, got sameNode=%v diffNode=%v", sameNode, diffNode)
		}
	})
}

// ==================== getPrimitiveResourceInfo ====================

func TestGetPrimitiveResourceInfo(t *testing.T) {
	xml := `<primitive class="ocf" id="dummy" provider="pacemaker" type="Dummy">
		<operations>
			<op id="dummy-monitor-interval-10s" interval="10s" name="monitor" timeout="20s"/>
			<op id="dummy-start-interval-0s" interval="0s" name="start" timeout="20s"/>
		</operations>
	</primitive>`

	doc := etree.NewDocument()
	if err := doc.ReadFromString(xml); err != nil {
		t.Fatalf("failed to parse XML: %v", err)
	}

	result := getPrimitiveResourceInfo(doc.Root())

	if result.Class != "ocf" {
		t.Errorf("Class = %q, want %q", result.Class, "ocf")
	}
	if result.ID != "dummy" {
		t.Errorf("ID = %q, want %q", result.ID, "dummy")
	}
	if result.Provider != "pacemaker" {
		t.Errorf("Provider = %q, want %q", result.Provider, "pacemaker")
	}
	if result.Type != "Dummy" {
		t.Errorf("Type = %q, want %q", result.Type, "Dummy")
	}
	// Note: getPrimitiveResourceInfo only processes <operations> elements that have an id attribute.
	// Standard Pacemaker XML <operations> has no id, so Operations will be empty.
}

// ==================== getResourceInfoFromXml ====================

func TestGetResourceInfoFromXml_Primitive(t *testing.T) {
	xml := `<primitive class="ocf" id="dummy" provider="pacemaker" type="Dummy">
		<operations>
			<op id="dummy-start-interval-0s" interval="0s" name="start" timeout="20s"/>
		</operations>
	</primitive>`

	doc := etree.NewDocument()
	if err := doc.ReadFromString(xml); err != nil {
		t.Fatalf("failed to parse XML: %v", err)
	}

	result, err := getResourceInfoFromXml("primitive", doc.Root())
	if err != nil {
		t.Fatalf("getResourceInfoFromXml() error = %v", err)
	}

	info, ok := result.(PrimitiveResource)
	if !ok {
		t.Fatalf("expected PrimitiveResource, got %T", result)
	}
	if info.ID != "dummy" {
		t.Errorf("ID = %q, want %q", info.ID, "dummy")
	}
	if info.Class != "ocf" {
		t.Errorf("Class = %q, want %q", info.Class, "ocf")
	}
}

func TestGetResourceInfoFromXml_Group(t *testing.T) {
	xml := `<group id="group1">
		<primitive class="ocf" id="dummy1" provider="pacemaker" type="Dummy"/>
		<primitive class="ocf" id="dummy2" provider="pacemaker" type="Dummy"/>
	</group>`

	doc := etree.NewDocument()
	if err := doc.ReadFromString(xml); err != nil {
		t.Fatalf("failed to parse XML: %v", err)
	}

	result, err := getResourceInfoFromXml("group", doc.Root())
	if err != nil {
		t.Fatalf("getResourceInfoFromXml() error = %v", err)
	}

	info, ok := result.(GroupResource)
	if !ok {
		t.Fatalf("expected GroupResource, got %T", result)
	}
	if info.ID != "group1" {
		t.Errorf("ID = %q, want %q", info.ID, "group1")
	}
	if len(info.Primitives) != 2 {
		t.Errorf("Primitives count = %d, want 2", len(info.Primitives))
	}
}

func TestGetResourceInfoFromXml_Clone(t *testing.T) {
	xml := `<clone id="sysinfo-clone">
		<primitive class="ocf" id="sysinfo" provider="pacemaker" type="SysInfo"/>
	</clone>`

	doc := etree.NewDocument()
	if err := doc.ReadFromString(xml); err != nil {
		t.Fatalf("failed to parse XML: %v", err)
	}

	result, err := getResourceInfoFromXml("clone", doc.Root())
	if err != nil {
		t.Fatalf("getResourceInfoFromXml() error = %v", err)
	}

	info, ok := result.(CloneResource)
	if !ok {
		t.Fatalf("expected CloneResource, got %T", result)
	}
	if info.ID != "sysinfo-clone" {
		t.Errorf("ID = %q, want %q", info.ID, "sysinfo-clone")
	}
	if len(info.Primitives) != 1 {
		t.Errorf("Primitives count = %d, want 1", len(info.Primitives))
	}
}

func TestGetResourceInfoFromXml_MetaAttributes(t *testing.T) {
	xml := `<meta_attributes id="dummy-meta_attributes">
		<nvpair id="dummy-target-role" name="target-role" value="Stopped"/>
		<nvpair id="dummy-priority" name="priority" value="10"/>
	</meta_attributes>`

	doc := etree.NewDocument()
	if err := doc.ReadFromString(xml); err != nil {
		t.Fatalf("failed to parse XML: %v", err)
	}

	result, err := getResourceInfoFromXml("meta", doc.Root())
	if err != nil {
		t.Fatalf("getResourceInfoFromXml() error = %v", err)
	}

	m, ok := result.(map[string]string)
	if !ok {
		t.Fatalf("expected map[string]string, got %T", result)
	}
	if m["target-role"] != "Stopped" {
		t.Errorf("target-role = %q, want %q", m["target-role"], "Stopped")
	}
	if m["priority"] != "10" {
		t.Errorf("priority = %q, want %q", m["priority"], "10")
	}
}

func TestGetResourceInfoFromXml_Operations(t *testing.T) {
	xml := `<operations>
		<op id="dummy-start-0" interval="0s" name="start" timeout="20s"/>
		<op id="dummy-monitor-10s" interval="10s" name="monitor" timeout="20s"/>
	</operations>`

	doc := etree.NewDocument()
	if err := doc.ReadFromString(xml); err != nil {
		t.Fatalf("failed to parse XML: %v", err)
	}

	result, err := getResourceInfoFromXml("operations", doc.Root())
	if err != nil {
		t.Fatalf("getResourceInfoFromXml() error = %v", err)
	}

	ops, ok := result.([]map[string]string)
	if !ok {
		t.Fatalf("expected []map[string]string, got %T", result)
	}
	if len(ops) != 2 {
		t.Errorf("operations count = %d, want 2", len(ops))
	}
}

// ==================== GetResourceInfoID ====================

func TestGetResourceInfoID_Primitive(t *testing.T) {
	xmlData := `<primitive class="ocf" id="dummy" provider="pacemaker" type="Dummy">
		<operations>
			<op id="dummy-start-0" interval="0s" name="start" timeout="20s"/>
		</operations>
	</primitive>`

	result, err := GetResourceInfoID("primitive", xmlData)
	if err != nil {
		t.Fatalf("GetResourceInfoID() error = %v", err)
	}

	if result["id"] != "dummy" {
		t.Errorf("id = %v, want %q", result["id"], "dummy")
	}
	if result["class"] != "ocf" {
		t.Errorf("class = %v, want %q", result["class"], "ocf")
	}
	if result["type"] != "Dummy" {
		t.Errorf("type = %v, want %q", result["type"], "Dummy")
	}
	if result["provider"] != "pacemaker" {
		t.Errorf("provider = %v, want %q", result["provider"], "pacemaker")
	}
}

func TestGetResourceInfoID_Group(t *testing.T) {
	xmlData := `<group id="group1">
		<primitive class="ocf" id="dummy1" provider="pacemaker" type="Dummy"/>
		<primitive class="ocf" id="dummy2" provider="pacemaker" type="Dummy"/>
	</group>`

	result, err := GetResourceInfoID("group", xmlData)
	if err != nil {
		t.Fatalf("GetResourceInfoID() error = %v", err)
	}

	if result["id"] != "group1" {
		t.Errorf("id = %v, want %q", result["id"], "group1")
	}
	rscs, ok := result["rscs"].([]string)
	if !ok {
		t.Fatalf("expected rscs to be []string, got %T", result["rscs"])
	}
	if len(rscs) != 2 {
		t.Errorf("rscs count = %d, want 2", len(rscs))
	}
}

func TestGetResourceInfoID_Clone(t *testing.T) {
	xmlData := `<clone id="sysinfo-clone">
		<primitive class="ocf" id="sysinfo" provider="pacemaker" type="SysInfo"/>
	</clone>`

	result, err := GetResourceInfoID("clone", xmlData)
	if err != nil {
		t.Fatalf("GetResourceInfoID() error = %v", err)
	}

	if result["id"] != "sysinfo-clone" {
		t.Errorf("id = %v, want %q", result["id"], "sysinfo-clone")
	}
	// Single primitive in clone -> rsc_id should be a string
	if result["rsc_id"] != "sysinfo" {
		t.Errorf("rsc_id = %v, want %q", result["rsc_id"], "sysinfo")
	}
}

func TestGetResourceInfoID_CloneMultiplePrimitives(t *testing.T) {
	xmlData := `<clone id="multi-clone">
		<primitive class="ocf" id="rsc1" provider="pacemaker" type="Dummy"/>
		<primitive class="ocf" id="rsc2" provider="pacemaker" type="Dummy"/>
	</clone>`

	result, err := GetResourceInfoID("clone", xmlData)
	if err != nil {
		t.Fatalf("GetResourceInfoID() error = %v", err)
	}

	// Multiple primitives -> rsc_id should be a list
	rscID, ok := result["rsc_id"].([]string)
	if !ok {
		t.Fatalf("expected rsc_id to be []string for multiple primitives, got %T", result["rsc_id"])
	}
	if len(rscID) != 2 {
		t.Errorf("rsc_id count = %d, want 2", len(rscID))
	}
}

func TestGetResourceInfoID_WithMetaAttributes(t *testing.T) {
	xmlData := `<primitive class="ocf" id="dummy" provider="pacemaker" type="Dummy">
		<meta_attributes id="dummy-meta">
			<nvpair id="dummy-target-role" name="target-role" value="Stopped"/>
		</meta_attributes>
	</primitive>`

	result, err := GetResourceInfoID("primitive", xmlData)
	if err != nil {
		t.Fatalf("GetResourceInfoID() error = %v", err)
	}

	meta, ok := result["meta_attributes"].(map[string]string)
	if !ok {
		t.Fatalf("expected meta_attributes to be map[string]string, got %T", result["meta_attributes"])
	}
	if meta["target-role"] != "Stopped" {
		t.Errorf("target-role = %q, want %q", meta["target-role"], "Stopped")
	}
}

func TestGetResourceInfoID_WithInstanceAttributes(t *testing.T) {
	xmlData := `<primitive class="ocf" id="dummy" provider="pacemaker" type="Dummy">
		<instance_attributes id="dummy-inst">
			<nvpair id="dummy-state" name="state" value="/var/run/dummy"/>
		</instance_attributes>
	</primitive>`

	result, err := GetResourceInfoID("primitive", xmlData)
	if err != nil {
		t.Fatalf("GetResourceInfoID() error = %v", err)
	}

	inst, ok := result["instance_attributes"].(map[string]string)
	if !ok {
		t.Fatalf("expected instance_attributes to be map[string]string, got %T", result["instance_attributes"])
	}
	if inst["state"] != "/var/run/dummy" {
		t.Errorf("state = %q, want %q", inst["state"], "/var/run/dummy")
	}
}

// ==================== GetResourceMetaAttributes ====================

func TestGetResourceMetaAttributes(t *testing.T) {
	t.Run("group category", func(t *testing.T) {
		result := GetResourceMetaAttributes("group")
		if result["action"] != true {
			t.Errorf("action = %v, want true", result["action"])
		}
		data, ok := result["data"].(map[string]map[string]interface{})
		if !ok {
			t.Fatalf("expected data to be map[string]map[string]interface{}, got %T", result["data"])
		}
		// Group should NOT have resource-stickiness or migration-threshold
		if _, exists := data["resource-stickiness"]; exists {
			t.Error("group should not have resource-stickiness")
		}
		if _, exists := data["migration-threshold"]; exists {
			t.Error("group should not have migration-threshold")
		}
		// Should have target-role, priority, is-managed
		if _, exists := data["target-role"]; !exists {
			t.Error("group should have target-role")
		}
	})

	t.Run("non-group category", func(t *testing.T) {
		result := GetResourceMetaAttributes("primitive")
		if result["action"] != true {
			t.Errorf("action = %v, want true", result["action"])
		}
		data, ok := result["data"].(map[string]map[string]interface{})
		if !ok {
			t.Fatalf("expected data to be map[string]map[string]interface{}, got %T", result["data"])
		}
		// Non-group should have resource-stickiness and migration-threshold
		if _, exists := data["resource-stickiness"]; !exists {
			t.Error("primitive should have resource-stickiness")
		}
		if _, exists := data["migration-threshold"]; !exists {
			t.Error("primitive should have migration-threshold")
		}
	})
}

// ==================== safeResourceName / safePeriod / safeIdentifier regex ====================

func TestSafeResourceName(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"dummy", true},
		{"my-resource", true},
		{"my_resource", true},
		{"rsc:0", true},
		{"rsc.node", true},
		{"rsc-1_v2:3.4", true},
		{"", false},
		{"has space", false},
		{"has/slash", false},
		{"has@symbol", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := safeResourceName.MatchString(tt.input)
			if got != tt.want {
				t.Errorf("safeResourceName.MatchString(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestSafePeriod(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"start", true},
		{"stop", true},
		{"promote0", true},
		{"", false},
		{"has-hyphen", false},
		{"has space", false},
		{"has_underscore", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := safePeriod.MatchString(tt.input)
			if got != tt.want {
				t.Errorf("safePeriod.MatchString(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestSafeIdentifier(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"ocf", true},
		{"pacemaker", true},
		{"my-type", true},
		{"my_type", true},
		{"type123", true},
		{"", false},
		{"has:colon", false},
		{"has.dot", false},
		{"has space", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := safeIdentifier.MatchString(tt.input)
			if got != tt.want {
				t.Errorf("safeIdentifier.MatchString(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
