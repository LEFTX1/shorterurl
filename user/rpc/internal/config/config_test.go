package config

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "./test.yaml", "the config file")

func TestConfig(t *testing.T) {
	flag.Parse()

	var c Config

	conf.MustLoad(*configFile, &c)
	assert.NotEmpty(t, c.DB.Host)
	assert.NotEmpty(t, c.DB.User)
	assert.NotEmpty(t, c.BizRedis.Host)
	assert.NotEmpty(t, c.BloomFilter.User.Name)
	assert.NotEmpty(t, c.Crypto.AESKey)
	t.Logf("config: %+v", c)
}
