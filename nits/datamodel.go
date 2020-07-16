package nits

// --------------------------------------------------------------------
type Person struct {
	Name string
}

func (p *Person) getLabel() string {
	return p.Name
}

// --------------------------------------------------------------------
type BrokenLegalRequirement struct {
	Description  string
	Persons      []*Person
	Consequences []*Event
	Explanation  *Explanation
	event        *Event
}

func (b *BrokenLegalRequirement) getLabel() string {
	return b.Description
}

// --------------------------------------------------------------------
type Duty struct {
	Description string
	OwedFrom    []*Person
	OwedTo      []*Person
	event       *Event
}

func (d *Duty) getLabel() string {
	return d.Description
}

// --------------------------------------------------------------------
type IrrelevantCause struct {
	Description string
	Explanation *Explanation
	event       *Event
}

func (i *IrrelevantCause) getLabel() string {
	return i.Description
}

// --------------------------------------------------------------------
type Event struct {
	Description       string
	Consequences      []*Event
	Duty              *Duty
	NegPerSe          *BrokenLegalRequirement
	IrrelevantCause   *IrrelevantCause
	InjuriesOrDamages []InjuryOrDamage
}

func (e *Event) getLabel() string {
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
	getLabel() string
}

type BodilyInjury struct {
	Description string
	Persons     []*Person
}

func (b *BodilyInjury) GetDescription() string {
	return b.Description
}

func (b *BodilyInjury) getLabel() string {
	return b.Description
}

func (b *BodilyInjury) GetPersons() []*Person {
	return b.Persons
}

type EmotionalHarm struct {
	Description string
	Persons     []*Person
}

func (e *EmotionalHarm) GetInjuryDescription() string {
	return e.Description
}

type PropertyDamage struct {
	Description string
	Persons     []*Person
}

func (p *PropertyDamage) GetDescription() string {
	return p.Description
}

func (p *PropertyDamage) getLabel() string {
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
