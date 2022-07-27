package diags

type Diagnostic interface {
	Severity() Severity
	Description() Description

	// ExtraInfo returns the raw extra information value. This is a low-level
	// API which requires some work on the part of the caller to properly
	// access associated information, so in most cases it'll be more convienient
	// to use the package-level ExtraInfo function to try to unpack a particular
	// specialized interface from this value.
	ExtraInfo() interface{}
}

type Severity rune

//go:generate go run golang.org/x/tools/cmd/stringer -type=Severity

const (
	Error   Severity = 'E'
	Warning Severity = 'W'
)

type Description struct {
	Address string
	Summary string
	Detail  string
}
