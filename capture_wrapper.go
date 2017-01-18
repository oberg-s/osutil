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

	"github.com/stretchr/testify/assert"
)

func testCapture(t *testing.T) {
	out, err := Capture(func() {
		fmt.Println("Go")
		C.printSomething()
	})
	assert.Nil(t, err)

	assert.Contains(t, string(out), "Go")
	assert.NotContains(t, string(out), "C")
}

func testCaptureWithCGo(t *testing.T) {
	out, err := CaptureWithCGo(func() {
		fmt.Println("Go")
		C.printSomething()
	})
	assert.Nil(t, err)

	assert.Contains(t, string(out), "Go")
	assert.Contains(t, string(out), "C")
}

func testCaptureStdoutWithCGo(t *testing.T) {
	out, err := CaptureStdoutWithCGo(func() {
//	out, err := CaptureWithCGo(func() {  // trigger test failure
		fmt.Println("Go")
		// if capturing both stderr and stdout the extra bring should would break this test
		C.printSomethingStderr()
		C.printSomething()
	})
	assert.Nil(t, err)

	assert.Contains(t, string(out), "Go")
	assert.Contains(t, string(out), "C")
	assert.NotContains(t, string(out), "E")
}
