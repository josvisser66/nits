package nits

import "sort"

// --------------------------------------------------------------------
type Concept struct {
	name        string
	shortName   string
	level       int
	explanation *Explanation
	related     []*Concept
	hints       []string
}

func (c *Concept) SortRelatedConcepts() {
	sort.Slice(c.related, func(i, j int) bool {
		return c.related[i].name < c.related[j].name
	})
}

func (c *Concept) GetReferenceText() string {
	return "Concept: " + c.name
}

var (
	Foreseeability = &Concept{
		name:      "foreseeability",
		shortName: "foreseeability",
		level:     1,
	}
	ComparativeNegligence = &Concept{
		name:      "comparative negligence",
		shortName: "compneg",
		level:     1,
		hints: []string{
			"Is the plaintiff negligent themselves?",
		},
		explanation: &Explanation{
			Text: []string{
				"Comparative negligence is a doctrine where the damages that the defendant is liable for " +
					"are reduced because the plaintiff is herself somewhat at fault for the injury or damage.",
			},
			References: []Reference{
				&URL{Url: "https://en.wikipedia.org/wiki/Comparative_negligence"},
			},
		},
	}
	ModifiedComparativeNegligence = &Concept{
		name:      "modified comparative negligence",
		shortName: "modcompneg",
		level:     1,
		explanation: &Explanation{Text: []string{
			"The doctrine of modified comparative negligence is a form of comparative negligence where there is " +
				"a threshold for the plaintiff's contribution to the injury or damage. There are two variants " +
				"of this doctrine, depending on whether an exact 50% culpability on the side of the plaintiff " +
				"bars recovery or not.",
		}},
		hints: []string{
			"What is the difference between this and pure comparative negligence?",
		},
	}
	PureComparativeNegligence = &Concept{
		name:      "pure comparative negligence",
		shortName: "purecompneg",
		level:     1,
		explanation: &Explanation{Text: []string{
			"In pure comparative negligence there is no threshold for barring the plaintiff for recovering " +
				"part of the damages, even though she is responsible for some (or a large) part of the " +
				"injury or property damages. For instance in pure comparative negligence you can recover 5% " +
				"of the damages if you yourself are 95% at fault.",
		}},
		hints: []string{
			"To what extent (percentage) is the plaintiff responsible for the injury or damage?",
			"Is there a threshold for the extent (percentage) that the plaintiff is responsible?",
		},
	}
	ContributoryNegligence = &Concept{
		name:      "contributory negligence",
		shortName: "contribneg",
		level:     1,
	}
	PreponderanceOfTheElements = &Concept{
		name:      "preponderance of the elements",
		shortName: "prepond",
		level:     1,
	}
	AssumptionOfRisk = &Concept{
		name:      "assumption of risk",
		shortName: "assumprisk",
		level:     1,
	}
	NegligencePerSe = &Concept{
		name:      "negligence per se",
		shortName: "negperse",
		level:     1,
		explanation: &Explanation{
			Text: []string{
				"In order for there to be negligence per se, the defendant must have been acting in violation of a " +
					"statute or regulation.",
				"To prove negligence per se the plaintiff must prove that the defendant was in violation of a " +
					"statute or regulation, the the statue or regulation was designed to prevent the kind of harm " +
					"that the plaintiff suffered, and that the plaintiff is in the class of people that the statute " +
					"or regulation sought to protect.",
			},
			References: []Reference{
				&Restatement{288},
			},
		},
	}
)

func initConcepts() {
	// Links all the contributory and comparative negligence concepts to
	// one another.
	ComparativeNegligence.related = []*Concept{
		ModifiedComparativeNegligence,
		PureComparativeNegligence,
		ContributoryNegligence,
	}
	ComparativeNegligence.SortRelatedConcepts()
	PureComparativeNegligence.related = []*Concept{
		ComparativeNegligence,
		ModifiedComparativeNegligence,
		ContributoryNegligence,
	}
	PureComparativeNegligence.SortRelatedConcepts()
	ModifiedComparativeNegligence.related = []*Concept{
		ComparativeNegligence,
		PureComparativeNegligence,
		ContributoryNegligence}
	ModifiedComparativeNegligence.SortRelatedConcepts()
	ContributoryNegligence.related = []*Concept{
		ComparativeNegligence,
		ModifiedComparativeNegligence,
		PureComparativeNegligence}
	ContributoryNegligence.SortRelatedConcepts()
}

// --------------------------------------------------------------------
func conceptMapToSlice(m map[*Concept]interface{}) []*Concept {
	keys := make([]*Concept, len(m))
	i := 0

	for k := range m {
		keys[i] = k
		i++
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].name < keys[j].name
	})

	return keys
}
