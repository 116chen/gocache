package lib

import (
	"github.com/hashicorp/raft"
	"io"
)

type MyFSM struct {
}

func (this *MyFSM) Apply(log *raft.Log) interface{} {
	return nil
}

func (this *MyFSM) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (this *MyFSM) Restore(closer io.ReadCloser) error {
	return nil
}
