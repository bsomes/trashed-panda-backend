package main

import "testing"

func TestIngredientLookup(t *testing.T) {
	allPossible := makeIDMap(getAllIngredients())
	ids := []string{"53", "54"}
	used := ingredientsForDrink(ids, allPossible)
	if used[0].ID != 53 || used[1].ID != 54 {
		t.Error("expected 53 and 54", "got", used[0].ID, used[1].ID)
	}
}

func TestIngredientLoad(t *testing.T) {
	ingredients := getAllIngredients()
	expectedName := "cream liqueur"
	if ingredients[0].Name != expectedName {
		t.Error("Expected the first ingredient to be", "instead it was", expectedName, ingredients[0].Name)
	}
}
