package govalidator

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
)

func TestValidator_SetDefaultRequired(t *testing.T) {
	v := New(Options{})
	v.SetDefaultRequired(true)
	if !v.Opts.RequiredDefault {
		t.Error("SetDefaultRequired failed")
	}
}

func TestValidator_Validate(t *testing.T) {
	var URL *url.URL
	URL, _ = url.Parse("http://www.example.com")
	params := url.Values{}
	params.Add("name", "John Doe")
	params.Add("username", "jhondoe")
	params.Add("email", "john@mail.com")
	params.Add("zip", "8233")
	URL.RawQuery = params.Encode()
	r, _ := http.NewRequest("GET", URL.String(), nil)
	rulesList := MapData{
		"name":  []string{"required"},
		"age":   []string{"between:5,16"},
		"email": []string{"email"},
		"zip":   []string{"digits:4"},
	}

	opts := Options{
		Request: r,
		Rules:   rulesList,
	}
	v := New(opts)
	validationError := v.Validate()
	if len(validationError) > 0 {
		t.Log(validationError)
		t.Error("Validate failed to validate correct inputs!")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Validate did not panic")
		}
	}()

	v1 := New(Options{Rules: MapData{}})
	v1.Validate()
}

func Benchmark_Validate(b *testing.B) {
	var URL *url.URL
	URL, _ = url.Parse("http://www.example.com")
	params := url.Values{}
	params.Add("name", "John Doe")
	params.Add("age", "27")
	params.Add("email", "john@mail.com")
	params.Add("zip", "8233")
	URL.RawQuery = params.Encode()
	r, _ := http.NewRequest("GET", URL.String(), nil)
	rulesList := MapData{
		"name":  []string{"required"},
		"age":   []string{"numeric_between:18,60"},
		"email": []string{"email"},
		"zip":   []string{"digits:4"},
	}

	opts := Options{
		Request: r,
		Rules:   rulesList,
	}
	v := New(opts)
	for n := 0; n < b.N; n++ {
		v.Validate()
	}
}

//============ validate json test ====================

func TestValidator_ValidateJSON(t *testing.T) {
	type User struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Address string `json:"address"`
		Age     int    `json:"age"`
		Zip     string `json:"zip"`
		Color   int    `json:"color"`
	}

	postUser := User{
		Name:    "",
		Email:   "inalid email",
		Address: "",
		Age:     1,
		Zip:     "122",
		Color:   5,
	}

	rules := MapData{
		"name":    []string{"required"},
		"email":   []string{"email"},
		"address": []string{"required", "between:3,5"},
		"age":     []string{"bool"},
		"zip":     []string{"len:4"},
		"color":   []string{"min:10"},
	}

	var user User

	body, _ := json.Marshal(postUser)
	req, _ := http.NewRequest("POST", "http://www.example.com", bytes.NewReader(body))

	opts := Options{
		Request: req,
		Data:    &user,
		Rules:   rules,
	}

	vd := New(opts)
	vd.SetTagIdentifier("json")
	validationErr := vd.ValidateJSON()
	if len(validationErr) != 5 {
		t.Error("ValidateJSON failed")
	}
}

func TestValidator_ValidateJSON_WithoutRequest(t *testing.T) {
	postUser := map[string]interface{}{
		"name":    "",
		"email":   "inalid email",
		"address": "",
		"age":     1,
		"zip":     "122",
		"color":   5,
	}

	rules := MapData{
		"name":    []string{"required"},
		"email":   []string{"email"},
		"address": []string{"required", "between:3,5"},
		"age":     []string{"bool"},
		"zip":     []string{"len:4"},
		"color":   []string{"min:10"},
	}

	//body, _ := json.Marshal(postUser)
	//req, _ := http.NewRequest("POST", "http://www.example.com", bytes.NewReader(body))

	opts := Options{
		Data:    postUser,
		Rules:   rules,
	}

	vd := New(opts)
	vd.SetTagIdentifier("json")
	validationErr := vd.Validate()
	if len(validationErr) != 5 {
		t.Error("ValidateJSON failed")
	}
}

func TestValidator_ValidateJSON_NULLValue(t *testing.T) {
	type User struct {
		Name   string `json:"name"`
		Count  Int    `json:"count"`
		Option Int    `json:"option"`
		Active Bool   `json:"active"`
	}

	rules := MapData{
		"name":   []string{"required"},
		"count":  []string{"required"},
		"option": []string{"required"},
		"active": []string{"required"},
	}

	postUser := map[string]interface{}{
		"name":   "John Doe",
		"count":  0,
		"option": nil,
		"active": nil,
	}

	var user User
	body, _ := json.Marshal(postUser)
	req, _ := http.NewRequest("POST", "http://www.example.com", bytes.NewReader(body))

	opts := Options{
		Request: req,
		Data:    &user,
		Rules:   rules,
	}

	vd := New(opts)
	vd.SetTagIdentifier("json")
	validationErr := vd.ValidateJSON()
	if len(validationErr) != 2 {
		t.Error("ValidateJSON failed")
	}
}

func TestValidator_ValidateStruct(t *testing.T) {
	type User struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Address string `json:"address"`
		Age     int    `json:"age"`
		Zip     string `json:"zip"`
		Color   int    `json:"color"`
	}

	postUser := User{
		Name:    "",
		Email:   "inalid email",
		Address: "",
		Age:     1,
		Zip:     "122",
		Color:   5,
	}

	rules := MapData{
		"name":    []string{"required"},
		"email":   []string{"email"},
		"address": []string{"required", "between:3,5"},
		"age":     []string{"bool"},
		"zip":     []string{"len:4"},
		"color":   []string{"min:10"},
	}

	opts := Options{
		Data:  &postUser,
		Rules: rules,
	}

	vd := New(opts)
	vd.SetTagIdentifier("json")
	validationErr := vd.ValidateStruct()
	if len(validationErr) != 5 {
		t.Error("ValidateStruct failed")
	}
}

func TestValidator_ValidateJSON_NoRules_panic(t *testing.T) {
	opts := Options{}

	assertPanicWith(t, errValidateArgsMismatch, func() {
		New(opts).ValidateJSON()
	})
}

func TestValidator_ValidateJSON_NonPointer_panic(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", nil)

	type User struct {
	}

	var user User
	opts := Options{
		Request: req,
		Data:    user,
		Rules: MapData{
			"name": []string{"required"},
		},
	}

	assertPanicWith(t, errRequirePtr, func() {
		New(opts).ValidateJSON()
	})
}

func TestValidator_ValidateStruct_NoRules_panic(t *testing.T) {
	opts := Options{}

	assertPanicWith(t, errRequireRules, func() {
		New(opts).ValidateStruct()
	})
}

func TestValidator_ValidateStruct_RequestProvided_panic(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", nil)
	opts := Options{
		Request: req,
		Rules: MapData{
			"name": []string{"required"},
		},
	}

	assertPanicWith(t, errRequestNotAccepted, func() {
		New(opts).ValidateStruct()
	})
}

func TestValidator_ValidateStruct_NonPointer_panic(t *testing.T) {
	type User struct {
	}

	var user User
	opts := Options{
		Data: user,
		Rules: MapData{
			"name": []string{"required"},
		},
	}

	assertPanicWith(t, errRequirePtr, func() {
		New(opts).ValidateStruct()
	})
}

func TestValidator_ValidateStruct_DataNil_panic(t *testing.T) {
	opts := Options{
		Rules: MapData{
			"name": []string{"required"},
		},
	}

	assertPanicWith(t, errRequireData, func() {
		New(opts).ValidateStruct()
	})
}

func TestHtmlTags(t *testing.T) {
	type TestCheck struct {
		Input  string
		HasTag bool
	}

	tests := []TestCheck {
		{"", false},
		{"Hello, World!", false},
		{"foo&amp;bar", false},
		{`Hello <a href="www.example.com/">World</a>!`, true},
		{"Foo <textarea>Bar</textarea> Baz", true},
		{"Foo <!-- Bar --> Baz", true},
		{"<", false},
		{"foo < bar", false},
		{`Foo<script type="text/javascript">alert(1337)</script>Bar`, true},
		{`Foo<div title="1>2">Bar`, true},
		{`I <3 Ponies!`, false},
		{`<script>foo()</script>`, true},
	}

	type User struct {
		Data   string `json:"data"`
	}

	rules := MapData{
		"data":   []string{"htmltag"},
	}

	//var result bool
	for _, test := range tests {
		postUser := User{
			Data:   test.Input,
		}

		opts := Options{
			Data:  &postUser,
			Rules: rules,
		}

		vd := New(opts)
		vd.SetTagIdentifier("json")
		err := vd.ValidateStruct()

		// give error if there is tag in string and still validator returns no error
		if test.HasTag == isEmpty(err) {
			t.Errorf("%s: expected %t, actual %v", test.Input, test.HasTag, err)
		}
	}
}
