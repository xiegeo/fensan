package store

import (
	"fmt"

	ht "github.com/xiegeo/fensan/hashtree"
)

type metaStore struct {
	minLevel  ht.Level
	ttlStore  KV
	hashStore LV
}

func (m *metaStore) InnerHashMinLevel() ht.Level {
	return m.minLevel
}
func (m *metaStore) GetInnerHashes(key HLKey, hs []byte, level ht.Level, off ht.Nodes) error {
	panic("not implemented")
}
func (m *metaStore) PutInnerHashes(key HLKey, hs []byte, level ht.Level, off ht.Nodes) (has ht.Nodes, complete bool, err error) {
	panic("not implemented")
}
func (m *metaStore) TTLGet(key HLKey) TTL {
	v := m.ttlStore.Get(key.FullBytes())
	return TTLFromBytes(v)
}
func (m *metaStore) TTLSetAtleast(key HLKey, freeFrom, until TTL) (byteMonth int64) {
	get := m.ttlStore.Get(key.FullBytes())
	old := TTLFromBytes(get)
	if old >= until {
		return 0
	}
	if freeFrom < old {
		freeFrom = old
	}
	buyMonth := freeFrom.MonthUntil(until)
	m.ttlStore.Set(key.FullBytes(), until.Bytes())
	return int64(buyMonth) * key.Length()
	//todo: check dedup
}

func (m *metaStore) Close() error {
	err := m.ttlStore.Close()
	err2 := m.hashStore.Close()
	if err != nil || err2 != nil {
		return fmt.Errorf("fail close meta store, ttlStore err: %v; hashStore err: %v", err, err2)
	}
	return nil
}
