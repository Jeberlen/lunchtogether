package hojdencrawler

import (
	"testing"

	"github.com/Jeberlen/lunchtogether/menu_items"
)

func TestKeepEveryThird(t *testing.T) {
	list := []menu_items.MenuItem{
		{
			ID:          "1",
			DayOfWeek:   "1",
			Type:        "meat",
			Name:        "Pannbiff 1",
			Description: "A great pannbiff",
			URL:         "localhost",
		},
		{
			ID:          "2",
			DayOfWeek:   "1",
			Type:        "meat",
			Name:        "Pannbiff 2",
			Description: "A great vegan pannbiff",
			URL:         "localhost",
		},
		{
			ID:          "3",
			DayOfWeek:   "1",
			Type:        "meat",
			Name:        "Pannbiff 3",
			Description: "A great pannbiff",
			URL:         "localhost",
		},
	}
	want := []menu_items.MenuItem{
		{
			ID:          "1",
			DayOfWeek:   "1",
			Type:        "meat",
			Name:        "Pannbiff 1",
			Description: "A great pannbiff",
			URL:         "localhost",
		},
	}
	actual := KeepEveryThird(list)

	if len(actual) != len(want) {
		t.Fatalf("Failed: Not same length - actual: %d | want: %d", len(actual), len(want))
	}
}

func TestGetDayOfWeek_Meat(t *testing.T) {
	text := "MEAT"
	want := true
	actual := HasTypePrefix(text)

	if actual != want {
		t.Fatalf(
			"Failed: Did not eval properly - actual: %t | want: %t",
			actual,
			want,
		)
	}
}

func TestGetDayOfWeek_Failure(t *testing.T) {
	text := "TEST"
	want := false
	actual := HasTypePrefix(text)

	if actual != want {
		t.Fatalf(
			"Failed: Did not eval properly - actual: %t | want: %t",
			actual,
			want,
		)
	}
}
