package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"sgridnext.com/src/config"
)

func Test_ReadJson(t *testing.T) {
	data := config.ReadJson("./test.json")
	t.Logf("data %s \n", data)
	assert.Equal(t, data["serverName"], "TestServer")
}
