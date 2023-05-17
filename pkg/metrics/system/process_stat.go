package system

// Stat represents process-stat data
type Stat struct {
	Utime     float64
	Stime     float64
	StartTime float64
	InOctets  float64
	OutOctets float64
}
