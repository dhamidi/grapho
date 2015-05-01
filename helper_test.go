package grapho

import "encoding/json"

func asJSON(want interface{}, got interface{}) ([]byte, []byte) {
	wantJSON, err := json.MarshalIndent(want, "", "  ")
	if err != nil {
		panic(err)
	}

	gotJSON, err := json.MarshalIndent(got, "", "  ")
	if err != nil {
		panic(err)
	}

	return wantJSON, gotJSON
}
