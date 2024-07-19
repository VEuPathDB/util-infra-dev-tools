package project

import (
	"bufio"
	"log"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"

	"vpdb-dev-tool/internal/lib/xio"
	"vpdb-dev-tool/internal/tools/ssh_compose/compose"
	"vpdb-dev-tool/internal/tools/ssh_compose/tunnel"
)

const (
	OutputComposeFileName = "docker-compose.ssh.yml"
	EnvFileName           = ".env"
	GitIgnoreFileName     = ".gitignore"
)

func WriteOutConfigs(configs tunnel.BuiltConfigs, tunnelConfigFile string) {
	writeOutComposeFile(&configs.Compose)
	writeOutEnvFile(configs.Hosts)
	writeOutGitIgnoreFile(path.Base(tunnelConfigFile))
}

func writeOutComposeFile(config *compose.Config) {
	file := requireFile(OutputComposeFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	defer file.Close()

	enc := yaml.NewEncoder(file)
	enc.SetIndent(2)

	err := enc.Encode(config)
	if err != nil {
		log.Fatalf("failed to write docker-compose config out to file: %s\n", err)
	}
}

func writeOutEnvFile(hosts map[string]string) {
	file := requireFile(EnvFileName, os.O_RDWR|os.O_CREATE)
	defer file.Close()

	patchEnvFile(xio.ReqRWFile{File: file}, hosts)
}

func writeOutGitIgnoreFile(tunnelConfigFile string) {
	_, err := os.Stat(GitIgnoreFileName)
	if err != nil {
		if os.IsNotExist(err) {
			createGitIgnoreFile(tunnelConfigFile)
			return
		}

		log.Fatalf("failed to stat file %s: %s\n", GitIgnoreFileName, err)
	}

	patchGitIgnoreFile(tunnelConfigFile)
}

func patchGitIgnoreFile(tunnelConfigFile string) {
	file := requireFile(GitIgnoreFileName, os.O_RDWR|os.O_APPEND)
	defer file.Close()

	entries := make(map[string]bool, 32)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if len(line) == 0 || line[0] == '#' {
			continue
		}

		entries[line] = true
	}

	if scanner.Err() != nil {
		log.Fatalf("encountered error while scanning %s: %s", GitIgnoreFileName, scanner.Err())
	}

	buff := bufio.NewWriter(file)

	if !entries[OutputComposeFileName] {
		requireWriteLn(buff, GitIgnoreFileName, OutputComposeFileName)
	}

	if !entries[tunnelConfigFile] {
		requireWriteLn(buff, GitIgnoreFileName, tunnelConfigFile)
	}

	if err := buff.Flush(); err != nil {
		log.Fatalf("failed to write buffer to file %s: %s", GitIgnoreFileName, err)
	}
}

func createGitIgnoreFile(tunnelConfigFile string) {
	file := requireFile(GitIgnoreFileName, os.O_WRONLY|os.O_CREATE|os.O_EXCL)
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
		OutputComposeFileName,
		tunnelConfigFile,
	}

	if hasGradle() {
		ignore = append(ignore, ".gradle/", "build/")
	}

	for _, entry := range ignore {
		requireWriteLn(buff, GitIgnoreFileName, entry)
	}

	if err := buff.Flush(); err != nil {
		log.Fatalf("failed to write buffer to file %s: %s", GitIgnoreFileName, err)
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
