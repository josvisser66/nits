package nits

import "testing"

func TestIsParentOf(t *testing.T) {
	pp := preprocess(test_case())
	plows := pp.findEvent("plows")
	carDies := pp.findEvent("car_dies")
	if !carDies.isParentOf(plows) {
		t.Error("carDies.isParentOf(plows); got:false, want:true")
	}
	if plows.isParentOf(carDies) {
		t.Error("plows.isParentOf(carDies); got:true, want:false")
	}
}
