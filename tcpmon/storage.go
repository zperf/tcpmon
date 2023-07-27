package tcpmon

type Datastore struct{}

func NewDatastore() *Datastore {
	return &Datastore{}
}

func (d *Datastore) Put() error {
	return nil
}
