package external

import (
	"math"
	"testing"
)

func TestRoundCoordinates(t *testing.T) {
	lat, lon := roundCoordinates(52.520833, 13.409444)
	if math.Abs(lat-52.5182) > 0.0001 {
		t.Errorf("expected approx lat to be 52.5182, got %f", lat)
	}
	if math.Abs(lon-13.4052) > 0.0001 {
		t.Errorf("expected approx lon to be 13.4052, got %f", lon)
	}
}
