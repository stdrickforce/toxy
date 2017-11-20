package xprocessor

import (
	"xlog"
	. "xprotocol"
	. "xthrift"
)

type Messenger struct {
	iprot Protocol
	oprot Protocol
}

func (m *Messenger) Reverse() {
	m.iprot, m.oprot = m.oprot, m.iprot
}

func (m *Messenger) ReadMessageBegin() (name string, mtype byte, seqid int32, err error) {
	return m.iprot.ReadMessageBegin()
}

func (m *Messenger) WriteMessageBegin(name string, mtype byte, seqid int32) error {
	return m.oprot.WriteMessageBegin(name, mtype, seqid)
}

func (m *Messenger) ForwardMessageBegin() error {
	name, mtype, seqid, err := m.iprot.ReadMessageBegin()
	if err != nil {
		return err
	}
	if err := m.oprot.WriteMessageBegin(name, mtype, seqid); err != nil {
		return err
	}
	return nil
}

func (m *Messenger) ForwardMessageEnd() error {
	if err := m.iprot.ReadMessageEnd(); err != nil {
		return err
	}
	if err := m.oprot.WriteMessageEnd(); err != nil {
		return err
	}
	if err := m.oprot.Flush(); err != nil {
		return err
	}
	return nil
}

func (m *Messenger) FastReply(seqid int32) (err error) {
	xlog.Debug("fast reply: ping")
	if err = m.iprot.Skip(T_STRUCT); err != nil {
		return
	}
	if err = m.iprot.ReadMessageEnd(); err != nil {
		return
	}
	if err = m.iprot.WriteMessageBegin("ping", T_REPLY, seqid); err != nil {
		return
	}
	if err = m.iprot.WriteByte(T_STOP); err != nil {
		return
	}
	if err = m.iprot.WriteMessageEnd(); err != nil {
		return
	}
	return
}

func (m *Messenger) Reply(name string, seqid int32) {
}

// TODO This only happened while iprot and oprot are same.
func (m *Messenger) FastForward(ftype byte) error {
	return nil
}

func (m *Messenger) Forward(ftype byte) error {
	if err := forward(m.iprot, m.oprot, ftype); err != nil {
		return err
	}
	return nil
}

func (m *Messenger) SetOutputProtocol(oprot Protocol) {
	m.oprot = oprot
}

func (m *Messenger) DelOutputProtocol() {
	m.oprot.Close()
	m.oprot = nil
}

func NewMessenger(protocol Protocol) *Messenger {
	return &Messenger{
		iprot: protocol,
	}
}