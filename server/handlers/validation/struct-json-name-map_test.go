package validation

import (
	"reflect"
	"testing"
)

// Test struct with nested fields and json tags
type Inner struct {
	FieldA string `json:"field_a"`
	FieldB int    `json:"field_b"`
}

type Outer struct {
	Name   string `json:"name"`
	Inner1 Inner  `json:"inner1"`
	Inner2 *Inner `json:"inner2"`
}

type Deep struct {
	OuterField Outer `json:"outer_field"`
}

func TestGetOrBuildFieldMap_Cache(t *testing.T) {
	typ := reflect.TypeOf(Outer{})
	// First call should build and cache
	m1 := GetOrBuildFieldMap(typ)
	if m1["Inner1.FieldA"] != "inner1.field_a" {
		t.Errorf("expected Inner1.FieldA -> inner1.field_a, got %s", m1["Inner1.FieldA"])
	}
	// Second call should hit cache (simulate by changing map and checking it persists)
	m1["test"] = "value"
	m2 := GetOrBuildFieldMap(typ)
	if m2["test"] != "value" {
		t.Errorf("expected cache to persist custom key, got %s", m2["test"])
	}
}

func TestBuildJSONFieldMap_NestedStructs(t *testing.T) {
	m := buildJSONFieldMap(reflect.TypeOf(Outer{}), "", "")
	expected := map[string]string{
		"Name":          "name",
		"Inner1":        "inner1",
		"Inner1.FieldA": "inner1.field_a",
		"Inner1.FieldB": "inner1.field_b",
		"Inner2":        "inner2",
		"Inner2.FieldA": "inner2.field_a",
		"Inner2.FieldB": "inner2.field_b",
	}
	for k, v := range expected {
		if m[k] != v {
			t.Errorf("expected %s -> %s, got %s", k, v, m[k])
		}
	}
}

func TestBuildJSONFieldMap_DeeplyNested(t *testing.T) {
	m := buildJSONFieldMap(reflect.TypeOf(Deep{}), "", "")
	expected := map[string]string{
		"OuterField":               "outer_field",
		"OuterField.Name":          "outer_field.name",
		"OuterField.Inner1":        "outer_field.inner1",
		"OuterField.Inner1.FieldA": "outer_field.inner1.field_a",
		"OuterField.Inner1.FieldB": "outer_field.inner1.field_b",
		"OuterField.Inner2":        "outer_field.inner2",
		"OuterField.Inner2.FieldA": "outer_field.inner2.field_a",
		"OuterField.Inner2.FieldB": "outer_field.inner2.field_b",
	}
	for k, v := range expected {
		if m[k] != v {
			t.Errorf("expected %s -> %s, got %s", k, v, m[k])
		}
	}
}
