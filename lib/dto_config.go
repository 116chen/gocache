package lib

import "github.com/hashicorp/raft"

var RaftNode *raft.Raft

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

type SysConfig struct {
	RaftConfig      RaftConfig      `yaml:"raft-config"`
	PeerNodesConfig PeerNodesConfig `yaml:"peer-nodes-config"`
}

type Servers struct {
	ID      int    `yaml:"id"`
	Address string `yaml:"address"`
}
type PeerNodesConfig struct {
	Port    string    `yaml:"port"`
	IsAPI   bool      `yaml:"is-api"`
	Servers []Servers `yaml:"servers"`
}

type RaftConfig struct {
	ServerName  string        `yaml:"server-name"`
	ServerID    string        `yaml:"server-id"`
	LogStore    string        `yaml:"log-store"`
	StableStore string        `yaml:"stable-store"`
	Transport   string        `yaml:"transport"`
	Role        string        `yaml:"role"`
	Servers     []raft.Server `yaml:"servers"`
}
