package io

import (
	"fmt"
	"github.com/RenatoGeh/gospn/sys"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// DownloadFromURL takes an URL u and a destination path p, downloading the contents of u to p.  If
// p is not a complete path (not a directory, contains extension), then the name of the file to be
// downloaded is copied as the new file's name. If override is set to true and p points to a file,
// this function overrides file p with the new download. Having said that, take extreme care when
// using override!
func DownloadFromURL(u, p string, override bool) error {
	var f string
	var isDir bool
	p = filepath.Clean(GetPath(p))
	if isDir = (filepath.Ext(p) == ""); isDir {
		t := strings.Split(u, "/")
		f = p + "/" + t[len(t)-1]
	} else {
		f = p
	}
	_, e := os.Stat(f)
	// File exists.
	if e == nil && !override {
		fmt.Println("File already exists and you chose not to override. Stopping download.")
		return e
		// File does not exist.
	} else if e != nil && !os.IsNotExist(e) {
		fmt.Printf("Could not stat file [%s].\n", p)
		fmt.Println(e)
		return e
	}
	out, e := os.Create(f)
	defer out.Close()
	if e != nil {
		fmt.Println("Error when trying to create file [%s].\n", f)
		return e
	}
	sys.Printf("Downloading from [%s] to {./%s}.\n", u, f)
	d, e := http.Get(u)
	defer d.Body.Close()
	if e != nil {
		fmt.Println("Error while downloading [%s].\nStopping download.", u)
		return e
	}
	_, e = io.Copy(out, d.Body)
	if e != nil {
		fmt.Println("Error while copying download to local directory.")
		return e
	}
	sys.Printf("Download finished.")
	return nil
}
