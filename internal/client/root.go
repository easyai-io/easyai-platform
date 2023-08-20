// Package client ...
/**
Copyright © 2023 shuaiyy
client: the CLI tool for easyai-platform
*/
package client

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/easyai-io/easyai-platform/internal/client/utils/stdlogger"
	"github.com/easyai-io/easyai-platform/pkg/auth"
	"github.com/easyai-io/easyai-platform/pkg/contextx"
)

var cfgFile string
var authWhiteList []string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "client",
	Short: "A CLI tool for easyai platform",
	Long: `client is a CLI client for easyai platform:

you can use client:
1. submit training jobs, and get/list/stop/delete jobs.
2. login a remote container of a job.
3. download files produced by a job.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.SetContext(context.Background())
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here, will be global for your application.
	// 全局flag
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.client/config.yaml)")

	// 全局 Auth 校验，替代 addAuthToCmd
	// 为部分命令添加全局auth, config,update不需要auth
	authWhiteList = []string{}
	rootCmd.PersistentPreRunE = globalAuthCmd
}

func globalAuthCmd(cmd *cobra.Command, _ []string) error {
	commandPath := cmd.CommandPath()
	stdlogger.Debugf("now path: %v", stdlogger.Yellow(commandPath))
	for _, whitePath := range authWhiteList {
		if commandPath == whitePath {
			return nil
		}
	}
	token := viper.GetString("user-token")
	uid, err := parseUserID(token)
	if err != nil {
		stdlogger.Error("解析用户token失败，请检查配置文件")
		stdlogger.Info("获取token方式：平台主页右上角; 配置token: client config init --help")
		return err
	}
	viper.Set("user-uid", uid)
	contextx.SetCmdUser(cmd, uid)
	stdlogger.Debugf("当前用户: %s", stdlogger.Yellow(uid))
	stdlogger.DebugSensitive(context.Background(), "my-token: %s", token)
	return nil
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else if v := os.Getenv("ClientConfig"); v != "" {
		viper.SetConfigFile(v)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cli" (without extension).
		viper.AddConfigPath(filepath.Join(home, ".client"))
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
	if viper.GetBool("debug") {
		stdlogger.SetDebugLevel()
	} else {
		stdlogger.SetInfoLevel()
	}
	// 	环境变量中获取的 CLIENT_TOKEN 和 EASYAI_PLATFORM_URL 优先级最高
	if token := os.Getenv("CLIENT_TOKEN"); token != "" {
		viper.Set("user-token", token)
	}
	if url := os.Getenv("EASYAI_PLATFORM_URL"); url != "" {
		viper.Set("host", url)
		if strings.HasPrefix(url, "http://") {
			viper.Set("port", 80)
			viper.Set("protocol", "http")
		} else {
			viper.Set("port", 443)
			viper.Set("protocol", "https")
		}
	}
}

func parseUserID(token string) (string, error) {
	signingKey := "easyai"
	var opts []auth.Option
	opts = append(opts, auth.SetExpired(3600*24*365*2))
	opts = append(opts, auth.SetSigningKey([]byte(signingKey)))
	opts = append(opts, auth.SetKeyfunc(func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, auth.ErrInvalidToken
		}
		return []byte(signingKey), nil
	}))
	var method jwt.SigningMethod = jwt.SigningMethodHS512
	opts = append(opts, auth.SetSigningMethod(method))

	auther := auth.New(nil, auth.NoCheckUserStatus, opts...)
	uid, _, _, _, err := auther.ParseUserInfo(context.Background(), token, true) // nolint:dogsled
	return uid, err
}
