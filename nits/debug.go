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
				Help:    "Displays all registered answers",
				Executor: func(strings []string) bool {
					displayAnswers(ui)
					return false
				},
			},
			{
				Aliases: []string{"done"},
				Help:    "Exits the debugger",
				Executor: func(strings []string) bool {
					return true
				},
			},
			{
				Aliases: []string{"train"},
				Help:    "Runs trainhmm",
				Executor: func(strings []string) bool {
					td, err := writeTrainhmmInput(content)
					ui.println("Temporary directory: td=%s; err=%v", td, err)
					if err != nil {
						return false
					}
					err = runTrainhmm(td)
					ui.println("Training: err=%v", err)
					err = readPrediction(td, content, ui)
					ui.println("Reading predictions: err=%v", err)
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
