# ledisdb

Ledisdb is a high performance nosql like redis based on leveldb written by go. It's supports some advanced data structure like kv, list, hash and zset.

## Build and Install

+ Create a workspace and checkout ledisdb source

        mkdir $WORKSPACE
        cd $WORKSPACE
        git clone git@github.com:siddontang/ledisdb.git src/github.com/siddontang/ledisdb

        cd src/github.com/siddontang/ledisdb

+ Then:

        . ./bootstap.sh 
        . ./dev.sh

        go install ./...

## Run

    ./ledis-server -config=/etc/ledis.json

    //another shell
    redis-cli -p 6380
    
    redis 127.0.0.1:6380> set a 1
    OK
    redis 127.0.0.1:6380> get a
    "1"

## Lib
    
    import "github.com/siddontang/ledisdb/ledis"
    l, _ := ledis.OpenWithConfig(cfg)
    db, _ := l.Select(0)

    db.Set(key, value)

    db.Get(key)


## Replication

set slaveof in config or dynamiclly

    redis-cli -p 6381 

    redis 127.0.0.1:6381> slaveof 127.0.0.1:6380
    OK

## Benchmark

See benchmark.md for more.

## Todo

+ Admin

## Thanks

Gamil: cenqichao@gmail.com

## Feedback

Gmail: siddontang@gmail.com