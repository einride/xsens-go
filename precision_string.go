// Code generated by "stringer -type Precision -trimprefix Precision -output precision_string.go precision.go"; DO NOT EDIT.

package xsens

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[PrecisionFloat32-0]
	_ = x[PrecisionFP1220-1]
	_ = x[PrecisionFP1632-2]
	_ = x[PrecisionFloat64-3]
}

const _Precision_name = "Float32FP1220FP1632Float64"

var _Precision_index = [...]uint8{0, 7, 13, 19, 26}

func (i Precision) String() string {
	if i >= Precision(len(_Precision_index)-1) {
		return "Precision(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Precision_name[_Precision_index[i]:_Precision_index[i+1]]
}
