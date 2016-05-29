package adapt

import (
	"bytes"
)

// HTTPServer interface
type HTTPServer interface {
    // GetDomain returns a domain name of new vhost (eg. example.com)
	GetDomain(name string) string
    
    // GetVhostConfig retruns a virtual host config string
    GetVhostConfig(name string, ip string, port int, docRoot string, ssl bool) bytes.Buffer
    
    // GetSSLPart returns a SSL part of virtual host config
    GetSSLPart() string
    
    // GetVhostPath returns a part to virtual host config file
    GetVhostPath(name string) string
    
    // IsExits tests if a virtual host already exists
    IsExists(name string) bool
    
    // GetReloadCommand func
    GetReloadCommand() string;
    
    // Enable func
    Enable(name string)
}