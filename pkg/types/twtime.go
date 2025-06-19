package types

import (
	"encoding/json"
	"strings"
	"time"
)

type TWTime time.Time

func (t *TWTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" {
		return nil
	}
	parsed, err := time.Parse("20060102T150405Z", s)
	if err != nil {
		return err
	}
	*t = TWTime(parsed)
	return nil
}

func (t TWTime) Time() time.Time {
	return time.Time(t)
}

func (t TWTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t).Format("20060102T150405Z"))
}