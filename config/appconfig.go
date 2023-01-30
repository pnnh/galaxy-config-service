package config

import (
	"context" 
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/appconfig"
	"github.com/sirupsen/logrus"
)

var awsAppConfigClient *appconfig.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-east-1"))
	if err != nil {
		logrus.Fatalf("unable to load SDK config, %v", err) 
	}
	svc := appconfig.NewFromConfig(cfg)
	awsAppConfigClient = svc

}

func LoadAwsConfig(fileName, env string) string {
	hostname, err := os.Hostname()
	if err != nil {
		logrus.Fatalln("获取主机名出错", err)
	}

	in := &appconfig.GetConfigurationInput{
		Application:   aws.String("multiverse.direct"),
		ClientId:      aws.String(hostname),
		Configuration: aws.String(fileName),
		Environment:   aws.String(env),
	}
	if Debug() {
		in.Application = aws.String("debug.multiverse.direct") 
	}
	out, err := awsAppConfigClient.GetConfiguration(context.Background(), in)
	if err != nil {
		logrus.Fatalln("获取配置出错", fileName, env, err)
	}
	content := string(out.Content)
	return content
}

func GetConfigurationMap() (map[string]string, error) {
	var cmdEnv []string

	awsConfig := LoadAwsConfig("main.config", "default")
	awsEnvs := strings.Split(awsConfig, "\n")
	for _, e := range awsEnvs {
		cmdEnv = append(cmdEnv, e)
	}
	// 系统环境变量可以覆盖掉默认配置
	osEnv := os.Environ()
	for _, e := range osEnv {
		cmdEnv = append(cmdEnv, e)
	}
	configMap := make(map[string]string)
	for _, e := range cmdEnv {
		index := strings.Index(e, "=")
		if index > 0 {
			configMap[e[:index]] = e[index+1:]
		}
	}

	return configMap, nil
}
