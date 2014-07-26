package ledis

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var m sync.Mutex

type adaptor struct {
	set    func([]byte, []byte) (int64, error)
	del    func([]byte) (int64, error)
	exists func([]byte) (int64, error)

	expire   func([]byte, int64) (int64, error)
	expireAt func([]byte, int64) (int64, error)
	ttl      func([]byte) (int64, error)

	showIdent func() string
}

func kvAdaptor(db *DB) *adaptor {
	adp := new(adaptor)
	adp.showIdent = func() string {
		return "kv-adptor"
	}

	adp.set = db.SetNX
	adp.exists = db.Exists
	adp.del = func(k []byte) (int64, error) {
		return db.Del(k)
	}

	adp.expire = db.Expire
	adp.expireAt = db.ExpireAt
	adp.ttl = db.TTL

	return adp
}

func listAdaptor(db *DB) *adaptor {
	adp := new(adaptor)
	adp.showIdent = func() string {
		return "list-adptor"
	}

	adp.set = func(k []byte, v []byte) (int64, error) {
		eles := make([][]byte, 0)
		for i := 0; i < 3; i++ {
			e := []byte(String(v) + fmt.Sprintf("_%d", i))
			eles = append(eles, e)
		}

		if n, err := db.LPush(k, eles...); err != nil {
			return 0, err
		} else {
			return n, nil
		}
	}

	adp.exists = func(k []byte) (int64, error) {
		if llen, err := db.LLen(k); err != nil || llen <= 0 {
			return 0, err
		} else {
			return 1, nil
		}
	}

	adp.del = db.LClear
	adp.expire = db.LExpire
	adp.expireAt = db.LExpireAt
	adp.ttl = db.LTTL

	return adp
}

func hashAdaptor(db *DB) *adaptor {
	adp := new(adaptor)
	adp.showIdent = func() string {
		return "hash-adptor"
	}

	adp.set = func(k []byte, v []byte) (int64, error) {
		datas := make([]FVPair, 0)
		for i := 0; i < 3; i++ {
			suffix := fmt.Sprintf("_%d", i)
			pair := FVPair{
				Field: []byte(String(k) + suffix),
				Value: []byte(String(v) + suffix)}

			datas = append(datas, pair)
		}

		if err := db.HMset(k, datas...); err != nil {
			return 0, err
		} else {
			return int64(len(datas)), nil
		}
	}

	adp.exists = func(k []byte) (int64, error) {
		if hlen, err := db.HLen(k); err != nil || hlen <= 0 {
			return 0, err
		} else {
			return 1, nil
		}
	}

	adp.del = db.HClear
	adp.expire = db.HExpire
	adp.expireAt = db.HExpireAt
	adp.ttl = db.HTTL

	return adp
}

func zsetAdaptor(db *DB) *adaptor {
	adp := new(adaptor)
	adp.showIdent = func() string {
		return "zset-adptor"
	}

	adp.set = func(k []byte, v []byte) (int64, error) {
		datas := make([]ScorePair, 0)
		for i := 0; i < 3; i++ {
			memb := []byte(String(k) + fmt.Sprintf("_%d", i))
			pair := ScorePair{
				Score:  int64(i),
				Member: memb}

			datas = append(datas, pair)
		}

		if n, err := db.ZAdd(k, datas...); err != nil {
			return 0, err
		} else {
			return n, nil
		}
	}

	adp.exists = func(k []byte) (int64, error) {
		if cnt, err := db.ZCard(k); err != nil || cnt <= 0 {
			return 0, err
		} else {
			return 1, nil
		}
	}

	adp.del = db.ZClear
	adp.expire = db.ZExpire
	adp.expireAt = db.ZExpireAt
	adp.ttl = db.ZTTL

	return adp
}

func allAdaptors(db *DB) []*adaptor {
	adps := make([]*adaptor, 4)
	adps[0] = kvAdaptor(db)
	adps[1] = listAdaptor(db)
	adps[2] = hashAdaptor(db)
	adps[3] = zsetAdaptor(db)
	return adps
}

///////////////////////////////////////////////////////

func TestExpire(t *testing.T) {
	db := getTestDB()
	m.Lock()
	defer m.Unlock()

	k := []byte("ttl_a")
	ek := []byte("ttl_b")

	dbEntrys := allAdaptors(db)
	for _, entry := range dbEntrys {
		ident := entry.showIdent()

		entry.set(k, []byte("1"))

		if ok, _ := entry.expire(k, 10); ok != 1 {
			t.Fatal(ident, ok)
		}

		//	err - expire on an inexisting key
		if ok, _ := entry.expire(ek, 10); ok != 0 {
			t.Fatal(ident, ok)
		}

		//	err - duration is zero
		if ok, err := entry.expire(k, 0); err == nil || ok != 0 {
			t.Fatal(ident, fmt.Sprintf("res = %d, err = %s", ok, err))
		}

		//	err - duration is negative
		if ok, err := entry.expire(k, -10); err == nil || ok != 0 {
			t.Fatal(ident, fmt.Sprintf("res = %d, err = %s", ok, err))
		}
	}
}

func TestExpireAt(t *testing.T) {
	db := getTestDB()
	m.Lock()
	defer m.Unlock()

	k := []byte("ttl_a")
	ek := []byte("ttl_b")

	dbEntrys := allAdaptors(db)
	for _, entry := range dbEntrys {
		ident := entry.showIdent()
		now := time.Now().Unix()

		entry.set(k, []byte("1"))

		if ok, _ := entry.expireAt(k, now+5); ok != 1 {
			t.Fatal(ident, ok)
		}

		//	err - expire on an inexisting key
		if ok, _ := entry.expireAt(ek, now+5); ok != 0 {
			t.Fatal(ident, ok)
		}

		//	err - expire with the current time
		if ok, err := entry.expireAt(k, now); err == nil || ok != 0 {
			t.Fatal(ident, fmt.Sprintf("res = %d, err = %s", ok, err))
		}

		//	err - expire with the time before
		if ok, err := entry.expireAt(k, now-5); err == nil || ok != 0 {
			t.Fatal(ident, fmt.Sprintf("res = %d, err = %s", ok, err))
		}
	}
}

func TestTTL(t *testing.T) {
	db := getTestDB()
	m.Lock()
	defer m.Unlock()

	k := []byte("ttl_a")
	ek := []byte("ttl_b")

	dbEntrys := allAdaptors(db)
	for _, entry := range dbEntrys {
		ident := entry.showIdent()

		entry.set(k, []byte("1"))
		entry.expire(k, 2)

		if tRemain, _ := entry.ttl(k); tRemain != 2 {
			t.Fatal(ident, tRemain)
		}

		//	err - check ttl on an inexisting key
		if tRemain, _ := entry.ttl(ek); tRemain != -1 {
			t.Fatal(ident, tRemain)
		}

		entry.del(k)
		if tRemain, _ := entry.ttl(k); tRemain != -1 {
			t.Fatal(ident, tRemain)
		}
	}
}

func TestExpCompose(t *testing.T) {
	db := getTestDB()
	m.Lock()
	defer m.Unlock()

	k0 := []byte("ttl_a")
	k1 := []byte("ttl_b")
	k2 := []byte("ttl_c")

	dbEntrys := allAdaptors(db)

	for _, entry := range dbEntrys {
		ident := entry.showIdent()

		entry.set(k0, k0)
		entry.set(k1, k1)
		entry.set(k2, k2)

		entry.expire(k0, 5)
		entry.expire(k1, 2)
		entry.expire(k2, 60)

		if tRemain, _ := entry.ttl(k0); tRemain != 5 {
			t.Fatal(ident, tRemain)
		}
		if tRemain, _ := entry.ttl(k1); tRemain != 2 {
			t.Fatal(ident, tRemain)
		}
		if tRemain, _ := entry.ttl(k2); tRemain != 60 {
			t.Fatal(ident, tRemain)
		}
	}

	// after 1 sec
	time.Sleep(1 * time.Second)

	for _, entry := range dbEntrys {
		ident := entry.showIdent()

		if tRemain, _ := entry.ttl(k0); tRemain != 4 {
			t.Fatal(ident, tRemain)
		}
		if tRemain, _ := entry.ttl(k1); tRemain != 1 {
			t.Fatal(ident, tRemain)
		}
	}

	// after 2 sec
	time.Sleep(2 * time.Second)

	for _, entry := range dbEntrys {
		ident := entry.showIdent()

		if tRemain, _ := entry.ttl(k1); tRemain != -1 {
			t.Fatal(ident, tRemain)
		}
		if exist, _ := entry.exists(k1); exist > 0 {
			t.Fatal(ident, false)
		}

		if tRemain, _ := entry.ttl(k0); tRemain != 2 {
			t.Fatal(ident, tRemain)
		}
		if exist, _ := entry.exists(k0); exist <= 0 {
			t.Fatal(ident, false)
		}

		// refresh the expiration of key
		if tRemain, _ := entry.ttl(k2); !(0 < tRemain && tRemain < 60) {
			t.Fatal(ident, tRemain)
		}

		if ok, _ := entry.expire(k2, 100); ok != 1 {
			t.Fatal(ident, false)
		}

		if tRemain, _ := entry.ttl(k2); tRemain != 100 {
			t.Fatal(ident, tRemain)
		}

		//	expire an inexisting key
		if ok, _ := entry.expire(k1, 10); ok == 1 {
			t.Fatal(ident, false)
		}
		if tRemain, _ := entry.ttl(k1); tRemain != -1 {
			t.Fatal(ident, tRemain)
		}
	}

	return
}
