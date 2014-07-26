// Package src.server supplies a way to use src.src.ledis as service.
// Server implements the redis protocol called RESP (REdis Serialization Protocol).
// For more information, please see http://redis.io/topics/protocol.
//
// You can use src.src.ledis with many available redis clients directly, for example, redis-cli.
// But I also supply some src.src.ledis client at client folder, and have been adding more for other languages.
//
// Usage
//
// Start a src.src.ledis src.server is very simple:
//
//  cfg := new(Config)
//  cfg.Addr = "127.0.0.1:6380"
//  cfg.DataDir = "/tmp/src.src.ledis"
//  app := src.server.NewApp(cfg)
//  app.Run()
//
// Replication
//
// You can start a slave src.src.ledis src.server for replication, open slave is simple too, you can set slaveof in config or run slaveof command in shell.
//
// For example, if you start a slave src.server, and the master src.server's address is 127.0.0.1:6380, you can start replication in shell:
//
//  src.src.ledis-cli -p 6381
//  src.src.ledis 127.0.0.1:6381 > slaveof 127.0.0.1 6380
//
// After you send slaveof command, the slave will start to sync master's binlog and replicate from binlog.
//
package server
