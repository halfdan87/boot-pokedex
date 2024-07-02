package pokeapi

import (
	"io"
    "time"
	"errors"
    "fmt"
	"net/http"
    "encoding/json"
    "github.com/halfdan87/boot-pokedex/internal/pokecache"
)

type Pagination struct {
  Next *string
  Prev *string
}

type Response struct {
    Count int
    Next *string
    Previous *string
    Results []Location
}

type Location struct {
    Name string
    Url string
}

var cache *pokecache.Cache = pokecache.NewCache(2 * time.Minute)

func GetLocations(url *string) ([]string, *Pagination, error) {
  var body []byte
  var err error

  if item, ok := cache.Get(*url); ok {
    body = item
  } else {
    res, err := http.Get(*url)
    if err != nil {
        return nil, nil, err 
    }

    body, err = io.ReadAll(res.Body)
    if err != nil {
        return nil, nil, err
    }
    res.Body.Close()
    if res.StatusCode > 299 {
        return nil, nil, errors.New(fmt.Sprintf("Response failed: %d", res.StatusCode))
    }

    cache.Add(*url, body)
  }

  response := Response{}
  err = json.Unmarshal(body, &response)
  if err != nil {
    return nil, nil, err
  }

  locationStrings := []string{}

  for _, loc := range response.Results {
    locationStrings = append(locationStrings, loc.Name)
  }

  pager := Pagination{
      Next : response.Next,
      Prev : response.Previous,
  }

  return locationStrings, &pager, nil
}

type ResponseLocation struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

var cacheLocationData *pokecache.Cache = pokecache.NewCache(2 * time.Minute)

func GetPokemons(loc string) ([]string, error) {
  var body []byte
  var err error

  url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s-area/", loc)

  fmt.Println(url)

  if item, ok := cacheLocationData.Get(url); ok {
    body = item
  } else {
    res, err := http.Get(url)
    if err != nil {
        return nil, err 
    }

    body, err = io.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }
    res.Body.Close()
    if res.StatusCode > 299 {
        return nil, errors.New(fmt.Sprintf("Response failed: %d", res.StatusCode))
    }

    cacheLocationData.Add(url, body)
  }

  fmt.Println(string(body))

  response := ResponseLocation{}

  err = json.Unmarshal(body, &response)
  if err != nil {
    return nil, err
  }

      fmt.Println(response)
  pokemones := []string{}

  for _, pok := range response.PokemonEncounters{
    pokemones = append(pokemones, pok.Pokemon.Name)
  }

  return pokemones, nil
}

type ResponsePokemon struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Cries          struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	Height    int `json:"height"`
	ID                     int    `json:"id"`
	IsDefault              bool   `json:"is_default"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			VersionGroup struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name          string `json:"name"`
	Order         int    `json:"order"`
	PastAbilities []any  `json:"past_abilities"`
	PastTypes     []any  `json:"past_types"`
	Species       struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}



var cachePokemonData *pokecache.Cache = pokecache.NewCache(2 * time.Minute)

func GetPokemon(name string) (ResponsePokemon, error) {
  var body []byte
  var err error
  response := ResponsePokemon{}

  url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", name)

  if item, ok := cacheLocationData.Get(url); ok {
    body = item
  } else {
    res, err := http.Get(url)
    if err != nil {
        return response, err 
    }

    body, err = io.ReadAll(res.Body)
    if err != nil {
        return response, err
    }
    res.Body.Close()
    if res.StatusCode > 299 {
        return response, errors.New(fmt.Sprintf("Response failed: %d", res.StatusCode))
    }

    cacheLocationData.Add(url, body)
  }


  err = json.Unmarshal(body, &response)
  if err != nil {
    return response, err
  }


  return response, nil
}


