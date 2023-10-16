package api

import "encoding/json"

const EXPOSED_ANUS = "EXPOSED_ANUS"
const EXPOSED_ARMPITS = "EXPOSED_ARMPITS"
const COVERED_BELLY = "COVERED_BELLY"
const EXPOSED_BELLY = "EXPOSED_BELLY"
const COVERED_BUTTOCKS = "COVERED_BUTTOCKS"
const EXPOSED_BUTTOCKS = "EXPOSED_BUTTOCKS"
const FACE_F = "FACE_F"
const FACE_M = "FACE_M"
const COVERED_FEET = "COVERED_FEET"
const EXPOSED_FEET = "EXPOSED_FEET"
const COVERED_BREAST_F = "COVERED_BREAST_F"
const EXPOSED_BREAST_F = "EXPOSED_BREAST_F"
const COVERED_GENITALIA_F = "COVERED_GENITALIA_F"
const EXPOSED_GENITALIA_F = "EXPOSED_GENITALIA_F"
const EXPOSED_BREAST_M = "EXPOSED_BREAST_M"
const EXPOSED_GENITALIA_M = "EXPOSED_GENITALIA_M"

type APIResponse struct {
	Unsafe  bool      `json:"unsafe" bson:"unsafe,omitempty"`
	Objects []Objects `json:"objects" bson:"objects,omitempty"`
}

type Objects struct {
	Box   []int   `json:"box" bson:"box,omitempty"`
	Score float32 `json:"score" bson:"score,omitempty"`
	Label string  `json:"label" bson:"label,omitempty"`
}

func ParseAPI(data []byte) (APIResponse, error) {
	var res APIResponse
	err := json.Unmarshal(data, &res)
	if err != nil {
		return res, err
	}
	return res, nil
}
