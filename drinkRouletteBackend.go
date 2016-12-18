package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
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
	ID   int      `json:"ID"`
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

func getEncodedIngredients(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	json.NewEncoder(w).Encode(getIngredients())
}

func makeDrinkFromList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idString := ps.ByName("ingredients")
	ids := strings.Split(idString, "-")
	ingredients := make([]Ingredient, len(ids))
	allIngredients := getIngredients()
	for ind, v := range ids {
		if i, err := strconv.Atoi(v); err == nil {
			ingredients[ind] = allIngredients[i]
		}
	}
	json.NewEncoder(w).Encode(Drink{
		Name:     "Test",
		Contents: uniform(ingredients),
		Size:     Buttchug,
	})
}

func main() {
	port := os.Getenv("PORT")

	router := httprouter.New()
	router.GET("/ingredients", getEncodedIngredients)
	router.GET("/makedrink/:ingredients", makeDrinkFromList)

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func getIngredients() []Ingredient {
	return []Ingredient{
		Ingredient{
			ID:   0,
			Name: "Vodka",
			Cat:  Alcohol,
		},
		Ingredient{
			ID:   1,
			Name: "Milk",
			Cat:  Mixer,
		},
		Ingredient{
			ID:   2,
			Name: "Kahlua",
			Cat:  Alcohol,
		},
		Ingredient{
			ID:   3,
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
