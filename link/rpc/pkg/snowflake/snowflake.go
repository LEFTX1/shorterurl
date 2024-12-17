package snowflake

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/snowflake"
)

var (
	node    *snowflake.Node
	once    sync.Once
	initErr error
	nodeID  int64 = 1 // 默认节点ID
)

// InitSnowflake 初始化雪花算法
func InitSnowflake() error {
	once.Do(func() {
		node, initErr = snowflake.NewNode(nodeID)
	})
	return initErr
}

// GetSnowflakeGenerator 获取ID生成器函数
func GetSnowflakeGenerator() (func() int64, error) {
	if node == nil {
		return nil, fmt.Errorf("snowflake node not initialized")
	}
	return func() int64 {
		return node.Generate().Int64()
	}, nil
}

// SetNodeID 设置节点ID
func SetNodeID(id int64) {
	nodeID = id
}
