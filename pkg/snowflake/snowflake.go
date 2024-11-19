package snowflake

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"os"
	"sync"
)

var (
	nodes []*snowflake.Node
	once  sync.Once
)

// InitSnowflake 初始化雪花算法节点池
func InitSnowflake() error {
	var initErr error
	once.Do(func() {
		nodes = make([]*snowflake.Node, 1024)
		for i := int64(0); i < 1024; i++ {
			node, err := snowflake.NewNode(i)
			if err != nil {
				initErr = fmt.Errorf("init snowflake node error: %w", err)
				return
			}
			nodes[i] = node
		}
	})
	return initErr
}

// GenerateNodeID 根据主机名生成节点ID
func GenerateNodeID(hostname string) int64 {
	hash := int64(0)
	for _, c := range hostname {
		hash = 31*hash + int64(c)
	}

	nodeID := hash % 1024
	if nodeID < 0 {
		nodeID += 1024
	}
	return nodeID
}

// GetSnowflakeGenerator 获取雪花算法生成器
func GetSnowflakeGenerator() (func() int64, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("get hostname failed: %w", err)
	}

	nodeID := GenerateNodeID(hostname)
	node := nodes[nodeID]

	return func() int64 {
		return node.Generate().Int64()
	}, nil
}
