package main

import (
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

type nameGenerator interface {
	NameWithIngredients(ingredients []int) string
}

type rnnNameGenerator struct {
	vocabulary []string
	sess       *tf.Session
}
