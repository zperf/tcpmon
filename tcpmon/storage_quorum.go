package tcpmon

import (
	"encoding/json"

	"github.com/dgraph-io/badger/v4"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func (d *Datastore) GetMembers() ([]KVPair, error) {
	return d.GetPrefix([]byte(PrefixMember), 0, false)
}

func (d *Datastore) GetMemberAddressList() ([]string, error) {
	r := make([]string, 0)
	members, err := d.GetMembers()
	if err != nil {
		return nil, err
	}

	for _, p := range members {
		r = append(r, p.Key)
	}
	return r, nil
}

func (d *Datastore) AddMember(member string, buf []byte) error {
	return d.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(KeyJoin(PrefixMember, member)), buf)
	})
}

func (d *Datastore) DeleteMember(member string) error {
	return d.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(KeyJoin(PrefixMember, member)))
	})
}

func (d *Datastore) UpdateMember(member string, buf []byte) error {
	return d.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(KeyJoin(PrefixMember, member)), buf)
	})
}

func (d *Datastore) GetMemberMeta(member string) (map[string]any, error) {
	var buf []byte

	err := d.db.View(func(txn *badger.Txn) error {
		it, err := txn.Get([]byte(KeyJoin(PrefixMember, member)))
		if err != nil {
			return err
		}
		buf, err = it.ValueCopy(nil)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	var m MemberInfo
	err = proto.Unmarshal(buf, &m)
	if err != nil {
		return nil, err
	}

	buf, err = protojson.Marshal(&m)
	if err != nil {
		return nil, err
	}

	out := make(map[string]any)
	err = json.Unmarshal(buf, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}
