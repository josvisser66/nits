package nits

func displayAnswers(ui *userInterface) {
	ui.println("Registered answers:")
	ui.newline()

	for _, answer := range answers {
		ui.println("%5t: %s", answer.correct, answer.questionShortName)
	}
}

func debug(ui *userInterface, content *Content) {
	ui.pushCommandContext(&CommandContext{
		Description: "NITS debugger",
		Commands: []*Command{
			{
				Aliases: []string{"answers"},
				Help:    "Displays all registered answers.",
				Executor: func(strings []string) bool {
					displayAnswers(ui)
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
					td, err := writeTrainhmmInput(content)
					ui.println("Temporary directory: td=%s; err=%v", td, err)
					if err != nil {
						return false
					}
					err = runTrainhmm(td)
					ui.println("Training: err=%v", err)
					err = readPrediction(td, content)
					ui.println("Reading predictions: err=%v", err)
					for shortName, skillLevel := range concepts {
						ui.println("%-20s: %f", shortName, skillLevel)
					}
					return false
				},
			},
			{
				Aliases: []string{"concepts"},
				Help: "Shows all concepts and the student's skill levels.",
				Executor: func([]string) bool {
					for shortName, skillLevel := range concepts {
						ui.println("%-20s: %f", shortName, skillLevel)
					}
					return false
				},
			},
			{
				Aliases: []string{"dot"},
				Help: "Generates dot file for content and concepts.",
				Executor: func([]string) bool {
					fname, err := makeDot(content, ui.yesNo("Augment concepts with skill level ratios"))
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
