package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func unmarshalRequestBody(r *http.Request, v interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}
