package hosts

import (
	"log"
	"os"
)

const (
	hostList = "host-list.yml"
)

// TODO: allow wildcards that match to service entries in other docker-compose files in the project.

func ReadHostsFile(path string) map[Host][]string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("encountered error while attempting to open host list file %s: %s\n", file, err)
	}
	defer file.Close()

	return ReadHosts(file)
}

func MakeHostsList() {
	file, err := os.OpenFile(hostList, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0664)
	if err != nil {
		log.Fatalf("failed to create file %s: %s\n", hostList, err)
	}
	defer file.Close()

	_, err = file.WriteString(`# Each entry in the 'hosts' map is an array of container definitions that depend
# on access to that host.  The entries in the list should be the names of the
# dependent containers as defined in the primary docker-compose 'services' block
# or blocks.
#
# The 'hosts' map keys MUST be "host:port" pairs.
hosts:
  ldap.foo.bar:389:
  - some-service-name
  - some-other-service
  s3.foo.bar:9000:
  - file-service
`)

	if err != nil {
		log.Fatalf("failed to write contents out to file %s: %s", hostList, err)
	}
}
