package osutil

/*
#include <stdio.h>
#include <stdlib.h>
*/
import "C"

import (
	"errors"
	"bytes"
	"io"
	"os"
	"unsafe"
)

// Capture captures stderr and stdout of a given function call.
func Capture(call func()) ([]byte, error) {
	originalStdErr, originalStdOut := os.Stderr, os.Stdout
	defer func() {
		os.Stderr, os.Stdout = originalStdErr, originalStdOut
	}()

	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	defer r.Close()

	os.Stderr, os.Stdout = w, w

	out := make(chan []byte)
	go func() {
		var b bytes.Buffer

		_, err := io.Copy(&b, r)
		if err != nil {
			panic(err)
		}

		out <- b.Bytes()
	}()

	call()

	err = w.Close()
	if err != nil {
		return nil, err
	}

	return <-out, err
}

// CaptureWithCGo captures stderr and stdout as well as stderr and stdout of C of a given function call.
func CaptureWithCGo(call func()) ([]byte, error) {
	return captureWithCGoImpl(true, call)
}

// CaptureStdoutWithCGo captures stdout as well as stdout of C of a given function call. stderr is not redirected or captured
func CaptureStdoutWithCGo(call func()) ([]byte, error) {
	return captureWithCGoImpl(false, call)
}
// TBD, a cleaner impl would be to use two Pipes, handle stderr and stdout indepedently and return separate []byte buffers
func captureWithCGoImpl(include_stderr bool, call func()) ([]byte, error) {
	originalStdErr, originalStdOut := os.Stderr, os.Stdout
	originalCStdErr, originalCStdOut := C.stderr, C.stdout
	defer func() {
		os.Stderr, os.Stdout = originalStdErr, originalStdOut
		C.stderr, C.stdout = originalCStdErr, originalCStdOut
	}()

	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	defer r.Close()

	cw := C.CString("w")
	defer C.free(unsafe.Pointer(cw))

	f := C.fdopen((C.int)(w.Fd()), cw)
	if (f == nil) {
		return nil, errors.New("fdopen returned nil")
	}
	defer C.fclose(f)

	os.Stdout = w
	C.stdout = f
	if include_stderr {
		os.Stderr = w
		C.stderr = f
	}

	out := make(chan []byte)
	go func() {
		var b bytes.Buffer

		_, err := io.Copy(&b, r)
		if err != nil {
			panic(err)
		}

		out <- b.Bytes()
	}()

	call()

	C.fflush(f)

	err = w.Close()
	if err != nil {
		return nil, err
	}

	return <-out, err
}
