package vdi

import (
	cli "github.com/Foxcapades/Argonaut"
	"github.com/Foxcapades/Argonaut/pkg/argo"

	"vpdb-dev-tool/internal/tools/vdi/docker"
)

func Init(builder argo.CommandTreeBuilder) {
	branch := cli.Branch("vdi")

	vdi_docker.Init(branch)

	builder.WithBranch(branch)
}
