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

type mockTest struct {
	w             io.Writer
	cleanupFunc   func()
	cleanupCalled bool
}

func NewMockTest(w io.Writer) *mockTest {
	return &mockTest{w: w}
}

func (t *mockTest) Errorf(format string, args ...any) {
	_, file, line, _ := runtime.Caller(1)
	fmt.Fprintf(t.w, "    %s:%d: ", file, line)
	fmt.Fprintf(t.w, format, args...)
	fmt.Fprintln(t.w)
}

func (t *mockTest) FailNow() {
	t.doCleanup()
	runtime.Goexit()
}

func (t *mockTest) Cleanup(f func()) {
	t.cleanupFunc = f
}

func (t *mockTest) Finish() {
	t.doCleanup()
}

func (t *mockTest) doCleanup() {
	if !t.cleanupCalled && t.cleanupFunc != nil {
		t.cleanupFunc()
		t.cleanupCalled = true
	}
}
