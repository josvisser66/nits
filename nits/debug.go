package nits

func debug(ui *userInterface) {
	ui.pushCommandContext(&CommandContext{
		Description: "NITS debugger",
		Commands: []*Command{
			{
				Aliases: []string{"answers"},
				Help:    "Displays all registered answers",
				Executor: func(strings []string) bool {
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
