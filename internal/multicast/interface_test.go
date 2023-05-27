//go:build !android
// +build !android

package multicast

import (
	"net"
	"reflect"
	"testing"
)

func TestInterfaces(t *testing.T) {
	list, err := interfaces()
	if err != nil {
		t.Fatalf("interfaces() failed: %s", err)
	}
	if len(list) == 0 {
		t.Error("interfaces() returns no interfaces")
	}
}

func TestInterafceProviders(t *testing.T) {
	want := []net.Interface{
		{Index: 123, Name: "Test#1"},
		{Index: 456, Name: "Test#2"},
		{Index: 789, Name: "Test#3"},
	}
	InterfacesProvider = func() []net.Interface {
		return want
	}
	defer func() { InterfacesProvider = nil }()
	got, err := interfaces()
	if err != nil {
		t.Fatalf("interfaces() failed: %s", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("unexpected interfaces:\nwant=%+v\ngot=%+v", want, got)
	}
}
