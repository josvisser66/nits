package nits

import (
	"math/rand"
	"os"
	"time"
)

func Run(content *Content) {
	// This initialization of the random generator is not cryptographically
	// secure, but it's good enough for our purpose.
	rand.Seed(time.Now().UnixNano())

	println("NITS 1.0 -- An ITS for negligence")
	println("            (c) Copyright 2020  Jos Visser <josvisser66@gmail.com>")
	println()
	println("Use ? to get help.")
	println()

	initBKT()
	initConcepts()
	ui := newUserInterface()
	defer ui.rl.Close()

	ui.pushCommandContext(&CommandContext{
		Description: "NITS core commands",
		Commands: []*Command{
			{
				Aliases: []string{"exit", "quit"},
				Global: true,
				Help:    "Exits NITS.",
				Executor: func(line []string) bool {
					if ui.yesNo("Are you sure you want to quit") {
						os.Exit(0)
					}
					return false
				},
			},
			{
				Aliases:  []string{"debug"},
				Global: true,
				Help:     "NITS debugging (internal)",
				Executor: func([]string) bool {
					debug(ui)
					return false
				},
			},
		},
	})
	defer ui.popCommandContext()

	content.Questions[2].ask(ui)
}
