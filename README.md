wall
-----

A library for populating data structures based on parsing an input string with a regular expression.
Each subcapture group is returned as its own field.

## Usage

The only exported function is `Parse(re *regexp.Regexp, input string, output interface{}) error`.
Returns an error if the input does not match the given regex. Also returns an error if the output
is not an acceptable type. Accepted types are `[]string` (matches returned in order, with the matching
input removed), `map[string]string`, and a tagged struct. Tagged struct fields must be of type `string`.
If subcatpture groups are not labeled, their labels default to the string representation of their
indices (i.e. `wall:"1"`, `wall:"3"`, etc.). Subcaptures with no matching struct fields will be
dropped without an error.

```
package main

import (
	"fmt"
	"regexp"

	"github.com/packrat386/wall"
)

type SemanticVersion struct {
	Major string `wall:"major"`
	Minor string `wall:"minor"`
	Patch string `wall:"patch"`
}

func main() {
	semverRegex := regexp.MustCompile(`^(?P<major>\d+)\.(?P<minor>\d+)\.(?P<patch>\d+)$`)
	fmt.Printf("Input a semantic version: ")

	var v string
	fmt.Scan(&v)

	var version SemanticVersion
	err := wall.Parse(semverRegex, v, &version)
	if err != nil {
		fmt.Printf("semantic verison is not valid\n")
		return
	}

	fmt.Printf("Major: %s\nMinor: %s\nPatch: %s\n", version.Major, version.Minor, version.Patch)
}
```

## Issues

This repo was mostly a way for me to explore using reflection and struct tags. I'll happily fix bugs
or look at pull requests, but I'm unlikely to be adding major functionality to this library.
