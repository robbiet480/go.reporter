package reporter_test

import (
	"fmt"

	"github.com/robbiet480/go.reporter"
)

// This example processes the given JSON as a single report for a single day.
func ExampleDecodeJSONString() {
	day, err := reporter.DecodeJSONString(`{"snapshots":[{"battery":0.9}]}`) // Truncated JSON
	if err != nil {
		fmt.Print(err)
	}
	fmt.Print(day)
}

// This example sets up a filesystem backend and returns the latest found report.
func ExampleNewFilesystemBackend() {
	backend, err := reporter.NewFilesystemBackend("")
	if err != nil {
		fmt.Print(err)
	}
	file, err := backend.GetLatestReport()
	if err != nil {
		fmt.Print(err)
	}
	day, err := reporter.DecodeFile(file)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Print(day)
}

// This example sets up a Dropbox backend and returns the latest found report.
func ExampleNewDropboxBackend() {
	backend, err := reporter.NewDropboxBackend("DROPBOX_ACCESS_TOKEN", "")
	if err != nil {
		fmt.Print(err)
	}
	file, err := backend.GetLatestReport()
	if err != nil {
		fmt.Print(err)
	}
	day, err := reporter.DecodeFile(file)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Print(day)
}
