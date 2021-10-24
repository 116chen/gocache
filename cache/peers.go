package cache

import pb "golearning/go-cache/cachepb"

// PeerPicker 通过key获取兄弟节点
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter 节点根据key获取[]byte数据
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
