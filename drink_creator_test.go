package main

import (
	"math"
	"testing"
)

type fakeData struct{}

func (f *fakeData) NumDrinksWithIngredient(ingredientID int) int {
	if ingredientID == 69 {
		return 10
	}
	return 20
}

func (f *fakeData) NumDrinksWithBothIngredients(firstID, secondID int) int {
	if firstID == 69 {
		return 4
	}
	return 5
}

func (f *fakeData) NumDrinksWithThreeIngredients(firstID, secondID, thirdID int) int {
	if firstID == 69 {
		return 2
	}
	return 1
}

func (f *fakeData) NumIngredientsInDrink(ingredientID int) (int, int) {
	return 4, 7
}

func Test_Utility(t *testing.T) {
	t.Parallel()
	inCommonWithLast := 3
	totalWithLast := 20
	inCommonWithLastTwo := 1
	totalWithLastTwo := 3
	expected := 2.4833333

	result := utility(inCommonWithLast, totalWithLast, inCommonWithLastTwo, totalWithLastTwo)
	if math.Abs(expected-result) > 1e-6 {
		t.Error("Expected utility was ", "Actual utility was ", expected, result)
	}
}

func Test_IngredientUtilityError(t *testing.T) {
	t.Parallel()
	creator := drinkCreator{&fakeData{}}
	_, err := creator.utilityOfIngredient(0, make([]int, 0))
	if err == nil {
		t.Error("Utility should throw error for empty included list")
	}
}

func Test_IngredientUtilityOneIncluded(t *testing.T) {
	t.Parallel()
	creator := drinkCreator{&fakeData{}}
	included := []int{1}
	expected := 1.25
	result, err := creator.utilityOfIngredient(0, included)
	if err != nil {
		t.Error(err)
	}
	if math.Abs(result-expected) > 1e-6 {
		t.Error("Expected ", "Actually was ", expected, result)
	}
}

func Test_IngredientUtilityMultiple(t *testing.T) {
	t.Parallel()
	creator := drinkCreator{&fakeData{}}
	included := []int{1, 2}
	expected := 2.05

	result, err := creator.utilityOfIngredient(0, included)
	if err != nil {
		t.Error(err)
	}
	if math.Abs(result-expected) > 1e-6 {
		t.Error("Expected ", expected, "Actually was ", result)
	}
}

func Test_AllCandidateUtility(t *testing.T) {
	t.Parallel()
	creator := drinkCreator{&fakeData{}}
	included := []int{1, 2}
	expected := []float64{2.05, 2.8}

	actual := creator.utilityOfAllCandidates([]int{4, 69}, included)
	for i, v := range actual {
		if math.Abs(expected[i]-v) > 1e-6 {
			t.Error("Expected: ", expected[i], " Actual: ", v)
		}
	}
}
