package svc

import (
	"fmt"
	"hash/crc32"
	"shorterurl/link/rpc/internal/config"
	"shorterurl/link/rpc/internal/model"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var testIdSequence int64 = 100000 // 使用一个足够大的起始值避免和已有数据冲突

// 测试专用ID生成器
func getTestID() int64 {
	return atomic.AddInt64(&testIdSequence, 1)
}

// 测试数据库连接和分片
func TestInitDBs(t *testing.T) {
	// 清理之前运行测试可能遗留的数据
	cleanupTestData(t)

	// 初始化配置
	c := config.Config{}
	c.DB.Host = "localhost"
	c.DB.Port = 3306
	c.DB.User = "root"
	c.DB.Password = "123456"
	c.DB.Database = "shorterurl"
	c.DB.Sharding.ShardingKey = "gid"
	c.DB.Sharding.NumberOfShards = 16

	// 使用测试ID生成器
	idGen := getTestID

	// 初始化数据库连接
	dbs, err := InitDBs(c, idGen)
	assert.NoError(t, err)
	assert.NotNil(t, dbs)

	// 确保所有数据库连接都成功创建
	assert.NotNil(t, dbs.Common)
	assert.NotNil(t, dbs.LinkDB)
	assert.NotNil(t, dbs.GotoLinkDB)
	assert.NotNil(t, dbs.GroupDB)
	assert.NotNil(t, dbs.UserDB)

	// 确保所有分片中间件都成功创建
	assert.NotNil(t, dbs.Shardings["t_link"])
	assert.NotNil(t, dbs.Shardings["t_link_goto"])
	assert.NotNil(t, dbs.Shardings["t_group"])
	assert.NotNil(t, dbs.Shardings["t_user"])

	// 测试各个分片表的操作
	t.Run("TestLinkShard", func(t *testing.T) {
		testLinkShard(t, dbs.LinkDB)
	})

	t.Run("TestLinkGotoShard", func(t *testing.T) {
		testLinkGotoShard(t, dbs.GotoLinkDB)
	})

	t.Run("TestGroupShard", func(t *testing.T) {
		testGroupShard(t, dbs.GroupDB)
	})

	t.Run("TestUserShard", func(t *testing.T) {
		testUserShard(t, dbs.UserDB)
	})
}

// 清理测试数据
func cleanupTestData(t *testing.T) {
	// 创建一个数据库连接来清理数据
	dsn := "root:123456@tcp(localhost:3306)/shorterurl?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Logf("打开数据库连接失败: %v", err)
		return
	}

	// 清理t_link表
	for i := 0; i < 16; i++ {
		tableName := fmt.Sprintf("t_link_%d", i)
		db.Exec(fmt.Sprintf("DELETE FROM %s WHERE gid LIKE 'test%%'", tableName))
	}

	// 清理t_link_goto表
	for i := 0; i < 16; i++ {
		tableName := fmt.Sprintf("t_link_goto_%d", i)
		db.Exec(fmt.Sprintf("DELETE FROM %s WHERE full_short_url LIKE 'test%%'", tableName))
	}

	// 清理t_group表
	for i := 0; i < 16; i++ {
		tableName := fmt.Sprintf("t_group_%d", i)
		db.Exec(fmt.Sprintf("DELETE FROM %s WHERE username LIKE 'test%%'", tableName))
	}

	// 清理t_user表
	for i := 0; i < 16; i++ {
		tableName := fmt.Sprintf("t_user_%d", i)
		db.Exec(fmt.Sprintf("DELETE FROM %s WHERE username LIKE 'test%%'", tableName))
	}
}

// 测试t_link表的分片
func testLinkShard(t *testing.T, db *gorm.DB) {
	// 先清理测试数据
	for i := 0; i < 16; i++ {
		tableName := fmt.Sprintf("t_link_%d", i)
		db.Exec(fmt.Sprintf("DELETE FROM %s WHERE gid LIKE 'test%%'", tableName))
	}

	// 创建测试数据
	for i := 0; i < 10; i++ {
		gid := fmt.Sprintf("test_gid_%d", i)
		link := model.Link{
			ID:            getTestID(),
			Gid:           gid,
			Domain:        "s.xleft.cn",
			ShortUri:      fmt.Sprintf("t%d", i), // 确保不超过8个字符
			FullShortUrl:  fmt.Sprintf("https://s.xleft.cn/t%d", i),
			OriginUrl:     fmt.Sprintf("http://example.com/%d", i),
			ClickNum:      0,
			EnableStatus:  0,
			CreatedType:   0,
			ValidDateType: 0,
			ValidDate:     time.Now(),
			Describe:      fmt.Sprintf("测试短链接_%d", i),
			TotalPv:       0,
			TotalUv:       0,
			TotalUip:      0,
			CreateTime:    time.Now(),
			UpdateTime:    time.Now(),
			DelTime:       0,
			DelFlag:       0,
		}

		result := db.Create(&link)
		assert.NoError(t, result.Error)
		assert.NotZero(t, link.ID)

		// 验证能够查询到
		var foundLink model.Link
		result = db.Where("gid = ?", gid).First(&foundLink)
		assert.NoError(t, result.Error)
		assert.Equal(t, gid, foundLink.Gid)
		assert.Equal(t, link.OriginUrl, foundLink.OriginUrl)
	}

	// 测试不同分片的记录分布
	var counts []struct {
		Table string
		Count int
	}

	for i := 0; i < 16; i++ {
		var count int64
		tableName := fmt.Sprintf("t_link_%d", i)
		db.Table(tableName).Where("gid LIKE 'test%'").Count(&count)
		counts = append(counts, struct {
			Table string
			Count int
		}{
			Table: tableName,
			Count: int(count),
		})
	}

	// 记录各分片表中的数据分布情况
	t.Logf("t_link分片数据分布: %+v", counts)
}

// 测试t_link_goto表的分片
func testLinkGotoShard(t *testing.T, db *gorm.DB) {
	// 先清理测试数据
	for i := 0; i < 16; i++ {
		tableName := fmt.Sprintf("t_link_goto_%d", i)
		db.Exec(fmt.Sprintf("DELETE FROM %s WHERE full_short_url LIKE 'test%%'", tableName))
	}

	// 创建测试数据
	for i := 0; i < 10; i++ {
		shortUrl := fmt.Sprintf("test_short_url_%d", i)
		linkGoto := model.LinkGoto{
			ID:           getTestID(),
			FullShortUrl: shortUrl,
			Gid:          fmt.Sprintf("test_gid_%d", i),
		}

		result := db.Create(&linkGoto)
		assert.NoError(t, result.Error)
		assert.NotZero(t, linkGoto.ID)

		// 验证能够查询到
		var foundLinkGoto model.LinkGoto
		result = db.Where("full_short_url = ?", shortUrl).First(&foundLinkGoto)
		assert.NoError(t, result.Error)
		assert.Equal(t, shortUrl, foundLinkGoto.FullShortUrl)
		assert.Equal(t, linkGoto.Gid, foundLinkGoto.Gid)
	}

	// 测试不同分片的记录分布
	var counts []struct {
		Table string
		Count int
	}

	for i := 0; i < 16; i++ {
		var count int64
		tableName := fmt.Sprintf("t_link_goto_%d", i)
		db.Table(tableName).Where("full_short_url LIKE 'test%'").Count(&count)
		counts = append(counts, struct {
			Table string
			Count int
		}{
			Table: tableName,
			Count: int(count),
		})
	}

	// 记录各分片表中的数据分布情况
	t.Logf("t_link_goto分片数据分布: %+v", counts)
}

// 测试t_group表的分片
func testGroupShard(t *testing.T, db *gorm.DB) {
	// 先清理测试数据
	for i := 0; i < 16; i++ {
		tableName := fmt.Sprintf("t_group_%d", i)
		db.Exec(fmt.Sprintf("DELETE FROM %s WHERE username LIKE 'test%%'", tableName))
	}

	// 创建测试数据
	for i := 0; i < 10; i++ {
		username := fmt.Sprintf("test_user_%d", i)
		group := model.Group{
			ID:         getTestID(),
			Username:   username,
			Name:       fmt.Sprintf("测试分组_%d", i),
			Gid:        fmt.Sprintf("test_gid_%d", i),
			SortOrder:  i,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
			DelFlag:    0,
		}

		result := db.Create(&group)
		assert.NoError(t, result.Error)
		assert.NotZero(t, group.ID)

		// 验证能够查询到
		var foundGroup model.Group
		result = db.Where("username = ?", username).First(&foundGroup)
		assert.NoError(t, result.Error)
		assert.Equal(t, username, foundGroup.Username)
		assert.Equal(t, group.Name, foundGroup.Name)
	}

	// 测试不同分片的记录分布
	var counts []struct {
		Table string
		Count int
	}

	for i := 0; i < 16; i++ {
		var count int64
		tableName := fmt.Sprintf("t_group_%d", i)
		db.Table(tableName).Where("username LIKE 'test%'").Count(&count)
		counts = append(counts, struct {
			Table string
			Count int
		}{
			Table: tableName,
			Count: int(count),
		})
	}

	// 记录各分片表中的数据分布情况
	t.Logf("t_group分片数据分布: %+v", counts)
}

// 测试t_user表的分片
func testUserShard(t *testing.T, db *gorm.DB) {
	// 先清理测试数据
	for i := 0; i < 16; i++ {
		tableName := fmt.Sprintf("t_user_%d", i)
		db.Exec(fmt.Sprintf("DELETE FROM %s WHERE username LIKE 'test%%'", tableName))
	}

	// 创建测试数据
	for i := 0; i < 10; i++ {
		username := fmt.Sprintf("test_user_%d", i)
		user := model.User{
			ID:         getTestID(),
			Username:   username,
			Password:   "password",
			RealName:   "",
			Mail:       fmt.Sprintf("test%d@example.com", i),
			Phone:      fmt.Sprintf("1380000%04d", i),
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
			DelFlag:    0,
		}

		result := db.Create(&user)
		assert.NoError(t, result.Error)
		assert.NotZero(t, user.ID)

		// 验证能够查询到
		var foundUser model.User
		result = db.Where("username = ?", username).First(&foundUser)
		assert.NoError(t, result.Error)
		assert.Equal(t, username, foundUser.Username)
		assert.Equal(t, user.Mail, foundUser.Mail)
	}

	// 测试不同分片的记录分布
	var counts []struct {
		Table string
		Count int
	}

	for i := 0; i < 16; i++ {
		var count int64
		tableName := fmt.Sprintf("t_user_%d", i)
		db.Table(tableName).Where("username LIKE 'test%'").Count(&count)
		counts = append(counts, struct {
			Table string
			Count int
		}{
			Table: tableName,
			Count: int(count),
		})
	}

	// 记录各分片表中的数据分布情况
	t.Logf("t_user分片数据分布: %+v", counts)
}

// 测试字符串分片键的分片算法效果
func TestShardingAlgorithm(t *testing.T) {
	// 模拟分片算法
	numberOfShards := 16
	values := []string{
		"test_user_1",
		"test_user_2",
		"test_gid_101",
		"test_short_url_abc",
		"username_admin",
		"gid_group1",
	}

	// 打印不同值映射到的分片
	for _, value := range values {
		hash := crc32.ChecksumIEEE([]byte(value))
		shard := int(hash % uint32(numberOfShards))
		t.Logf("值 '%s' 映射到分片: %d", value, shard)
	}
}
