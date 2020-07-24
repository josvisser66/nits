// Package nits contains the implementation of NITS, an ITS for negligence.
package nits

// This file contains the outermost function of NITS.

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"
)

// Run runs NITS on some content.
func Run(content *Content) {
	content.check()

	// This initialization of the random generator is not cryptographically
	// secure, but it's good enough for our purpose.
	rand.Seed(time.Now().UnixNano())

	println("NITS 1.0 -- An ITS for negligence")
	println("            (c) Copyright 2020  Jos Visser <josvisser66@gmail.com>")
	println()
	println("Use ? to get help.")
	println()
	println("Running on", runtime.GOOS)

	initBKT()
	initConcepts()
	state := newStudentState(content)
	ui := newUserInterface()
	defer ui.rl.Close()

	// Pushes the outermost command context.
	ui.pushCommandContext(&CommandContext{
		description: "NITS core commands",
		commands: []*Command{
			{
				aliases: []string{"exit", "quit"},
				global:  true,
				help:    "Exits NITS. Use nosave argument *not* to save the user state.",
				executor: func(words []string) bool {
					if sure, _ := ui.yesNo("Are you sure you want to quit"); sure {
						if len(words) == 1 || (len(words) > 1 && words[1] != "nosave") {
							state.saveUserData()
						}
						os.Exit(0)
					}
					return false
				},
			},
			{
				aliases: []string{"reset"},
				global:  true,
				help:    "Resets internal student state",
				executor: func(words []string) bool {
					state.reset()
					return false
				},
			},{
				aliases: []string{"debug"},
				global:  true,
				help:    "NITS debugging (internal)",
				executor: func(words []string) bool {
					debug(ui, state, words)
					return false
				},
			},
			{
				aliases: []string{"load"},
				global:  true,
				help:    "Load student data",
				executor: func([]string) bool {
					if err := state.loadUserData(); err != nil {
						ui.error("Loading failed: %s", err)
					}
					return false
				},
			},
			{
				aliases: []string{"save"},
				global:  true,
				help:    "Save student data",
				executor: func([]string) bool {
					if err := state.saveUserData(); err != nil {
						ui.error("Saving failed: %s", err)
					}
					return false
				},
			},
		},
	})
	defer ui.popCommandContext()

	if err := state.loadUserData(); err != nil {
		ui.error("User data *not* loaded: %s", err)
	} else {
		ui.println("User data restored.")
	}

	ui.newline()

	for {
		if next := state.selectQuestion(); next != nil {
			next.ask(ui, state)
		} else {
			ui.println("We are out of questions!")
			break
		}
	}

	state.saveUserData()
}

// check checks the content. Mostly delegates to the check methods
// of each of the questions.
func (c *Content) check() {
	m := make(map[string]interface{})

	for _, q := range c.Questions {
		name := q.getShortName()
		if _, ok := m[name]; ok {
			panic(fmt.Sprintf("Duplicate question short name: %s", name))
		}
		m[name] = nil
		q.check()
	}
}

// Helper functions. Go people don't like this CHECK function, but I do.
func CHECK(b bool, s string, args ...interface{}) {
	if !b {
		panic(fmt.Sprintf(s, args...))
	}
}
