// Copyright © 2017 Alex Oleshkevich <alex.oleshkevich@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var debugMode bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "vhost",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

// Execute command
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/vhost/config.yaml)")
	RootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "Enable debug output.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if debugMode {
		log.SetLevel(log.DebugLevel)
	}

	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName("config")       // name of config file (without extension)
	viper.AddConfigPath(".")            // adding current directory as search path
	viper.AddConfigPath("$HOME/.vhost") // adding home directory as search path
	viper.AddConfigPath("/etc/vhost")   // adding global directory as search path
	viper.AutomaticEnv()                // read in environment variables that match

	viper.SetDefault("sites-enabled", "/etc/nginx/sites-enabled")
	viper.SetDefault("domain", ".lan")
	viper.SetDefault("mysql-host", "127.0.0.1")
	viper.SetDefault("mysql-port", "3306")
	viper.SetDefault("mysql-user", "root")
	viper.SetDefault("mysql-pass", "")
	viper.SetDefault("postgres-host", "127.0.0.1")
	viper.SetDefault("postgres-port", "5432")
	viper.SetDefault("postgres-user", "postgres")
	viper.SetDefault("postgres-pass", "")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debugf("Using config file: %s", viper.ConfigFileUsed())
	}
}
