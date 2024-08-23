package use

import (
	"errors"
	"fmt"

	cli "github.com/Foxcapades/Argonaut"
	"github.com/Foxcapades/Argonaut/pkg/argo"

	"vpdb-dev-tool/internal/lib/xos"
)

const (
	leafDesc = "Updates the local .env file to pin the stack image versions to" +
		" a specific set of images."

	argDesc = "Environment to mimic.\n\n" +
		"May be one of:\n" +
		"- latest\n" +
		"- qa\n" +
		"- prod"

	fileFlagDesc = "Specifies a docker compose file containing images whose" +
		" versions should be pinned.\n" +
		"May be provided more than once.\n\n" +
		"If unused, then 'docker-compose.yml' will be assumed."

	backupFlagDesc = "Backup .env file (if exist) before writing modifications.\n\n" +
		"May optionally be used to specify the backup name if desired."

	noBackupIndicator = "I, the user of this cli tool, hereby certify that I do" +
		" not want a backup made of my '.env' file."
)

type options struct {
	target string
	files  []string
	backup string
}

func Init(branch argo.CommandBranchBuilder) {
	var opts options

	branch.WithLeaf(cli.Leaf("use").
		WithDescription(leafDesc).
		WithArgument(cli.Argument().
			WithName("version").
			WithDescription(argDesc).
			WithValidator(func(raw string) error {
				if raw == "qa" || raw == "prod" || raw == "latest" {
					return nil
				} else {
					return errors.New(`invalid stack version specified, must be one of "qa", "prod", or "latest"`)
				}
			}).
			WithBinding(&opts.target).
			Require()).
		WithFlag(cli.ComboFlag('b', "make-backup").
			WithDescription(backupFlagDesc).
			WithBindingAndDefault(&opts.backup, noBackupIndicator, false)).
		WithFlag(cli.ComboFlag('f', "compose-file").
			WithDescription(fileFlagDesc).
			WithArgument(cli.Argument().
				WithName("path").
				WithBinding(&opts.files).
				WithValidator(func(path string) error {
					if ok, err := xos.PathExists(path); err != nil {
						return err
					} else if !ok {
						return fmt.Errorf(`specified compose file "%s" does not exist`, path)
					} else {
						return nil
					}
				}).
				Require())).
		WithCallback(func(_ argo.CommandLeaf) { run(opts) }))
}
