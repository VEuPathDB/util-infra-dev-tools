#!/usr/bin/env sh

repoUrl="https://github.com/VEuPathDB/util-vpdb-dev-tool"
releaseSuffix="/releases/latest"

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