package camunda

import (
	"strings"
	"time"
)

// Time a custom time format
type Time struct {
	time.Time
}

// UnmarshalJSON
func (t *Time) UnmarshalJSON(b []byte) (err error) {
	t.Time, err = time.Parse(DefaultDateTimeFormat, strings.Trim(string(b), "\""))
	return
}

// MarshalJSON
func (t *Time) MarshalJSON() ([]byte, error) {
	timeStr := t.Time.Format(DefaultDateTimeFormat)
	return []byte("\"" + timeStr + "\""), nil
}

// toCamundaTime return time formatted for camunda
func toCamundaTime(dt time.Time) string {
	if dt.IsZero() {
		return ""
	}

	return dt.Format(DefaultDateTimeFormat)
}
