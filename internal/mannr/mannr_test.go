package mannr

import "testing"

func TestSN2MN(t *testing.T) {
	mnr := Serial2Manufacturer("M2XTEST0061")
	if mnr != "29B472FBF02F6A5552C7FEE2" {
		t.Errorf("got invalid manufacturer nr: %s", mnr)
	}
}
