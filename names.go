package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

type nameGenerator interface {
	NameWithIngredients(ingredients []int) string
}

type rnnNameGenerator struct {
	vocabulary []string
	model      *tf.SavedModel
	buckets    []bucket
}

type bucket struct {
	id         int
	inputSize  int
	outputSize int
}

func buckets() []bucket {
	return []bucket{
		bucket{id: 0, inputSize: 6, outputSize: 3},
		bucket{id: 1, inputSize: 4, outputSize: 7},
		bucket{id: 2, inputSize: 9, outputSize: 5}}
}

func (rnn *rnnNameGenerator) NameWithIngredients(ingredients []int) string {
	bucket, err := rnn.findBucket(ingredients)
	if err != nil {
		log.Fatal(err)
	}
	inputs := rnn.makeInputs(ingredients, bucket)
	outFeed, er := rnn.outputFeed(bucket.id)
	if err != nil {
		log.Fatal(er)
	}
	return rnn.run(inputs, outFeed)
}

func (rnn *rnnNameGenerator) makeInputs(inputs []int, bucket *bucket) map[tf.Output]*tf.Tensor {
	inputDict := make(map[tf.Output]*tf.Tensor)
	for ind := 0; ind < bucket.inputSize; ind++ {
		var v int32
		if ind < len(inputs) {
			v = int32(inputs[ind])
		} else {
			v = 0
		}
		vals := [64]int32{}
		for i := range vals {
			vals[i] = v
		}
		val, err := tf.NewTensor(vals)
		if err != nil {
			log.Fatal(err)
		}
		op := rnn.model.Graph.Operation("encoder" + strconv.Itoa(ind)).Output(0)
		inputDict[op] = val
	}
	for ind := 0; ind <= bucket.outputSize; ind++ {
		vals := [64]int32{}
		for i := range vals {
			if i == 0 {
				vals[i] = 1
			} else {
				vals[i] = 0
			}
		}
		val, err := tf.NewTensor(vals)
		if err != nil {
			log.Fatal(err)
		}
		inputDict[rnn.model.Graph.Operation("decoder"+strconv.Itoa(ind)).Output(0)] = val
		val, err = tf.NewTensor([64]float32{})
		if err != nil {
			log.Fatal(err)
		}
		inputDict[rnn.model.Graph.Operation("weight"+strconv.Itoa(ind)).Output(0)] = val
	}
	return inputDict

}

func (rnn *rnnNameGenerator) outputFeed(bucketId int) ([]tf.Output, error) {
	switch bucketId {
	case 0:
		return []tf.Output{
			//rnn.model.Graph.Operation("model_with_buckets/sequence_loss/truediv").Output(0),
			rnn.model.Graph.Operation("add").Output(0),
			rnn.model.Graph.Operation("add_1").Output(0),
			rnn.model.Graph.Operation("add_2").Output(0),
		}, nil
	case 1:
		return []tf.Output{
			//rnn.model.Graph.Operation("model_with_buckets/sequence_loss_1/truediv").Output(0),
			rnn.model.Graph.Operation("add_3").Output(0),
			rnn.model.Graph.Operation("add_4").Output(0),
			rnn.model.Graph.Operation("add_5").Output(0),
			rnn.model.Graph.Operation("add_6").Output(0),
			rnn.model.Graph.Operation("add_7").Output(0),
			rnn.model.Graph.Operation("add_8").Output(0),
			rnn.model.Graph.Operation("add_9").Output(0),
		}, nil
	case 2:
		return []tf.Output{
			//rnn.model.Graph.Operation("model_with_buckets/sequence_loss_2/truediv").Output(0),
			rnn.model.Graph.Operation("add_10").Output(0),
			rnn.model.Graph.Operation("add_11").Output(0),
			rnn.model.Graph.Operation("add_12").Output(0),
			rnn.model.Graph.Operation("add_13").Output(0),
			rnn.model.Graph.Operation("add_14").Output(0),
		}, nil
	default:
		return nil, errors.New("BucketId was outside range of all buckets")
	}

}

func (rnn *rnnNameGenerator) run(input map[tf.Output]*tf.Tensor, output []tf.Output) string {
	out, err := rnn.model.Session.Run(input, output, nil)
	if err != nil {
		log.Fatal(err)
	}
	return rnn.makeName(findHighestScoringIndices(out))
}

func (rnn *rnnNameGenerator) findBucket(ingredients []int) (*bucket, error) {
	if len(ingredients) > 9 {
		return nil, errors.New("input length was longer than longest bucket")
	}
	var bucketToUse bucket
	for _, v := range rnn.buckets {
		if len(ingredients) <= v.inputSize {
			bucketToUse = v
			break
		}
	}
	return &bucketToUse, nil
}

func (rnn *rnnNameGenerator) Close() {
	rnn.model.Session.Close()
}

func (rnn *rnnNameGenerator) makeName(indices []int) string {
	words := make([]string, 0)
	for i, v := range indices {
		//Hack to prevent repeated words
		if i == 0 || indices[i-1] != indices[i] {
			words = append(words, rnn.vocabulary[v])
		}
	}
	return strings.Join(words, " ")
}

func findHighestScoringIndices(outputs []*tf.Tensor) []int {
	wordIndices := make([]int, 0)
	for _, output := range outputs {
		scores := output.Value().([][]float32)[0]
		wordIndices = append(wordIndices, argmax(scores...))
	}
	return wordIndices
}

func argmax(vals ...float32) int {
	max := vals[0]
	var index int
	for i, v := range vals {
		if v > max {
			max = v
			index = i
		}
	}
	return index
}

func makeRnnNameGenerator(vocabPath string, modelPath string) *rnnNameGenerator {
	model, err := tf.LoadSavedModel(modelPath, []string{"serve"}, nil)
	if err != nil {
		log.Fatal("Failed to load rnn model", err)
	}
	vocab := loadVocabulary(vocabPath)
	return &rnnNameGenerator{
		vocabulary: vocab,
		model:      model,
		buckets:    buckets(),
	}
}

func loadVocabulary(path string) []string {
	const (
		PAD = "PAD"
		GO  = "GO"
		EOS = "EOS"
	)
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("Failed to open vocabulary at path " + path)
	}
	defer f.Close()
	var words []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return append([]string{PAD, GO, EOS}, words...)
}
