package daytime

import (
	"time"
)

type Clockfn func() time.Time

func Get(now Clockfn) string {
	const format = "Monday, January 2, 2006 15:04:05-MST"
	return now().Format(format)
}
