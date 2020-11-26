package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// structs for each pokemon response from pokeapi
type responseKind struct {
	Types []types `json:"types"`
}
type types struct {
	Type typestruct `json:"type"`
}
type typestruct struct {
	Name string `json:"name"`
}

// data of the pokemons simplify from pokeapi
type responsePokemons struct {
	Results []results `json:"results"`
}
type result struct {
	Name string `json:"name"`
}

// final body response for pokedex
type bodyResponse struct {
	Result []results `json:"result"`
}
type results struct {
	Name string   `json:"name"`
	Type []string `json:"type"`
}

//global vars to pokeApi urls
var pokeAPIKind string = "https://pokeapi.co/api/v2/pokemon/"
var pokeAPIList string = "https://pokeapi.co/api/v2/pokemon?offset=0&limit=10"

// this methos is being used to handle the POST request for /pokemon
func yourHandlerTs(w http.ResponseWriter, r *http.Request) {
	pokemons, err := getPokemons()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	fmt.Fprintf(w, "%v\n", pokemons)
}

// main method to execute the flow on golang
func main() {
	r := mux.NewRouter()
	// Routes r is the way to route path or handle funcions using gorilla/Mux
	r.Path("/pokemon").HandlerFunc(yourHandlerTs).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", r))
}

func getPokemons() (string, error) {
	body, err := makeRequest(pokeAPIList, "")
	if err != nil {
		return "", err
	}
	var responsePoke responsePokemons
	json.Unmarshal(body, &responsePoke)
	result := &results{}
	bodyRs := &bodyResponse{}
	dataBytes := []byte{}
	for _, j := range responsePoke.Results {

		data, err := getKindPokemons(j.Name)
		if err != nil {
			return "", err
		}
		result = &results{Name: j.Name, Type: data}
		bodyRs.Result = append(bodyRs.Result, *result)
	}
	dataBytes, err = json.Marshal(bodyRs)
	if err != nil {
		return "", err
	}
	return string(dataBytes), nil
}

func getKindPokemons(id string) ([]string, error) {
	body, err := makeRequest(pokeAPIKind, id)
	if err != nil {
		return []string{}, err
	}
	var responseTypes responseKind
	json.Unmarshal(body, &responseTypes)
	types := []string{}
	for _, j := range responseTypes.Types {
		types = append(types, j.Type.Name)
	}
	return types, nil
}

// method to do request on pokeapi
func makeRequest(url string, id string) ([]byte, error) {
	resp, err := http.Get(url + id)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
