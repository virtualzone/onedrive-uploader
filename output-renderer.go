package main

import (
	"fmt"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

type OutputRenderer struct {
	ProgressBar *progressbar.ProgressBar
	Quiet       bool
}

func (r *OutputRenderer) initProgressBar(maxBytes int64, desc string) {
	if r.Quiet {
		return
	}
	r.ProgressBar = progressbar.NewOptions64(
		maxBytes,
		progressbar.OptionSetDescription(desc),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(20),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() { fmt.Printf("\n") }),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionThrottle(50*time.Millisecond),
	)
}

func (r *OutputRenderer) initSpinner(desc string) {
	if r.Quiet {
		return
	}
	r.ProgressBar = progressbar.NewOptions64(
		-1,
		progressbar.OptionSetDescription(desc),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetWidth(20),
		progressbar.OptionOnCompletion(func() {}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionClearOnFinish(),
	)
	// Spin in an asynchronous thread
	go func() {
		for !r.ProgressBar.IsFinished() {
			r.ProgressBar.Add(1)
			time.Sleep(50 * time.Millisecond)
		}
	}()
}

func (r *OutputRenderer) stopSpinner() {
	if r.Quiet {
		return
	}
	r.ProgressBar.Finish()
}

func (r *OutputRenderer) updateProgressBar(totalBytes int64) {
	if r.Quiet {
		return
	}
	r.ProgressBar.Set64(totalBytes)
}
