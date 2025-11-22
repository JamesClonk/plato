package main

import "github.com/JamesClonk/plato/cmd"

var (
	version    = "0.0.0-dev.0"
	buildstamp = "now"
	githash    = ""
)

func main() {
	cmd.Execute(version, buildstamp, githash)
}
