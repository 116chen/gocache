package lib

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"golearning/go-cache/cache"
	"golearning/go-cache/common"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type Listener func(chan interface{})

func loadSysConfig(path string) (*SysConfig, error) {
	config := new(SysConfig)
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func Boot(path string) error {
	// 加载配置文件
	sysConfig, err := loadSysConfig(path)
	if err != nil {
		return err
	}

	// 启动raft服务
	leaderNotifyCh, err := bootRaft(&sysConfig.RaftConfig)
	if err != nil {
		return err
	}

	// 建立和对等节点的连接
	//if sysConfig.RaftConfig.Role == common.RoleTypeMaster.String() {
	//	err = bootPeerNodes(&sysConfig.PeerNodesConfig)
	//	if err != nil {
	//		return err
	//	}
	//}

	//for true {
	//	select {
	//	case isLeader := <-leaderNotifyCh:
	//		if isLeader {
	//			// raft集群中的主节点
	//			sysConfig.RaftConfig.Role = common.RoleTypeMaster.String()
	//			err := bootPeerNodes(&sysConfig.PeerNodesConfig)
	//			if err != nil {
	//				panic(err)
	//			}
	//		}
	//	}
	//	log.Println("current leader: ", RaftNode.Leader())
	//}

	RegisterListener(sysConfig, leaderNotifyCh)

	return nil
}

// TODO 线程池实现注册监听器，暂时只实现了切换主从节点的监听器
func RegisterListener(sysConfig *SysConfig, ch <-chan bool) {
	go func() {
		for true {
			select {
			//  TODO isLeader有问题....
			case isLeader := <-ch:
				log.Println("进入该分支！！！ isLeader=", isLeader)
				if isLeader {
					// raft集群中的主节点
					sysConfig.RaftConfig.Role = common.RoleTypeMaster.String()
					//err := bootPeerNodes(&sysConfig.PeerNodesConfig)
					//if err != nil {
					//	panic(err)
					//}
				}
			}
			log.Println("current leader: ", RaftNode.Leader())
		}
	}()
}

func createGroup() *cache.Group {
	return cache.NewGroup("test", cache.GetterFunc(func(key string) ([]byte, error) {
		log.Println("slowDb load. key =", key)
		if value, ok := db[key]; ok {
			return []byte(value), nil
		}
		return nil, fmt.Errorf("%s is not found", key)
	}), 1<<20)
}

func startCacheServer(addr string, addrs []string, gee *cache.Group) {
	peers := cache.NewHttpPool(addr)
	peers.Set(addrs...)
	gee.RegisterPeers(peers)
	log.Println("gocache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}

func startAPIServer(apiAddr string, gee *cache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := gee.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())

		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))

}

func bootPeerNodes(peerNodesConfig *PeerNodesConfig) error {
	var addrs []string
	addrMap := make(map[string]string, 0)
	// TODO 当主从关系变化时，修改配置类的做法是否正确？？如果不正确，应该怎么操作会比较好一点
	for _, v := range peerNodesConfig.Servers {
		addrs = append(addrs, v.Address)
		params := strings.Split(v.Address, ":")
		addrMap[params[0]] = params[1]
	}

	gee := createGroup()
	if peerNodesConfig.IsAPI {
		startAPIServer(addrs[0], gee)
	} else {
		startCacheServer(addrMap[peerNodesConfig.Port], addrs, gee)
	}

	// TODO 广播告知对等节点raft主节点是谁....

	return nil
}

func bootRaft(raftConfig *RaftConfig) (<-chan bool, error) {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(raftConfig.ServerID)
	config.Logger = hclog.New(&hclog.LoggerOptions{
		Name:   raftConfig.ServerName,
		Level:  hclog.LevelFromString("DEBUG"),
		Output: os.Stderr,
	})
	leaderNotifyCh := make(chan bool, 1)
	config.NotifyCh = leaderNotifyCh

	dir, _ := os.Getwd()
	root := strings.Replace(dir, "\\", "/", -1)
	logStore, err := raftboltdb.NewBoltStore(root + raftConfig.LogStore)
	if err != nil {
		return nil, err
	}

	stableStore, err := raftboltdb.NewBoltStore(root + raftConfig.StableStore)
	if err != nil {
		return nil, err
	}

	snapshotStore := raft.NewDiscardSnapshotStore()

	addr, err := net.ResolveTCPAddr("tcp", raftConfig.Transport)

	transport, err := raft.NewTCPTransport(addr.String(), addr, 5, time.Second*10, os.Stdout)
	if err != nil {
		return nil, err
	}

	fsm := &MyFSM{}

	// 这里不能使用:=，如果使用，则RaftNode作用域仅限该函数内，外部使用RaftNode==nil
	RaftNode, err = raft.NewRaft(config, fsm, logStore, stableStore, snapshotStore, transport)
	if err != nil {
		return nil, err
	}

	configuration := raft.Configuration{
		Servers: raftConfig.Servers,
	}

	RaftNode.BootstrapCluster(configuration)
	return leaderNotifyCh, nil
}
