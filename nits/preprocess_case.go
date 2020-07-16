package nits

type preprocessedCase struct {
	events                 map[*Event]interface{}
	persons                map[*Person]interface{}
	duties                 map[*Duty]interface{}
	irrelevantCauses       map[*IrrelevantCause]interface{}
	injuriesOrDamages      map[InjuryOrDamage]interface{}
	brokenLegalRequirement map[*BrokenLegalRequirement]interface{}
}

func (p *preprocessedCase) ppIrrelevantCause(event *Event, cause *IrrelevantCause) {
	cause.event = event
	p.irrelevantCauses[cause] = nil
}

func (p *preprocessedCase) ppPersons(persons []*Person) {
	for _, person := range persons {
		p.persons[person] = nil
	}
}

func (p *preprocessedCase) ppInjuriesOrDamages(event *Event, injuriesOrDamages []InjuryOrDamage) {
	for _, id := range injuriesOrDamages {
		if _, ok := p.injuriesOrDamages[id]; ok {
			continue
		}

		p.ppPersons(id.GetPersons())
	}
}

func (p *preprocessedCase) ppBrokenLegalRequirement(event *Event, negperse *BrokenLegalRequirement) {
	negperse.event = event
	p.ppEvents(negperse.Consequences)
	p.ppPersons(negperse.Persons)
}

func (p *preprocessedCase) ppDuty(event *Event, duty *Duty) {
	duty.event = event
	p.duties[duty] = nil
	p.ppPersons(duty.OwedFrom)
	p.ppPersons(duty.OwedTo)
}

func (p *preprocessedCase) ppEvent(e *Event) {
	if _, ok := p.events[e]; ok {
		return
	}

	p.events[e] = nil
	p.ppEvents(e.Consequences)

	if e.InjuriesOrDamages != nil {
		p.ppInjuriesOrDamages(e, e.InjuriesOrDamages)
	}

	if e.IrrelevantCause != nil {
		p.ppIrrelevantCause(e, e.IrrelevantCause)
	}

	if e.NegPerSe != nil {
		p.ppBrokenLegalRequirement(e, e.NegPerSe)
	}

	if e.Duty != nil {
		p.ppDuty(e, e.Duty)
	}
}

func (p *preprocessedCase) ppEvents(events []*Event) {
	for _, e := range events {
		p.ppEvent(e)
	}
}

func preprocess(c *Case) *preprocessedCase {
	p := &preprocessedCase{
		events:                 make(map[*Event]interface{}),
		persons:                make(map[*Person]interface{}),
		duties:                make(map[*Duty]interface{}),
		irrelevantCauses:       make(map[*IrrelevantCause]interface{}),
		injuriesOrDamages:      make(map[InjuryOrDamage]interface{}),
		brokenLegalRequirement: make(map[*BrokenLegalRequirement]interface{}),
	}
	p.ppEvents(c.RootEvents)
	return p
}
