// Package version holds the application version, injected at build time by
// GoReleaser via -ldflags.
package version

// Version defaults to "dev" for local builds and is overwritten at release
// time via -X github.com/mzeahmed/noticeal/internal/version.Version=....
var Version = "dev"
