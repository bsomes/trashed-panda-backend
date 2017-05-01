package main

import (
	"database/sql"
	"log"
)

type drinkDataSQL struct {
	db *sql.DB
}

func (d *drinkDataSQL) GetBaseIngredient(id int) Ingredient {
	var (
		name string
		cat  Category
	)
	err := d.db.QueryRow("select name, category from BaseIngredients where id = $1", id).Scan(&name, &cat)
	if err != nil {
		log.Fatal(err)
	}
	return Ingredient{
		ID:   id,
		Name: name,
		Cat:  cat,
	}
}

func (d *drinkDataSQL) NumDrinksWithIngredient(ingredientID int) int {
	var count int
	err := d.db.QueryRow("select count(drinkid) from contains where baseid = $1", ingredientID).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count
}

func (d *drinkDataSQL) NumDrinksWithBothIngredients(firstID int, secondID int) int {
	var count int
	err := d.db.QueryRow(
		"select count(c.drinkid) from contains c join contains s on c.drinkid = s.drinkid where c.baseid = $1 and s.baseid = $2",
		firstID, secondID).
		Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count
}

func (d *drinkDataSQL) NumDrinksWithThreeIngredients(firstID int, secondID int, thirdID int) int {
	var count int
	err := d.db.QueryRow(
		"select count(c.drinkid) from contains c join contains s on c.drinkid = s.drinkid join contains t on s.drinkid = t.drinkid where c.baseid = $1 and s.baseid = $2 and t.baseid = $3",
		firstID, secondID, thirdID).
		Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count
}

func (d *drinkDataSQL) NumIngredientsInDrink(ingredientID int) (int, int) {
	var (
		mean float64
		max  int
	)
	statement := `select AVG(d.count), MAX(d.count) from
(select COUNT(*)
from contains c
JOIN (select drinkid from contains where baseid = $1) ids
on c.drinkid = ids.drinkid
group by c.drinkid) d`

	d.db.QueryRow(statement, ingredientID).Scan(&mean, &max)
	rounded := int(mean + 0.5)
	return rounded, max
}
