package simulator

import (
    "os"
    "strconv"
    "testing"
    "reflect"
    "sort"
    "github.com/pastet89/alien-invasion/utils"
    "github.com/pastet89/alien-invasion/city"
)


func (s *Simulator) resetData() {
	s.Map.Cities = make(city.Cities)
	s.aliensNamesLocations = make(aliensNamesLocations)
}

func muteOutput() {
	os.Stdout,_ = os.Open(os.DevNull)
}

func getSampleMapPath() string {
	var sampleMapPath string = utils.GetRootPath() + "/world/testdata/map.txt"
	return sampleMapPath
}

func getDummySimulator() *Simulator {
	muteOutput()
	path := getSampleMapPath()
	s:= ConstructSimulator(5, path)
	s.resetData()
	return s
}

func TestConstructSimulator(t *testing.T) {
	s := getDummySimulator()
	if  s.configVars.aliensCount != 5 {
		t.Error("Wrong alien count")
	}
	if  s.Map.MapFilePath != getSampleMapPath() {
		t.Error("Wrong map file path")
	}
}

func TestAlienIsTrappedDueToNotExistingMainCity(t *testing.T) {
	s := getDummySimulator()
	alienIsTrappedDueToNotExistingMainCity := s.alienIsTrapped("SomeCrazyNonExistingCity")
	if  !alienIsTrappedDueToNotExistingMainCity {
		t.Error("Not detecting trapped alien in TestAlienIsTrappedDueToNotExistingMainCity")
	}
}

func TestAlienIsTrappedDueToNoCityExits(t *testing.T) {
	s := getDummySimulator()
	s.Map.Cities["Tokyo"] = &city.City{
			DestinationCitiesDirections: make(city.DestinationCitiesDirections),
	}
	alienIsTrappedDueToNotExistingMainCity := s.alienIsTrapped("Tokyo")
	if  !alienIsTrappedDueToNotExistingMainCity {
		t.Error("Not detecting trapped alien in TestAlienIsTrappedDueToNoCityExits")
	}
}

func TestCreateAliens(t *testing.T) {
	s := getDummySimulator()
	config := utils.GetConfig()
	prefix := utils.GetStringConfigVar(config, "game/alien_name_prefix")
	s.configVars.aliensCount = 11
	s.createAliens()
	expected := make(aliensNamesLocations)
	for i := 1; i < s.configVars.aliensCount + 1; i++ {
		alienName := prefix + strconv.Itoa(i)
		expected[alienName] = ""
	}

	if  !reflect.DeepEqual(expected, s.aliensNamesLocations) {
		t.Error("Error setting init aliens")
	}
}

func TestGetCitiesNamesSlice(t *testing.T) {
	s := getDummySimulator()

	s.Map.Cities["Sofia"] = &city.City{
			DestinationCitiesDirections: make(city.DestinationCitiesDirections),
	}
	s.Map.Cities["Moscow"] = &city.City{
			DestinationCitiesDirections: make(city.DestinationCitiesDirections),
	}
	result := s.getCitiesNamesSlice()
	expected := []string{"Moscow", "Sofia"}
	sort.Strings(result)
	if  !reflect.DeepEqual(result, expected)  {
		t.Error("Error generating a slice with city names")
	}

}

func TestInitiallyDistributeAliensDestroyingACity(t *testing.T) {
	s := getDummySimulator()

	s.Map.Cities["Sofia"] = &city.City{
			DestinationCitiesDirections: make(city.DestinationCitiesDirections),
	}
	s.Map.Cities["Sofia"].DestinationCitiesDirections["west"] = "Belgrade"

	s.aliensNamesLocations["X1"] = ""
	s.aliensNamesLocations["X2"] = ""

	s.initiallyDistributeAliens()

	if  len(s.Map.Cities) > 0 || len(s.aliensNamesLocations) > 0 {
		t.Error("Cities and aliens not destroyed as expected during initial alien distribution")
	}

}


func TestInitiallyDistributeAliensAllocatesARealCity(t *testing.T) {
	s := getDummySimulator()

	s.Map.Cities["Sofia"] = &city.City{
			DestinationCitiesDirections: make(city.DestinationCitiesDirections),
	}
	s.Map.Cities["Sofia"].DestinationCitiesDirections["west"] = "Belgrade"

	s.aliensNamesLocations["X1"] = ""

	s.initiallyDistributeAliens()

	expectedAliensNamesLocations := make(aliensNamesLocations)
	expectedAliensNamesLocations["X1"] = "Sofia"

	if  !reflect.DeepEqual(expectedAliensNamesLocations, s.aliensNamesLocations)  {
		t.Error("Error in aliensNamesLocations data during initial alien distribution")
	}
}

func TestProcessAlienTravelResultingInFight(t *testing.T) {
	s := getDummySimulator()

	s.Map.Cities["Belgrade"] = &city.City{
			DestinationCitiesDirections: make(city.DestinationCitiesDirections),
	}
	s.Map.Cities["Belgrade"].DestinationCitiesDirections["west"] = "Munich"
	s.Map.Cities["Belgrade"].AliensInCity = append(s.Map.Cities["Belgrade"].AliensInCity, "X1")

	s.aliensNamesLocations["X1"] = "Belgrade"
	s.aliensNamesLocations["X2"] = "Earth"

	s.processAlienTravel("X2", "Belgrade")
	if  len(s.aliensNamesLocations) > 0 {
		t.Error("Aliens not destroyed during travel expected to result in fight")
	}
	if  len(s.Map.Cities) > 0 {
		t.Error("Cities not destroyed during travel expected to result in fight")
	}
}


func TestProcessAlienTravelResultingInPeace(t *testing.T) {
	s := getDummySimulator()

	s.Map.Cities["Belgrade"] = &city.City{
			DestinationCitiesDirections: make(city.DestinationCitiesDirections),
	}
	s.Map.Cities["Belgrade"].DestinationCitiesDirections["west"] = "Munich"

	s.Map.Cities["Moon"] = &city.City{
			DestinationCitiesDirections: make(city.DestinationCitiesDirections),
	}
	s.Map.Cities["Moon"].DestinationCitiesDirections["west"] = "Jupiter"

	s.Map.Cities["Belgrade"].AliensInCity = append(s.Map.Cities["Belgrade"].AliensInCity, "X1")

	s.aliensNamesLocations["X1"] = "Belgrade"
	s.aliensNamesLocations["X2"] = "Earth"

	s.processAlienTravel("X2", "Moon")

	expectedAliensNamesLocations := make(aliensNamesLocations)
	expectedAliensNamesLocations["X1"] = "Belgrade"
	expectedAliensNamesLocations["X2"] = "Moon"

	expectedCities := make(city.Cities)
	expectedCities["Belgrade"] = &city.City{
			DestinationCitiesDirections: make(city.DestinationCitiesDirections),
	}
	expectedCities["Belgrade"].DestinationCitiesDirections["west"] = "Munich"
	expectedCities["Belgrade"].AliensInCity = append(expectedCities["Belgrade"].AliensInCity, "X1")

	expectedCities["Moon"] = &city.City{
			DestinationCitiesDirections: make(city.DestinationCitiesDirections),
	}
	expectedCities["Moon"].DestinationCitiesDirections["west"] = "Jupiter"
	expectedCities["Moon"].AliensInCity = append(expectedCities["Moon"].AliensInCity, "X2")

	if  !reflect.DeepEqual(expectedAliensNamesLocations, s.aliensNamesLocations) {
		t.Error("Unexpected s.aliensNamesLocations during an Alien Travel not resulting in a fight")
	}
	if  !reflect.DeepEqual(expectedCities, s.Map.Cities) {
		t.Error("Unexpected s.Map.Cities during an Alien Travel not resulting in a fight")
	}
}

func TestProcessAlienFight(t *testing.T) {
	s := getDummySimulator()

	s.Map.Cities["Belgrade"] = &city.City{
			DestinationCitiesDirections: make(city.DestinationCitiesDirections),
	}
	s.Map.Cities["Belgrade"].DestinationCitiesDirections["west"] = "Munich"

	s.Map.Cities["Belgrade"].AliensInCity = append(s.Map.Cities["Belgrade"].AliensInCity, "X1")

	s.aliensNamesLocations["X1"] = "Belgrade"
	s.aliensNamesLocations["X2"] = "Earth"

	s.processAlienFight("X2", "Belgrade")

	if  len(s.aliensNamesLocations) > 0 {
		t.Error("Aliens not destroyed during alien fight")
	}
	if  len(s.Map.Cities) > 0 {
		t.Error("Cities not destroyed during alien fight")
	}
}

func TestKillAlien(t *testing.T) {
	s := getDummySimulator()

	s.aliensNamesLocations["X2"] = "Earth"
	s.killAlien("X2")

	if len(s.aliensNamesLocations) > 0 {
		t.Error("Unable to kill alien")
	}
}


func TestProcessAliensDeath(t *testing.T) {
	s := getDummySimulator()

	s.Map.Cities["Belgrade"] = &city.City{
			DestinationCitiesDirections: make(city.DestinationCitiesDirections),
	}
	s.Map.Cities["Belgrade"].AliensInCity = append(s.Map.Cities["Belgrade"].AliensInCity, "X1")

	s.aliensNamesLocations["X1"] = "Belgrade"
	s.aliensNamesLocations["X2"] = "Earth"

	s.processAliensDeath("X2", "Belgrade")

	if  len(s.aliensNamesLocations) > 0 {
		t.Error("Aliens not destroyed during alien death")
	}
}

func TestGetRandomCityReturnsExistingCity(t *testing.T) {
	s := getDummySimulator()
	s.Map.Cities["Belgrade"] = &city.City{
			DestinationCitiesDirections: make(city.DestinationCitiesDirections),
	}
	randomCity := s.getRandomCityName([]string{"Belgrade"})
	if  randomCity != "Belgrade" {
		t.Error("Not returning expected city from getRandomCityName")
	}
}

func TestGetRandomCityReturnsRandomCity(t *testing.T) {
	s := getDummySimulator()
	var slice []string
	for i := 0; i <1000 ; i++ {
		name := "Belgrade" + strconv.Itoa(i)
		s.Map.Cities[name] = &city.City{
			DestinationCitiesDirections: make(city.DestinationCitiesDirections),
		}
		slice = append(slice, name)
	}

	uniqueResults := make(map[string]string)
	
	for i := 0; i < 1000; i++ {
		randomCity := s.getRandomCityName(slice)
		uniqueResults[randomCity] = ""
	}
	
	if  len(uniqueResults) < 2 {
		t.Error("getRandomCityName returning same values each time!")
	}
}