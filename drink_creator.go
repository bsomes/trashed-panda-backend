package main

import (
	"errors"
	"log"
	"math"
	"math/rand"
	"sync"
)

type drinkData interface {
	//Gets count of drinks containing this base ingredient
	NumDrinksWithIngredient(ingredientID int) int
	//Gets count of drinks containing first and second base ingredients
	NumDrinksWithBothIngredients(firstID int, secondID int) int
	//Gets count of drinks containing first, second, and third base ingredients
	NumDrinksWithThreeIngredients(firstID int, secondID int, thirdID int) int
	NumIngredientsInDrink(ingredientID int) (int, int)
}

type drinkCreator struct {
	data      drinkData
	nameMaker nameGenerator
}

func (c *drinkCreator) utilityOfIngredient(candidateID int, currentIncludedIds []int) (float64, error) {
	if len(currentIncludedIds) == 0 {
		return 0, errors.New("Tried to calculate utility without any ingredients already included")
	}

	inCommonWithLast := c.data.NumDrinksWithBothIngredients(candidateID, currentIncludedIds[len(currentIncludedIds)-1])

	totalWithLast := c.data.NumDrinksWithIngredient(currentIncludedIds[len(currentIncludedIds)-1])

	if len(currentIncludedIds) < 2 {
		return utility(inCommonWithLast, totalWithLast, 0, 0), nil
	}

	inCommonWithLastTwo := c.data.NumDrinksWithThreeIngredients(
		candidateID,
		currentIncludedIds[len(currentIncludedIds)-1],
		currentIncludedIds[len(currentIncludedIds)-2])
	totalWithLastTwo := c.data.NumDrinksWithBothIngredients(
		currentIncludedIds[len(currentIncludedIds)-1],
		currentIncludedIds[len(currentIncludedIds)-2])
	return utility(inCommonWithLast, totalWithLast, inCommonWithLastTwo, totalWithLastTwo), nil
}

func (c *drinkCreator) utilityOfAllCandidates(candidateIDs []int, currentlyIncludedIds []int) []float64 {
	var wg sync.WaitGroup
	wg.Add(len(candidateIDs))
	utilities := make([]float64, len(candidateIDs))
	for i, v := range candidateIDs {
		go func(i int, val int) {
			defer wg.Done()
			var err error
			utilities[i], err = c.utilityOfIngredient(val, currentlyIncludedIds)
			if err != nil {
				log.Fatal(err)
			}
		}(i, v)
	}
	wg.Wait()
	return utilities
}

func continueDrawing(utilities []float64) bool {
	firstVal := utilities[0]
	for _, v := range utilities {
		if v != firstVal {
			return true
		}
	}
	draw := rand.Float64()
	if draw < 0.5 {
		return false
	}
	return true
}

func (c *drinkCreator) drawNextIngredient(candidateIDs []int, currentlyIncludedIDs []int) int {
	utilities := c.utilityOfAllCandidates(candidateIDs, currentlyIncludedIDs)
	if continueDrawing(utilities) {
		sum := 0.0
		for _, v := range utilities {
			sum += v
		}
		probabilities := make([]float64, len(utilities))
		cumulative := 0.0
		for i := range probabilities {
			cumulative += utilities[i] / sum
			probabilities[i] = cumulative
		}
		draw := rand.Float64()
		for i, v := range probabilities {
			if draw < v {
				return candidateIDs[i]
			}
		}
	}
	return -1
}

func (c *drinkCreator) makeIngredientList(candidateBaseIDs []int, firstIngredient int) []int {
	averageNumIngredients, maxNumIngredients := c.data.NumIngredientsInDrink(firstIngredient)
	if maxNumIngredients == 0 {
		//If ingredient has no drinks in database make an arbitrary cocktail
		averageNumIngredients = 4
		maxNumIngredients = 6
	}
	numIngredients := numIngredients(averageNumIngredients, maxNumIngredients)
	ingList := []int{firstIngredient}
	for i := 0; i < numIngredients; i++ {
		if len(candidateBaseIDs) == 0 {
			return ingList
		}
		nextIngredient := c.drawNextIngredient(candidateBaseIDs, ingList)
		if nextIngredient == -1 {
			return ingList
		}
		ingList = append(ingList, nextIngredient)
		for i, v := range candidateBaseIDs {
			if v == nextIngredient {
				candidateBaseIDs[i] = candidateBaseIDs[len(candidateBaseIDs)-1]
				candidateBaseIDs = candidateBaseIDs[:len(candidateBaseIDs)-1]
				break
			}
		}
	}
	return ingList
}

func (c *drinkCreator) makeDrink(ingredients []Ingredient) Drink {
	//Get unique base ids
	baseSet := make(map[int]struct{})
	for _, v := range ingredients {
		baseSet[v.BaseID] = struct{}{}
	}
	baseIDs := make([]int, len(baseSet))
	i := 0
	for k := range baseSet {
		baseIDs[i] = k
		i++
	}
	firstIngredient := drawAlcohol(ingredients)
	for i, v := range baseIDs {
		if v == firstIngredient.BaseID {
			baseIDs[i] = baseIDs[len(baseIDs)-1]
			baseIDs = baseIDs[:len(baseIDs)-1]
		}
	}
	ingredientBaseIDs := c.makeIngredientList(baseIDs, firstIngredient.BaseID)
	finalIngredients := make([]Ingredient, 0)
	for _, v := range ingredients {
		if baseIDIncluded(ingredientBaseIDs, v.BaseID) {
			finalIngredients = append(finalIngredients, v)
		}
	}
	return Drink{
		Name:     c.nameMaker.NameWithIngredients(getIds(finalIngredients)),
		Contents: uniform(finalIngredients),
		Size:     Buttchug,
	}
}

func getIds(ingredients []Ingredient) []int {
	ids := make([]int, 0)
	for _, v := range ingredients {
		ids = append(ids, v.BaseID)
	}
	return ids
}

func baseIDIncluded(ids []int, baseID int) bool {
	for _, v := range ids {
		if v == baseID {
			return true
		}
	}
	return false
}

func drawAlcohol(ingredients []Ingredient) Ingredient {
	alcohols := getAlcohols(ingredients)
	return alcohols[rand.Intn(len(alcohols))]
}

func getAlcohols(ingredients []Ingredient) []Ingredient {
	result := make([]Ingredient, 0)
	for _, v := range ingredients {
		if v.Cat == Alcohol {
			result = append(result, v)
		}
	}
	return result
}

//Draws a rounded value from a triangular distribution between 1 and the max
//With the average as the mode
func numIngredients(averageNumIngredients, maxNumIngredients int) int {
	min := 1.0
	max := float64(maxNumIngredients)
	mode := float64(averageNumIngredients)
	if mode > 1 {
		mode /= 2
	}

	f := (mode - min) / (max - min)
	draw := rand.Float64()
	if draw < f {
		return int(min + math.Sqrt(draw*(max-min)*(mode-min)) + 0.5)
	}
	return int(max - math.Sqrt((1-draw)*(max-min)*(max-mode)) + 0.5)
}

func utility(inCommonWithLast, totalWithLast, inCommonWithLastTwo, totalWithLastTwo int) float64 {
	const (
		a float64 = 3.0
		b float64 = 10.0
		c float64 = 0.5
	)

	if totalWithLast == 0 {
		return 0
	} else if totalWithLastTwo == 0 {
		return a*(float64(inCommonWithLast)/float64(totalWithLast)) + c
	} else {
		return a*(float64(inCommonWithLast)/float64(totalWithLast)) +
			b*(float64(inCommonWithLastTwo)/float64(totalWithLastTwo)) + c
	}
}
