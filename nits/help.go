package nits

// --------------------------------------------------------------------
type HelpItem struct {
	Text string
}

type Help struct {
	Hints []string
	Items []*HelpItem
}

type helpContext struct {
	help *Help
}
