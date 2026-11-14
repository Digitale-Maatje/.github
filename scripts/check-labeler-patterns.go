// Compiles every pattern in .github/labeler.yml with Go's regexp, the RE2
// engine srvaroa/labeler matches with. The labeler drops a pattern it cannot
// compile without a word (it ignores regexp.Match's error), so a lookahead
// such as ^(?!examples/) silently never matches. This turns that into a
// failed check.
//
// Reads labeler.yml as JSON on stdin (yq -o=json converts it):
//
//	yq -o=json .github/labeler.yml | go run scripts/check-labeler-patterns.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

// The rule conditions srvaroa/labeler (v1.14) matches as regular expressions.
var stringConditions = []string{"title", "branch", "base-branch", "body"}

type config struct {
	Version int              `json:"version"`
	Labels  []map[string]any `json:"labels"`
}

func main() {
	var cfg config
	if err := json.NewDecoder(os.Stdin).Decode(&cfg); err != nil {
		fmt.Fprintln(os.Stderr, "labeler.yml: cannot read:", err)
		os.Exit(1)
	}
	if len(cfg.Labels) == 0 {
		fmt.Fprintln(os.Stderr, "labeler.yml: no labels: the input is not a labeler configuration")
		os.Exit(1)
	}

	failed, checked := 0, 0
	check := func(i int, rule map[string]any, key, pattern string) {
		checked++
		if _, err := regexp.Compile(pattern); err != nil {
			failed++
			fmt.Fprintf(os.Stderr, "labels[%d] (%v) %s: %v\n", i, rule["label"], key, err)
		}
	}
	for i, rule := range cfg.Labels {
		for _, key := range stringConditions {
			if v, ok := rule[key]; ok {
				s, isString := v.(string)
				if !isString {
					failed++
					fmt.Fprintf(os.Stderr, "labels[%d] (%v) %s: not a string\n", i, rule["label"], key)
					continue
				}
				check(i, rule, key, s)
			}
		}
		if v, ok := rule["files"]; ok {
			files, isList := v.([]any)
			if !isList {
				failed++
				fmt.Fprintf(os.Stderr, "labels[%d] (%v) files: not a list\n", i, rule["label"])
				continue
			}
			for _, f := range files {
				s, isString := f.(string)
				if !isString {
					failed++
					fmt.Fprintf(os.Stderr, "labels[%d] (%v) files: %v is not a string\n", i, rule["label"], f)
					continue
				}
				check(i, rule, "files", s)
			}
		}
	}

	if failed > 0 {
		fmt.Fprintf(os.Stderr, "%d of %d patterns fail to compile as RE2\n", failed, checked)
		os.Exit(1)
	}
	fmt.Printf("%d patterns in %d rules compile as RE2\n", checked, len(cfg.Labels))
}
