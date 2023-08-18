package tcpmon

import (
	"github.com/cockroachdb/errors"
	"github.com/dgraph-io/badger/v4"
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

func (d *Datastore) AddMember(member string) error {
	m := &MemberInfo{}
	buf, err := proto.Marshal(m)
	if err != nil {
		return errors.WithStack(err)
	}

	return d.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(KeyJoin(PrefixMember, member)), buf)
	})
}

func (d *Datastore) DeleteMember(member string) error {
	return d.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(KeyJoin(PrefixMember, member)))
	})
}
