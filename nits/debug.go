package nits

func displayAnswers(ui *userInterface, state *studentState) {
	ui.println("Registered answers:")
	ui.newline()

	for _, answer := range state.answers {
		ui.println("%5t: %s", answer.correct, answer.questionShortName)
	}
}

func debug(ui *userInterface, state *studentState) {
	ui.pushCommandContext(&CommandContext{
		Description: "NITS debugger",
		Commands: []*Command{
			{
				Aliases: []string{"answers"},
				Help:    "Displays all registered answers.",
				Executor: func(strings []string) bool {
					displayAnswers(ui, state)
					return false
				},
			},
			{
				Aliases: []string{"done"},
				Help:    "Exits the debugger.",
				Executor: func(strings []string) bool {
					return true
				},
			},
			{
				Aliases: []string{"train"},
				Help:    "Runs trainhmm.",
				Executor: func([]string) bool {
					td, err := state.writeTrainhmmInput()
					ui.println("Temporary directory: td=%s; err=%v", td, err)
					if err != nil {
						return false
					}
					err = state.runTrainhmm(td)
					ui.println("Training: err=%v", err)
					err = state.readPrediction(td)
					ui.println("Reading predictions: err=%v", err)
					for concept, skillLevel := range state.scores {
						ui.println("%-20s: %f", concept.shortName, skillLevel)
					}
					return false
				},
			},
			{
				Aliases: []string{"scores"},
				Help:    "Shows all concepts and the student's skill levels.",
				Executor: func([]string) bool {
					for concept, skillLevel := range state.scores {
						ui.println("%-20s: %f", concept.shortName, skillLevel)
					}
					return false
				},
			},
			{
				Aliases: []string{"dot"},
				Help:    "Generates dot file for content and concepts.",
				Executor: func([]string) bool {
					fname, err := makeDot(state, ui.yesNo("Augment concepts with skill level ratios"))
					if err != nil {
						ui.println("error: %s", err)
					} else {
						ui.println("now run: dot -Tpdf %s >/tmp/aap.pdf", fname)
					}
					return false
				},
			},
		},
	})
	ui.pushPrompt("debug> ")
	defer ui.popPrompt()
	defer ui.popCommandContext()

	for {
		_, ret := ui.getInput()
		if ret {
			return
		}
		ui.println("Please enter one of the debugger's commands. Use ? for help.")
	}
}
