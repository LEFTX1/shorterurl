package consumer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

// RedisStream Redis Stream操作包装
type RedisStream struct {
	client *redis.Redis
}

// StreamMessage Redis Stream 消息结构
type StreamMessage struct {
	ID     string
	Fields map[string]string
}

// NewRedisStream 创建Redis Stream操作包装
func NewRedisStream(client *redis.Redis) *RedisStream {
	return &RedisStream{
		client: client,
	}
}

// Exists 检查 Stream 是否存在
func (r *RedisStream) Exists(key string) (bool, error) {
	result, err := r.client.Eval(`
		return redis.call('EXISTS', KEYS[1])
	`, []string{key})
	if err != nil {
		return false, err
	}
	if count, ok := result.(int64); ok {
		return count > 0, nil
	}
	return false, nil
}

// Xadd 添加消息到 Stream
func (r *RedisStream) Xadd(key string, id string, values map[string]string) (string, error) {
	// 构建 XADD 命令的参数
	args := []string{key, "*"} // 使用 * 让 Redis 自动生成 ID
	for k, v := range values {
		args = append(args, k, v)
	}

	// 使用简单的 Redis 命令
	result, err := r.client.Eval(`
		local key = KEYS[1]
		local args = {'XADD', key, '*'}
		
		for i = 1, #ARGV, 2 do
			table.insert(args, ARGV[i])
			table.insert(args, ARGV[i+1])
		end
		
		return redis.call(unpack(args))
	`, []string{key}, args[2:]) // 跳过 key 和 id 参数

	if err != nil {
		return "", err
	}

	if resultStr, ok := result.(string); ok {
		return resultStr, nil
	}
	return fmt.Sprintf("%v", result), nil
}

// Xgroup 创建消费者组
func (r *RedisStream) Xgroup(key, group, start string, mkStream bool) (string, error) {
	// 使用简单的 Redis 命令
	result, err := r.client.Eval(`
		local key = KEYS[1]
		local group = ARGV[1]
		
		-- 检查 Stream 是否存在
		local exists = redis.call('EXISTS', key)
		if exists == 0 then
			-- 如果不存在且需要创建，则创建一个空消息
			redis.call('XADD', key, '*', 'init', 'init')
		end
		
		-- 创建消费者组
		return redis.call('XGROUP', 'CREATE', key, group, '0', 'MKSTREAM')
	`, []string{key}, []string{group})

	if err != nil {
		return "", err
	}

	if resultStr, ok := result.(string); ok {
		return resultStr, nil
	}
	return fmt.Sprintf("%v", result), nil
}

// Xinfo 获取Stream或消费者组信息
func (r *RedisStream) Xinfo(subcommand, key string) (string, error) {
	// 执行XINFO命令
	result, err := r.client.Eval(fmt.Sprintf(`
		return redis.call('XINFO', '%s', KEYS[1])
	`, subcommand), []string{key}, []string{})

	if err != nil {
		return "", err
	}

	// 将结果转换为字符串
	if resultStr, ok := result.(string); ok {
		return resultStr, nil
	}

	return fmt.Sprintf("%v", result), nil
}

// Xreadgroup 从消费者组读取消息
func (r *RedisStream) Xreadgroup(stream, group, consumer string, start string, count int, block int) ([]StreamMessage, error) {
	result, err := r.client.Eval(`
		local messages = redis.call('XREADGROUP', 'GROUP', ARGV[1], ARGV[2], 'COUNT', ARGV[3], 'BLOCK', ARGV[4], 'STREAMS', KEYS[1], ARGV[5])
		if messages == false then
			return {}
		end
		return messages[1][2]
	`, []string{stream}, []string{group, consumer, strconv.Itoa(count), strconv.Itoa(block), start})
	if err != nil {
		return nil, err
	}

	// 解析结果
	messages := make([]StreamMessage, 0)
	if result == nil {
		return messages, nil
	}

	// 将结果转换为 StreamMessage 切片
	if msgs, ok := result.([]interface{}); ok {
		for _, msg := range msgs {
			if msgSlice, ok := msg.([]interface{}); ok && len(msgSlice) >= 2 {
				if id, ok := msgSlice[0].(string); ok {
					if fields, ok := msgSlice[1].([]interface{}); ok {
						message := StreamMessage{
							ID:     id,
							Fields: make(map[string]string),
						}
						for i := 0; i < len(fields); i += 2 {
							if key, ok := fields[i].(string); ok {
								if value, ok := fields[i+1].(string); ok {
									message.Fields[key] = value
								}
							}
						}
						messages = append(messages, message)
					}
				}
			}
		}
	}

	return messages, nil
}

// Xack 确认消息
func (r *RedisStream) Xack(key, group string, ids ...string) (string, error) {
	// 构建XACK命令
	allArgs := []string{group}
	allArgs = append(allArgs, ids...)

	evalScript := `
		local args = {'XACK', KEYS[1], ARGV[1]}
		for i=2,#ARGV do
			table.insert(args, ARGV[i])
		end
		return redis.call(unpack(args))
	`

	result, err := r.client.Eval(evalScript, []string{key}, allArgs)

	if err != nil {
		return "", err
	}

	// 将结果转换为字符串
	if resultStr, ok := result.(string); ok {
		return resultStr, nil
	}

	return fmt.Sprintf("%v", result), nil
}

// Xdel 删除消息
func (r *RedisStream) Xdel(key string, ids ...string) (string, error) {
	// 构建XDEL命令
	evalScript := `
		local args = {'XDEL', KEYS[1]}
		for i=1,#ARGV do
			table.insert(args, ARGV[i])
		end
		return redis.call(unpack(args))
	`

	result, err := r.client.Eval(evalScript, []string{key}, ids)

	if err != nil {
		return "", err
	}

	// 将结果转换为字符串
	if resultStr, ok := result.(string); ok {
		return resultStr, nil
	}

	return fmt.Sprintf("%v", result), nil
}

// Xtrim 裁剪Stream
func (r *RedisStream) Xtrim(key string, maxLen int) (string, error) {
	// 执行XTRIM命令
	result, err := r.client.Eval(`
		return redis.call('XTRIM', KEYS[1], 'MAXLEN', '~', ARGV[1])
	`, []string{key}, []string{strconv.Itoa(maxLen)})

	if err != nil {
		return "", err
	}

	// 将结果转换为字符串
	if resultStr, ok := result.(string); ok {
		return resultStr, nil
	}

	return fmt.Sprintf("%v", result), nil
}

// 辅助函数：构建参数对的字符串
func buildArgPairs(count int) string {
	if count == 0 {
		return ""
	}

	pairs := make([]string, count)
	for i := 0; i < count; i++ {
		argIndex := i*2 + 2 // 参数索引，第一个是ID
		pairs[i] = fmt.Sprintf("ARGV[%d], ARGV[%d]", argIndex, argIndex+1)
	}

	return strings.Join(pairs, ", ")
}

// 辅助函数：将map转换为数组
func mapToArray(m map[string]string) []string {
	result := make([]string, 0, len(m)*2)
	for k, v := range m {
		result = append(result, k, v)
	}
	return result
}
