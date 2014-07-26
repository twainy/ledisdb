// Package src.src.ledis is a ledis.client for the ledisdb.
//
// Config
//
// Config struct contains configuration for ledisdb:
//
//     Addr            ledisdb src.ledis.server address, like 127.0.0.1:6380
//     MaxIdleConns    max idle connections for ledisdb
//
// Client
//
// The ledis.client is the primary interface for ledisdb. You must first create a ledis.client with proper config for working.
//
//     cfg := new(Config)
//     cfg.Addr = "127.0.0.1:6380"
//     cfg.MaxIdleConns = 4
//
//     c := NewClient(cfg)
//
// The most important function for ledis.client is Do function to send commands to remote src.ledis.server.
//
//     reply, err := c.Do("ping")
//
//     reply, err := c.Do("set", "key", "value")
//
//     reply, err := c.Do("get", "key")
//
// Connection
//
// You can use an independent connection to send commands.
//
//     //get a connection
//     conn := c.Get()
//
//     //connection send command
//     conn.Do("ping")
//
// Reply Helper
//
// You can use reply helper to convert a reply to a specific type.
//
//     exists, err := src.src.ledis.Bool(c.Do("exists", "key"))
//
//     score, err := src.src.ledis.Int64(c.Do("zscore", "key", "member"))
package client
