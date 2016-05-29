package adapt

import (
	"bytes"
	"io/ioutil"
	"text/template"
    "path"
	"os"

	"github.com/alex-oleshkevich/vhost/util"
    log "github.com/Sirupsen/logrus"
)

var config = util.LoadConfig()

// ApacheTemplateData struct
type ApacheTemplateData struct {
	IP       string
	Port     int
	Domain   string
	SitesDir string
	DocRoot  string
	SSL      string
}

// ApacheSSLData struct
type ApacheSSLData struct {
	CertFile string
	KeyFile  string
}

// ApacheServer struct
type ApacheServer struct {
}

// GetVhostConfig func
func (s ApacheServer) GetVhostConfig(name string, ip string, port int, docRoot string, ssl bool) bytes.Buffer {
	apacheVhostTemplate, err := util.FindFile("share/apache.conf")
	if err != nil {
		panic("Apache virtual host template was not found.")
	}

	sshPart := ""
	if ssl {
		sshPart = s.GetSSLPart()
	}

	data := ApacheTemplateData{
		ip,
		port,
		s.GetDomain(name),
		config.General.SitesDir,
		docRoot,
		sshPart,
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

// GetDomain func
func (s ApacheServer) GetDomain(name string) string {
	return name + config.Apache.Domain
}

// GetSSLPart func
func (s ApacheServer) GetSSLPart() string {
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

// GetVhostPath func
func (s ApacheServer) GetVhostPath(name string) string {
	return config.Apache.DirSitesAvailable + "/" + name + config.Apache.VhostFileSuffix
}

// IsExists func
func (s ApacheServer) IsExists(name string) bool {
	var _, err = os.Stat(s.GetVhostPath(name))
	return err == nil
}

// GetReloadCommand func
func (s ApacheServer) GetReloadCommand() string {
	return config.Apache.RestartCommand
}

// Enable func
func (s ApacheServer) Enable(name string) {
    source := s.GetVhostPath(name)
    target := config.Apache.DirSitesEnabled + "/" + path.Base(source)
    
    log.Info(source + " -> " + target)
    err := os.Symlink(source, target)
    if err != nil {
        log.Error(err)
    }
}
