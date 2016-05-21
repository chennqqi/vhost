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
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"

	log "github.com/Sirupsen/logrus"
	"github.com/alex-oleshkevich/vhost/util"
	"github.com/spf13/cobra"
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

// ApacheTemplateData struct
type ApacheTemplateData struct {
	IP       string
	Port     int
	Domain   string
	SitesDir string
	SubDir   string
	SSL      string
}

// ApacheSSLData struct
type ApacheSSLData struct {
	CertFile string
	KeyFile  string
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

		// handle Apache virtual host
		if httpApache {
			if dumpVhostOnly {
				log.Info("Generated Apache vhost config:")
				buffer := generateApacheVhost(name)
				fmt.Print("\n", buffer.String(), "\n\n")
			} else {
				if isApacheVhostExists(name) && !force {
					log.Warn("Apache vhost with name \"" + name + "\" already exists. Do you want to enable it?")
					return
				}
				makeApacheVhost(name)
				makeDirs(name, httpSubDir)
				addToHosts(vhostIPAddress, makeDomain(name))
				reloadApache()
			}
		}
		// end Apache virtual host
	},
}

func init() {
	RootCmd.AddCommand(createCmd)

	createCmd.Flags().BoolVarP(&createMysqlDb, "mysql", "m", false, "Create MySQL database.")
	createCmd.Flags().StringVar(&mysqlDbName, "dbname", "", "MySQL database name to use.")
	createCmd.Flags().StringVar(&vhostIPAddress, "ip", "127.0.0.1", "IP address to use.")
	createCmd.Flags().IntVar(&vhostPort, "port", 80, "Port number to use.")
	createCmd.Flags().StringVar(&httpSubDir, "subdir", "", "Point document root into this subdirectory.")
	createCmd.Flags().BoolVarP(&httpApache, "apache", "a", false, "Virtual host for Apache.")
	createCmd.Flags().BoolVar(&dumpVhostOnly, "dump", false, "Dump generated vhost config and exit. Does not writes anything.")
	createCmd.Flags().BoolVar(&enableSSL, "ssl", false, "Enable SSL.")
	createCmd.Flags().BoolVar(&force, "force", false, "Force action.")
}

func isApacheVhostExists(name string) bool {
	var _, err = os.Stat(makeApacheVhostPath(name))
	return err == nil
}

func makeApacheVhostPath(name string) string {
	return config.Apache.DirSitesAvailable + "/" + name + config.Apache.VhostFileSuffix
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

func makeApacheVhost(name string) {
	buffer := generateApacheVhost(name)
	ioutil.WriteFile(makeApacheVhostPath(name), buffer.Bytes(), 0755)
}

func generateApacheVhost(name string) bytes.Buffer {
	apacheVhostTemplate, err := util.FindFile("share/apache.conf")
	if err != nil {
		panic("Apache virtual host template was not found.")
	}

	if httpSubDir != "" {
		httpSubDir = "/" + httpSubDir
	}

	data := ApacheTemplateData{
		vhostIPAddress, vhostPort, makeDomain(name), config.General.SitesDir, httpSubDir, makeApacheSSLPart(),
	}

	templateData, err := ioutil.ReadFile(apacheVhostTemplate)
	if err != nil {
		panic(err)
	}
	tmpl, err := template.New("apache").Parse(string(templateData))
	if err != nil {
		panic(err)
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, data)
	if err != nil {
		panic(err)
	}
	return buffer
}

func makeApacheSSLPart() string {
	if !enableSSL {
		return ""
	}

	templateFile, err := util.FindFile("share/apache-ssl.conf")
	if err != nil {
		panic("SSL config for Apache was not found")
	}

	data := ApacheSSLData{
		config.Ssl.CertFile, config.Ssl.KeyFile,
	}

	content, err := ioutil.ReadFile(templateFile)
	if err != nil {
		panic(err)
	}
	tmpl, err := template.New("apache-ssl").Parse(string(content))
	if err != nil {
		panic(err)
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, data)
	return buffer.String()
}

func makeDomain(name string) string {
	return name + config.General.Domain
}

func addToHosts(ip string, domain string) {
	hosts := util.LoadHosts(config.General.HostsFile)
	
	if hosts.Has(domain) && !force {
		log.Warn("Domain \"" + domain + "\" is already in \"" + config.General.HostsFile + "\".")
		return
	}
	
	hosts.Add(ip, domain)
	hosts.Write()
}

func reloadApache() {
	var arguments []string
	var output bytes.Buffer
	
	command := strings.Split(config.Apache.RestartCommand, " ")
	if len(command) > 1 {
		arguments = command[1:]
	}
	cmd := exec.Command(command[0], strings.Join(arguments, " "))
	cmd.Stdout = &output
	err := cmd.Run()
	if err != nil {
		log.Error(err)
	}
	log.Info(output.String())
}