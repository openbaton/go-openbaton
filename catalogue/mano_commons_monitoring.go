package catalogue

type PerceivedSeverity string

const (
	SeverityIndeterminate = PerceivedSeverity("INDETERMINATE")
	SeverityWarning       = PerceivedSeverity("WARNING")
	SeverityMinor         = PerceivedSeverity("MINOR")
	SeverityMajor         = PerceivedSeverity("MAJOR")
	SeverityCritical      = PerceivedSeverity("CRITICAL")
)

