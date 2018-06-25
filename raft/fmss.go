package raft

import (
	"github.com/hashicorp/raft"
)

type rFSMSnapshot struct {
}

func (f *rFSMSnapshot) Persist(sink raft.SnapshotSink) error {
	return nil
}

func (f *rFSMSnapshot) Release() {
}
