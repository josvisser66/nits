package nits

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

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

	initBKT()
	initConcepts()
	ui := newUserInterface()
	defer ui.rl.Close()

	ui.pushCommandContext(&CommandContext{
		Description: "NITS core commands",
		Commands: []*Command{
			{
				Aliases: []string{"exit", "quit"},
				Global:  true,
				Help:    "Exits NITS.",
				Executor: func(line []string) bool {
					if ui.yesNo("Are you sure you want to quit") {
						os.Exit(0)
					}
					return false
				},
			},
			{
				Aliases: []string{"debug"},
				Global:  true,
				Help:    "NITS debugging (internal)",
				Executor: func([]string) bool {
					debug(ui)
					return false
				},
			},
			{
				Aliases: []string{"load"},
				Global:  true,
				Help:    "Load user data",
				Executor: func([]string) bool {
					if err := loadUserData(); err != nil {
						ui.println("Loading failed: %s", err)
					}
					return false
				},
			},
			{
				Aliases: []string{"save"},
				Global:  true,
				Help:    "Save user data",
				Executor: func([]string) bool {
					if err := saveUserData(); err != nil {
						ui.println("Saving failed: %s", err)
					}
					return false
				},
			},
		},
	})
	defer ui.popCommandContext()

	if err := loadUserData(); err != nil {
		ui.println("User data *not* loaded: %s", err)
	} else {
		ui.println("User data restored.")
	}

	ui.newline()

	for {
		selectQuestion(content).ask(ui)
	}
}

func (c *Content) check() {
checkConcepts()
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
