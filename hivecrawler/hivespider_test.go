package hivecrawler

import (
	"testing"
)

func TestRemoveStrings(t *testing.T) {
	inputStrings := []string{"this", "is", "a", "test"}
	toRemove := "test"
	want := []string{"this", "is", "a"}

	actual := RemoveStrings(inputStrings, toRemove)

	if len(actual) != len(want) {
		t.Fatal("Failed: want and actual are different lenghts")
	}

	for i, str := range actual {
		if str != want[i] {
			t.Fatal("Failed: Incorrect values in actual")
		}
	}
}

func TestContains(t *testing.T) {
	slice := []string{"1", "2", "3", "4"}
	str := "3"
	want := true
	actual := Contains(slice, str)

	if actual != want {
		t.Fatal("Failed: actual should be true")
	}
}

func TestGetDayOfWeek_WeekDay(t *testing.T) {
	day := "MONDAY"
	wantedDay := "1"
	wantedUrl := "https://thehivefoodmarket.se/#monday"

	actualDay, actualUrl := GetDayOfWeek(day)

	if actualDay != wantedDay || actualUrl != wantedUrl {
		t.Fatalf(
			"Failed: Day or url is not equal %s =? %s || %s =? %s",
			actualDay,
			wantedDay,
			actualUrl,
			wantedUrl,
		)
	}
}

func TestGetDayOfWeek_WeekEnd(t *testing.T) {
	day := "SATURDAY"
	wantedDay := "1"
	wantedUrl := "https://thehivefoodmarket.se/"

	actualDay, actualUrl := GetDayOfWeek(day)

	if actualDay != wantedDay || actualUrl != wantedUrl {
		t.Fatalf(
			"Failed: Day or url is not equal %s =? %s || %s =? %s",
			actualDay,
			wantedDay,
			actualUrl,
			wantedUrl,
		)
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
