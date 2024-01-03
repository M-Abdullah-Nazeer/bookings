package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form creates custom form struct, embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

func (f *Form) Valid() bool {

	return len(f.Errors) == 0
}

// New initializes a form struct
func New(data url.Values) *Form {

	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required checks for required fields
func (f *Form) Required(fields ...string) {

	for _, field := range fields {
		value := f.Get(field)

		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be empty")
		}

	}

}

// Has checks if form field is in post and not empty, incase it is required
func (f *Form) Has(field string) bool {

	x := f.Get(field)

	if x == "" {
		f.Errors.Add(field, "This field cannot be blank")
		return false
	}
	return true

}

// MinLength checks for string min length
func (f *Form) MinLength(field string, length int) bool {

	x := f.Get(field)

	if len(x) < length {

		f.Errors.Add(field, fmt.Sprintf("Field must contain at least %d letters", length))
		return false
	}
	return true

}

// IsEmail validates email address
func (f *Form) IsEmail(field string) bool {

	if !govalidator.IsEmail(f.Get(field)) { //IsEmail is built-in func of govalidator

		f.Errors.Add(field, "Invalid Email Address")
		return false
	}

	return true

}
