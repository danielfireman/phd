package load

import "time"

type DurationMetric struct {
	Max        time.Duration
	Min        time.Duration
	NumSamples int64
	Avg        time.Duration
	Sum        time.Duration
}

func (m *DurationMetric) sample(s time.Duration) {
	m.NumSamples++
	m.Sum += s
	if s > m.Max {
		m.Max = s
	} else if s < m.Min || m.Min == 0 {
		m.Min = s
	}
}

type Summary struct {
	Duration    time.Duration
	DurationStr string
	NumErrors   int64
	QPS         float64
	Latency     DurationMetric
}

func summarize(reqs <-chan RequestResult, start time.Time) *Summary {
	results := &Summary{}
	for r := range reqs {
		switch {
		case r.Err != nil:
			results.NumErrors++
		default:
			results.Latency.sample(r.Duration)
		}
	}
	results.Duration = time.Now().Sub(start)
	results.DurationStr = results.Duration.String()
	results.QPS = float64(results.NumErrors+results.Latency.NumSamples) / results.Duration.Seconds()
	results.Latency.Avg = time.Duration(float64(results.Latency.Sum) / float64(results.Latency.NumSamples))
	return results
}
