package city

import (
		"time"
		"math/rand"
)

type City struct {
    DestinationCitiesDirections DestinationCitiesDirections
    AliensInCity []string
    IsMainCityWithRoadsOut bool
}

type DestinationCitiesDirections map[string]string

type Cities map[string]*City

type ReverseRoadsToCities map[string][]string

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func (dcd DestinationCitiesDirections) GetRandomDestination() string {
	randomI := rand.Intn(len(dcd))
	i := 0
	for destination, _ := range dcd {
		if i == randomI {
			return destination
		}
		i++
	}
	return ""
}

func (c Cities) CityHasRoadsOut(cityName string) bool {
	cityObj, ok := c[cityName]
	return ok && len(cityObj.DestinationCitiesDirections) > 0
}

func (c Cities) getCityFromCities(cityName string) *City {
	cityObj, _ := c[cityName]
	return cityObj
}

func (c Cities) GetAllDestinations(cityName string) DestinationCitiesDirections {
	var result = make(DestinationCitiesDirections)
	cityObj := c.getCityFromCities(cityName)
	if cityObj != nil {
		return cityObj.DestinationCitiesDirections
	}
	return result
}

func (c Cities) GetAliensInCity(cityName string) []string {
	var result []string
	cityObj := c.getCityFromCities(cityName)
	if cityObj != nil {
		return cityObj.AliensInCity
	}
	return result
}