// Package config implements configuration for the wash executable using
// https://github.com/spf13/viper.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Contains all the keys for Wash's shared config
const (
	SocketKey = "socket"
)

// Socket is the path to the Wash server's UNIX
// socket
var Socket string

// Load Wash's config.
func Load() error {
	// Set any defaults
	cdir, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	viper.SetDefault(SocketKey, filepath.Join(cdir, "wash", "wash-api.sock"))

	// Tell viper that the config. can be read from WASH_<entry>
	// environment variables
	viper.SetEnvPrefix("WASH")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Set the default config file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// TODO: Should we error here, or log a warning and continue? UserHomeDir()
		// checks environment variables
		return fmt.Errorf("could not determine the home directory: %v", homeDir)
	}
	viper.SetConfigName("wash")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Join(homeDir, ".puppetlabs", "wash"))
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}
	// Load the shared config
	Socket = viper.GetString(SocketKey)

	return nil
}
