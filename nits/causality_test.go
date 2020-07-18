package nits

import "testing"

func TestIsParentOf(t *testing.T) {
	pp := preprocess(DefaultCase())
	plows := pp.findEvent("plows")
	carDies := pp.findEvent("car_dies")
	if !isParentOf(carDies, plows) {
		t.Error("isParentOf(carDies, plows); got:false, want:true")
	}
	if isParentOf(plows, carDies) {
		t.Error("isParentOf(plows, carDies); got:true, want:false")
	}
}
