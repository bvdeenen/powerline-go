package main

import (
	"os"

	pwl "github.com/bvdeenen/powerline-go/powerline"
)

// segmentDirenv shows the direnv state: the contents of DIRENV_MESG when set,
// otherwise a "Δ" marker when a direnv environment is loaded (DIRENV_DIR set).
func segmentDirenv(p *powerline) []pwl.Segment {
	var content string
	if mesg := os.Getenv("DIRENV_MESG"); mesg != "" {
		content = "⊂" + mesg + "⊃"
	} else if os.Getenv("DIRENV_DIR") != "" {
		content = "Δ"
	} else {
		return []pwl.Segment{}
	}

	return []pwl.Segment{{
		Name:       "direnv",
		Content:    content,
		Foreground: p.theme.DirenvFg,
		Background: p.theme.DirenvBg,
	}}
}
