package utils

import "testing"

func TestContains(t *testing.T) {
	slice := []string{"1", "2", "3", "4"}
	str := "3"
	want := true
	actual := Contains(slice, str)

	if actual != want {
		t.Fatal("Failed: actual should be true")
	}
}

func TestGetFoodTypeFromName_Meat(t *testing.T) {
	name := "good pork"
	want := "meat"
	actualType := GetFoodTypeFromName(name)
	if want != actualType {
		t.Fatalf(
			"Failed: actual (%s) not equal to want (%s)",
			actualType,
			want,
		)
	}
}

func TestGetFoodTypeFromName_Fish(t *testing.T) {
	name := "catched fish"
	want := "fish"
	actualType := GetFoodTypeFromName(name)
	if want != actualType {
		t.Fatalf(
			"Failed: actual (%s) not equal to want (%s)",
			actualType,
			want,
		)
	}
}

func TestGetFoodTypeFromName_Vegetarian(t *testing.T) {
	name := "Vegan beef patty"
	want := "vegetarian"
	actualType := GetFoodTypeFromName(name)
	if want != actualType {
		t.Fatalf(
			"Failed: actual (%s) not equal to want (%s)",
			actualType,
			want,
		)
	}
}

func TestGetFoodTypeFromName_Other(t *testing.T) {
	name := "This is not food"
	want := "other"
	actualType := GetFoodTypeFromName(name)
	if want != actualType {
		t.Fatalf(
			"Failed: actual (%s) not equal to want (%s)",
			actualType,
			want,
		)
	}
}

func TestGetFoodTypeFromDescription_Meat(t *testing.T) {
	description := "Salad, peas, ham, onion"
	want := "meat"
	actualType := GetFoodTypeFromDescription(description)
	if want != actualType {
		t.Fatalf(
			"Failed: actual (%s) not equal to want (%s)",
			actualType,
			want,
		)
	}
}

func TestGetFoodTypeFromDescription_Fish(t *testing.T) {
	description := "fresh cod caught, brussle sprouts"
	want := "fish"
	actualType := GetFoodTypeFromDescription(description)
	if want != actualType {
		t.Fatalf(
			"Failed: actual (%s) not equal to want (%s)",
			actualType,
			want,
		)
	}
}

func TestGetFoodTypeFromDescription_Vegetarian(t *testing.T) {
	description := "onion, garlic, vegan patty"
	want := "vegetarian"
	actualType := GetFoodTypeFromDescription(description)
	if want != actualType {
		t.Fatalf(
			"Failed: actual (%s) not equal to want (%s)",
			actualType,
			want,
		)
	}
}

func TestGetFoodTypeFromDescription_Other(t *testing.T) {
	description := "This is not food"
	want := "other"
	actualType := GetFoodTypeFromName(description)
	if want != actualType {
		t.Fatalf(
			"Failed: actual (%s) not equal to want (%s)",
			actualType,
			want,
		)
	}
}
