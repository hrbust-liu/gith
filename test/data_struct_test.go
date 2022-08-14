package test

import (
    "testing"

    "github.com/hrbust-liu/gith/src/common"
)

func TestSliceSetOpt(t *testing.T) {
    s1 := []string{"A", "B", "C", "C"}
    s2 := []string{"A", "B", "B", "D"}
    onlySlice1, intersectSlice, onlySlice2 := common.SliceSetOpt(s1, s2)
    common.AssertEqual(t, onlySlice1, []string{"C"})
    common.AssertEqual(t, intersectSlice, []string{"A", "B"})
    common.AssertEqual(t, onlySlice2, []string{"D"})
}

func TestMapSetOpt(t *testing.T) {
    m1 := map[string]string{"A1": "B1", "A2": "B2", "A3": "B3"}
    m2 := map[string]string{"A1": "B1", "A2": "B3", "A4": "B4"}
    onlyMap1, commonMapValueSame, commonMapValueDiff, onlyMap2 := common.MapSetOpt(m1, m2)
    common.AssertEqual(t, onlyMap1, []string{"A3"})
    common.AssertEqual(t, commonMapValueSame, []string{"A1"})
    common.AssertEqual(t, commonMapValueDiff, []string{"A2"})
    common.AssertEqual(t, onlyMap2, []string{"A4"})
}
