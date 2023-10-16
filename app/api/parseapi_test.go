package api

import "testing"

func TestParseAPI(t *testing.T) {
	data := []byte(`{"unsafe":true,"objects":[{"box":[52,88,287,299],"score":0.6513780355453491,"label":"FACE_F"}]}`)
	
	apires, err := ParseAPI(data)
	if err != nil {
		t.Fatal(err)
	}

	if apires.Unsafe != true {
		t.Error("Incorrect Unsafe != true")
	}
	
	if len(apires.Objects) != 1 {
		t.Fatalf("Invalid number of objects: %d", len(apires.Objects))
	}

	if apires.Objects[0].Label != "FACE_F" {
		t.Errorf("Incorrect object label: %s != %s", apires.Objects[0].Label, "FACE_F")
	}
}