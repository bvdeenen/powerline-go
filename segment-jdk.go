package main

import (
	"os"
	"path/filepath"
	"strings"

	pwl "github.com/bvdeenen/powerline-go/powerline"
)

// segmentJdk shows the JDK selected by JAVA_HOME. When JAVA_HOME is a
// symlink (e.g. /usr/lib/jvm/default-jdk -> openjdk11) the link target is
// shown, cut at the first '-'; otherwise the directory's base name is shown.
// Nothing is shown when JAVA_HOME is unset.
func segmentJdk(p *powerline) []pwl.Segment {
	javaHome := os.Getenv("JAVA_HOME")
	if javaHome == "" {
		return []pwl.Segment{}
	}
	name := filepath.Base(javaHome)
	if target, err := os.Readlink(javaHome); err == nil {
		name = strings.SplitN(target, "-", 2)[0]
	}

	return []pwl.Segment{{
		Name:       "jdk",
		Content:    "𝄢" + name,
		Foreground: p.theme.JdkFg,
		Background: p.theme.JdkBg,
	}}
}
