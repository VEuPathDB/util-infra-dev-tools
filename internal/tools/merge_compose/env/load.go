package env

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"vpdb-dev-tool/internal/lib/xbytes"
	"vpdb-dev-tool/internal/lib/xos"
)

var emptyBytes = [0]byte{}

func ProcessUserEnv(into map[string]string) {
	for _, value := range os.Environ() {
		pos := strings.IndexByte(value, '=')

		if pos == -1 {
			into[value] = ""
		} else {
			into[value[:pos]] = value[pos+1:]
		}
	}
}

func LoadMany(paths []string, into map[string]string) {
	for _, path := range paths {
		Load(path, into)
	}
}

func Load(path string, into map[string]string) {
	file := xos.MustOpen(path, os.O_RDONLY, 0644)
	defer xos.MustClose(file)

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		processLine(scanner.Bytes(), lineNum, path, into)
	}

	if err := scanner.Err(); err != nil {
		logrus.Fatalf("error while scanning env file %s: %s", path, err.Error())
	}
}

func processLine(line []byte, lineNum int, path string, into map[string]string) {
	line = skipLeadingSpace(line)

	// If the line is empty, or starts with a pound sign
	if len(line) == 0 || line[0] == '#' {
		return
	}

	key, value := splitEnvLine(line, lineNum, path)

	if len(key) == 0 {
		logrus.Fatalf("no key found for value on line %d in file %s", lineNum, path)
	}

	into[string(key)] = string(value)
}

func skipLeadingSpace(line []byte) []byte {
	for i, b := range line {
		if !xbytes.IsWhitespace(b) {
			return line[i:]
		}
	}

	return emptyBytes[:]
}

func splitEnvLine(line []byte, lineNum int, path string) (name []byte, value []byte) {
	for i, b := range line {
		if !xbytes.IsWordChar(b) {
			if b == '=' {
				return line[:i], line[i+1:]
			}
			logrus.Fatalf("encountered illegal character on line %d in file %s", lineNum, path)
		}
	}

	logrus.Fatalf("no value set for env var on line %d in file %s", lineNum, path)
	return line, nil // unreachable
}
