package world

import (
	"strings"
	"testing"
	"reflect"
	"sort"
	"github.com/pastet89/alien-invasion/utils"
	"github.com/pastet89/alien-invasion/city"
)

var sampleMapPath string = utils.GetRootPath() + "/world/testdata/map.txt"

func getDummyMap() *Map {
	m := &Map{MapFilePath: sampleMapPath}
	m.Cities = make(city.Cities)
	m.ReverseRoadsToCities = make(city.ReverseRoadsToCities)
	return m
}

func getExpectedParsedData(setMainCityWithRoadsOutFlags bool) (city.Cities, city.ReverseRoadsToCities) {
	expectedReversedRoads := make(city.ReverseRoadsToCities)
	expectedCities := make(city.Cities)
	for _, cityName := range []string{"Moscow", "Sofia", "Tokyo", "London", "Paris", "Beijing"} {
		destinations := make(city.DestinationCitiesDirections)
		expectedCities[cityName] = &city.City{
				DestinationCitiesDirections: destinations,
				IsMainCityWithRoadsOut: false,
		}
	}

	expectedCities["Moscow"].DestinationCitiesDirections["London"] = "west"
	expectedCities["Moscow"].DestinationCitiesDirections["Tokyo"] = "east"
	
	expectedCities["Sofia"].DestinationCitiesDirections["Paris"] = "west"
	expectedCities["Sofia"].DestinationCitiesDirections["Tokyo"] = "east"

	expectedCities["Tokyo"].DestinationCitiesDirections["Beijing"] = "east"

	if setMainCityWithRoadsOutFlags {
		expectedCities["Moscow"].IsMainCityWithRoadsOut = true
		expectedCities["Sofia"].IsMainCityWithRoadsOut = true
		expectedCities["Tokyo"].IsMainCityWithRoadsOut = true

	}

	expectedReversedRoads["Tokyo"] = append(expectedReversedRoads["Tokyo"], "Moscow")
	expectedReversedRoads["Tokyo"] = append(expectedReversedRoads["Tokyo"], "Sofia")


	expectedReversedRoads["London"] = append(expectedReversedRoads["London"], "Moscow")
	expectedReversedRoads["Paris"] = append(expectedReversedRoads["Paris"], "Sofia")
	expectedReversedRoads["Beijing"] = append(expectedReversedRoads["Beijing"], "Tokyo")

	return expectedCities, expectedReversedRoads
}

func confirmParsedData(testName string, m *Map, t *testing.T, confirmWithFlags bool) {
	expectedCities, expectedReversedRoads := getExpectedParsedData(confirmWithFlags)

	for startCity, _ := range m.ReverseRoadsToCities {
		sort.Strings(m.ReverseRoadsToCities[startCity])
	}

	if  !reflect.DeepEqual(m.Cities, expectedCities) {
		t.Errorf("Wrong Cities data in test in %s", testName)
	}

	if  !reflect.DeepEqual(m.ReverseRoadsToCities, expectedReversedRoads){
		t.Errorf("Wrong Reversed Roads to Cities data in %s", testName)
	}
}

func TestParseCitiesRoadsData(t *testing.T) {
	m := getDummyMap()

	testData := make(map[string][]string)
	testData["Moscow"] = []string{"west=London", "east=Tokyo"}
	testData["Sofia"] = []string{"west=Paris", "east=Tokyo"}
	testData["Tokyo"] = []string{"east=Beijing"}
	for cityName, cityRoadsData := range testData {
		m.Cities[cityName] = &city.City{
			DestinationCitiesDirections: make(city.DestinationCitiesDirections),
		}
		m.parseCitiesRoadsData(cityName, cityRoadsData)
	}
	confirmWithFlags := false
	confirmParsedData("TestParseCitiesRoadsData", m, t, confirmWithFlags)
}

func TestParseCities(t *testing.T) {
	m := getDummyMap()
	testData := []string{
		"Moscow west=London east=Tokyo",
		"Sofia west=Paris east=Tokyo",
		"Tokyo east=Beijing",
	}
	m.parseCities(testData)
	confirmWithFlags := true
	confirmParsedData("TestParseCities", m, t, confirmWithFlags)
}

func TestParseMap(t *testing.T) {
	m := getDummyMap()
	m.ParseMap()
	confirmWithFlags := true
	confirmParsedData("TestParseMap", m, t, confirmWithFlags)
}

func TestEncodeMap(t *testing.T) {
	m := getDummyMap()
	m.ParseMap()
	encoded := m.getEncodedMap()
	m.Cities = make(city.Cities)
	m.ReverseRoadsToCities = make(city.ReverseRoadsToCities)
	m.parseCities(strings.Split(encoded, "\n"))
	confirmWithFlags := true
	confirmParsedData("TestEncodeMap", m, t, confirmWithFlags)
}