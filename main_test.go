package main

import "testing"

func TestIngredientLookup(t *testing.T) {
	allPossible := getAllIngredients()
	ids := []string{"0", "1"}
	used := ingredientsForDrink(ids, allPossible)
	if used[0].ID != 0 || used[1].ID != 1 {
		t.Error("expected 0 and 1", "got", used[0].ID, used[1].ID)
	}
}
