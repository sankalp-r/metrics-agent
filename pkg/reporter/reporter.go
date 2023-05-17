package reporter

import (
	"bufio"
	"go.uber.org/zap"
	"io"
	"os"
)

// Reporter of collected metrics
type Reporter interface {
	Report()
	Write([]byte)
	Close()
}

// NewReporter returns new reporter
func NewReporter(fileName string) (Reporter, error) {
	targetFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	return &fileReporter{
		targetFile: targetFile,
		metrics:    make(chan []byte, 100),
		buffer:     bufio.NewWriter(targetFile),
	}, nil
}

// fileReporter collects metrics in a file
type fileReporter struct {
	targetFile io.WriteCloser
	metrics    chan []byte
	buffer     *bufio.Writer
}

// Report metrics in a buffer
func (f *fileReporter) Report() {
	for metric := range f.metrics {
		_, err := f.buffer.Write(metric)
		if err != nil {
			zap.L().Error(err.Error())
			f.buffer.Reset(f.targetFile)
		}
	}
}

// Close the files and channels
func (f *fileReporter) Close() {
	close(f.metrics)

	if err := f.buffer.Flush(); err != nil {
		zap.L().Error(err.Error())
	}

	if err := f.targetFile.Close(); err != nil {
		zap.L().Error(err.Error())
	}
}

// Write to reporter
func (f *fileReporter) Write(metric []byte) {
	f.metrics <- metric
}
