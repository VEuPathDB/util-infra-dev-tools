package project

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"ssh-compose/internal/compose"
	"ssh-compose/internal/env"
	"ssh-compose/internal/tunnel"
)

const (
	outputComposeFileName = "docker-compose.ssh.yml"
	envFileName           = ".env"
	gitIgnoreFileName     = ".gitignore"
)

func WriteOutConfigs(configs tunnel.BuiltConfigs, tunnelConfigFile string) {
	writeOutComposeFile(&configs.Compose)
	writeOutEnvFile(configs.Hosts)
	writeOutGitIgnoreFile(tunnelConfigFile)
}

func writeOutComposeFile(config *compose.Config) {
	file := requireFile(outputComposeFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	defer file.Close()

	enc := yaml.NewEncoder(file)
	enc.SetIndent(2)

	err := enc.Encode(config)
	if err != nil {
		log.Fatalf("failed to write docker-compose config out to file: %s\n", err)
	}
}

func writeOutEnvFile(hosts map[string]string) {
	file := requireFile(envFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND)
	defer file.Close()

	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		log.Fatalf("failed to fast-forward to end of file %s: %s\n", envFileName, err)
	}

	if _, err := file.WriteString(env.BuildNewEnv(hosts)); err != nil {
		log.Fatalf("failed to write env vars: %s", err)
	}
}

func writeOutGitIgnoreFile(tunnelConfigFile string) {
	_, err := os.Stat(gitIgnoreFileName)
	if err != nil {
		if os.IsNotExist(err) {
			createGitIgnoreFile(tunnelConfigFile)
			return
		}

		log.Fatalf("failed to stat file %s: %s\n", gitIgnoreFileName, err)
	}

	patchGitIgnoreFile(tunnelConfigFile)
}

func patchGitIgnoreFile(tunnelConfigFile string) {
	file := requireFile(gitIgnoreFileName, os.O_RDWR|os.O_APPEND)
	defer file.Close()

	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		log.Fatalf("failed to fast-forward to end of file %s: %s\n", gitIgnoreFileName, err)
	}

	buff := bufio.NewWriter(file)

	requireWriteLn(buff, gitIgnoreFileName, outputComposeFileName)
	requireWriteLn(buff, gitIgnoreFileName, tunnelConfigFile)
}

func createGitIgnoreFile(tunnelConfigFile string) {
	file := requireFile(gitIgnoreFileName, os.O_WRONLY|os.O_CREATE|os.O_EXCL)
	defer file.Close()

	buff := bufio.NewWriter(file)

	ignore := []string{
		".idea/",
		".tools/",
		".bin/",
		".settings/",
		"/bin/",
		".env",
		".envrc",
		"*.tmp",
		"*.swp",
		"*.jar",
		"!gradle-wrapper.jar",
		"*.bak",
		".classpath",
		".project",
		".settings",
		"*~",
		".*~",
		outputComposeFileName,
		tunnelConfigFile,
	}

	if hasGradle() {
		ignore = append(ignore, ".gradle/", "build/")
	}

	for _, entry := range ignore {
		requireWriteLn(buff, gitIgnoreFileName, entry)
	}
}

func hasGradle() bool {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to lookup CWD")
	}

	entries, err := os.ReadDir(cwd)
	if err != nil {
		log.Fatalf("failed to list CWD directory contents")
	}

	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "build.gradle") {
			return true
		}
	}

	return false
}

func requireFile(path string, flags int) *os.File {
	file, err := os.OpenFile(path, flags, 0664)
	if err != nil {
		log.Fatalf("failed to open file %s for writing: %s\n", path, err)
	}

	return file
}

func requireWriteLn(b *bufio.Writer, file, line string) {
	if _, err := b.WriteString(line); err != nil {
		log.Fatalf("encountered error while writing to file %s: %s\n", file, err)
	}
	if err := b.WriteByte('\n'); err != nil {
		log.Fatalf("encountered error while writing to file %s: %s\n", file, err)
	}
}
