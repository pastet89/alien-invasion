package simulator

import (
	"fmt"
	"strconv"
)

func (s *Simulator) createAliens() {
	s.aliensNamesLocations = make(aliensNamesLocations)
	for i := 1; i < s.configVars.aliensCount + 1; i++ {
		alienName, initLocation := s.generateAlienName(i), ""
		s.aliensNamesLocations[alienName] = initLocation
	}
}

func (s *Simulator) initiallyDistributeAliens() {
	citiesArray := s.getCitiesNamesSlice()
	for alienName, _ := range s.aliensNamesLocations {
		cityName := s.getRandomCityName(citiesArray)
		s.processAlienTravel(alienName, cityName)
	}
}

func (s Simulator) alienIsTrapped(alienLocation string) bool {
	return !s.Map.Cities.CityHasRoadsOut(alienLocation)
}

func (s Simulator) generateAlienName(id int) string {
	return s.configVars.alienNamePrefix + strconv.Itoa(id)
}

func (s *Simulator) processAlienTravel(alienName string, destinationCityName string) {
	aliensInCity := s.Map.Cities.GetAliensInCity(destinationCityName)
	aliensMustFight := len(aliensInCity) > 0
	if aliensMustFight {
		s.processAlienFight(alienName, destinationCityName)
	} else {
		s.addAlienToCity(alienName, destinationCityName)
		s.setAlienLocation(alienName, destinationCityName)
	}
}

func (s *Simulator) processAlienFight(lastArrivingAlien string, cityName string) {
	s.processAliensDeath(lastArrivingAlien, cityName)
	s.Map.DestroyCity(cityName)
}

func (s *Simulator) processAliensDeath(lastArrivingAlien string, cityName string) {
	aliensInCity := s.Map.Cities.GetAliensInCity(cityName)
	firstArrivingAlien := aliensInCity[0]
	s.killAlien(firstArrivingAlien)
	s.killAlien(lastArrivingAlien)
	fmt.Printf(s.configVars.fightFormatMsg, cityName, firstArrivingAlien, lastArrivingAlien)
}

func (s *Simulator) addAlienToCity(alienName string, cityName string) {
	s.Map.Cities[cityName].AliensInCity = append(s.Map.Cities[cityName].AliensInCity, alienName)
}

func (s *Simulator) setAlienLocation(alienName string, cityName string) {
	s.aliensNamesLocations[alienName] = cityName
}

func (s *Simulator) killAlien(alienName string) {
	delete(s.aliensNamesLocations, alienName)
}