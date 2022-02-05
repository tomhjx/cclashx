package core

import (
	"net"
	"testing"
)

func TestAddProxy(t *testing.T) {
	ip0, err := net.ResolveIPAddr("ip", "baidu.com")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(ip0)

	ip1, err := net.ResolveIPAddr("ip", "cm-jm.okvpn.xyz")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(ip1)

	ip2, err := net.ResolveIPAddr("ip", "66.42.103.1")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(ip2)

	ip3, err := net.ResolveIPAddr("ip", "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(ip3)

	_, err2 := net.ResolveIPAddr("ip", "abc")
	if err2 == nil {
		t.Fatal("abc hava a ip??")
	}
	t.Log(err2)

}
