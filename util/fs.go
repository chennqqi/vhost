package util

import (
    "os"
    "path/filepath"
)

var paths = []string {
    "./",
    os.Getenv("HOME") + "/.vhost/",
    "/etc/vhost/",
}

// FindFile finds filename in paths
func FindFile(filename string) (pathname string, found bool) {
    for _, path := range paths {
        var pathname = path + filename
        if _, err := os.Stat(pathname); err == nil {
            abspath, err := filepath.Abs(pathname)
            if err != nil {
                panic(err)
            } 
            return abspath, true
        }
    }
    return "", false
}