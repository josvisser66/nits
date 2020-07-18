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
	Consequences []Event
	Explanation  *Explanation
	event        Event
}

func (b *BrokenLegalRequirement) getLabel() string {
	return b.Description
}

// --------------------------------------------------------------------
type Duty struct {
	Description string
	OwedFrom    []*Person
	OwedTo      []*Person
	event       Event
}

func (d *Duty) getLabel() string {
	return d.Description
}

// --------------------------------------------------------------------
type IrrelevantCause struct {
	Description string
	Explanation *Explanation
	event       Event
}

func (i *IrrelevantCause) getLabel() string {
	return i.Description
}

// --------------------------------------------------------------------
type Event interface {
	getLabel() string // For dot drawings.
	getShortName() string
	getDescription() string
	getConsequences() []Event
	getDuty() *Duty
	getNegPerSe() *BrokenLegalRequirement
	getIrrelevantCause() *IrrelevantCause
	getInjuriesOrDamages() []InjuryOrDamage
	getCauses() []Event
	addCause(e Event)
}

type PassiveEvent struct {
	shortName         string
	Description       string
	Consequences      []Event
	Duty              *Duty
	NegPerSe          *BrokenLegalRequirement
	IrrelevantCause   *IrrelevantCause
	InjuriesOrDamages []InjuryOrDamage
	causes            []Event
}

func (pe *PassiveEvent) getLabel() string {
	return pe.Description
}

func (pe *PassiveEvent) addCause(event Event) {
	if event == nil {
		return
	}
	if pe.causes == nil {
		pe.causes = make([]Event, 0, 1)
	}
	pe.causes = append(pe.causes, event)
}

func (pe *PassiveEvent) getShortName() string {
	return pe.shortName
}

func (pe *PassiveEvent) getDescription() string {
	return pe.Description
}

func (pe *PassiveEvent) getConsequences() []Event {
	return pe.Consequences
}

func (pe *PassiveEvent) getDuty() *Duty {
	return pe.Duty
}

func (pe *PassiveEvent) getNegPerSe() *BrokenLegalRequirement {
	return pe.NegPerSe
}

func (pe *PassiveEvent) getIrrelevantCause() *IrrelevantCause {
	return pe.IrrelevantCause
}

func (pe *PassiveEvent) getInjuriesOrDamages() []InjuryOrDamage {
	return pe.InjuriesOrDamages
}

func (pe *PassiveEvent) getCauses() []Event {
	return pe.causes
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
	getCauses() []Event
	addCause(event Event)
}

type BodilyInjury struct {
	Description string
	Persons     []*Person
	causes      []Event
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

func (b *BodilyInjury) addCause(event Event) {
	if b.causes == nil {
		b.causes = make([]Event, 0, 1)
	}
	b.causes = append(b.causes, event)
}

func (b *BodilyInjury) getCauses() []Event {
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
	causes      []Event
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

func (p *PropertyDamage) addCause(event Event) {
	if p.causes == nil {
		p.causes = make([]Event, 0, 1)
	}
	p.causes = append(p.causes, event)
}

func (p *PropertyDamage) getCauses() []Event {
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
