#!/bin/bash

REPO_PATH=$PIGEON_REPO_PATH

if [ -z "$REPO_PATH" ]; then
  echo 'Error: please set REPO_PATH environment variable'
  exit 1
fi

pushd ${REPO_PATH}
RESULT_PULL="$(hg pull 2>&1) | grep added"
echo ${RESULT_PULL}

if [ -z "$RESULT_PULL" ]; then
  echo "Nothing new. exit."
  exit 1
fi

hg update -C
sudo cp /scripts/stool-pigeon.sh /usr/local/bin
popd
