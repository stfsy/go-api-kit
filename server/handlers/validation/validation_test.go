package validation

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

// Nested struct for testing
type Address struct {
	City    string `json:"c" validate:"required"`
	ZipCode string `json:"z" validate:"len=5"`
}

type UserWithAddress struct {
	Name    string   `json:"name" validate:"required"`
	Address *Address `json:"a" validate:"required"`
}

func TestValidateStruct_NestedStruct(t *testing.T) {
	u := UserWithAddress{
		Name:    "Valid Name",
		Address: &Address{City: "", ZipCode: "123"},
	}
	errors := ValidateStruct(u)
	if errors["a.c"].Validator != "required" {
		t.Errorf("expected required error for a.c, got %v", errors["a.c"])
	}
	if errors["a.z"].Validator != "len" {
		t.Errorf("expected len error for a.z, got %v", errors["a.z"])
	}
}

// Slice validation
type Emails struct {
	Emails []string `validate:"required,dive,email"`
}

func TestValidateStruct_SliceValidation(t *testing.T) {
	e := Emails{Emails: []string{"good@email.com", "bad-email"}}
	errors := ValidateStruct(e)
	if len(errors) == 0 {
		t.Errorf("expected errors for invalid email in slice")
	}
	found := false
	for k, v := range errors {
		if v.Validator == "email" && k == "emails[1]" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected email validation error for emails.1, got %v", errors)
	}
}

// Custom validator example
func customFoo(fl validator.FieldLevel) bool {
	return fl.Field().String() == "foo"
}

type CustomTagStruct struct {
	Foo string `json:"foo" validate:"foo"`
}

func TestValidateStruct_CustomValidator(t *testing.T) {
	err := validate.RegisterValidation("foo", customFoo)
	if err != nil {
		t.Fatalf("failed to register custom validator: %v", err)
	}
	s := CustomTagStruct{Foo: "bar"}
	errors := ValidateStruct(s)
	if errors["foo"].Validator != "foo" {
		t.Errorf("expected custom validator 'foo', got %v", errors["foo"])
	}
}

// TestStruct for testing validation
type TestStruct struct {
	Name  string `json:"name" validate:"required,min=2,max=10"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age" validate:"min=18,max=100"`
}

func TestValidateStruct_ValidatorTagAndMessage(t *testing.T) {
	data := TestStruct{
		Name:  "J", // too short
		Email: "invalid-email",
		Age:   17, // too young
	}
	errors := ValidateStruct(data)
	if errors["name"].Validator != "min" {
		t.Errorf("expected validator 'min' for name, got '%s'", errors["name"].Validator)
	}
	if errors["name"].Message == "" {
		t.Errorf("expected non-empty message for name")
	}
	if errors["email"].Validator != "email" {
		t.Errorf("expected validator 'email' for email, got '%s'", errors["email"].Validator)
	}
	if errors["email"].Message == "" {
		t.Errorf("expected non-empty message for email")
	}
	if errors["age"].Validator != "min" {
		t.Errorf("expected validator 'min' for age, got '%s'", errors["age"].Validator)
	}
	if errors["age"].Message == "" {
		t.Errorf("expected non-empty message for age")
	}
}

func TestValidateStruct_ValidData(t *testing.T) {
	data := TestStruct{
		Name:  "John",
		Email: "john@example.com",
		Age:   25,
	}

	errors := ValidateStruct(data)
	if len(errors) != 0 {
		t.Errorf("expected no validation errors, got %d: %v", len(errors), errors)
	}
}

func TestValidateStruct_InvalidData(t *testing.T) {
	data := TestStruct{
		Name:  "J", // too short
		Email: "invalid-email",
		Age:   17, // too young
	}

	errors := ValidateStruct(data)
	if len(errors) != 3 {
		t.Errorf("expected 3 validation errors, got %d: %v", len(errors), errors)
	}

	// Check for specific error fields
	expectedFields := []string{"name", "email", "age"}
	for _, expected := range expectedFields {
		if _, ok := errors[expected]; !ok {
			t.Errorf("expected validation error for field %s, but not found in: %v", expected, errors)
		}
	}
}

// Edge case: struct with no validation tags
type NoValidationStruct struct {
	Field string
}

func TestValidateStruct_NoValidationTags(t *testing.T) {
	s := NoValidationStruct{Field: "anything"}
	errors := ValidateStruct(s)
	if len(errors) != 0 {
		t.Errorf("expected no errors, got %v", errors)
	}
}

// Edge case: struct with pointer fields
type PointerStruct struct {
	Name *string `json:"name" validate:"required"`
}

func TestValidateStruct_PointerFieldNil(t *testing.T) {
	s := PointerStruct{Name: nil}
	errors := ValidateStruct(s)
	if len(errors) != 1 || errors["name"].Validator != "required" {
		t.Errorf("expected required error for name, got %v", errors)
	}
}

func TestValidateStruct_PointerFieldNotNil(t *testing.T) {
	name := "ok"
	s := PointerStruct{Name: &name}
	errors := ValidateStruct(s)
	if len(errors) != 0 {
		t.Errorf("expected no errors, got %v", errors)
	}
}
