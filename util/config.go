package util

import (
    "io/ioutil"
    "gopkg.in/yaml.v2"
    log "github.com/Sirupsen/logrus"
)

// Apache struct
type Apache struct {
    DirSitesEnabled string `yaml:"dir_sites_enabled"`
    DirSitesAvailable string `yaml:"dir_sites_available"`
    VhostFileSuffix string `yaml:"vhost_file_suffix"`
    RestartCommand string `yaml:"restart_command"`
    Domain string `yaml:"domain"`
}

// General struct
type General struct {
    Domain string
    SitesDir string `yaml:"sites_dir"`
    HostsFile string `yaml:"hosts_file"`
    User string
    Group string
}

// Mysql struct
type Mysql struct {
    User string
    Password string
    Host string
    Charset string
}

// Ssl struct
type Ssl struct {
    CertFile string `yaml:"cert_file"`
    KeyFile string `yaml:"key_file"`
}

// Config structure
type Config struct {
    Apache Apache
    General General
    Mysql Mysql
    Ssl Ssl
}

// LoadConfig config file into object
func LoadConfig() Config {
    configFile, err := FindFile("config.yaml")
    if err != nil {
        panic("Config file was not found")
    }
    
    log.Infoln("Using config from: " + configFile)
    
    data, err := ioutil.ReadFile(configFile)
    if err != nil {
        panic(err)
    }
    
    config := Config{}
    err = yaml.Unmarshal([]byte(data), &config)
    if err != nil {
        panic(err)
    }
    
    return config
}