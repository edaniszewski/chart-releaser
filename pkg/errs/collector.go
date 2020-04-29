package errs

import "strings"

// Collector collects multiple errors, allowing for deferred error
// handling and propagation.
type Collector struct {
	errors []error
}

// NewCollector creates a new error Collector.
func NewCollector() *Collector {
	return &Collector{}
}

// Count the number of errors currently collected.
func (c *Collector) Count() int {
	return len(c.errors)
}

// HasErrors checks whether the collector has any errors added.
func (c *Collector) HasErrors() bool {
	return len(c.errors) > 0
}

// Add an error to the Collector.
//
// If an error collector is being added to this one, all errors which
// it has accumulated are copied into this collector.
func (c *Collector) Add(err error) {
	if col, ok := err.(*Collector); ok {
		c.errors = append(c.errors, col.errors...)
	} else {
		c.errors = append(c.errors, err)
	}
}

// Error returns the error string for the collector.
func (c *Collector) Error() string {
	if c.HasErrors() {
		var errors []string
		for _, e := range c.errors {
			errors = append(errors, e.Error())
		}
		return "\nErrors:\n • " + strings.Join(errors, "\n • ") + "\n\n"
	}
	return ""
}
