package util

import (
	"bufio"
	"net"
	"os"
	"strings"
)

// CommentChar constant
const CommentChar string = "#"

// HostLine represents a line from /etc/hosts file
type HostLine struct {
	IP    string
	Hosts []string
	Raw   string
}

// IsComment tests if a line of /etc/hosts if a comment
func (hl *HostLine) IsComment() bool {
	trimmed := strings.TrimSpace(hl.Raw)
	return strings.HasPrefix(trimmed, CommentChar)
}

// MakeLine returns a line prepared for writing
func (hl *HostLine) MakeLine() string {
    if hl.IsComment() {
        return hl.Raw + "\n"
    } 
    
    return hl.IP + "         " + strings.Join(hl.Hosts, " ") + "\n"
}

// Hosts represents a /etc/hosts file
type Hosts struct {
	Path  string
	Lines []HostLine
}

// IsWritable tests if /etc/hosts is writable by current user
func (h *Hosts) IsWritable() bool {
	_, err := os.OpenFile(h.Path, os.O_WRONLY, 0660)
	if err != nil {
		return false
	}

	return true
}

// Add a new item to /etc/hosts
func (h *Hosts) Add(ip string, domain string) {
    line := HostLine{
        IP: ip,
        Hosts: []string{domain},
        Raw: ip + " " + domain,
    }
    h.Lines = append(h.Lines, line)
}

// Has tests if domain is in the hosts
func (h *Hosts) Has(domain string) bool {
    for _, line := range h.Lines {
        if strings.Contains(line.Raw, domain) {
            return true
        }
    }
    return false
}

func (h *Hosts) Write() {
    file, err := os.Create(h.Path)
    if err != nil {
        panic(err)
    }
    
    handle := bufio.NewWriter(file)
    defer handle.Flush()
    for _, line := range h.Lines {
        handle.Write([]byte(line.MakeLine()))
    }
}

// LoadHosts loads /etc/hosts file
func LoadHosts(filename string) Hosts {
	hosts := Hosts{
		Path: filename,
	}

	var lines []HostLine

	file, err := os.Open(hosts.Path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := makeHostLine(scanner.Text())
		lines = append(lines, line)
	}
	hosts.Lines = lines
	return hosts
}

func makeHostLine(rawLine string) HostLine {
	line := HostLine{
		Raw: rawLine,
	}

	parts := strings.Fields(rawLine)
	if len(parts) == 0 {
		return line
	}

	if !line.IsComment() {
		ip := parts[0]
		if net.ParseIP(ip) == nil {
			panic("Invalid hosts line: " + rawLine)
		}

		line.IP = ip
		line.Hosts = parts[1:]
	}
	return line
}
