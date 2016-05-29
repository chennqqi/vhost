// Copyright © 2016 Alex Oleshkevich <alex.oleshkevich@gmail.com>
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
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/alex-oleshkevich/vhost/util"
	"github.com/alex-oleshkevich/vhost/adapt"
	"github.com/spf13/cobra"
	"io/ioutil"
)

var (
	createMysqlDb  bool
	mysqlDbName    string
	httpApache     bool
	httpNginx      bool
	httpSubDir     string
	vhostIPAddress string
	vhostPort      int
	dumpVhostOnly  bool
	enableSSL      bool
	force          bool
)

var config = util.LoadConfig()

func init() {
	RootCmd.AddCommand(createCmd)

	createCmd.Flags().BoolVarP(&createMysqlDb, "mysql", "m", false, "Create MySQL database.")
	createCmd.Flags().StringVar(&mysqlDbName, "dbname", "", "MySQL database name to use.")
	createCmd.Flags().StringVar(&vhostIPAddress, "ip", "*", "IP address to use.")
	createCmd.Flags().IntVar(&vhostPort, "port", 80, "Port number to use.")
	createCmd.Flags().StringVar(&httpSubDir, "subdir", "", "Point document root into this subdirectory.")
	createCmd.Flags().BoolVarP(&httpApache, "apache", "a", false, "Virtual host for Apache.")
	createCmd.Flags().BoolVar(&dumpVhostOnly, "dump", false, "Dump generated vhost config and exit. Does not writes anything.")
	createCmd.Flags().BoolVar(&enableSSL, "ssl", false, "Enable SSL.")
	createCmd.Flags().BoolVar(&force, "force", false, "Force action.")
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Aliases: []string{"c"},
	Use:     "create",
	Short:   "Create a new virtual host",
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Errorln("No vhost name provided. Exit.")
			return
		}
		var name = args[0]
		if mysqlDbName == "" {
			mysqlDbName = strings.Replace(name, "-", "_", -1)
		}

		log.Infoln("Creating vhost: " + name)
		if createMysqlDb {
			log.Infoln(" --> add database: " + mysqlDbName)
		}

        var httpServer adapt.HTTPServer
		if httpApache {
            httpServer = adapt.ApacheServer{}
        } else {
            log.Fatal("No webserver specified in flags. Use --apache or --nginx.")
            return
        }
        
        if httpSubDir != "" {
            httpSubDir = "/" + httpSubDir
        }
        docRoot := config.General.SitesDir + "/" + httpServer.GetDomain(name) + "/www" + httpSubDir
        
        vhostConfig := httpServer.GetVhostConfig(name, vhostIPAddress, vhostPort, docRoot, enableSSL)
        
        if dumpVhostOnly {
            log.Info("Generated vhost config:")
            fmt.Print("\n", vhostConfig.String(), "\n\n")
        } else {
            if httpServer.IsExists(name) && !force {
                log.Warn("Virtual host with name '" + name + "' already exists in '" + httpServer.GetVhostPath(name) + "'. Did you want to enable it?")
                return
            }
            
            ioutil.WriteFile(httpServer.GetVhostPath(name), vhostConfig.Bytes(), 0755)
            makeDirs(name, httpSubDir)
            addToHosts(vhostIPAddress, httpServer.GetDomain(name))
            enable(name, httpServer)
            reload(httpServer)
        }

	},
}

func makeDirs(name string, httpSubDir string) {
	var rootDir = config.General.SitesDir
	var domain = name + config.General.Domain
	var siteRoot = rootDir + "/" + domain
	var docRoot = siteRoot + "/www"
	if httpSubDir != "" {
		docRoot += "/" + httpSubDir
	}
	var dirs = []string{
		docRoot,
		siteRoot + "/log",
		siteRoot + "/tmp",
	}

	for _, dir := range dirs {
		log.Info("--> create: " + dir)
		os.MkdirAll(dir, 0755)
	}
}

func addToHosts(ip string, domain string) {
    log.Info("Adding domain '" + domain + "' to ", config.General.HostsFile)
	hosts := util.LoadHosts(config.General.HostsFile)
	
	if hosts.Has(domain) && !force {
		log.Warn("Domain \"" + domain + "\" is already in \"" + config.General.HostsFile + "\".")
		return
	}
	
	hosts.Add(ip, domain)
	hosts.Write()
}

func reload(server adapt.HTTPServer) {
    command := server.GetReloadCommand()
    log.Info("Reloading Apache webserver")
    log.Info("Executing: /bin/sh -c " + command)
    out, err := exec.Command("/bin/sh", "-c", command).Output()
    if err != nil {
        log.Error("Error restarting apache: ", err)
    } else {
        fmt.Println(out)
    }
}

func enable(name string, server adapt.HTTPServer) {
    log.Info("Enable vhost")
    server.Enable(name)
}