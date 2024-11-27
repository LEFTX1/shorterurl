package user

import (
	"context"                                // 导入上下文包，用于控制操作的生命周期
	"crypto/rand"                            // 导入随机数生成包，用于生成随机用户名
	"encoding/hex"                           // 导入十六进制编码包，用于将随机字节编码为十六进制字符串
	"errors"                                 // 导入错误处理包
	"flag"                                   // 导入命令行标志解析包
	"shorterurl/admin/internal/types/errorx" // 导入自定义错误类型包
	"testing"                                // 导入测试包
	"time"                                   // 导入时间包

	"github.com/stretchr/testify/require" // 导入 testify 包中的 require，用于断言

	"github.com/stretchr/testify/assert"     // 导入 testify 包中的 assert，用于断言
	"github.com/zeromicro/go-zero/core/conf" // 导入 go-zero 配置包
	"github.com/zeromicro/go-zero/core/logx" // 导入 go-zero 日志包

	"shorterurl/admin/internal/config" // 导入项目配置包
	"shorterurl/admin/internal/svc"    // 导入服务上下文包
	"shorterurl/admin/internal/types"  // 导入项目类型包
)

// generateTestUsername 生成一个随机的测试用户名
func generateTestUsername() string {
	bytes := make([]byte, 16)                   // 创建一个长度为16的字节切片
	if _, err := rand.Read(bytes); err != nil { // 生成随机字节
		return "test_" + hex.EncodeToString([]byte(time.Now().String())) // 如果生成随机字节失败，使用当前时间生成用户名
	}
	return "test_" + hex.EncodeToString(bytes) // 返回生成的随机用户名
}

// setupTest 设置测试环境
func setupTest(t *testing.T) (*svc.ServiceContext, context.Context) {
	configFile := flag.String("f", "../../../etc/admin-api.yaml", "配置文件路径") // 定义配置文件路径标志
	flag.Parse()                                                            // 解析命令行标志

	var c config.Config            // 创建配置结构体
	conf.MustLoad(*configFile, &c) // 加载配置文件

	// 创建服务上下文
	svcCtx, err := svc.NewServiceContext(c)
	require.NoError(t, err, "创建服务上下文失败") // 断言没有错误

	// 禁用测试日志
	logx.Disable()

	return svcCtx, context.Background() // 返回服务上下文和背景上下文
}

// TestUserRegisterLogic 测试用户注册逻辑
func TestUserRegisterLogic(t *testing.T) {
	svcCtx, ctx := setupTest(t)                // 设置测试环境
	logic := NewUserRegisterLogic(ctx, svcCtx) // 创建用户注册逻辑实例

	t.Run("成功注册新用户", func(t *testing.T) {
		req := &types.UserRegisterReq{
			Username: generateTestUsername(), // 生成测试用户名
			Password: "password123",          // 设置密码
			RealName: "Test User",            // 设置真实姓名
			Phone:    "13800138000",          // 设置电话号码
			Mail:     "test@example.com",     // 设置邮箱
		}

		resp, err := logic.UserRegister(req) // 调用用户注册方法
		require.NoError(t, err, "注册应该成功")    // 断言没有错误
		require.NotNil(t, resp, "响应不应为空")    // 断言响应不为空

		assert.Equal(t, req.Username, resp.Username) // 断言用户名相等
		assert.NotEmpty(t, resp.CreateTime)          // 断言创建时间不为空
		assert.Equal(t, "注册成功", resp.Message)        // 断言注册成功消息

		// 验证布隆过滤器
		exists, err := svcCtx.BloomFilters.UserExists(ctx, req.Username)
		require.NoError(t, err, "布隆过滤器检查失败")    // 断言没有错误
		assert.True(t, exists, "用户应该存在于布隆过滤器中") // 断言用户存在于布隆过滤器中
	})

	t.Run("重复注册用户名", func(t *testing.T) {
		username := generateTestUsername() // 生成测试用户名
		firstReq := &types.UserRegisterReq{
			Username: username,           // 设置用户名
			Password: "password123",      // 设置密码
			RealName: "Test User",        // 设置真实姓名
			Phone:    "13800138000",      // 设置电话号码
			Mail:     "test@example.com", // 设置邮箱
		}

		// 第一次注册
		resp1, err := logic.UserRegister(firstReq) // 调用用户注册方法
		require.NoError(t, err, "第一次注册应该成功")       // 断言没有错误
		require.NotNil(t, resp1, "第一次注册响应不应为空")    // 断言响应不为空

		// 尝试重复注册
		resp2, err := logic.UserRegister(firstReq) // 再次调用用户注册方法
		require.Error(t, err, "重复注册应该失败")          // 断言有错误
		assert.Nil(t, resp2, "重复注册响应应该为空")         // 断言响应为空

		// 检查错误类型
		var userErr *errorx.UserError
		ok := errors.As(err, &userErr)                                          // 检查错误是否为用户错误
		assert.True(t, ok, "应该返回用户错误")                                          // 断言是用户错误
		assert.Equal(t, errorx.UserNameExistError, userErr.Code, "应该是用户名已存在错误") // 断言错误码相等

		// 验证布隆过滤器
		exists, err := svcCtx.BloomFilters.UserExists(ctx, username)
		require.NoError(t, err, "布隆过滤器检查失败")    // 断言没有错误
		assert.True(t, exists, "用户应该存在于布隆过滤器中") // 断言用户存在于布隆过滤器中
	})

	t.Run("并发注册相同用户名", func(t *testing.T) {
		username := generateTestUsername() // 生成测试用户名
		concurrency := 10                  // 设置并发数
		type result struct {
			resp *types.UserRegisterResp // 注册响应
			err  error                   // 错误
		}
		results := make(chan result, concurrency) // 创建结果通道

		// 并发注册
		for i := 0; i < concurrency; i++ {
			go func() {
				req := &types.UserRegisterReq{
					Username: username,           // 设置用户名
					Password: "password123",      // 设置密码
					RealName: "Test User",        // 设置真实姓名
					Phone:    "13800138000",      // 设置电话号码
					Mail:     "test@example.com", // 设置邮箱
				}
				resp, err := logic.UserRegister(req) // 调用用户注册方法
				results <- result{resp, err}         // 将结果发送到通道
			}()
		}

		// 统计结果
		successCount := 0     // 成功计数
		userErrorCount := 0   // 用户错误计数
		systemErrorCount := 0 // 系统错误计数

		for i := 0; i < concurrency; i++ {
			r := <-results // 从通道接收结果
			switch {
			case r.err == nil:
				successCount++ // 成功计数加一
			case errorx.IsUserError(r.err):
				userErrorCount++ // 用户错误计数加一
			default:
				systemErrorCount++ // 系统错误计数加一
			}
		}

		// 验证结果
		assert.Equal(t, 1, successCount, "应该只有一次注册成功")               // 断言只有一次注册成功
		assert.Equal(t, concurrency-1, userErrorCount, "其他应该都是用户错误") // 断言其他都是用户错误
		assert.Equal(t, 0, systemErrorCount, "不应该有系统错误")             // 断言没有系统错误

		// 验证布隆过滤器
		exists, err := svcCtx.BloomFilters.UserExists(ctx, username)
		require.NoError(t, err, "布隆过滤器检查失败")    // 断言没有错误
		assert.True(t, exists, "用户应该存在于布隆过滤器中") // 断言用户存在于布隆过滤器中
	})
}
