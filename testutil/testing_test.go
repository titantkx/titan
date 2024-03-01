package testutil_test

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tokenize-titan/titan/testutil"
)

func TestMockTestErrorf(t *testing.T) {
	var output bytes.Buffer

	mt := testutil.NewMockTest(&output)
	defer mt.Finish()

	_, file, line, _ := runtime.Caller(0)
	mt.Errorf("Hello %s", "World")

	expectedOutput := fmt.Sprintf("    %s:%d: Hello World\n", file, line+1)

	require.Equal(t, expectedOutput, output.String())
}

func TestMockTestCleanupIsCalledAfterFinished(t *testing.T) {
	cleanupCalled := 0
	defer func() {
		require.Equal(t, 1, cleanupCalled)
	}()

	mt := testutil.NewMockTest(os.Stderr)
	defer mt.Finish()

	mt.Cleanup(func() {
		cleanupCalled++
	})
}

func TestMockTestCleanupIsCalledAfterFailed(t *testing.T) {
	ch := make(chan struct{})

	go func() {
		defer func() {
			ch <- struct{}{}
		}()

		cleanupCalled := 0
		defer func() {
			require.Equal(t, 1, cleanupCalled)
		}()

		mt := testutil.NewMockTest(os.Stderr)
		defer mt.Finish()

		mt.Cleanup(func() {
			cleanupCalled++
		})

		mt.FailNow()
	}()

	<-ch
}
