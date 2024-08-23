package use

import (
	"bytes"
	"fmt"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"vpdb-dev-tool/internal/lib/env"
	"vpdb-dev-tool/internal/lib/gh"
	"vpdb-dev-tool/internal/lib/util"
	"vpdb-dev-tool/internal/lib/xos"
)

const (
	tmpEnvFileName = "HOLY_CRIPES_DO_NOT_COMMIT_THIS_FILE"
	defaultCompose = "docker-compose.yml"
)

func run(opts options) {
	if len(opts.files) == 0 {
		ensureDefaultCompose()
		opts.files = []string{defaultCompose}
	}

	// if the user specified just "-b" with no followup arg, make a default path
	if opts.backup == "" {
		opts.backup = makeDefaultBackupPath()
	}

	// fetch a map of image names to env var names
	images := scanConfigs(opts.files)

	logrus.Debugf("will attempt to set env vars for %d image tags", len(images))

	mods := env.NewEditor()

	// If the user requested "latest" then remove any env settings?
	if opts.target == "latest" {
		for _, v := range images {
			mods.AddOrReplace(v, "latest")
		}

		fmt.Println("\n\nRemember to pull down any project images to avoid stale 'latest'!\n")
	} else {
		versions := loadTaggerFile(opts.target)

		for image, envKey := range images {
			if version, ok := versions[image]; ok {
				mods.AddOrReplace(envKey, version)
			} else {
				logrus.Warningf("No version found in tagger for image %s:${%s}", image, envKey)
			}
		}
	}

	if opts.backup != noBackupIndicator {
		if xos.MustPathExists(env.DotEnvFileName) {
			xos.MustCopyFile(env.DotEnvFileName, opts.backup)
			logrus.Infof("created backup of .env file: %s", opts.backup)
		} else {
			logrus.Info("no pre-existing .env file was found, so no backup was made")
		}
	} else {
		logrus.Debugf("no .env file backup flag was specified, not making a backup")
	}

	envFile := env.GetOrCreateDotEnvFile()
	defer xos.MustClose(envFile)
	tmpFile := xos.MustCreateFile(tmpEnvFileName)
	defer xos.MustClose(tmpFile)

	mods.ApplyEdits(envFile, tmpFile)

	xos.MustClose(envFile)
	xos.MustClose(tmpFile)

	util.Must(os.Rename(tmpEnvFileName, env.DotEnvFileName))
}

func loadTaggerFile(target string) (versions map[string]string) {
	ghCreds := gh.RequireCredentials()
	taggerFileBytes := util.MustReturn(gh.GetPrivateFileContents("tagger", "versions.yml", target, ghCreds.Token))

	decoder := yaml.NewDecoder(bytes.NewReader(taggerFileBytes))
	util.Must(decoder.Decode(&versions))

	return
}

func makeDefaultBackupPath() string {
	nameBase := path.Join(util.MustReturn(os.UserHomeDir()), path.Base(util.MustReturn(os.Getwd()))+".env.backup.%d")
	distinct := 1

	for {
		attempt := fmt.Sprintf(nameBase, distinct)

		if !xos.MustPathExists(attempt) {
			return attempt
		}

		distinct++
	}
}

func ensureDefaultCompose() {
	if ok, err := xos.PathExists(defaultCompose); err != nil {
		logrus.Fatalf("failed to stat %s: %s", defaultCompose, err)
		panic(err) // unreachable
	} else if !ok {
		logrus.Fatal("no docker compose file was specified and docker-compose.yml could not be found")
		panic(err) // unreachable
	}
}
