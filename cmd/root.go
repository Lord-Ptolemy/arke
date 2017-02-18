package cmd // import "iris.arke.works/forum/cmd"

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"iris.arke.works/forum/http"
	"net"
)

var log *zap.Logger

// RootCmd exports the main arke command interface, including all subcommands
var RootCmd = &cobra.Command{
	Use:   "arke",
	Short: "TBD",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		lvl := zap.NewAtomicLevel()
		// Set Logging Level
		if viper.GetBool("verbose") {
			lvl.SetLevel(zap.DebugLevel)
		} else if viper.IsSet("log.level") {
			switch viper.GetString("log.level") {
			case "debug":
				lvl.SetLevel(zap.DebugLevel)
			case "warn":
				lvl.SetLevel(zap.WarnLevel)
			case "dpanic":
				lvl.SetLevel(zap.DPanicLevel)
			case "error":
				lvl.SetLevel(zap.ErrorLevel)
			case "fatal":
				lvl.SetLevel(zap.FatalLevel)
			case "info":
				lvl.SetLevel(zap.InfoLevel)
			case "panic":
				lvl.SetLevel(zap.PanicLevel)
			default:
				return errors.New("Unknown Logging Level")
			}
		}
		var err error
		conf := zap.NewDevelopmentConfig()
		conf.Level = lvl
		conf.DisableStacktrace = true
		conf.DisableCaller = true
		log, err = conf.Build()
		return err
	},
	Run: server,
}

func init() {
	// We must set the persistent flags *before* we start up the configuration init function
	RootCmd.PersistentFlags().String("log.level", "warn", "Set logging level")
	RootCmd.PersistentFlags().Bool("verbose", false, "Enable Verbose Mode. This overrides the logging level setting")

	RootCmd.Flags().String("listen", "[::1]:80", "Interface and Port to listen on")

	initConf()
}

func server(cmd *cobra.Command, args []string) {
	iface, err := cmd.Flags().GetString("listen")
	if err != nil {
		log.Error("Error while getting interface", zap.Error(err))
	}
	ip, err := net.ResolveTCPAddr("tcp", iface)
	if err != nil {
		log.Error("Error parsing interfaace", zap.Error(err))
	}
	err = http.StartBlocking(ip, log)
	if err != nil {
		log.Error("Error on HTTP Server", zap.Error(err))
	}
}
