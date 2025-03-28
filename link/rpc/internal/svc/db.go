package svc

import (
	"fmt"
	"hash/crc32"
	"shorterurl/link/rpc/internal/config"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/sharding"
)

// DBs 包含所有数据库连接
type DBs struct {
	Common     *gorm.DB                      // 用于没有分片的表
	LinkDB     *gorm.DB                      // t_link表的分片DB
	GotoLinkDB *gorm.DB                      // t_link_goto表的分片DB
	GroupDB    *gorm.DB                      // t_group表的分片DB
	UserDB     *gorm.DB                      // t_user表的分片DB
	Shardings  map[string]*sharding.Sharding // 保存所有分片中间件
}

// InitDBs 初始化所有数据库连接
func InitDBs(c config.Config, idGen func() int64) (*DBs, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DB.User,
		c.DB.Password,
		c.DB.Host,
		c.DB.Port,
		c.DB.Database,
	)

	// 初始化普通DB连接
	commonDB, err := initCommonDB(dsn)
	if err != nil {
		return nil, err
	}

	// 初始化各个分片DB
	linkDB, linkSharding, err := initLinkDB(dsn, c, idGen)
	if err != nil {
		return nil, fmt.Errorf("init link db error: %v", err)
	}

	gotoLinkDB, gotoLinkSharding, err := initGotoLinkDB(dsn, c, idGen)
	if err != nil {
		return nil, fmt.Errorf("init goto link db error: %v", err)
	}

	groupDB, groupSharding, err := initGroupDB(dsn, c, idGen)
	if err != nil {
		return nil, fmt.Errorf("init group db error: %v", err)
	}

	userDB, userSharding, err := initUserDB(dsn, c, idGen)
	if err != nil {
		return nil, fmt.Errorf("init user db error: %v", err)
	}

	shardings := map[string]*sharding.Sharding{
		"t_link":      linkSharding,
		"t_link_goto": gotoLinkSharding,
		"t_group":     groupSharding,
		"t_user":      userSharding,
	}

	return &DBs{
		Common:     commonDB,
		LinkDB:     linkDB,
		GotoLinkDB: gotoLinkDB,
		GroupDB:    groupDB,
		UserDB:     userDB,
		Shardings:  shardings,
	}, nil
}

// 初始化通用数据库连接
func initCommonDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open common database failed: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB failed: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// 初始化t_link表的分片DB
func initLinkDB(dsn string, c config.Config, idGen func() int64) (*gorm.DB, *sharding.Sharding, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("open link database failed: %v", err)
	}

	configDBPool(db)

	// 使用gid作为分片键
	shardingInstance := sharding.Register(sharding.Config{
		ShardingKey:    "gid", // 分组标识作为分片键
		NumberOfShards: uint(c.DB.Sharding.NumberOfShards),
		ShardingAlgorithm: func(value any) (suffix string, err error) {
			var shard int
			switch v := value.(type) {
			case string:
				hash := crc32.ChecksumIEEE([]byte(v))
				shard = int(hash % uint32(c.DB.Sharding.NumberOfShards))
			case int, int32, int64:
				var id int64
				switch val := v.(type) {
				case int:
					id = int64(val)
				case int32:
					id = int64(val)
				case int64:
					id = val
				}
				shard = int(id % int64(c.DB.Sharding.NumberOfShards))
			default:
				return "", fmt.Errorf("unsupported type for sharding key")
			}
			return "_" + strconv.Itoa(shard), nil
		},
		ShardingSuffixs:     genShardingSuffixs(c.DB.Sharding.NumberOfShards),
		PrimaryKeyGenerator: sharding.PKCustom,
		PrimaryKeyGeneratorFn: func(tableIdx int64) int64 {
			return idGen()
		},
	}, "t_link")

	if err = db.Use(shardingInstance); err != nil {
		return nil, nil, fmt.Errorf("register t_link sharding middleware failed: %v", err)
	}

	return db, shardingInstance, nil
}

// 初始化t_link_goto表的分片DB
func initGotoLinkDB(dsn string, c config.Config, idGen func() int64) (*gorm.DB, *sharding.Sharding, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("open goto link database failed: %v", err)
	}

	configDBPool(db)

	// 使用full_short_url作为分片键
	shardingInstance := sharding.Register(sharding.Config{
		ShardingKey:    "full_short_url", // 完整短链接作为分片键
		NumberOfShards: uint(c.DB.Sharding.NumberOfShards),
		ShardingAlgorithm: func(value any) (suffix string, err error) {
			var shard int
			switch v := value.(type) {
			case string:
				hash := crc32.ChecksumIEEE([]byte(v))
				shard = int(hash % uint32(c.DB.Sharding.NumberOfShards))
			default:
				return "", fmt.Errorf("unsupported type for sharding key")
			}
			return "_" + strconv.Itoa(shard), nil
		},
		ShardingSuffixs:     genShardingSuffixs(c.DB.Sharding.NumberOfShards),
		PrimaryKeyGenerator: sharding.PKCustom,
		PrimaryKeyGeneratorFn: func(tableIdx int64) int64 {
			return idGen()
		},
	}, "t_link_goto")

	if err = db.Use(shardingInstance); err != nil {
		return nil, nil, fmt.Errorf("register t_link_goto sharding middleware failed: %v", err)
	}

	return db, shardingInstance, nil
}

// 初始化t_group表的分片DB
func initGroupDB(dsn string, c config.Config, idGen func() int64) (*gorm.DB, *sharding.Sharding, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("open group database failed: %v", err)
	}

	configDBPool(db)

	// 使用username作为分片键
	shardingInstance := sharding.Register(sharding.Config{
		ShardingKey:    "username", // 用户名作为分片键
		NumberOfShards: uint(c.DB.Sharding.NumberOfShards),
		ShardingAlgorithm: func(value any) (suffix string, err error) {
			var shard int
			switch v := value.(type) {
			case string:
				hash := crc32.ChecksumIEEE([]byte(v))
				shard = int(hash % uint32(c.DB.Sharding.NumberOfShards))
			default:
				return "", fmt.Errorf("unsupported type for sharding key")
			}
			return "_" + strconv.Itoa(shard), nil
		},
		ShardingSuffixs:     genShardingSuffixs(c.DB.Sharding.NumberOfShards),
		PrimaryKeyGenerator: sharding.PKCustom,
		PrimaryKeyGeneratorFn: func(tableIdx int64) int64 {
			return idGen()
		},
	}, "t_group")

	if err = db.Use(shardingInstance); err != nil {
		return nil, nil, fmt.Errorf("register t_group sharding middleware failed: %v", err)
	}

	return db, shardingInstance, nil
}

// 初始化t_user表的分片DB
func initUserDB(dsn string, c config.Config, idGen func() int64) (*gorm.DB, *sharding.Sharding, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("open user database failed: %v", err)
	}

	configDBPool(db)

	// 使用username作为分片键
	shardingInstance := sharding.Register(sharding.Config{
		ShardingKey:    "username", // 用户名作为分片键
		NumberOfShards: uint(c.DB.Sharding.NumberOfShards),
		ShardingAlgorithm: func(value any) (suffix string, err error) {
			var shard int
			switch v := value.(type) {
			case string:
				hash := crc32.ChecksumIEEE([]byte(v))
				shard = int(hash % uint32(c.DB.Sharding.NumberOfShards))
			default:
				return "", fmt.Errorf("unsupported type for sharding key")
			}
			return "_" + strconv.Itoa(shard), nil
		},
		ShardingSuffixs:     genShardingSuffixs(c.DB.Sharding.NumberOfShards),
		PrimaryKeyGenerator: sharding.PKCustom,
		PrimaryKeyGeneratorFn: func(tableIdx int64) int64 {
			return idGen()
		},
	}, "t_user")

	if err = db.Use(shardingInstance); err != nil {
		return nil, nil, fmt.Errorf("register t_user sharding middleware failed: %v", err)
	}

	return db, shardingInstance, nil
}

// 配置数据库连接池
func configDBPool(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get sql.DB failed: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return nil
}

// 生成分片后缀数组
func genShardingSuffixs(numberOfShards int) func() []string {
	return func() (suffixs []string) {
		for i := 0; i < numberOfShards; i++ {
			suffixs = append(suffixs, "_"+strconv.Itoa(i))
		}
		return
	}
}
