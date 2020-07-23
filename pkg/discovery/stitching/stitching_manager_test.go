package stitching

import (
	"testing"
)

func TestNewStitchingManager(t *testing.T) {

	items := []StitchingPropertyType{UUID, IP, UUID, IP}

	for _, ptype := range items {
		m := NewStitchingManager(ptype)
		if m.GetStitchType() != ptype {
			t.Errorf("Stitching type wrong: %v Vs. %v", m.stitchType, ptype)
		}
	}
}
