package world

import "github.com/pastet89/alien-invasion/city"


type Map struct {
    Cities city.Cities
    ReverseRoadsToCities city.ReverseRoadsToCities
    MapFilePath string
}

func (m *Map) DestroyCity(cityName string) {
	m.destroyRoadsToCity(cityName)
	delete(m.Cities, cityName)
}

func (m *Map) destroyRoadsToCity(cityName string) {
	for _, cityLeadingToDeletedCity := range m.ReverseRoadsToCities[cityName] {
		delete(m.Cities.GetAllDestinations(cityLeadingToDeletedCity), cityName)
	}
	delete(m.ReverseRoadsToCities, cityName)
}
