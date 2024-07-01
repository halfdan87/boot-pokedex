package pokeapi

import (
	"io"
	"errors"
    "fmt"
	"net/http"
    "encoding/json"
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

func GetLocations(url *string) ([]string, *Pagination, error) {
  res, err := http.Get(*url)
  if err != nil {
    return nil, nil, err 
  }

  body, err := io.ReadAll(res.Body)
  res.Body.Close()
  if res.StatusCode > 299 {
    return nil, nil, errors.New(fmt.Sprintf("Response failed: %d", res.StatusCode))
  }
  if err != nil {
    return nil, nil, err
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
