package nits

// This file contains the core data model of NITS.

// --------------------------------------------------------------------
// Person is a person that is involved in a case.
type Person struct {
	Name    string
	damages map[InjuryOrDamage]interface{}
}

func (p *Person) getLabel() string {
	return p.Name
}

// --------------------------------------------------------------------
// BrokenLegalRequirement is the fact that one or more persons are in
// violation of a statute ir regulation (negligence per se).
type BrokenLegalRequirement struct {
	Description  string
	Persons      []*Person
	Consequences []Event
	Explanation  *Explanation
	event        Event // Back link to event that points to this BrokenLegalRequirement.
}

func (b *BrokenLegalRequirement) getLabel() string {
	return b.Description
}

// --------------------------------------------------------------------
// Duty is a legal obligation.
type Duty struct {
	Description string
	OwedFrom    []*Person
	OwedTo      []*Person
	event       Event // Back link to the event that points to this Duty.
}

func (d *Duty) getLabel() string {
	return d.Description
}

// --------------------------------------------------------------------
// Claim is an (irrelevant) claim that a person is making for an event.
type Claim struct {
	Person      *Person
	Description string
	Explanation *Explanation
	event       Event // Back link to the event that points to this Claim.
}

func (c *Claim) getLabel() string {
	return c.Description
}

// --------------------------------------------------------------------
// Event is something that happened.
type Event interface {
	getLabel() string // For dot drawings.
	getShortName() string
	getDescription() string
	getConsequences() []Event
	getDuty() *Duty
	getNegPerSe() *BrokenLegalRequirement
	getInjuriesOrDamages() []InjuryOrDamage
	getDirectCauses() []Event
	getClaims() []*Claim
	addCause(e Event)
}

// PassiveEvent is an event that just happens, it is not an Act.
type PassiveEvent struct {
	shortName         string // for internal (testing) use only.
	Description       string
	Consequences      []Event
	Duty              *Duty
	NegPerSe          *BrokenLegalRequirement
	InjuriesOrDamages []InjuryOrDamage
	Claims            []*Claim
	directCauses      []Event // Back links to the events that this event is a consequence of.
}

func (pe *PassiveEvent) getLabel() string {
	return pe.Description
}

// addCause adds an event that is the cause of this event.
func (pe *PassiveEvent) addCause(event Event) {
	if event == nil {
		return
	}
	if pe.directCauses == nil {
		pe.directCauses = make([]Event, 0, 1)
	}
	pe.directCauses = append(pe.directCauses, event)
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

func (pe *PassiveEvent) getInjuriesOrDamages() []InjuryOrDamage {
	return pe.InjuriesOrDamages
}

func (pe *PassiveEvent) getClaims() []*Claim {
	return pe.Claims
}

func (pe *PassiveEvent) getDirectCauses() []Event {
	return pe.directCauses
}

// --------------------------------------------------------------------
// Act is an event that was a willful act by a person.
type Act struct {
	shortName         string // Internal use (testing) only.
	Person            *Person
	Description       string
	Consequences      []Event
	Duty              *Duty
	NegPerSe          *BrokenLegalRequirement
	InjuriesOrDamages []InjuryOrDamage
	Claims            []*Claim
	directCauses      []Event // Back links to the events that inspired this act.
}

func (a *Act) getLabel() string {
	return a.Description
}

func (a *Act) addCause(event Event) {
	if event == nil {
		return
	}
	if a.directCauses == nil {
		a.directCauses = make([]Event, 0, 1)
	}
	a.directCauses = append(a.directCauses, event)
}

func (a *Act) getShortName() string {
	return a.shortName
}

func (a *Act) getDescription() string {
	return a.Description
}

func (a *Act) getConsequences() []Event {
	return a.Consequences
}

func (a *Act) getDuty() *Duty {
	return a.Duty
}

func (a *Act) getNegPerSe() *BrokenLegalRequirement {
	return a.NegPerSe
}

func (a *Act) getInjuriesOrDamages() []InjuryOrDamage {
	return a.InjuriesOrDamages
}

func (a *Act) getClaims() []*Claim {
	return a.Claims
}

func (a *Act) getDirectCauses() []Event {
	return a.directCauses
}

// --------------------------------------------------------------------
// InjuryOrDamage; speaks for itself :-)
type InjuryOrDamage interface {
	GetDescription() string
	GetPersons() []*Person
	getLabel() string
	getDirectCauses() []Event
	addCause(event Event)
}

// BodilyInjury is a bodily injury suffered by one or more persons.
type BodilyInjury struct {
	Description  string
	Persons      []*Person
	directCauses []Event // Back links to the events that directly caused this injury.
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
	if b.directCauses == nil {
		b.directCauses = make([]Event, 0, 1)
	}
	b.directCauses = append(b.directCauses, event)
}

func (b *BodilyInjury) getDirectCauses() []Event {
	return b.directCauses
}

// PropertyDamage is damage to somebody's property.
type PropertyDamage struct {
	Description  string
	Persons      []*Person
	directCauses []Event // Back links to the events that directly caused this damage.
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
	if p.directCauses == nil {
		p.directCauses = make([]Event, 0, 1)
	}
	p.directCauses = append(p.directCauses, event)
}

func (p *PropertyDamage) getDirectCauses() []Event {
	return p.directCauses
}

// --------------------------------------------------------------------

// Content is the question content that NITS operates on.
type Content struct {
	Questions []Question
}

// findQuestion finds a question by short name.
func (c *Content) findQuestion(shortName string) Question {
	for _, q := range c.Questions {
		if q.getShortName() == shortName {
			return q
		}
	}

	return nil
}
