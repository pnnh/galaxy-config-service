package config

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var ACCOUNT_DB_DSN = ""
var GINMODE = "release"
var ISSUER = "sfx.xyz" // TOTP发行机构
var JWTRealm = "sfx.xyz"
var JWTKey = ""
var CSRFToken = ""

func init() {

	mode := os.Getenv("MODE")
	if len(mode) > 0 {
		GINMODE = mode
	}
	configMap, err := GetConfigurationMap()
	if err != nil {
		logrus.Fatalln("获取appconfig配置出错: %w", err)
	}

	ACCOUNT_DB_DSN = configMap["ACCOUNT_DB"]
	if len(ACCOUNT_DB_DSN) < 1 {
		logrus.Fatalln("数据库未配置")
	}
}

func Debug() bool {
	return GINMODE == gin.DebugMode
}

func Test() bool {
	return GINMODE == gin.TestMode
}

func Release() bool {
	return GINMODE != gin.DebugMode && GINMODE != gin.TestMode
}

func LoadConfig(fileName, env string) (string, error) {
	var awsConfig string
	var err error

	if Debug() {
		awsConfig, err = LoadDebugConfig("main.config", "default")
	} else {
		awsConfig, err = LoadAwsConfig("main.config", "default")
	}
	return awsConfig, err
}