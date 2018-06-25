package raft

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"
	rbdb "github.com/tidwall/raft-boltdb"
)

// Node raft节点
type Node struct {
	dataDir   string
	buildAddr string
	localID   raft.ServerID
	timeout   time.Duration
	raft      *raft.Raft
}

// NewNode 实例化一个raft节点
func NewNode(buildAddr string) *Node {
	n := new(Node)
	n.dataDir = "/tmp/raft-example"
	n.buildAddr = buildAddr
	n.localID = getLocalID(buildAddr)
	n.timeout = 2 * time.Second
	return n
}

// 为了方便以节点地址生成ID，实际上要做更多的保障以确保节点ID再集群中唯一
func getLocalID(address string) raft.ServerID {
	return raft.ServerID(address)
}

// GetRaftInstince 获取raft实例
func (n *Node) GetRaftInstince() *raft.Raft {
	return n.raft
}

// Start 启动节点
// isStartSingle 是否以单节点启动
func (n *Node) Start(isStartSingle bool) error {
	// 默认配置启动
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(n.localID)

	// 绑定地址
	addr, err := net.ResolveTCPAddr("tcp", n.buildAddr)
	if err != nil {
		return err
	}
	transport, err := raft.NewTCPTransport(n.buildAddr, addr, 3, n.timeout, os.Stderr)
	if err != nil {
		return err
	}

	// 快照存储区
	snapshot, err := raft.NewFileSnapshotStore(n.dataDir, 3, os.Stderr)
	if err != nil {
		return err
	}

	// 日志存储
	logStore, err := rbdb.NewBoltStore(filepath.Join(n.dataDir, "raft.db"))
	if err != nil {
		return err
	}

	fsm := new(rFSM)
	// 实例化节点
	r, err := raft.NewRaft(config, fsm, logStore, logStore, snapshot, transport)
	if err != nil {
		return err
	}
	n.raft = r

	//单节点启动
	if isStartSingle {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		// 启动  这个方法只能在集群启动的时候调用
		r.BootstrapCluster(configuration)
	}
	return nil
}

// Join 将一个节点加入集群
func (n *Node) Join(address string) {
	//不是领导者，则需要将请求转发至领导者
	if n.raft.State() != raft.Leader {
		leader := n.raft.Leader()
		fmt.Println("当前领导者是:", leader)
		return
	}
	serverID := getLocalID(address)
	serverAddr := raft.ServerAddress(address)

	fmt.Printf("开始连接节点: NodeID: %s, NodeAddr: %s \n", serverID, serverAddr)
	future := n.raft.AddVoter(serverID, serverAddr, 0, 0)
	if future.Error() != nil {
		fmt.Println("加入集群失败", future.Error())
	}
}
