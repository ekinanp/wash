package cmd

import (
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Benchkram/errz"
	"github.com/puppetlabs/wash/cmd/internal/server"
	cmdutil "github.com/puppetlabs/wash/cmd/util"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func serverCommand() *cobra.Command {
	serverCmd := &cobra.Command{
		Use:   "server <mountpoint>",
		Short: "Sets up the Wash daemon (API and FUSE servers)",
		Long: `Initializes all of the plugins, then sets up the Wash daemon (its API and FUSE servers).
To stop it, make sure you're not using the filesystem at <mountpoint>, then enter Ctrl-C.`,
		Args:   cobra.MinimumNArgs(1),
		PreRun: bindServerArgs,
		RunE:   toRunE(serverMain),
	}
	addServerArgs(serverCmd)

	return serverCmd
}

func serverMain(cmd *cobra.Command, args []string) exitCode {
	mountpoint := args[0]
	mountpoint, err := filepath.Abs(mountpoint)
	if err != nil {
		cmdutil.ErrPrintf("Could not compute the absolute path of the mountpoint %v: %v", mountpoint, err)
		return exitCode{1}
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
	})

	// On Ctrl-C, trigger the clean-up. This consists of shutting down the API
	// server and unmounting the FS.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	srv := server.New(mountpoint, serverOptsFromFlags())
	if err := srv.Start(); err != nil {
		log.Warn(err)
		return exitCode{1}
	}
	srv.Wait(sigCh)
	return exitCode{0}
}

func addServerArgs(cmd *cobra.Command) {
	cmd.Flags().String("loglevel", "info", "Set the logging level")
	cmd.Flags().String("logfile", "", "Set the log file's location. Defaults to stdout")
	cmd.Flags().String("cpuprofile", "", "Write cpu profile to file")
}

func bindServerArgs(cmd *cobra.Command, args []string) {
	// Only bind config lookup when invoking the specific command as viper bindings are global.
	errz.Fatal(viper.BindPFlag("loglevel", cmd.Flags().Lookup("loglevel")))
	errz.Fatal(viper.BindPFlag("logfile", cmd.Flags().Lookup("logfile")))
	errz.Fatal(viper.BindPFlag("cpuprofile", cmd.Flags().Lookup("cpuprofile")))
}

// OptsFromFlags returns server.Opts as set by command-line flags.
func serverOptsFromFlags() server.Opts {
	return server.Opts{
		CPUProfilePath:      viper.GetString("cpuprofile"),
		LogFile:             viper.GetString("logfile"),
		LogLevel:            viper.GetString("loglevel"),
	}
}
