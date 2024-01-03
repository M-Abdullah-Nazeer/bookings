package forms

import (
	"net/http"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {

	r, err := http.NewRequest("POST", "/my-url", nil)
	if err != nil {

		t.Error("Error creating request")
	}

	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("Form is not valid")
	}

}

// PostForm is a field in the Request struct. It represents form data in the body of an HTTP POST request. This field is a map of string slices (map[string][]string), where each key corresponds to a form field name, and the associated value is a slice of strings containing the values for that field.
func TestRequired(t *testing.T) {

	r, err := http.NewRequest("POST", "/my-url", nil)
	if err != nil {

		t.Error("Error creating request")
	}

	// creating a form object using some custom constructor (New) and passing in the PostForm data from the request

	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("Shows valid while required fields are missing")
	}

	//  Adds values to the form data for fields "a," "b," and "c."
	postData := url.Values{}
	postData.Add("a", "a")
	postData.Add("b", "a")
	postData.Add("c", "a")

	r, err = http.NewRequest("POST", "/my-url", nil)
	if err != nil {

		t.Error("Error creating 2nd request")
	}

	// Sets the form data created earlier as the PostForm of the new request.
	r.PostForm = postData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("Shows required feilds are missing when it does not")
	}

}

func TestForm_MinLength(t *testing.T) {

	postData := url.Values{}
	form := New(postData)

	form.MinLength("x", 3)
	if form.Valid() {
		t.Error("Shows no error for non-existent field minlen")
	}

	IsError := form.Errors.Get("x")
	if IsError == "" {
		t.Error("Should have an error but did not get one")
	}

	postData = url.Values{}
	postData.Add("a", "abdullah")
	form = New(postData)

	form.MinLength("a", 3)
	if !form.Valid() {
		t.Error("Got an error where minlength is satisfied")
	}

	IsError = form.Errors.Get("a")
	if IsError != "" {
		t.Error("Should not have an error but did got one")
	}

	postData = url.Values{}
	postData.Add("b", "al")
	form = New(postData)
	form.MinLength("b", 3)
	if form.Valid() {
		t.Error("Shows minlen 3 is met but its not met")
	}
}

func TestForm_Has(t *testing.T) {

	postData := url.Values{}
	form := New(postData)

	form.Has("a")
	if form.Valid() {
		t.Error("Should give has error for non existent field but flagged valid")
	}

	postData = url.Values{}
	postData.Add("a", "abdullah")
	form = New(postData)

	form.Has("a")
	if !form.Valid() {
		t.Error("Should not give has error for provided field")
	}

}

func TestForm_IsEmail(t *testing.T) {

	postData := url.Values{}
	form := New(postData)
	form.IsEmail("a")
	if form.Valid() {
		t.Error("Should give email error for non existent field but flagged valid")
	}

	postData = url.Values{}
	postData.Add("a", "abdulla")
	form = New(postData)
	form.IsEmail("a")
	if form.Valid() {
		t.Error("Should give invalid email error but flagged valid")
	}

	postData = url.Values{}
	postData.Add("a", "abdullah@gmail.com")
	form = New(postData)
	form.IsEmail("a")
	if !form.Valid() {
		t.Error("Should not show error, email is valid")
	}

}
