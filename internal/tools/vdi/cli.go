package vdi

import (
	cli "github.com/Foxcapades/Argonaut"
	"github.com/Foxcapades/Argonaut/pkg/argo"

	"vpdb-dev-tool/internal/tools/vdi/docker"
)

const branchDesc = "VDI related utility operations."

func Init(builder argo.CommandTreeBuilder) {
	branch := cli.Branch("vdi").
		WithDescription(branchDesc)

	vdi_docker.Init(branch)

	builder.WithBranch(branch)
}
