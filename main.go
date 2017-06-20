//The main package contains all logic for handling http requests to trashed-panda-backend
package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/stdlib"
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
	Frac float64    `json:"Fraction"`
}

//Drink The overall representation of a finished drink.
type Drink struct {
	Name     string       `json:"Name"`
	Contents []Proportion `json:"Contents"`
	Size     Scale        `json:"Size"`
}

// allows GET requests from all external URLs
func setDefaultHeader(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func getEncodedIngredients(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	setDefaultHeader(w)
	ingData := ingredientData{
		data: db,
	}
	json.NewEncoder(w).Encode(ingData.getAllIngredients())
}

func makeDrinkFromList(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var ids []int
	json.NewDecoder(r.Body).Decode(&ids)

	ingData := ingredientData{
		data: db,
	}
	ingredients := ingData.getAllIngredientsWithIDs(ids)
	drinkData := drinkDataSQL{
		db: db,
	}
	creator := drinkCreator{
		data: &drinkData,
	}
	drink := creator.makeDrink(ingredients)
	setDefaultHeader(w)
	json.NewEncoder(w).Encode(drink)
}

var db *sql.DB

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databaseConnection := os.Getenv("DATABASE_URL")
	if databaseConnection == "" {
		databaseConnection = "postgresql://localhost:5432/briansomes"
	}

	data, err := sql.Open("pgx", databaseConnection)
	if err != nil {
		log.Fatal(err)
	}
	db = data
	router := httprouter.New()
	router.GET("/ingredients", getEncodedIngredients)
	router.POST("/makedrink", makeDrinkFromList)

	log.Fatal(http.ListenAndServe(":"+port, router))
}

func uniform(ingredients []Ingredient) []Proportion {
	proportions := make([]Proportion, len(ingredients))
	for i := range proportions {
		proportions[i] = Proportion{Ing: ingredients[i], Frac: 1 / float64(len(ingredients))}
	}
	return proportions
}
