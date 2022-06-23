package simulator

import "math/rand"

func (s Simulator) cityExistsOnMap(cityName string) bool {
	_, ok := s.Map.Cities[cityName]
	return ok
}

func (s *Simulator) getRandomCityName(citiesArray []string) string {
	for {
		randIndex := rand.Intn(len(citiesArray))
		cityName := citiesArray[randIndex]
		if s.cityExistsOnMap(cityName) {
			return cityName
		}
	}
}

func (s Simulator) getCitiesNamesSlice() []string {
	citiesArray := make([]string, len(s.Map.Cities))
	i := 0
	for cityName, _ := range(s.Map.Cities) {
		citiesArray[i] = cityName
		i++
	}
	return citiesArray
}