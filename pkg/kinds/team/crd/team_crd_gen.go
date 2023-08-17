// Code generated - EDITING IS FUTILE. DO NOT EDIT.
//
// Generated by:
//     kinds/gen.go
// Using jennies:
//     CRDTypesJenny
//
// Run 'make gen-cue' from repository root to regenerate.

package crd

import (
	_ "embed"

	"github.com/grafana/grafana/pkg/kinds/team"
	"github.com/grafana/kindsys/k8ssys"
)

// The CRD YAML representation of the Team kind.
//
//go:embed team.crd.yml
var CRDYaml []byte

// Team is the Go CRD representation of a single Team object.
// It implements [runtime.Object], and is used in k8s scheme construction.
type Team struct {
	k8ssys.Base[team.Team]
}

// TeamList is the Go CRD representation of a list Team objects.
// It implements [runtime.Object], and is used in k8s scheme construction.
type TeamList struct {
	k8ssys.ListBase[team.Team]
}
