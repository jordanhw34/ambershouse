package forms

type errors map[string][]string

// Add => adds error message for a given form field
func (err errors) Add(field, message string) {
	err[field] = append(err[field], message)
}

// Get => returns first error message
func (err errors) Get(field string) string {
	errStr := err[field]
	if len(errStr) == 0 {
		return ""
	}
	return errStr[0]
}
