package protohelper

import "fmt"

type DummyMessage struct {
	payload []byte
}

func (dm *DummyMessage) Reset()                   { dm.payload = dm.payload[:0] }
func (dm *DummyMessage) String() string           { return fmt.Sprintf("%q", dm.payload) }
func (dm *DummyMessage) ProtoMessage()            {}
func (dm *DummyMessage) Marshal() ([]byte, error) { return dm.payload, nil }
func (dm *DummyMessage) Unmarshal(in []byte) error {
	dm.payload = append(dm.payload[:0], in...)
	return nil
}
