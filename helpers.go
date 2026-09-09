package main

import (
	"runtime/debug"
	"strings"
)

// dep is one entry of the build's module graph.
type dep struct {
	Name    string
	Version string
}

// moduleDeps returns the modules linked into this binary, so --version can
// report the exact bubbletea build in use.
func moduleDeps() []dep {
	var out []dep
	if bi, ok := debug.ReadBuildInfo(); ok {
		for _, m := range bi.Deps {
			out = append(out, dep{Name: m.Path, Version: m.Version})
		}
	}
	return out
}

// wrap hard-wraps text to width for plain-text (non-TTY) output.
func wrap(s string, width int) string {
	if width < 20 {
		width = 20
	}
	var out strings.Builder
	for i, para := range strings.Split(s, "\n") {
		if i > 0 {
			out.WriteString("\n")
		}
		lineLen := 0
		for j, word := range strings.Fields(para) {
			switch {
			case j == 0:
				lineLen = len([]rune(word))
				out.WriteString(word)
			case lineLen+1+len([]rune(word)) > width:
				out.WriteString("\n")
				out.WriteString(word)
				lineLen = len([]rune(word))
			default:
				out.WriteString(" ")
				out.WriteString(word)
				lineLen += 1 + len([]rune(word))
			}
		}
	}
	return out.String()
}
