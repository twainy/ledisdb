#!/bin/bash

export VTTOP=$(pwd)
export VTROOT="${VTROOT:-${VTTOP/\/src\/github.com\/siddontang\/ledisdb/}}"
# VTTOP sanity check
if [[ "$VTTOP" == "${VTTOP/\/src\/github.com\/siddontang\/ledisdb/}" ]]; then
    echo "WARNING: VTTOP($VTTOP) does not contain src/github.com/siddontang/ledisdb"
    exit 1
fi


#default snappy and leveldb install path
#you may change yourself

# SNAPPY_DIR=/usr/local/snappy
# LEVELDB_DIR=/usr/local/leveldb

function add_path()
{
  # $1 path variable
  # $2 path to add
  if [ -d "$2" ] && [[ ":$1:" != *":$2:"* ]]; then
    echo "$1:$2"
  else
    echo "$1"
  fi
}

export GOPATH=$(add_path $GOPATH $VTROOT)
export GOPATH=$(add_path $VTROOT/_vendor $GOPATH)
