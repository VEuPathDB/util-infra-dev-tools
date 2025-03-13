#!/usr/bin/env bash

function __vpdb_get_help() {
  vpdb $* -h \
    | awk '
      BEGIN { mode = 0; }
      {
        if ($1 == "Flags") {
          mode = 1;
        } else if ($1 == "Commands") {
          mode = 2;
        } else if (mode == 1 && $0 ~ /^  -/) {
          if ($3 == "|") {
            print $1, gensub(/=.+/, "", "g", $4);
          } else if ($3 ~ /^--/) {
            print $1, gensub(/=.+/, "", "g", $3);
          } else {
            print $1;
          }
        } else if (mode == 2 && $0 ~ /^  [^ ]/) {
          print $1;
        }
      }' \
    | tr '\n' ' '
}

function __vpdb_completion_ssh_compose() {
  # --gen-example is a complete command
  if grep -q "\--gen-example" <<< "${REMAINING_LINE}"; then
    return
  fi

  # If they aren't looking for a flag, use files
  if [[ "${COMP_WORDS[$((COMP_CWORD-1))]}" =~ ^-(s|-ssh-home)$ ]]; then
    COMPREPLY=($(compgen -f -- "${CURRENT_WORD}"))
    return
  fi

  if [[ "${COMP_WORDS[$((COMP_CWORD-1))]}" =~ ^-(i|-image)$ ]]; then
    if [[ -n "${CURRENT_WORD}" ]]; then
      COMPREPLY=($(docker search --limit 5 "${CURRENT_WORD}" --format '{{.Name}} {{.StarCount}}' | sort -rk2 | cut -f1 -d' ' | tr '\n' ' '))
    fi
    return
  fi

  if [[ "${CURRENT_WORD}" == -* ]]; then
    if ! grep -q '\-i\|--image' <<< "${COMP_LINE}"; then
      COMPREPLY+=(-i --image)
    fi

    if ! grep -q '\-s\|--ssh-home' <<< "${COMP_LINE}"; then
      COMPREPLY+=(-s --ssh-home)
    fi

    if [ "${VERBOSE_LEVEL}" -lt 3 ]; then
      COMPREPLY+=(-V)
    fi

    return
  fi

  if ! grep -q '.*\.ya\?ml' <<< "${COMP_LINE}"; then
    COMPREPLY=($(compgen -C "find . -maxdepth 1 '(' -name '*.yml' -o -name '*.yaml' ')' -printf '%P '" -- "${CURRENT_WORD}"))
    return
  fi

  COMPREPLY=(-h --help -v --version)

  if ! grep -q '\-i\|--image' <<< "${COMP_LINE}"; then
    COMPREPLY+=(-i --image)
  fi

  if ! grep -q '\-s\|--ssh-home' <<< "${COMP_LINE}"; then
    COMPREPLY+=(-s --ssh-home)
  fi

  if [ "${VERBOSE_LEVEL}" -lt 3 ]; then
    COMPREPLY+=(-V)
  fi
}

function __vpdb_completion_vdi() {
  local -a COMMANDS=(vdi)

  if grep -q "gen-tagger" <<< "${REMAINING_LINE}"; then
    if ! grep -q "\-w\|--write-version"; then
      return
    fi

    COMMANDS+=(gen-tagger)
  fi

  COMPREPLY=($(__vpdb_get_help ${COMMANDS[@]}))
  if [ "${VERBOSE_LEVEL}" -gt 2 ]; then
    COMPREPLY=($(sed 's/-V\+//g' <<< "${COMPREPLY[@]}"))
  fi
}

function __vpdb_completion_stack() {
  COMPREPLY=($(__vpdb_get_help stack))
}

function __vpdb_completion_merge_compose() {
  COMPREPLY=($(__vpdb_get_help merge-compose))
}

function __vpdb_completion() {
  COMPREPLY=()

  if [ ${#COMP_WORDS[@]} -eq 1 ]; then
    return
  fi

  # There is no valid operation after a help or version flag.
  if grep -q '\-h\|--help\|-v\|--version' <<< "${COMP_LINE}"; then
    return
  fi

  declare -ra TOP_LEVEL_COMMANDS=(
    "ssh-compose"
    "vdi"
    "stack"
    "merge-compose"
  )

  declare -ri VERBOSE_LEVEL="$(grep -o '\-V\+' <<< "${COMP_LINE}" | sed 's/-//' | tr -d '\n' | wc -c)"

  declare -r CURRENT_WORD="${COMP_WORDS[${COMP_CWORD}]}"

  for COMMAND in "${TOP_LEVEL_COMMANDS[@]}"; do
    if grep -q "${COMMAND}" <<< "${COMP_LINE}"; then
      local -r REMAINING_LINE="$(sed "s/vpdb\|${COMMAND}\|-V\+//g" <<< "${COMP_LINE}")"
      __vpdb_completion_$(sed 's/-/_/g' <<< "${COMMAND}")
      return
    fi
  done

  COMPREPLY=($(compgen -W "$(__vpdb_get_help)"))

  if [ ${VERBOSE_LEVEL} -gt 2 ]; then
    COMPREPLY=($(sed 's/-V\+//g' <<< "${COMPREPLY[@]}"))
  fi
}

complete -F __vpdb_completion vpdb