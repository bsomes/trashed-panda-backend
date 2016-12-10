package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Category int

const (
	Alcohol Category = iota
	Mixer
	Other
)

type Scale int

const (
	Shot Scale = iota
	Glass
	Cup
	Buttchug
)

type Ingredient struct {
	Name string   `json:"Name"`
	Cat  Category `json:"Category"`
}

type Proportion struct {
	Ing  Ingredient `json:"Ingredient"`
	Frac float32    `json:"Fraction"`
}

type Drink struct {
	Name     string       `json:"Name"`
	Contents []Proportion `json:"Contents"`
	Size     Scale        `json:"Size"`
}

func main() {
	/*http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})*/

	http.HandleFunc("/ingredients", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(getIngredients())
	})

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(getTestDrink())
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getIngredients() []Ingredient {
	return []Ingredient{
		Ingredient{
			Name: "Vodka",
			Cat:  Alcohol,
		},
		Ingredient{
			Name: "Milk",
			Cat:  Mixer,
		},
		Ingredient{
			Name: "Kahlua",
			Cat:  Alcohol,
		},
		Ingredient{
			Name: "Cocoa powder",
			Cat:  Other,
		},
	}
}

func uniform(ingredients []Ingredient) []Proportion {
	proportions := make([]Proportion, len(ingredients))
	for i := range proportions {
		proportions[i] = Proportion{Ing: ingredients[i], Frac: 1 / float32(len(ingredients))}
	}
	return proportions
}

func getTestDrink() Drink {
	return Drink{
		Name:     "Mulatto",
		Contents: uniform(getIngredients()),
		Size:     Buttchug,
	}
}
