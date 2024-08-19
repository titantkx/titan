package testutil

import (
	"fmt"
	"io"
	"runtime"
)

type TestingT interface {
	Errorf(format string, args ...any)
	FailNow()
	Cleanup(f func())
}

type MockTest struct {
	w             io.Writer
	cleanupFunc   func()
	cleanupCalled bool
}

func NewMockTest(w io.Writer) *MockTest {
	return &MockTest{w: w}
}

func (t *MockTest) Errorf(format string, args ...any) {
	_, file, line, _ := runtime.Caller(1)
	fmt.Fprintf(t.w, "    %s:%d: ", file, line)
	fmt.Fprintf(t.w, format, args...)
	fmt.Fprintln(t.w)
}

func (t *MockTest) FailNow() {
	t.doCleanup()
	runtime.Goexit()
}

func (t *MockTest) Cleanup(f func()) {
	t.cleanupFunc = f
}

func (t *MockTest) Finish() {
	t.doCleanup()
}

func (t *MockTest) doCleanup() {
	if !t.cleanupCalled && t.cleanupFunc != nil {
		t.cleanupFunc()
		t.cleanupCalled = true
	}
}
