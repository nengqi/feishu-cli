package cmd

import (
	"encoding/json"
	"testing"
)

func TestParseBoardExportNodesArray(t *testing.T) {
	nodes, err := parseBoardExportNodes(json.RawMessage(`[{"id":"n1","type":"svg"}]`))
	if err != nil {
		t.Fatalf("parseBoardExportNodes() error = %v", err)
	}
	if len(nodes) != 1 || nodes[0]["id"] != "n1" {
		t.Fatalf("nodes = %#v", nodes)
	}
}

func TestParseBoardExportNodesMap(t *testing.T) {
	nodes, err := parseBoardExportNodes(json.RawMessage(`{"n1":{"type":"svg"}}`))
	if err != nil {
		t.Fatalf("parseBoardExportNodes() error = %v", err)
	}
	if len(nodes) != 1 || nodes[0]["id"] != "n1" || nodes[0]["type"] != "svg" {
		t.Fatalf("nodes = %#v", nodes)
	}
}
