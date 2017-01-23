package osutil

/*
#include <stdio.h>

void printSomething() {
	printf("C\n");
}
void printSomethingStderr() {
	fprintf(stderr, "E\n");
}
*/
import "C"

import (
	"fmt"
	"testing"
	"os"
	"github.com/stretchr/testify/assert"
)

func printAll() {
	fmt.Println("Go")
	fmt.Fprintln(os.Stderr, "e")
	C.printSomething()
	C.printSomethingStderr()
}

func testCapture(t *testing.T) {
	out, err := Capture(func() {
		printAll()
	})
	assert.Nil(t, err)

	assert.Contains(t, string(out), "Go")
	assert.Contains(t, string(out), "e")
	assert.NotContains(t, string(out), "C")
	assert.NotContains(t, string(out), "E")
}

func testCaptureStdout(t *testing.T) {
//	out, err := Capture(func() { // trigger test failure
	out, err := CaptureStdout(func() {
		printAll()
	})
	assert.Nil(t, err)

	assert.Contains(t, string(out), "Go")
	assert.NotContains(t, string(out), "e")
	assert.NotContains(t, string(out), "C")
	assert.NotContains(t, string(out), "E")
}

func testCaptureStderr(t *testing.T) {
//	out, err := Capture(func() { // trigger test failure
	out, err := CaptureStderr(func() {
		printAll()
	})
	assert.Nil(t, err)

	assert.NotContains(t, string(out), "Go")
	assert.Contains(t, string(out), "e")
	assert.NotContains(t, string(out), "C")
	assert.NotContains(t, string(out), "E")
}

func testCaptureWithCGo(t *testing.T) {
	out, err := CaptureWithCGo(func() {
		printAll()
	})
	assert.Nil(t, err)

	assert.Contains(t, string(out), "Go")
	assert.Contains(t, string(out), "e")
	assert.Contains(t, string(out), "C")
	assert.Contains(t, string(out), "E")
}

func testCaptureStdoutWithCGo(t *testing.T) {
//	out, err := CaptureWithCGo(func() {  // trigger test failure
	out, err := CaptureStdoutWithCGo(func() {
		printAll()
	})
	assert.Nil(t, err)

	assert.Contains(t, string(out), "Go")
	assert.NotContains(t, string(out), "e")
	assert.Contains(t, string(out), 	"C")
	assert.NotContains(t, string(out), "E")
}
