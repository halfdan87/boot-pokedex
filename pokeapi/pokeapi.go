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

  fmt.Println("End Unmarshal")

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
