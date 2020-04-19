package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

type PokemonSpecies struct {
	BaseHappiness int `json:"base_happiness"`
	CaptureRate   int `json:"capture_rate"`
	Color         struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"color"`
	EggGroups []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"egg_groups"`
	EvolutionChain struct {
		URL string `json:"url"`
	} `json:"evolution_chain"`
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var ch = make(chan PokemonSpecies)

func fetchPokemon(id int) {
	// format link without printing
	pokemonLink := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon-species/%d", id)
	response, err := http.Get(pokemonLink)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	//retrieve response body
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// unmarshall response to PokemonSpecies struct
	ps := PokemonSpecies{}
	err = json.Unmarshal(responseData, &ps)
	if err != nil {
		fmt.Println(err)
	}

	defer response.Body.Close()

	//send PokemonSpecies to PokemonSpecies channel
	ch <- ps
}

func main() {
	//start a WaitGroup to wait for all goroutines to run
	var wg sync.WaitGroup
	wg.Add(20)
	for i := 1; i <= 10; i++ {
		go fetchPokemon(i)
		wg.Done()
	}

	for i := 1; i <= 10; i++ {
		go func() {
			ps := <-ch
			fmt.Printf("%+v\n", ps)
			wg.Done()
		}()
	}
	wg.Wait()
}
