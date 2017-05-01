package main

import (
	"database/sql"
	"log"
)

type ingredientData struct {
	data *sql.DB
}

func (i *ingredientData) getAllIngredients() []Ingredient {
	var (
		id     int
		name   string
		color  sql.NullString
		baseID int
		ing    Ingredient
	)
	ingredients := make([]Ingredient, 0)
	rows, err := i.data.Query("select id, name, color, baseid from ingredients i join baseingredients b on i.id = b.id")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name, &color, &baseID)
		if err != nil {
			log.Fatal(err)
		}
		if !color.Valid {
			ing = Ingredient{
				ID:     id,
				Name:   name,
				Color:  "",
				BaseID: baseID,
			}
		} else {
			ing = Ingredient{
				ID:     id,
				Name:   name,
				Color:  color.String,
				BaseID: baseID,
			}
		}
		ingredients = append(ingredients, ing)
	}
	return ingredients
}

func (i *ingredientData) getAllIngredientsWithIDs(ids []string) []Ingredient {
	var (
		id     int
		name   string
		color  string
		baseID int
	)
	statement := "select id, name, color, baseid from ingredients i join baseIngredients b on i.id = b.id  where i.id = any($1::integer[])"
	args := "{"
	for _, v := range ids {
		args += v + ","
	}
	args = args[:len(args)-1] + "}"
	ingredients := make([]Ingredient, 0)
	rows, err := i.data.Query(statement, args)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name, &color, &baseID)
		if err != nil {
			log.Fatal(err)
		}
		ingredients = append(ingredients, Ingredient{
			ID:     id,
			Name:   name,
			Color:  color,
			BaseID: baseID,
		})
	}
	return ingredients
}
