#!/usr/bin/env sh

repoUrl="https://github.com/VEuPathDB/util-vpdb-dev-tool"
releaseSuffix="/releases/latest"

checkPrereqs() {
  if ! command -v curl; then
    echo "curl is required to run the automatic install script" 1>&2
    exit 1
  fi
}

getReleases() {
  curl -L https://api.github.com/repos/VEuPathDB/util-vpdb-dev-tool/releases
}

case $(uname | tr '[:upper:]' '[:lower:]') in
  linux*)

    ;;
  darwin*)
    ;;
  *)
    echo "unsupported os" 1>&2
    exit 1
    ;;
esac