package catalogue

//go:generate stringer -type=PerceivedSeverity
type PerceivedSeverity string

const (
	SeverityIndeterminate = PerceivedSeverity("INDETERMINATE")
	SeverityWarning       = PerceivedSeverity("WARNING")
	SeverityMinor         = PerceivedSeverity("MINOR")
	SeverityMajor         = PerceivedSeverity("MAJOR")
	SeverityCritical      = PerceivedSeverity("CRITICAL")
)
