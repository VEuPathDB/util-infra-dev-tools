package project

import (
	"os"
	"vpdb-dev-tool/internal/lib/must"

	E "vpdb-dev-tool/internal/lib/env"
	"vpdb-dev-tool/internal/lib/xos"
	"vpdb-dev-tool/internal/tools/ssh_compose/env"
)

func patchEnvFile(file *os.File, hosts map[string]string) {
	editor := E.NewEditor().
		AddIfAbsent(env.TunnelHost, "").
		AddIfAbsent(env.TunnelPort, "").
		AddIfAbsent(env.TunnelUser, "")

	tmpFile := xos.MustCreateFile("/tmp/env-backup")
	defer xos.MustClose(tmpFile)

	for key := range hosts {
		editor.AddIfAbsent(key, "")
	}

	editor.ApplyEdits(file, tmpFile)
	xos.MustClose(tmpFile)
	xos.MustDelete(file.Name())

	must.NotError(xos.CopyFile(tmpFile.Name(), file.Name()))
	xos.MustDelete(tmpFile.Name())
}
