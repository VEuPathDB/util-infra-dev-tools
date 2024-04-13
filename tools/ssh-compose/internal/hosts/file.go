package hosts

import (
	"log"
	"os"
)

func ReadHostsFile(path string) map[Host][]string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("encountered error while attempting to open host list file %s: %s\n", file, err)
	}
	defer file.Close()

	return ReadHosts(file)
}
