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
	shortName		  string
	Description       string
	Consequences      []*Event
	Duty              *Duty
	NegPerSe          *BrokenLegalRequirement
	IrrelevantCause   *IrrelevantCause
	InjuriesOrDamages []InjuryOrDamage
	causes            []*Event
}

func (e *Event) getLabel() string {
	return e.Description
}

func (e *Event) addCause(event *Event) {
	if event == nil {
		return
	}
	if e.causes == nil {
		e.causes = make([]*Event,0,1)
	}
	e.causes = append(e.causes, event)
}

func (e *Event) getCauses() []*Event {
	return e.causes
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
	getCauses() []*Event
	addCause(event *Event)
}

type BodilyInjury struct {
	Description string
	Persons     []*Person
	causes      []*Event
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

func (b *BodilyInjury) addCause(event *Event) {
	if b.causes == nil {
		b.causes = make([]*Event, 0, 1)
	}
	b.causes = append(b.causes, event)
}

func (b *BodilyInjury) getCauses() []*Event {
	return b.causes
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
	causes      []*Event
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

func (p *PropertyDamage) addCause(event *Event) {
	if p.causes == nil {
		p.causes = make([]*Event, 0, 1)
	}
	p.causes = append(p.causes, event)
}

func (p *PropertyDamage) getCauses() []*Event {
	return p.causes
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
