package nits

import "fmt"

// --------------------------------------------------------------------
type HelpItem struct {
	Text string
}

type Help struct {
	Hints []string
	Items []*HelpItem
}

type Reference interface {
	GetReferenceText() string
}

type Explanation struct {
	Text []string
	References []Reference
}

type Restatement struct {
	Paragraph int
}

func (r *Restatement) GetReferenceText() string {
	return fmt.Sprintf("Restatement, Second, Torts, §%d", r.Paragraph)
}

type URL struct {
	Url string
}

func (u *URL) GetReferenceText() string {
	return u.Url
}

type helpContext struct {
	help *Help
}
