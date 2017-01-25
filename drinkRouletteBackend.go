package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	ID     int      `json:"ID"`
	Name   string   `json:"Name"`
	Cat    Category `json:"Category"`
	Color  string   `json:"Color"`
	Brands []string `json:"Brands"`
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
	json.NewEncoder(w).Encode(getAllIngredients())
}

func makeDrinkFromList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	idString := ps.ByName("ingredients")
	ids := strings.Split(idString, "-")

	allIngredients := getAllIngredients()
	ingredients := ingredientsForDrink(ids, makeIDMap(allIngredients))
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

func getAllIngredients() []Ingredient {
	goPath := os.Getenv("GOPATH")
	println(goPath)
	file, err := ioutil.ReadFile(goPath + "/src/github.com/bsomes/trashed-panda-backend/inputs/classified-ingredients.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	var ings []Ingredient
	marshalError := json.Unmarshal(file, &ings)
	if marshalError != nil {
		println(marshalError.Error())
	}
	return ings
}

func makeIDMap(ings []Ingredient) map[int]Ingredient {
	result := make(map[int]Ingredient, len(ings))
	for _, v := range ings {
		result[v.ID] = v
	}
	return result
}

func ingredientsForDrink(ids []string, available map[int]Ingredient) []Ingredient {
	ingredients := make([]Ingredient, len(ids))
	for ind, v := range ids {
		if i, err := strconv.Atoi(v); err == nil {
			ingredients[ind] = available[i]
		}
	}
	return ingredients
}

func uniform(ingredients []Ingredient) []Proportion {
	proportions := make([]Proportion, len(ingredients))
	for i := range proportions {
		proportions[i] = Proportion{Ing: ingredients[i], Frac: 1 / float32(len(ingredients))}
	}
	return proportions
}
