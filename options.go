package bishamon

func WithFieldsFromMapFunc(fieldsFromMapFunc FieldsFromMapFunc) RedactorOption {
	return func(r *Redactor) {
		r.fieldsFromMapFunc = fieldsFromMapFunc
	}
}

func WithRedactFuncs(redactFuncs ...RedactFieldFunc) RedactorOption {
	return func(r *Redactor) {
		r.redactorFuncs.redactFieldsFuncs = append(r.redactorFuncs.redactFieldsFuncs, redactFuncs...)
	}
}

func WithRedactMapFuncs(redactMapFuncs ...RedactMapFunc) RedactorOption {
	return func(r *Redactor) {
		r.redactorFuncs.redactMapFuncs = append(r.redactorFuncs.redactMapFuncs, redactMapFuncs...)
	}
}

func WithRedactListFuncs(redactListFuncs ...RedactListFunc) RedactorOption {
	return func(r *Redactor) {
		r.redactorFuncs.redactListFuncs = append(r.redactorFuncs.redactListFuncs, redactListFuncs...)
	}
}
