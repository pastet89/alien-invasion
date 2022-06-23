# Alien Invasion

![image](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)

## Overview
Alien Invasion is a backend game simulating the invasion of aliens invading the Earth. 
The game starts by passing two CLI arguments: a `N` number of aliens and a path to a file with a world map. 

Each line of the map file contains a city leading to other cities. Each city can lead to 1-4 other cities. The format is as follows:
```
Foo north=Bar west=Baz south=Qu-ux
Bar south=Foo west=Bee
```

The main city is the first city on each line and is separated by a whitespace by the destination cities. Each destination city and the direction leading to it are separated by `=` and each destination city-direction pair is separated by a whitespace as well. 


The aliens are distributed in random cities around and once the game starts they move around the world by visiting other cities. They move by choosing a random destination going out of their city. If two aliens arrive at the same city, they fight and kill each other. In this case the city is destroyed and deleted from the map. In this case the game generates a destroy message. The roads from other cities leading to it are also destroyed. That can lead to certain aliens getting trapped in isolated cities with no exit.

The game stops when any of the following conditions is met:

- all remaining aliens are trapped
- each alien has moved at least 10,000 times
- there are no more aliens left

At the end of the game, the remaining city map is printed out in the same format it had been received initially.

## Requirements

- Go `1.16+`
- [Gonfig](https://github.com/creamdog/gonfig)

## Installation

- Download the project:
`git clone https://github.com/pastet89/alien-invasion.git`
- Install the dependencies:
`make deps`
- Run the tests:
`make test`
- Build:
`make build`

## Usage
From the project root folder run the build file corresponding to your OS with the following parameters:

`./bin/{BINARY_NAME} N path_to_map`

where:
- `N` is the number of aliens
- `path_to_map` is the path to the file containing the world map

You can also modify various game parameters such as the alien name prefix from the `config.json` file.

## Assumptions

- The city names can contain only latin alphabet characters, hyphen `-`, and numbers
- The city names have length `>= 2` and `<= 20`
- A direction going out of a city can not lead to the same city
- Each line of the map file represents a unique city
- Each city has between `1` and `4` roads going out of it
- The directions leading to the cities can have only 4 values: `west`, `east`, `north`, `south`
- The alien count is `>= 2`
- The city count is `>= 4` and `>=2 * ` the alien count

## Game structure

The game has three main modules: `simulator`, `city` and `world`.
The `city` module defines the city types and functionality.

The `world` module manages the parsing, encoding and functionality of the map.

The `simulator` module is the main game engine and commands the alien and cities actions.

There is no `alien` module because aliens are tracked just with their string names in `Simulator`. Also, aliens cannot perform actions on their own - they move only within the scope of the simulation dependent on other game objects - cities.

The `main()` function simply constructs and runs `Simulator`.
`Simulator.Run()` executes 10,000 loops in which each alien tries to travel and if it can travel to another city, performs the actions mentioned above. If before the 10,000th loop there are either no remaining aliens or no aliens able to travel, the game is ended.

## Algorithmic and data structures considerations

**World map**

The world map is parsed into a nested map `Cities map[string]*City` with city name keys pointing to `City` structs with the following fields:

```
type City struct {
    DestinationCitiesDirections DestinationCitiesDirections
    AliensInCity []string
    IsMainCityWithRoadsOut bool
}
```

The `City` structs contain data about the roads going out of the city and the aliens present in the city. By storing each `City` with its name as a key in the `Cities` map the following time complexity is achieved:

- `O(1)` deletion of a city
- `O(1)` access to the roads going out of the city
- `O(1)` access to the aliens present in the city

When deleting a city, the roads to the city must be deleted as well. If only the `Cities` map was present, this would mean that all cities should have been iterated and checked for remaining roads to the deleted city, which represents a time complexity of `O(n)`.

Theoretically each pair of aliens can contribute to a city being destroyed. This means that for each alien pair death a full cities iteration loop may be needed. That would result in a `O(m * n)` which while not being a quadratic solution may still come closer to quadratic than to a linear solution. If the aliens and cities numbers are high, that might significantly slow down the program execution.

For this reason, a reverse path to each city is kept from the cities to which it leads:
```
ReverseRoadsToCities map[string][]string
```
where the map key is the destination city name pointing to a list of all cities leading to this destination city.

E.g. If London and Paris lead to Berlin, in a separate map data structure Berlin will be stored as a key pointing to these two cities:
```
package city

type ReverseRoadsToCities map[string][]string
```
```
package world

ReverseRoadsToCities city.ReverseRoadsToCities
```
```
ReverseRoadsToCities["Berlin"] = []string{"London", "Paris"}
```
In this case the deletion of all remaining roads leading to Berlin becomes `O(1)` as well. A `O(1)` check provides that London and Paris lead to Berlin. And a `O(1)` deletion of `Cities["London"].DestinationCitiesDirections["Berlin"]` and `Cities["Paris"].DestinationCitiesDirections["Berlin"]` finishes the process.

**Isolated cities**

When parsing the map there will be isolated destination cities:
```
Foo north=Bar west=Baz south=Qu-ux
Bar south=Foo west=Bee
```
`Qu-ux` is an example for such a city.

Aliens can arrive in such cities, fight and get trapped. These actions must be tracked and these cities are also stored with their names as keys in the `Cities` map.
But when encoding back the remaining map to a string at the end of the game these cities must not be encoded as main cities with roads out - in order to keep the original map format. For these reasons a  `IsMainCityWithRoadsOut` flag inside each `City` struct is used to store and check for this when encoding and decoding the map.

**Aliens and cities tracking**

The aliens are stored in a map data structure `Simulator.aliensNamesLocations` with their names pointing to their city name locations.

Whenever an alien is to be transferred to a new city, a `O(1)` check is done  to find out how many aliens are there:
```
len(Cities[City].AliensInCity)
```
If there are no aliens, the alien is added to the `Cities[City].AliensInCity` slice.

If there is already another alien in the city, both aliens are deleted from the `aliensNamesLocations` map which is `O(1)` operation. The city and the roads to it are also each deleted with `O(1)` operations using the `ReverseRoadsToCities` and `Cities` data structures as discussed above.  

This way, the total game time complexity is kept down to approximately `10,000 * n` where `n` is the number of aliens.

**Generating random cities**

One disadvantage of maps as a data structure, however, is the manner in which a random item can be selected from them.

Golang does not iterate over the map items in a predefined sequence and each iteration can happen in a different order. For this reason, [some people](https://stackoverflow.com/questions/23482786/get-an-arbitrary-key-item-from-a-map) use a hacky way to get a `O(1)` random key from a map:

```
func get_some_key(m map[int]int) int {
    for k := range m {
        return k
    }
    return 0
}
```

This, however, should not be expected to provide randomness as from a PRNG. The fact that the iteration does not happen in the same order each time does not mean that it happens in a random order.

The real random map key algorithm has an `O(n)` time complexity. It selects a number `n` in the range of the size of the map and then iterates `n` times before returning the corresponding key.

For the random selection of a direction out of the city this is not a problem: there the range of options is between `1` and `4`. That is negligible and in practice is close to `O(1)`.

During the initial alien distribution, however, for each `n` alien a random city must be selected. As discussed above, if using the common random map key generator algorithm this may still result in a close to quadratic solution. For this reason, during the initial alien distribution, a slice containing all cities is created. This allows getting a random city with `O(1)` time complexity. This is done by generating a random index number and accessing the slice element by its index. That list is created within the inner scope of a function and released from the memory by the garbage collector as soon as the alien distribution is done.

## License

[MIT](https://github.com/pastet89/alien-invasion/blob/main/LICENSE.md)