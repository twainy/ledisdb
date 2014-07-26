# LedisDB

Ledisdb is a high performance NoSQL like Redis based on LevelDB written by go. It supports some advanced data structure like kv, list, hash and zset, and may be alternative for Redis.

![travis_img](https://api.travis-ci.org/twainy/ledisdb.svg)

## Features

+ Rich advanced data structure: KV, List, Hash, ZSet, Bit.
+ Uses leveldb to store lots of data, over the memory limit. 
+ Supports expiration and ttl.
+ Redis clients, like redis-cli, are supported directly.
+ Multi client API supports, including Golang, Python, Lua(Openresty). 
+ Easy to embed in Golang application. 
+ Replication to guarantee data safe.
+ Supplies tools to load, dump, repair database. 

## Build and Install

+ Create a workspace and checkout ledisdb source

        mkdir $WORKSPACE
        cd $WORKSPACE
        git clone git@github.com:siddontang/ledisdb.git src/github.com/siddontang/ledisdb

        cd src/github.com/siddontang/ledisdb

+ Install leveldb and snappy, if you have installed, skip.

    LedisDB supplies a simple shell to install leveldb and snappy: 

        sh build_leveldb.sh

    It will default install leveldb at /usr/local/leveldb and snappy at /usr/local/snappy.

    LedisDB use the modified LevelDB for better performance, see [here](https://github.com/siddontang/ledisdb/wiki/leveldb-source-modification).

+ Set LEVELDB_DIR and SNAPPY_DIR to the actual install path in dev.sh.

+ Then:

        . ./bootstap.sh 
        . ./dev.sh

        go install ./...

## Server Example

    ./ledis-server -config=/etc/ledis.json

    //another shell
    ledis-cli -p 6380
    
    ledis 127.0.0.1:6380> set a 1
    OK
    ledis 127.0.0.1:6380> get a
    "1"

## Package Example
    
    import "github.com/siddontang/ledisdb/ledis"
    l, _ := ledis.Open(cfg)
    db, _ := l.Select(0)

    db.Set(key, value)

    db.Get(key)


## Replication Example

Set slaveof in config or dynamiclly

    ledis-cli -p 6381 

    ledis 127.0.0.1:6381> slaveof 127.0.0.1 6380
    OK

## Benchmark

See benchmark.md for more.

## Todo

See [Issues todo](https://github.com/siddontang/ledisdb/issues?labels=todo&page=1&state=open)


## Links

+ [Official Website](http://ledisdb.com)
+ [Author's Chinese Blog](http://blog.csdn.net/siddontang/article/category/2264003)
+ [GoDoc](https://godoc.org/github.com/siddontang/ledisdb)
+ [Server Commands](https://github.com/siddontang/ledisdb/wiki/Commands)


## Thanks

Gmail: cenqichao@gmail.com

Gmail: chendahui007@gmail.com

Gmail: cppgohan@gmail.com

Gmail: tiaotiaoyly@gmail.com

Gmail: wyk4true@gmail.com


## Feedback

Gmail: siddontang@gmail.com
