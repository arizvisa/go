// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package cgo contains runtime support for code generated
by the cgo tool.  See the documentation for the cgo command
for details on using cgo.
*/
package cgo

/*

#cgo darwin,!arm,!arm64 LDFLAGS: -lpthread
#cgo darwin,arm LDFLAGS: -framework CoreFoundation
#cgo darwin,arm64 LDFLAGS: -framework CoreFoundation
#cgo dragonfly LDFLAGS: -lpthread
#cgo freebsd LDFLAGS: -lpthread
#cgo android LDFLAGS: -llog
#cgo !android,linux LDFLAGS: -lpthread
#cgo netbsd LDFLAGS: -lpthread
#cgo openbsd LDFLAGS: -lpthread

#cgo windows,386 CFLAGS: -target i686-w64-mingw32
#cgo windows,386 LDFLAGS: -target i686-w64-mingw32
#cgo windows,amd64 CFLAGS: -target x86_64-w64-mingw64
#cgo windows,amd64 LDFLAGS: -target x86_64-w64-mingw64

#cgo CFLAGS: -Wall -Werror

#cgo solaris CPPFLAGS: -D_POSIX_PTHREAD_SEMANTICS

*/
import "C"
