package nits

// --------------------------------------------------------------------
type Person struct {
	Name string
}

// --------------------------------------------------------------------
type BrokenLegalRequirement struct {
	Description string
	Persons []*Person
	Consequences []*Event
	Explanation *Explanation
}

// --------------------------------------------------------------------
type Duty interface {
	GetDutyDescription() string
}

type ActiveDuty struct {
	Description string
	OwedFrom []*Person
	OwedTo[] *Person
}

func (a *ActiveDuty) GetDutyDescription() string {
	return a.Description
}

type NonDuty struct {
	Description string
}

func (n *NonDuty) GetDutyDescription() string {
	return n.Description
}

// --------------------------------------------------------------------
type IrrelevantCause struct {
	Description string
	Explanation *Explanation
}

// --------------------------------------------------------------------
type Event struct {
	Description string
	Consequences []*Event
	Duty Duty
	NegPerSe BrokenLegalRequirement
	IrrelevantCause *IrrelevantCause
	InjuriesOrDamages []InjuryOrDamage
}

func (e *Event) GetCauseDescription() string {
	return e.Description
}

// --------------------------------------------------------------------
type Cause interface {
	GetCauseDescription() string
}

// --------------------------------------------------------------------
type InjuryOrDamage interface {
	GetDescription() string
	GetPersons() []*Person
}

type BodilyInjury struct {
	Description string
	Persons []*Person
}

func (b *BodilyInjury) GetDescription() string {
	return b.Description
}

func (b *BodilyInjury) GetPersons() []*Person {
	return b.Persons
}

type EmotionalHarm struct {
	Description string
	Persons []*Person
}

func (e *EmotionalHarm) GetInjuryDescription() string {
	return e.Description
}

type PropertyDamage struct {
	Description string
	Persons []*Person
}

func (p *PropertyDamage) GetDescription() string {
	return p.Description
}

func (p *PropertyDamage) GetPersons() []*Person {
	return p.Persons
}

// --------------------------------------------------------------------
type Content struct {
	Questions []Question
}

func (c *Content) findQuestion(shortName string) Question {
	for _, q := range c.Questions {
		if q.getShortName() == shortName {
			return q
		}
	}

	return nil
}
