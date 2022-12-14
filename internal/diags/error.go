package diags

// nativeError is a Diagnostic implementation that wraps a normal Go error
type nativeError struct {
	err error
}

var _ Diagnostic = nativeError{}

func (e nativeError) Severity() Severity {
	return Error
}

func (e nativeError) Description() Description {
	return Description{
		Summary: e.err.Error(),
	}
}

func (e nativeError) ExtraInfo() interface{} {
	// Native errors don't carry any "extra information".
	return nil
}
