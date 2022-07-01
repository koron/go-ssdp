package ssdplog_test

import (
	"bytes"
	"log"
	"testing"

	"github.com/koron/go-ssdp/internal/ssdplog"
)

func TestLoggerProvider(t *testing.T) {
	b := &bytes.Buffer{}
	logger := log.New(b, "", 0)

	ssdplog.Printf("never output")
	if s := b.String(); s != "" {
		t.Errorf("unexpected log #1:\nwant=(empty)\n got=%q", s)
	}

	// provide LoggerProvider
	ssdplog.LoggerProvider = func() *log.Logger { return logger }
	ssdplog.Printf("foo")
	if s := b.String(); s != "foo\n" {
		t.Errorf("unexpected log #1:\nwant=%q\n got=%q", "foo\n", s)
	}

	// disable LoggerProvider, output buffer never changed
	ssdplog.LoggerProvider = func() *log.Logger { return nil }
	ssdplog.Printf("bar")
	if s := b.String(); s != "foo\n" {
		t.Errorf("unexpected log #1:\nwant=%q\n got=%q", "foo\n", s)
	}
}
