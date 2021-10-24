package cache

import pb "golearning/go-cache/cachepb"

// PeerPicker 通过key获取实例节点
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter 实例节点根据key获取[]byte数据
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
