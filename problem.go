package problem

import (
	"encoding/json"
	"net/http"
)

// An Option configures a Problem using the functional options paradigm
// popularized by Rob Pike.
type Option interface {
	apply(*Problem)
}

type optionFunc func(*Problem)

func (f optionFunc) apply(problem *Problem) { f(problem) }

type Problem struct {
	data map[string]interface{}
}

// JSON returns the Problem as json bytes
func (p Problem) JSON() []byte {
	b, _ := p.MarshalJSON()
	return b
}

// MarshalJSON implements the json.Marshaler interface
func (p Problem) MarshalJSON() ([]byte, error) {
	return json.Marshal(&p.data)
}

// JSONString returns the Problem as json string
func (p Problem) JSONString() string {
	return string(p.JSON())
}

// WriteTo writes the Problem to a http Response Writer
func (p Problem) WriteTo(w http.ResponseWriter) (int, error) {
	w.Header().Set("Content-Type", "application/problem+json")
	if statuscode, ok := p.data["status"]; ok {
		if statusint, ok := statuscode.(int); ok {
			w.WriteHeader(statusint)
		}
	}
	return w.Write(p.JSON())
}

// New generates a new Problem
func New(opts ...Option) *Problem {
	problem := &Problem{}
	problem.data = make(map[string]interface{})
	for _, opt := range opts {
		opt.apply(problem)
	}
	return problem
}

// Append an Option to a existing Problem
func (p *Problem) Append(opts ...Option) *Problem {
	for _, opt := range opts {
		opt.apply(p)
	}
	return p
}

// Type sets the type URI (typically, with the "http" or "https" scheme) that identifies the problem type.
// When dereferenced, it SHOULD provide human-readable documentation for the problem type
func Type(uri string) Option {
	return optionFunc(func(problem *Problem) {
		problem.data["type"] = uri
	})
}

// Title sets a title that appropriately describes it (think short)
// Written in english and readable for engineers (usually not suited for
// non technical stakeholders and not localized); example: Service Unavailable
func Title(title string) Option {
	return optionFunc(func(problem *Problem) {
		problem.data["title"] = title
	})
}

// Status sets the HTTP status code generated by the origin server for this
// occurrence of the problem.
func Status(status int) Option {
	return optionFunc(func(problem *Problem) {
		problem.data["status"] = status
	})
}

// Detail A human readable explanation specific to this occurrence of the problem.
func Detail(detail string) Option {
	return optionFunc(func(problem *Problem) {
		problem.data["detail"] = detail
	})
}

// Instance an absolute URI that identifies the specific occurrence of the
// problem.
func Instance(uri string) Option {
	return optionFunc(func(problem *Problem) {
		problem.data["instance"] = uri
	})
}

// Custom sets a custom key value
func Custom(key string, value interface{}) Option {
	return optionFunc(func(problem *Problem) {
		problem.data[key] = value
	})
}
