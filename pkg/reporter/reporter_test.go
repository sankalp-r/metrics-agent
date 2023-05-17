package reporter

import (
	"bufio"
	"bytes"
	"io"
	"testing"
)

func TestReporter(t *testing.T) {
	testFile := &mockFile{}

	testReporter := fileReporter{
		targetFile: testFile,
		metrics:    make(chan []byte, 1),
		buffer:     bufio.NewWriterSize(testFile, 2),
	}

	go testReporter.Report()
	testReporter.Write([]byte("test_metric_1 1\n"))
	testReporter.Write([]byte("test_metric_2 2"))

	testReporter.Close()
	result := testFile.buffer.String()
	expected := "test_metric_1 1\ntest_metric_2 2"
	if result != expected {
		t.Errorf("expeted: %s, \ngot: %s", expected, result)
	}

}

type mockFile struct {
	io.WriteCloser
	buffer bytes.Buffer
}

func (m *mockFile) Write(b []byte) (n int, err error) {
	return m.buffer.Write(b)
}

func (m *mockFile) Close() error {
	return nil
}
