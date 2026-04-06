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

// WithUTC sets whether or not the time should use UTC.
// By default, utc is false and it uses the defined timezone.
func WithUTC(utc bool) Option {
	return func(e *Encoder) {
		e.utc = utc
	}
}

// WithCRLF sets the line ending to CRLF (\r\n) instead of LF (\n).
// This is useful for Windows compatibility.
func WithCRLF() Option {
	return func(e *Encoder) {
		e.crlf = true
	}
}

// WithDelimiter sets a custom field delimiter.
// By default, a tab character is used.
func WithDelimiter(delimiter rune) Option {
	return func(e *Encoder) {
		e.delimiter = delimiter
	}
}
