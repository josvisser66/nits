package nits

import (
	"math/rand"
	"time"
)

func Run(content *Content) {
	// This initialization of the random generator is not cryptographically
	// secure, but it's good enough for our purpose.
	rand.Seed(time.Now().UnixNano())

	println("NITS 1.0 -- An ITS for negligence")
	println("            (c) Copyright 2020  Jos Visser <josvisser66@gmail.com>")
	println()

	ui := newUserInterface()
	defer ui.rl.Close()

	content.Questions[1].ask(ui)
}
