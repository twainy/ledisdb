package server

import (
    ledis_client "ledis/client"
	"testing"
)

func TestKV(t *testing.T) {
	c := getTestConn()
	defer c.Close()

	if ok, err := ledis_client.String(c.Do("set", "a", "1234")); err != nil {
		t.Fatal(err)
	} else if ok != OK {
		t.Fatal(ok)
	}

	if n, err := ledis_client.Int(c.Do("setnx", "a", "123")); err != nil {
		t.Fatal(err)
	} else if n != 0 {
		t.Fatal(n)
	}

	if n, err := ledis_client.Int(c.Do("setnx", "b", "123")); err != nil {
		t.Fatal(err)
	} else if n != 1 {
		t.Fatal(n)
	}

	if v, err := ledis_client.String(c.Do("get", "a")); err != nil {
		t.Fatal(err)
	} else if v != "1234" {
		t.Fatal(v)
	}

	if v, err := ledis_client.String(c.Do("getset", "a", "123")); err != nil {
		t.Fatal(err)
	} else if v != "1234" {
		t.Fatal(v)
	}

	if v, err := ledis_client.String(c.Do("get", "a")); err != nil {
		t.Fatal(err)
	} else if v != "123" {
		t.Fatal(v)
	}

	if n, err := ledis_client.Int(c.Do("exists", "a")); err != nil {
		t.Fatal(err)
	} else if n != 1 {
		t.Fatal(n)
	}

	if n, err := ledis_client.Int(c.Do("exists", "empty_key_test")); err != nil {
		t.Fatal(err)
	} else if n != 0 {
		t.Fatal(n)
	}

	if _, err := ledis_client.Int(c.Do("del", "a", "b")); err != nil {
		t.Fatal(err)
	}

	if n, err := ledis_client.Int(c.Do("exists", "a")); err != nil {
		t.Fatal(err)
	} else if n != 0 {
		t.Fatal(n)
	}

	if n, err := ledis_client.Int(c.Do("exists", "b")); err != nil {
		t.Fatal(err)
	} else if n != 0 {
		t.Fatal(n)
	}
}

func TestKVM(t *testing.T) {
	c := getTestConn()
	defer c.Close()

	if ok, err := ledis_client.String(c.Do("mset", "a", "1", "b", "2")); err != nil {
		t.Fatal(err)
	} else if ok != OK {
		t.Fatal(ok)
	}

	if v, err := ledis_client.MultiBulk(c.Do("mget", "a", "b", "c")); err != nil {
		t.Fatal(err)
	} else if len(v) != 3 {
		t.Fatal(len(v))
	} else {
		if vv, ok := v[0].([]byte); !ok || string(vv) != "1" {
			t.Fatal("not 1")
		}

		if vv, ok := v[1].([]byte); !ok || string(vv) != "2" {
			t.Fatal("not 2")
		}

		if v[2] != nil {
			t.Fatal("must nil")
		}
	}
}

func TestKVIncrDecr(t *testing.T) {
	c := getTestConn()
	defer c.Close()

	if n, err := ledis_client.Int64(c.Do("incr", "n")); err != nil {
		t.Fatal(err)
	} else if n != 1 {
		t.Fatal(n)
	}

	if n, err := ledis_client.Int64(c.Do("incr", "n")); err != nil {
		t.Fatal(err)
	} else if n != 2 {
		t.Fatal(n)
	}

	if n, err := ledis_client.Int64(c.Do("decr", "n")); err != nil {
		t.Fatal(err)
	} else if n != 1 {
		t.Fatal(n)
	}

	if n, err := ledis_client.Int64(c.Do("incrby", "n", 10)); err != nil {
		t.Fatal(err)
	} else if n != 11 {
		t.Fatal(n)
	}

	if n, err := ledis_client.Int64(c.Do("decrby", "n", 10)); err != nil {
		t.Fatal(err)
	} else if n != 1 {
		t.Fatal(n)
	}
}
