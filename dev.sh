#!/bin/bash

export VTTOP=$(pwd)
export VTROOT=$(pwd)
# VTTOP sanity check

#default snappy and leveldb install path
#you may change yourself

SNAPPY_DIR=/usr/local/snappy
LEVELDB_DIR=/usr/local/leveldb

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

export CGO_CFLAGS="-I$LEVELDB_DIR/include -I$SNAPPY_DIR/include"
export CGO_CXXFLAGS="-I$LEVELDB_DIR/include -I$SNAPPY_DIR/include"
export CGO_LDFLAGS="-L$LEVELDB_DIR/lib -L$SNAPPY_DIR/lib -lsnappy"

#for linux, use LD_LIBRARY_PATH
export LD_LIBRARY_PATH=$(add_path $LD_LIBRARY_PATH $SNAPPY_DIR/lib)
export LD_LIBRARY_PATH=$(add_path $LD_LIBRARY_PATH $LEVELDB_DIR/lib)

#for macos, use DYLD_LIBRARY_PATH
export DYLD_LIBRARY_PATH=$(add_path $DYLD_LIBRARY_PATH $SNAPPY_DIR/lib)
export DYLD_LIBRARY_PATH=$(add_path $DYLD_LIBRARY_PATH $LEVELDB_DIR/lib)
