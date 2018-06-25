package raft

import (
	"encoding/binary"
	"hash/adler32"
	"io"

	"github.com/hashicorp/raft"
)

type logHash struct {
	lastHash []byte
}

func (l *logHash) Add(d []byte) {
	hasher := adler32.New()
	hasher.Write(l.lastHash)
	hasher.Write(d)
	l.lastHash = hasher.Sum(nil)
}

type applyItem struct {
	index uint64
	term  uint64
	data  []byte
}

func (a *applyItem) set(l *raft.Log) {
	a.index = l.Index
	a.term = l.Term
	a.data = make([]byte, len(l.Data))
	copy(a.data, l.Data)
}

type rFSM struct {
	logHash
	lastTerm  uint64
	lastIndex uint64
	applied   []applyItem
}

// 日志写入完成时回调此方法
func (f *rFSM) Apply(l *raft.Log) interface{} {
	return nil
}

func (f *rFSM) Snapshot() (raft.FSMSnapshot, error) {
	s := new(rFSMSnapshot)
	return s, nil
}

func (f *rFSM) Restore(r io.ReadCloser) error {
	err := binary.Read(r, binary.LittleEndian, &f.lastTerm)
	if err == nil {
		err = binary.Read(r, binary.LittleEndian, &f.lastIndex)
	}
	if err == nil {
		f.lastHash = make([]byte, adler32.Size)
		_, err = r.Read(f.lastHash)
	}
	return err
}
