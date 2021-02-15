package httpcli

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetResponse send http request with given url then decode response and fill the fields of given dest
func GetResponse(url string, dest interface{}, disallowUnknowFields bool) error {
	// #nosec
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("invalid url %s", url)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf(resp.Status)
	}

	// disallow unknow fileds to filter out error messages from wechat server when no err is returned
	decoder := json.NewDecoder(resp.Body)
	if disallowUnknowFields {
		decoder.DisallowUnknownFields()
	}

	if err := decoder.Decode(dest); err != nil {
		return err
	}

	return nil
}
