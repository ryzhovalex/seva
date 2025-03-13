// Time is in milliseconds, unless other is clearly specified.
package bone

import "time"

func Utc() int64 {
	return time.Now().Unix()
}

// Formats timestamp to a date.
func Date(ms int, format string) string {
	return time.Unix(int64(ms), 0).Format(format)
}

func Sleep(duration int64) {
	time.Sleep(time.Duration(duration) * time.Millisecond)
}
