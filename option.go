package tsv

// Option configures the Encoder.
type Option func(*Encoder)

// WithTimeFormat sets a custom time format.
// By default, time.Time is encoded as a Unix epoch timestamp.
func WithTimeFormat(format string) Option {
	return func(e *Encoder) {
		e.timeFormat = format
	}
}
