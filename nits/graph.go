package nits

// This file contains some generic routines to walk to graph of a case.

// IsParentOf tests if a particular event is an indirect cause for another
// event.
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

// IsCauseInFact checks if an event lead directly or indirectly to
// some damage.
func isCauseInFact(dam InjuryOrDamage, event Event) bool {
	for _, e := range dam.getDirectCauses() {
		if isParentOf(event, e) {
			return true
		}
	}
	return false
}

// intersectPersons returns the slice of all persons that
// are both in a and b.
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

// findDuties finds a slice of breached duties that led to some
// damage. It walks the graph up from the direct causes of the
// damage and collects all the duties it finds on the way.
func findDuties(dam InjuryOrDamage) []*Duty {
	causes := dam.getDirectCauses()
	// seen is a set that will prevent us from recursing infinitely.
	seen := make(map[Event]interface{})
	result := make([]*Duty, 0)
	for len(causes) > 0 {
		// next will contain the direct causes of the events we are looking at now.
		next := make([]Event, 0)
		for _, e := range causes {
			// Registers this event as seen.
			seen[e] = nil
			// Goes through the direct causes of this event and maybe adds them
			// to the set of next events to process.
			for _, cause := range e.getDirectCauses() {
				if _, ok := seen[cause]; !ok {
					// If we haven't seen this direct cause yet we'll add it to the
					// set of events to investigate next.
					next = append(next, cause)
				}
			}
			// If this event has a duty and that duty was owed to the people that
			// suffered the damage then append the duty to the result set.
			if e.getDuty() == nil {
				continue
			}
			p := intersectPersons(dam.GetPersons(), e.getDuty().OwedTo)
			if len(p) > 0 {
				result = append(result, e.getDuty())
			}
		}
		// Done with this level of the tree. Move one level up.
		causes = next
	}
	return result
}

// findDamages collects all damages that are in the direct and indirect
// consequences of a set of events.
// Note: This method does not have protection against infinite recursion.
func findDamages(events []Event) []InjuryOrDamage {
	// Collect them in a set, to prevent duplicates.
	result := make(map[InjuryOrDamage]interface{}, 0)
	for len(events) > 0 {
		// Keep track of the next level down to process next.
		next := make([]Event, 0)

		// Collects all the damages directly linked to the events
		// at this level in the tree into the result set.
		for _, e := range events {
			if dams := e.getInjuriesOrDamages(); dams != nil {
				for _, d := range dams {
					result[d] = nil
				}
			}

			// Adds all the consequences of this event to the next
			// level of the tree to process.
			next = append(next, e.getConsequences()...)
		}

		// Done here, process the next level down.
		events = next
	}

	// Collects the set into a slice and returns it.
	ret := make([]InjuryOrDamage, 0, len(result))

	for dam := range result {
		ret = append(ret, dam)
	}

	return ret
}

// Collects all persons that owe the duties in the slice.
func collectPersonsFromDuties(duties []*Duty) map[*Person]interface{} {
	m := make(map[*Person]interface{})
	for _, d := range duties {
		for _, p := range d.OwedFrom {
			m[p] = nil
		}
	}
	return m
}
