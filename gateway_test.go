package discgo

import "testing"

func TestGatewayEndpoint_Get(t *testing.T) {
	e := c.gateway()
	url, err := e.get()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(url)
}
