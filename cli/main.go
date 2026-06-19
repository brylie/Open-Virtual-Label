package main

import (
	"os"

	"github.com/open-virtual-label/ovl/cmd"
)

func main() {
	os.Exit(cmd.Execute())
}
