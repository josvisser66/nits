package nits

// This file contains the preprocessing code for a case.

import "math/rand"

// preprocessedCase is a struct that contains sets for all the individual
// objects of a given type in a case.
// In the act of creating this struct we also fill back links in the objects.
type preprocessedCase struct {
	events                 map[Event]interface{}
	persons                map[*Person]interface{}
	duties                 map[*Duty]interface{}
	claims                 map[*Claim]interface{}
	injuriesOrDamages      map[InjuryOrDamage]interface{}
	brokenLegalRequirement map[*BrokenLegalRequirement]interface{}
}

func (p *preprocessedCase) ppClaims(event Event, claims []*Claim) {
	for _, claim := range claims {
		claim.event = event
		p.claims[claim] = nil
	}
}

func (p *preprocessedCase) ppPersons(persons []*Person) {
	for _, person := range persons {
		p.persons[person] = nil
	}
}

func (p *preprocessedCase) ppInjuriesOrDamages(event Event, injuriesOrDamages []InjuryOrDamage) {
	for _, dam := range injuriesOrDamages {
		if _, ok := p.injuriesOrDamages[dam]; ok {
			continue
		}
		p.injuriesOrDamages[dam] = nil
		p.ppPersons(dam.GetPersons())
		dam.addCause(event)

		for _, person := range dam.GetPersons() {
			if person.damages == nil {
				person.damages = make(map[InjuryOrDamage]interface{}, 0)
			}
			person.damages[dam] = nil
		}
	}
}

func (p *preprocessedCase) ppBrokenLegalRequirement(event Event, negperse *BrokenLegalRequirement) {
	negperse.event = event
	p.ppEvents(nil, negperse.Consequences)
	p.ppPersons(negperse.Persons)
	p.brokenLegalRequirement[negperse] = nil
}

func (p *preprocessedCase) ppDuty(event Event, duty *Duty) {
	duty.event = event
	p.duties[duty] = nil
	p.ppPersons(duty.OwedFrom)
	p.ppPersons(duty.OwedTo)
}

func (p *preprocessedCase) ppEvent(parent Event, e Event) {
	if _, ok := p.events[e]; ok {
		return
	}

	e.addCause(parent)
	p.events[e] = nil
	p.ppEvents(e, e.getConsequences())

	if e.getInjuriesOrDamages() != nil {
		p.ppInjuriesOrDamages(e, e.getInjuriesOrDamages())
	}

	if e.getClaims() != nil {
		p.ppClaims(e, e.getClaims())
	}

	if e.getNegPerSe() != nil {
		p.ppBrokenLegalRequirement(e, e.getNegPerSe())
	}

	if e.getDuty() != nil {
		p.ppDuty(e, e.getDuty())
	}
}

func (p *preprocessedCase) ppEvents(parent Event, events []Event) {
	for _, e := range events {
		p.ppEvent(parent, e)
	}
}

// preprocess preprocesses a case. It starts at the root events of the
// case and collects all the individual objects into sets of their type.
// Along the way the back links in various objects are also filled in.
func (c *Case) preprocess() *preprocessedCase {
	if c.preproc != nil {
		return c.preproc
	}
	p := &preprocessedCase{
		events:                 make(map[Event]interface{}),
		persons:                make(map[*Person]interface{}),
		duties:                 make(map[*Duty]interface{}),
		claims:                 make(map[*Claim]interface{}),
		injuriesOrDamages:      make(map[InjuryOrDamage]interface{}),
		brokenLegalRequirement: make(map[*BrokenLegalRequirement]interface{}),
	}
	p.ppEvents(nil, c.RootEvents)
	c.preproc = p
	return p
}

// findEvent finds an event by short name. Used for unit testing.
func (p *preprocessedCase) findEvent(shortName string) Event {
	for event := range p.events {
		if event.getShortName() == shortName {
			return event
		}
	}
	return nil
}

// randomAct finds a random act.
func (p *preprocessedCase) randomAct() *Act {
	acts := make([]*Act, 0, len(p.events))
	for event := range p.events {
		if act, ok := event.(*Act); ok {
			acts = append(acts, act)
		}
	}

	return acts[rand.Int()%len(acts)]
}

// randomInjuryOrDamage finds a random injury or damage. Surprising, I know.
func (p *preprocessedCase) randomInjuryOrDamage() InjuryOrDamage {
	dams := make([]InjuryOrDamage, 0, len(p.injuriesOrDamages))
	for dam := range p.injuriesOrDamages {
		dams = append(dams, dam)
	}

	return dams[rand.Int()%len(dams)]
}
