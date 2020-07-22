package nits

func isParentOf(suspectedParent, child Event) bool {
	if suspectedParent == child {
		return true
	}
	for _, cause := range child.getDirectCauses() {
		if isParentOf(suspectedParent, cause) {
			return true
		}
	}
	return false
}

func isCauseInFact(dam InjuryOrDamage, event Event) bool {
	for _, e := range dam.getDirectCauses() {
		if isParentOf(event, e) {
			return true
		}
	}
	return false
}

func intersectPersons(a, b []*Person) []*Person {
	m := make(map[*Person]int)

	for _, p := range a {
		m[p]++
	}

	for _, p := range b {
		m[p]++
	}

	result := make([]*Person, 0)

	for p, n := range m {
		if n > 1 {
			result = append(result, p)
		}
	}

	return result
}

func findDuties(dam InjuryOrDamage, causes []Event) []*Duty {
	seen := make(map[Event]interface{})
	result := make([]*Duty, 0)
	for len(causes) > 0 {
		next := make([]Event, 0)
		for _, e := range causes {
			seen[e] = nil
			for _, cause := range e.getDirectCauses() {
				if _, ok := seen[cause]; !ok {
					next = append(next, cause)
				}
			}
			if e.getDuty() == nil {
				continue
			}
			p := intersectPersons(dam.GetPersons(), e.getDuty().OwedTo)
			if len(p) > 0 {
				result = append(result, e.getDuty())
			}
		}
		causes = next
	}
	return result
}

func findDamages(events []Event) []InjuryOrDamage {
	result := make(map[InjuryOrDamage]interface{}, 0)
	for len(events) > 0 {
		next := make([]Event, 0)

		for _, e := range events {
			if dams := e.getInjuriesOrDamages(); dams != nil {
				for _, d := range dams {
					result[d] = nil
				}
			}

			next = append(next, e.getConsequences()...)
		}

		events = next
	}

	ret := make([]InjuryOrDamage, 0, len(result))

	for dam := range result {
		ret = append(ret, dam)
	}

	return ret
}

func collectPersonsFromDuties(duties []*Duty) map[*Person]interface{} {
	m := make(map[*Person]interface{})
	for _, d := range duties {
		for _, p := range d.OwedFrom {
			m[p] = nil
		}
	}
	return m
}
