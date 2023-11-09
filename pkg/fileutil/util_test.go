package fileutil

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"
)

var origFileContent = `
Abbatxt7095
`

var updatedFileContent = `
efgwht2033
`

type testStruct struct {
	content         string
	expectedContent string
	mutex           sync.Mutex
}

func (a *testStruct) CallBackForFileLoad(ctx context.Context, dynamicContent []byte) error {
	a.mutex.Lock()
	a.expectedContent = string(dynamicContent)
	defer a.mutex.Unlock()
	return nil
}

func (a *testStruct) CallBackForFileDeletion(ctx context.Context) error {
	a.mutex.Lock()
	a.expectedContent = ""
	defer a.mutex.Unlock()
	return nil
}

func TestLoadDynamicFile(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{
			"abcde",
			"abcde",
		},
		{
			"fghijk",
			"fghijk",
		},
		{
			"xyzopq",
			"xyzopq",
		},
		{
			"eks:test",
			"eks:test",
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	testA := &testStruct{}
	StartLoadDynamicFile(ctx, "/tmp/util_test.txt", testA)
	time.Sleep(2 * time.Second)
	os.WriteFile("/tmp/util_test.txt", []byte("test"), 0777)
	for {
		time.Sleep(1 * time.Second)
		testA.mutex.Lock()
		if testA.expectedContent == "test" {
			t.Log("read to test")
			testA.mutex.Unlock()
			break
		}
		testA.mutex.Unlock()
	}
	for _, c := range cases {
		updateFile(testA, c.input, t)
		testA.mutex.Lock()
		if testA.expectedContent != c.want {
			t.Errorf(
				"Unexpected result: TestLoadDynamicFile: got: %s, wanted %s",
				testA.expectedContent,
				c.want,
			)
		}
		testA.mutex.Unlock()
	}
}

func updateFile(testA *testStruct, origFileContent string, t *testing.T) {
	testA.content = origFileContent
	data := []byte(origFileContent)
	err := os.WriteFile("/tmp/util_test.txt", data, 0600)
	if err != nil {
		t.Errorf("failed to create a local file /tmp/util_test.txt")
	}
	time.Sleep(1 * time.Second)
}

func TestDeleteDynamicFile(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	testA := &testStruct{}
	StartLoadDynamicFile(ctx, "/tmp/delete.txt", testA)
	time.Sleep(2 * time.Second)
	os.WriteFile("/tmp/delete.txt", []byte("test"), 0777)
	for {
		time.Sleep(1 * time.Second)
		testA.mutex.Lock()
		if testA.expectedContent == "test" {
			t.Log("read to test")
			testA.mutex.Unlock()
			break
		}
		testA.mutex.Unlock()
	}
	os.Remove("/tmp/delete.txt")
	time.Sleep(2 * time.Second)
	testA.mutex.Lock()
	if testA.expectedContent != "" {
		t.Errorf("failed in TestDeleteDynamicFile")
	}
	testA.mutex.Unlock()
}
