// +build ignore

package main

import (
	"log"

	"github.com/shurcooL/vfsgen"
	"github.com/yobro/introvert/web"
)

func main() {
	err := vfsgen.Generate(web.Assets, vfsgen.Options{
		PackageName:  "web",
		BuildTags:    "builtinassets",
		VariableName: "Assets",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
