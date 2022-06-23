package simulator

import (
	"fmt"
	"time"
	"math/rand"
	"github.com/pastet89/alien-invasion/world"
)

type aliensNamesLocations map[string]string

type Simulator struct {
    Map world.Map
    aliensNamesLocations aliensNamesLocations
    configVars configVars
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func ConstructSimulator(aliensCount int, mapFilePath string) *Simulator {
	s := &Simulator{
		Map: world.Map{MapFilePath: mapFilePath},
		aliensNamesLocations: make(aliensNamesLocations),
	}
	s.setMainConfigVars()
	s.configVars.aliensCount = aliensCount
	s.createAliens()
	s.Map.ParseMap()
	s.initiallyDistributeAliens()
	return s
}

func (s Simulator) Run() {
	for i := 0; i < s.configVars.maxTravelsPerAlien; i++ {
		aliensCanTravel := false
		for alienName, alienLocation := range s.aliensNamesLocations {
			if !s.alienIsTrapped(alienLocation) {
				aliensCanTravel = true
				destinations := s.Map.Cities.GetAllDestinations(alienLocation)
				destinationCity := destinations.GetRandomDestination()
				s.processAlienTravel(alienName, destinationCity)
			}
		}
		if !aliensCanTravel && len(s.aliensNamesLocations) > 0 {
			fmt.Println(s.configVars.allAliensAreTrappedMessage)
			break
		}
	}
	fmt.Println("========= THE END =========")
	s.Map.PrintMap()
}