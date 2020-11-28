package main

import (
	"github.com/markbates/pkger"
	"github.com/zchase/stevie/cmd"
)

func main() {
	pkger.Include("/pkg/application/file_templates")
	pkger.Include("/lib")
	cmd.Execute()
}
