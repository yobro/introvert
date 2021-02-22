// +build !builtinassets

package web

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Assets is a filesystem containing the web app assets
var Assets = func() http.FileSystem {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	switch filepath.Base(wd) {
	case "introvert":
		return http.Dir("web/static")
	case "web":
		return http.Dir("static")
	}

	panic(fmt.Sprintf("unable to create asset filesystem: unknown origin %s", wd))
}()

func init() {
	log.Println("debug: using local files in web/static directory")
}
