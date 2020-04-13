package mock

import "github.com/wolffcm/flux/values"

// AscendingTimeProvider provides ascending timestamps every nanosecond
// starting from Start.
type AscendingTimeProvider struct {
	Start int64
}

func (atp *AscendingTimeProvider) CurrentTime() values.Time {
	t := values.Time(atp.Start)
	atp.Start++
	return t
}
