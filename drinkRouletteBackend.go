//The main package contains all logic for handling http requests to trashed-panda-backend
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

//Category This corresponds to the category that a particular ingredient falls under.
type Category int

const (
	//Alcohol Anything that taking a shot of wouldn't be totally ridiculous
	Alcohol Category = iota
	//Mixer Any ingredient that wouldn't be totally ridiculous filling a 12 oz cup with
	Mixer
	//Other All other ingredients
	Other
)

//Scale The approximate size of the final drink
type Scale int

const (
	//Shot The size of a shot glass. Should be reserved for drinks that are mostly alcohol
	Shot Scale = iota
	//Glass Roughly the size of a small cocktail, like White Russian or Old Fashioned size
	Glass
	//Cup About a 12 oz cup full
	Cup
	//Buttchug For the adventrous.
	Buttchug
)

//Ingredient The fundamental building block of a drink.
type Ingredient struct {
	ID     int      `json:"ID"`
	Name   string   `json:"Name"`
	Cat    Category `json:"Category"`
	Color  string   `json:"Color"`
	BaseID int      `json:"BaseID"`
}

//Proportion This describes how much of a drink is commposed of this particular Ingredient
type Proportion struct {
	Ing  Ingredient `json:"Ingredient"`
	Frac float32    `json:"Fraction"`
}

//Drink The overall representation of a finished drink.
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
	path := http.Dir("./inputs")
	file, err := path.Open("classified-ingredients.json")
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	stats, statsErr := file.Stat()
	if statsErr != nil {
		return nil
	}
	var size = stats.Size()

	data := make([]byte, size)
	file.Read(data)
	var ings []Ingredient
	marshalError := json.Unmarshal(data, &ings)
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
