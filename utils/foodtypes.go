package utils

import "strings"

func GetMeatTypes() []string {
	return []string{"ham",
		"entrecôte", "turkey", "lamb",
		"veal", "rib", "fläsk", "pork",
		"beef", "chicken", "pannbiff",
		"bacon", "sirloin", "steak", "burger",
		"brisket", "kebab", "forno"}
}

func GetFishTypes() []string {
	return []string{
		"fish", "shrimp", "prawn", "clam",
		"clams", "crab", "salmon", "tuna",
		"tilapia", "cod", "snapper", "sardine",
		"sardines", "herring", "haddock",
		"flounder", "trout", "pollock", "bass",
		"halibut", "pike", "mackerel", "catch"}
}

func GetVegTypes() []string {
	return []string{"vegetarian", "vegan", "veggie", "halloumi", "pannoumi", "aranccini"}
}

func Contains(slice []string, str string) bool {

	for _, item := range slice {
		if item == str {
			return true
		}
	}

	return false
}

func GetFoodTypeFromName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ToLower(name)
	slicesOfName := strings.Split(name, " ")

	for _, slice := range slicesOfName {

		wordInName := strings.TrimSpace(slice)
		wordInName = strings.ToLower(wordInName)

		if Contains(GetFishTypes(), wordInName) {
			return "fish"
		} else if Contains(GetMeatTypes(), wordInName) {
			return "meat"
		} else if Contains(GetVegTypes(), wordInName) {
			return "vegetarian"
		}
	}
	return "other"
}

func GetFoodTypeFromDescription(description string) string {
	description = strings.ToLower(description)
	description = strings.ReplaceAll(description, ",", "")
	sliceOfDesc := strings.Split(description, " ")

	for _, slice := range sliceOfDesc {
		desc := strings.TrimSpace(slice)
		desc = strings.ToLower(desc)

		if Contains(GetMeatTypes(), desc) {
			return "meat"
		} else if Contains(GetFishTypes(), desc) {
			return "fish"
		} else if Contains(GetVegTypes(), desc) {
			return "vegetarian"
		}
	}
	return "other"
}
