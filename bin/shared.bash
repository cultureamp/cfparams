#!/usr/bin/env bash

function inline_link() {
  LINK=$(printf "url='%s'" "$1")

  if [ $# -gt 1 ]; then
    LINK=$(printf "$LINK;content='%s'" "$2")
  fi

  printf '\033]1339;%s\a\n' "$LINK"
}

function finish() {
  # Did the previous command fail? Then make Buildkite
  # auto-expand the build log for it.
  if [ "$?" -gt 0 ]; then
    echo "^^^ +++"
  fi
}