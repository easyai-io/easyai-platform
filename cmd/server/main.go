package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/easyai-io/easyai-platform/internal/server"
	"github.com/easyai-io/easyai-platform/internal/server/metrics"
	"github.com/easyai-io/easyai-platform/pkg/logger"
)

// VERSION Usage: go build -ldflags "-X main.Version=x.x.x -X main.Revision=xxxx -X main.Branch=xxxx -X main.BuildUser=xxxx -X main.BuildDate=xxxx"
var (
	Version   = "0.0.1"
	Revision  = "a.b.c"
	Branch    = "test"
	BuildUser = "shuaiyy"
	BuildDate = "2023-08-01"
)

var conf string
var www string

// @title easyai-platform
// @version 0.0.1
// @description a ml/dl platform build with golang + gin, based on k8s and cloud-native projects.
// @schemes http https
// @basePath /
// @contact.name shuaiyy
// @contact.email admin@easyai.io
// @contact.url https://easyai.io

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-Authorization-Token

func main() {
	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   "easyai-platform",
		Short: "the backend server for easyai machine learning platform",
		Long:  `a machine learning platform build with: golang + k8s`,
	}
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "print version",
		Long:  "print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("version: %s, revision: %s, branch: %s, buildUser: %s, buildDate: %s\n", Version, Revision, Branch, BuildUser, BuildDate)
		},
	}
	rootCmd.AddCommand(versionCmd)

	var webCmd = &cobra.Command{
		Use:     "web",
		Short:   "start web server",
		Long:    "start web server",
		Example: "./easyai-platform web --conf ./local.toml",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := logger.NewTagContext(cmd.Context(), "__easyai__")
			metrics.SetBuildInfo(Version, Revision, Branch, BuildUser, BuildDate)
			return server.Run(ctx,
				server.SetConfigFile(conf),
				server.SetWWWDir(www),
				server.SetVersion(Version))
		},
	}
	webCmd.Flags().StringVarP(&conf, "conf", "c", "", "App configuration file(*.toml)")
	webCmd.Flags().StringVarP(&www, "www", "w", "", "frontend static file directory")
	rootCmd.AddCommand(webCmd)

	if err := rootCmd.Execute(); err != nil {
		logger.Errorf("server.Run error: %v", err)
	}
}
