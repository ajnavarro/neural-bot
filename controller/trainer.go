package main

import (
	"encoding/csv"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/goml/gobrain"
)

type TrainerOptions struct {
	Path string `short:"p" long:"path" description:"Path where the data file is stored" default:"out.csv"`
}

func (o *TrainerOptions) Execute(args []string) error {
	log.Println("TRAIN")
	f, err := os.Open(o.Path)
	if err != nil {
		log.Fatalln("Error opening file", err)
	}

	r := csv.NewReader(f)

	// set the random seed to 0
	rand.Seed(0)

	// create the XOR representation patter to train the network
	patterns := [][][]float64{}
	count := 0
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		values := []float64{}
		for _, field := range record {
			fv, err := strconv.ParseFloat(field, 64)
			if err != nil {
				return err
			}

			values = append(values, fv)
		}

		sm := values[2]
		d := values[3]
		lm := values[0]
		rm := values[1]
		patterns = append(patterns, [][]float64{{sm, d}, {lm, rm}})

		count++
	}

	log.Println("Patterns added", count)

	// instantiate the Feed Forward
	ff := &gobrain.FeedForward{}

	// initialize the Neural Network;
	// the networks structure will contain:
	// 2 inputs, 5 hidden nodes and 2 outputs.
	ff.Init(2, 5, 2)

	// train the network using the XOR patterns
	// the training will run for 1000 epochs
	// the learning rate is set to 0.6 and the momentum factor to 0.4
	// use true in the last parameter to receive reports about the learning error
	ff.Train(patterns, 1000, 0.6, 0.4, true)

	ff.Test(patterns)

	log.Println("Input nodes:", ff.NInputs, "Input activations:", ff.InputActivations, "Input weights:", ff.InputWeights)
	log.Println("Hidden nodes:", ff.NHiddens, "Hidden activations:", ff.HiddenActivations)
	log.Println("Output nodes:", ff.NOutputs, "Output activations:", ff.OutputActivations, "Output weights:", ff.OutputWeights)

	return nil
}
