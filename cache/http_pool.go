package cache

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	pb "golearning/go-cache/cachepb"
	"golearning/go-cache/consistenthash"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/gocache/"
	defaultReplicas = 50
)

type HTTPPool struct {
	// 本机{ip:port}
	self string
	// e.g. "/gocache/"
	basePath string
	mu       sync.Mutex
	// 哈希一致性的实现
	peers *consistenthash.Map
	// keyed by e.g. "http://10.0.0.2:8888"
	httpGetters map[string]*httpGetter
}

func NewHttpPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

func (this *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", this.self, fmt.Sprintf(format, v...))
}

// ServeHTTP 接收HTTP的请求并返回响应
func (this *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 参数校验
	if !strings.HasPrefix(r.URL.Path, this.basePath) {
		http.Error(w, "check path. ps: /gocache/{groupName}/{key}", http.StatusNotFound)
		return
	}
	//this.Log("method=%s path=%s", r.Method, r.URL.Path)
	params := strings.SplitN(r.URL.Path[len(this.basePath):], "/", 2)
	if len(params) != 2 {
		http.Error(w, "param format err. ps: /gocache/{groupName}/{key}", http.StatusNotFound)
		return
	}

	groupName, key := params[0], params[1]
	this.Log("groupName=%s key=%s", groupName, key)

	// 获取group
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "group not found. groupName="+groupName, http.StatusOK)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusOK)
		return
	}

	// proto.Marshal
	body, err := proto.Marshal(&pb.Response{Value: view.ByteSlice()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	_, _ = w.Write(body)
}

// Set 注册兄弟节点
func (this *HTTPPool) Set(peers ...string) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.peers = consistenthash.NewMap(nil, defaultReplicas)
	this.peers.Add(peers...)
	this.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		this.httpGetters[peer] = &httpGetter{baseUrl: peer + this.basePath}
	}
}

func (this *HTTPPool) PickPeer(key string) (peer PeerGetter, ok bool) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if peer := this.peers.Get(key); peer != "" && peer != this.self {
		return this.httpGetters[peer], true
	}
	return
}

type httpGetter struct {
	baseUrl string
}

func (this *httpGetter) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", this.baseUrl, fmt.Sprintf(format, v...))
}

// Get 请求远端节点获取数据的Get 通过HTTP协议交互，Protobuf序列化....
func (this *httpGetter) Get(in *pb.Request, out *pb.Response) error {
	this.Log("peer receive a request... group = %s, key = %s", in.GetGroup(), in.GetKey())
	u := fmt.Sprintf("%v%v/%v", this.baseUrl, url.QueryEscape(in.GetGroup()), url.QueryEscape(in.GetKey()))
	resp, err := http.Get(u) // ignore_security_alert
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %v", resp.StatusCode)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("server read err: %v", err.Error())
	}

	if err = proto.Unmarshal(bytes, out); err != nil {
		return fmt.Errorf("decoding response body: %v", err)
	}
	return nil
}
