package docx

import "fmt"

// DocxError is the base error type for the go-docx package.
type DocxError struct {
	msg string
}

// Error implements the error interface.
func (e *DocxError) Error() string { return e.msg }

// NewDocxError creates a new DocxError with the given message.
func NewDocxError(msg string) *DocxError {
	return &DocxError{msg: msg}
}

// InvalidXmlError is raised when invalid XML is encountered, such as on attempt
// to access a missing required child element.
type InvalidXmlError struct {
	DocxError
}

// NewInvalidXmlError creates a new InvalidXmlError.
func NewInvalidXmlError(format string, args ...any) *InvalidXmlError {
	return &InvalidXmlError{DocxError{msg: fmt.Sprintf(format, args...)}}
}

// PackageNotFoundError is raised when a package (file) cannot be found.
type PackageNotFoundError struct {
	DocxError
}

// NewPackageNotFoundError creates a new PackageNotFoundError.
func NewPackageNotFoundError(format string, args ...any) *PackageNotFoundError {
	return &PackageNotFoundError{DocxError{msg: fmt.Sprintf(format, args...)}}
}

// InvalidSpanError is raised when an invalid merge region is specified
// in a request to merge table cells.
type InvalidSpanError struct {
	DocxError
}

// NewInvalidSpanError creates a new InvalidSpanError.
func NewInvalidSpanError(format string, args ...any) *InvalidSpanError {
	return &InvalidSpanError{DocxError{msg: fmt.Sprintf(format, args...)}}
}
