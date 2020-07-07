package nits

// --------------------------------------------------------------------
type Person struct {
	Name string
	Injuries []*Injury
	Duties []*Duty
	Acts []*Act
}

type Victim interface {
	GetPerson() *Person
	GetInjuries() []*Injury
}

func (p *Person) GetInjuries() []*Injury {
	return p.Injuries
}

type Actor interface {
	GetPerson() *Person
	GetActs() []*Act
}

func (p *Person) GetActs() []*Act {
	return p.Acts
}

func (p *Person) GetPerson() *Person {
	return p
}

// --------------------------------------------------------------------
type Act interface {
	GetActors() []*Actor
	GetConsequences() []*Event
}

type ActiveAct struct {
	Description string
	Actors []*Actor
	Consequences []*Event
}

func (a *ActiveAct) GetCauseDescription() string {
	return a.Description
}

func (a *ActiveAct) GetConsequences() []*Event {
	return a.Consequences
}

func (a *ActiveAct) GetActors() []*Actor {
	return a.Actors
}

type NonAct struct {
	Description string
	Actors []*Actor
	Consequences []*Event
}

func (n *NonAct) GetCauseDescription() string {
	return n.Description
}

func (n *NonAct) GetConsequences() []*Event {
	return n.Consequences
}

func (n *NonAct) GetActors() []*Actor {
	return n.Actors
}

// --------------------------------------------------------------------
type Duty interface {
	GetDutyDescription() string
}

type ActiveDuty struct {
	Description string
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
type Event struct {
	Description string
	Consequences []*Event
}

func (e *Event) GetCauseDescription() string {
	return e.Description
}

// --------------------------------------------------------------------
type Cause interface {
	GetCauseDescription() string
}

// --------------------------------------------------------------------
type Injury interface {
	GetInjuryDescription() string
	GetDirectCauses() []*Cause
}

type BodilyInjury struct {
	Description string
	Causes []*Cause
}

func (b *BodilyInjury) GetInjuryDescription() string {
	return b.Description
}

func (b *BodilyInjury) GetDirectCauses() []*Cause {
	return b.Causes
}

type EmotionalHarm struct {
	Description string
	Causes []*Cause
}

func (e *EmotionalHarm) GetInjuryDescription() string {
	return e.Description
}

func (e *EmotionalHarm) GetDirectCauses() []*Cause {
	return e.Causes
}

type PropertyDamage struct {
	Description string
	Causes []*Cause
}

func (p *PropertyDamage) GetInjuryDescription() string {
	return p.Description
}

func (p *PropertyDamage) GetDirectCauses() []*Cause {
	return p.Causes
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
