package main

import (
	"database/sql"
	"log"
	"strconv"
)

type ingredientData struct {
	data *sql.DB
}

func (i *ingredientData) getAllIngredients() []Ingredient {
	var (
		id       int
		name     string
		color    sql.NullString
		baseID   int
		category Category
		ing      Ingredient
	)
	ingredients := make([]Ingredient, 0)
	rows, err := i.data.Query("select i.id, i.name, i.color, i.baseid, b.category from ingredients i join baseingredients b on i.baseid = b.id")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name, &color, &baseID, &category)
		if err != nil {
			log.Fatal(err)
		}
		if !color.Valid {
			ing = Ingredient{
				ID:     id,
				Name:   name,
				Color:  "",
				BaseID: baseID,
				Cat:    category,
			}
		} else {
			ing = Ingredient{
				ID:     id,
				Name:   name,
				Color:  color.String,
				BaseID: baseID,
				Cat:    category,
			}
		}
		ingredients = append(ingredients, ing)
	}
	return ingredients
}

func (i *ingredientData) getAllIngredientsWithIDs(ids []int) []Ingredient {
	var (
		id       int
		name     string
		color    sql.NullString
		baseID   int
		category Category
	)
	statement := "select i.id, i.name, i.color, i.baseid, b.category from ingredients i join baseIngredients b on i.baseid = b.id  where i.id = any($1::integer[])"
	args := "{"
	for _, v := range ids {
		args += strconv.Itoa(v) + ","
	}
	args = args[:len(args)-1] + "}"
	ingredients := make([]Ingredient, 0)
	rows, err := i.data.Query(statement, args)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name, &color, &baseID, &category)
		if err != nil {
			log.Fatal(err)
		}
		if color.Valid {
			ingredients = append(ingredients, Ingredient{
				ID:     id,
				Name:   name,
				Color:  color.String,
				BaseID: baseID,
				Cat:    category,
			})
		} else {
			ingredients = append(ingredients, Ingredient{
				ID:     id,
				Name:   name,
				Color:  "",
				BaseID: baseID,
				Cat:    category,
			})
		}
	}
	return ingredients
}
