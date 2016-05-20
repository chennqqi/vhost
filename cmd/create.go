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
	"io/ioutil"
	"strings"
	"text/template"

	log "github.com/Sirupsen/logrus"
	"github.com/alex-oleshkevich/vhost/util"
	"github.com/spf13/cobra"
	"os"
)

var (
	createMysqlDb bool
	mysqlDbName   string
	httpApache    bool
	httpNginx     bool
	httpSubDir    string
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
			makeDirs(name, httpSubDir)
			makeApacheVhost(name)
		}
		// end Apache virtual host
	},
}

func init() {
	RootCmd.AddCommand(createCmd)

	createCmd.Flags().BoolVarP(&createMysqlDb, "mysql", "m", false, "Create MySQL database")
	createCmd.Flags().StringVar(&mysqlDbName, "dbname", "", "MySQL database name to use")
	createCmd.Flags().StringVar(&httpSubDir, "subdir", "", "Point document root into this subdirectory")
	createCmd.Flags().BoolVarP(&httpApache, "apache", "a", false, "Virtual host for Apache")
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
	apacheVhostTemplate, found := util.FindFile("share/apache.conf")
	if !found {
		panic("Apache virtual host template was not found.")
	}

	data := ApacheTemplateData{
		"127.0.0.1", 80, name + config.General.Domain, config.General.SitesDir, "/", "",
	}

	templateData, err := ioutil.ReadFile(apacheVhostTemplate)
	tmpl, err := template.New("apache").Parse(string(templateData))
	if err != nil {
		panic(err)
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, data)
	if err != nil {
		panic(err)
	}

	log.Debugln("Apache template:")
	log.Debugln(buffer.String())

	var outfile = config.Apache.DirSitesAvailable + "/" + name + config.Apache.VhostFileSuffix
	ioutil.WriteFile(outfile, buffer.Bytes(), 0755)
}
