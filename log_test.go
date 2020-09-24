package ssdp

import (
	"bytes"
	"log"
	"testing"
)

func TestLogger(t *testing.T) {
	b := &bytes.Buffer{}
	Logger = log.New(b, "", 0)
	logf("foo")
	if s := b.String(); s != "foo\n" {
		t.Errorf("unexpected log #1:\nwant=%q\n got=%q", "foo\n", s)
	}
	Logger = nil
	logf("bar")
	if s := b.String(); s != "foo\n" {
		t.Errorf("unexpected log #2:\nwant=%q\n got=%q", "foo\n", s)
	}
}
