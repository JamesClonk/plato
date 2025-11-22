package main

import "github.com/JamesClonk/plato/cmd"

var (
	version = "0.0.0"
	commit  = "-"
	date    = "now"
)

func main() {
	cmd.Execute(version, commit, date)
}
