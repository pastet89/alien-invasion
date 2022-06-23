package city

import "testing"

func TestGetRandomDestination(t *testing.T) {
	dcd := make(DestinationCitiesDirections)
	cities := []string{"London", "Nairobi", "Tokyo", "Moscow"}
	for _, city := range cities {
		dcd[city] = "west"
	}

	destination := dcd.GetRandomDestination()
	existingCity := false
	for _, city := range cities {
		if city == destination {
			existingCity = true
		}
	}
	if !existingCity {
		t.Error("GetRandomDestination does not return existing city")
	}

}
