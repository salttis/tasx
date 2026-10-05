package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/salttis/tasx/internal/cli"
)

func main() {
	name := strings.TrimSuffix(filepath.Base(os.Args[0]), filepath.Ext(os.Args[0]))
	if err := cli.Execute(strings.EqualFold(name, "tx")); err != nil {
		_, _ = os.Stderr.WriteString("tasx: " + err.Error() + "\n")
		os.Exit(1)
	}
}
