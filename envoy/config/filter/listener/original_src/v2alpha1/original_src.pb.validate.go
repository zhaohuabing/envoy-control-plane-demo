//go:build !disable_pgv
// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: envoy/config/filter/listener/original_src/v2alpha1/original_src.proto

package v2alpha1

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on OriginalSrc with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *OriginalSrc) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on OriginalSrc with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in OriginalSrcMultiError, or
// nil if none found.
func (m *OriginalSrc) ValidateAll() error {
	return m.validate(true)
}

func (m *OriginalSrc) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for BindPort

	// no validation rules for Mark

	if len(errors) > 0 {
		return OriginalSrcMultiError(errors)
	}

	return nil
}

// OriginalSrcMultiError is an error wrapping multiple validation errors
// returned by OriginalSrc.ValidateAll() if the designated constraints aren't met.
type OriginalSrcMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m OriginalSrcMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m OriginalSrcMultiError) AllErrors() []error { return m }

// OriginalSrcValidationError is the validation error returned by
// OriginalSrc.Validate if the designated constraints aren't met.
type OriginalSrcValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e OriginalSrcValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e OriginalSrcValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e OriginalSrcValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e OriginalSrcValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e OriginalSrcValidationError) ErrorName() string { return "OriginalSrcValidationError" }

// Error satisfies the builtin error interface
func (e OriginalSrcValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sOriginalSrc.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = OriginalSrcValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = OriginalSrcValidationError{}
