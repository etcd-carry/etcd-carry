package main

import (
	"fmt"
	"orcastack.io/etcd-mirror/cmd/app"
	"os"
)

func main() {
	cmd := app.NewEtcdMirrorCommand()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
