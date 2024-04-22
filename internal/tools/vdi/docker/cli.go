package vdi_docker

import (
	cli "github.com/Foxcapades/Argonaut"
	"github.com/Foxcapades/Argonaut/pkg/argo"
)

const (
	genTagDesc = "Generates YAML map entries for the latest VDI docker image tags and prints them on STDOUT."

	writeFlagDesc = "Write versions out to versions.yml file.  Command will fail if versions.yml file does not already " +
		"exist in the current working directory."
)

type cliOpts struct {
	writeToVersionsFile bool
}

func Init(branch argo.CommandBranchBuilder) {
	var opts cliOpts

	branch.WithLeaf(cli.Leaf("gen-tagger").
		WithDescription(genTagDesc).
		WithFlag(cli.ComboFlag('w', "write-versions").
			WithDescription(writeFlagDesc).
			WithBinding(&opts.writeToVersionsFile, true)).
		WithCallback(func(_ argo.CommandLeaf) { mainGenTagger(opts) }))
}
