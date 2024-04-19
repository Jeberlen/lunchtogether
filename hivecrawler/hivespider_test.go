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
