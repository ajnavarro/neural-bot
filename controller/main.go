package main

import (
	"os"

	"github.com/jessevdk/go-flags"
)

func main() {
	parser := flags.NewParser(nil, flags.Default)

	if _, err := parser.AddCommand("gamepad", "Gamepad controller",
		"Use a joystic to control the robot", &GamepadOptions{}); err != nil {
		panic(err)
	}

	if _, err := parser.AddCommand("train", "Train the neural network",
		"Using the output generated by the gamepad, train a neural network", &TrainerOptions{}); err != nil {
		panic(err)
	}

	if _, err := parser.Parse(); err != nil {
		if err, ok := err.(*flags.Error); ok {
			if err.Type == flags.ErrHelp {
				os.Exit(0)
			}

			parser.WriteHelp(os.Stdout)
		}

		os.Exit(1)
	}

}