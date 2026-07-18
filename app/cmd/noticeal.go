package main

import "fmt"

// Injected at build time via .goreleaser.yaml ldflags.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	fmt.Printf("noticeal %s (%s) built on %s\n", version, commit, date)
}
