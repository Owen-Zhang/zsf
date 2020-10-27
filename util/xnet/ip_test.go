package xnet

import "testing"

func Test_GetLocalMainIP(t *testing.T) {
	ip, _, err := GetLocalMainIP()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("ip is:%s", ip)
}
