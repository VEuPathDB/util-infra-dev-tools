package vdi_docker

import (
	cli "github.com/Foxcapades/Argonaut"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

const (
	genTagDesc = "Generates YAML map entries for the latest VDI docker image tags and prints them on STDOUT."
)

func Init(branch argo.CommandBranchBuilder) {
	branch.WithLeaf(cli.Leaf("gen-tagger").
		WithDescription(genTagDesc).
		WithCallback(func(_ argo.CommandLeaf) { mainGenTagger() }))
}
