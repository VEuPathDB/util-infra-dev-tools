package xyml

import (
	"gopkg.in/yaml.v3"
	"strings"
	"vpdb-dev-tool/internal/lib/must"
	"vpdb-dev-tool/internal/lib/xos"
)

func Stringify(node *yaml.Node) string {
	buf := strings.Builder{}
	buf.Grow(32768)

	enc := yaml.NewEncoder(&buf)
	defer xos.MustClose(enc)

	enc.SetIndent(2)

	must.NotError(enc.Encode(node))

	return buf.String()
}
