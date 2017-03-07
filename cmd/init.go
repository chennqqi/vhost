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
	"bufio"
	"fmt"
	"os"
	"strings"

	"text/template"

	"path/filepath"

	"os/exec"

	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/alex-oleshkevich/vhost/models"
	"github.com/alex-oleshkevich/vhost/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TemplateVars struct
type TemplateVars struct {
	IP           string
	Port         string
	ProjectPath  string
	Domain       string
	DomainSuffix string
}

var (
	httpAddress string
	httpPort    string
)

var templateData TemplateVars
var dbCreate bool
var dbName string
var dbType string

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init <name> <directory>",
	Short: "Initialize a new project",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// defer func() {
		// 	if r := recover(); r != nil {
		// 		if templateData.ProjectPath != "" && utils.DirectoryExists(templateData.ProjectPath) {
		// 			os.RemoveAll(templateData.ProjectPath)
		// 		}
		// 	}
		// }()

		lockfile := models.Lock{}

		if len(args) == 0 {
			log.Fatalf("Project name was not specified")
		}
		projectName := args[0]

		targetDirectory, err := os.Getwd()
		if err != nil {
			log.Panicln(err)
		}
		if len(args) >= 2 {
			targetDirectory = args[1]
		}

		log.Infof("Initializing project %s in %s", projectName, targetDirectory)
		if utils.DirectoryExists(targetDirectory) {
			log.Warnf("Directory %s already exists.", targetDirectory)
			lfname := path.Join(targetDirectory, models.LockFilename)
			if utils.FileExists(lfname) {
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Seems there is another project in the target directory. Overwrite? y,n ")
				text, _ := reader.ReadString('\n')
				if strings.TrimSpace(text) != "y" {
					log.Infoln("Cancel.")
					return
				}
			}
		}

		if dbName != "" {
			if !isDbSupported(dbType) {
				log.Fatalf("Database driver not supported: %s", dbType)
			}

			log.Infof("Create %s database %s", dbType, dbName)
			switch dbType {
			case "mysql":
				err = createMySQLDatabase(dbName)
			case "postgres":
				err = createPostgresDatabase(dbName)
			}

			if err != nil {
				log.Fatalln(err)
			}
		}

		err = os.MkdirAll(targetDirectory, 0755)
		if err != nil {
			log.Fatalln(err)
		}
		os.Chdir(targetDirectory)

		for _, dir := range []string{"etc", "log", "tmp", "www"} {
			if utils.DirectoryExists(dir) {
				log.Warnf("Directory %s exists.", dir)
			} else {
				os.Mkdir(dir, 0755)
			}
		}

		nginxTemplatePath, err := utils.FindFileInApp("templates/vhost.tpl")
		if err != nil {
			log.Errorln("Could not find nginx template in application directories")
			log.Fatalln(err)
		}

		log.Debugf("Using template file: %s", nginxTemplatePath)
		t := template.New("vhost.tpl")
		t, err = t.ParseFiles(nginxTemplatePath)
		if err != nil {
			log.Fatalln(err)
		}
		writer, err := os.Create("etc/vhost.conf")
		if err != nil {
			log.Fatalln(err)
		}

		projectPath, err := filepath.Abs(targetDirectory)
		if err != nil {
			log.Panicln(err)
		}

		templateData.ProjectPath = projectPath
		templateData.Domain = projectName
		templateData.DomainSuffix = viper.GetString("domain")

		err = t.Execute(writer, templateData)

		if err != nil {
			log.Panicln(err)
		}

		realConfDest := fmt.Sprintf("%s/%s.conf", viper.GetString("sites-enabled"), projectName)
		realConfDest, err = filepath.Abs(realConfDest)
		if err != nil {
			log.Fatalln(err)
		}
		if utils.FileExists(realConfDest) {
			log.Warnf("Vhost is already enabled in %s.", realConfDest)
			os.Remove(realConfDest)
		}

		realConfSrc, err := filepath.Abs("etc/vhost.conf")
		if err != nil {
			log.Panicln(err)
		}

		log.Debugf("Link %s -> %s", realConfSrc, realConfDest)
		err = os.Symlink(realConfSrc, realConfDest)
		if err != nil {
			log.Fatalln(err)
		}

		shellCmd := exec.Command("/usr/bin/systemctl", "reload", "nginx")
		err = shellCmd.Start()
		if err != nil {
			log.Panicln(err)
		}
		log.Debugln(shellCmd.CombinedOutput())

		lockfile.DbName = dbName
		lockfile.DbType = dbType
		lockfile.VhostLinkPath = realConfDest

		err = lockfile.Write()
		if err != nil {
			log.Fatalln(err)
		}
	},
}

func isDbSupported(dbtype string) bool {
	return utils.InArray(dbType, []string{"mysql", "postgres"})
}

func createMySQLDatabase(dbname string) error {
	db, err := utils.GetMySQLDb()
	if err != nil {
		return err
	}

	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbname)
	_, err = db.Exec(sql)
	return err
}

func createPostgresDatabase(dbname string) error {
	db, err := utils.GetPostgresDb()
	if err != nil {
		return err
	}

	sql := fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname = '%s'", dbname)
	rows, err := db.Query(sql)
	if err != nil {
		return err
	}
	if rows.Next() {
		log.Warningln("Postgres database already exists.")
		return nil
	}

	_, err = db.Query(fmt.Sprintf("CREATE DATABASE %s", dbname))
	if err != nil {
		return err
	}
	return nil
}

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&templateData.IP, "ip", "127.0.0.1", "IP address to bind to")
	initCmd.Flags().StringVar(&templateData.Port, "port", "80", "IP address to bind to")
	initCmd.Flags().StringVar(&dbType, "dbtype", "mysql", "Database type: mysql, postgres")
	initCmd.Flags().StringVar(&dbName, "dbname", "", "Database name")
}
