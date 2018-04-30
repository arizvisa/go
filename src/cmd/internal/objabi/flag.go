// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package objabi

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func Flagcount(name, usage string, val *int) {
	flag.Var((*count)(val), name, usage)
}

func Flagfn1(name, usage string, f func(string)) {
	flag.Var(fn1(f), name, usage)
}

func Flagprint(fd int) {
	if fd == 1 {
		flag.CommandLine.SetOutput(os.Stdout)
	}
	flag.PrintDefaults()
}

// This is a near exact copy of gcc's libiberty/argv.c buildargv
func buildargv(data []byte) []string {
	var result []string

	// current index into data
	var di int

	// string parsing states
	var squote, dquote, bsquote bool

	// simple check if a char is whitespace
	ISSPACE := func(ch byte) bool {
		return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\v' || ch == '\f' || ch == '\r'
	}

	// simple consume_whitespace implementation
	consume_whitespace := func(idx int, data []byte) int {
		for i := idx; i < len(data); i++ {
			if !ISSPACE(data[i]) {
				return i
			}
		}
		return len(data)
	}

	// consume initial whitespace
	di = consume_whitespace(0, data)
	if di >= len(data) {
		return result
	}

	// argument loop
	for di < len(data) {
		var arg string

		// scan each individual argument
		arg = ""
		for di < len(data) {
			if ISSPACE(data[di]) && !squote && !dquote && !bsquote {
				break
			}

			// backslash
			if bsquote {
				bsquote = false
				arg += string(data[di])
			} else if data[di] == '\\' {
				bsquote = true

			// single-quote
			} else if squote {
				if data[di] == '\'' {
					squote = false
				} else {
					arg += string(data[di])
				}

			// double-quote
			} else if dquote {
				if data[di] == '"' {
					dquote = false
				} else {
					arg += string(data[di])
				}

			// state entries
			} else {
				if data[di] == '\'' {
					squote = true
				} else if data[di] == '"' {
					dquote = true
				} else {
					arg += string(data[di])
				}
			}

			// process the next byte
			di += 1
		}

		// add the current arg to the results
		result = append(result, arg)
		di = consume_whitespace(di, data)
	}

	// ...and we're done!
	return result
}

func Flagparse(usage func()) {
	flag.Usage = usage

	// Expand any response files that were specified at the commandline. Anything
	// that is not a response file (file not found, zero-length arg, etc) gets
	// blindly added to the arg as we assume that the user knows what they're doing.

	// FIXME: I think response files are recursive, so if that's true then this
	// 		  code should be refactored to support recursive response files. Probably
	//		  with a channel or something.
	var args []string
	for _, arg := range os.Args {
		// Check that response file prefix doesn't exist or that arg is zero-length
		if !strings.HasPrefix(arg, "@") || len(arg) < 1 {
			args = append(args, arg)
			continue
		}

		// Check to see if the @-prefixed file is non-existent
		file, err := os.Open(arg[1:])
		if os.IsNotExist(err) {
			log.Printf("Unable to open response file (%s): %#v\n", arg[1:], err)
			args = append(args, arg)
			continue
		}

		// Okay, so now we have a file with args. So expand the file contents
		contents, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatalf("Unable to read contents of response file (%s): %#v", arg[1:], err)
		}

		// Now we can add each arg from the file
		for _, row := range buildargv(contents) {
			args = append(args, row)
		}

		// We're done. Close it and move on
		if err := file.Close(); err != nil {
			log.Fatalf("Unable to close response file (%s): %#v", arg[1:], err)
		}
	}

	os.Args = args
	flag.Parse()
}

func AddVersionFlag() {
	flag.Var(versionFlag{}, "V", "print version and exit")
}

var buildID string // filled in by linker

type versionFlag struct{}

func (versionFlag) IsBoolFlag() bool { return true }
func (versionFlag) Get() interface{} { return nil }
func (versionFlag) String() string   { return "" }
func (versionFlag) Set(s string) error {
	name := os.Args[0]
	name = name[strings.LastIndex(name, `/`)+1:]
	name = name[strings.LastIndex(name, `\`)+1:]
	name = strings.TrimSuffix(name, ".exe")
	p := Expstring()
	if p == DefaultExpstring() {
		p = ""
	}
	sep := ""
	if p != "" {
		sep = " "
	}

	// The go command invokes -V=full to get a unique identifier
	// for this tool. It is assumed that the release version is sufficient
	// for releases, but during development we include the full
	// build ID of the binary, so that if the compiler is changed and
	// rebuilt, we notice and rebuild all packages.
	if s == "full" && strings.HasPrefix(Version, "devel") {
		p += " buildID=" + buildID
	}
	fmt.Printf("%s version %s%s%s\n", name, Version, sep, p)
	os.Exit(0)
	return nil
}

// count is a flag.Value that is like a flag.Bool and a flag.Int.
// If used as -name, it increments the count, but -name=x sets the count.
// Used for verbose flag -v.
type count int

func (c *count) String() string {
	return fmt.Sprint(int(*c))
}

func (c *count) Set(s string) error {
	switch s {
	case "true":
		*c++
	case "false":
		*c = 0
	default:
		n, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("invalid count %q", s)
		}
		*c = count(n)
	}
	return nil
}

func (c *count) Get() interface{} {
	return int(*c)
}

func (c *count) IsBoolFlag() bool {
	return true
}

func (c *count) IsCountFlag() bool {
	return true
}

type fn0 func()

func (f fn0) Set(s string) error {
	f()
	return nil
}

func (f fn0) Get() interface{} { return nil }

func (f fn0) String() string { return "" }

func (f fn0) IsBoolFlag() bool {
	return true
}

type fn1 func(string)

func (f fn1) Set(s string) error {
	f(s)
	return nil
}

func (f fn1) String() string { return "" }
