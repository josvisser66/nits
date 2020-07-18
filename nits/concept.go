package nits

import (
	"fmt"
	"sort"
)

// --------------------------------------------------------------------
type Concept struct {
	name        string
	shortName   string
	level       int
	explanation *Explanation
	related     []*Concept
	hints       []string
}

var allConcepts = make([]*Concept, 0)

func (c *Concept) add() *Concept {
	allConcepts = append(allConcepts, c)
	return c
}
func (c *Concept) sortRelatedConcepts() {
	sort.Slice(c.related, func(i, j int) bool {
		return c.related[i].name < c.related[j].name
	})
}

func (c *Concept) GetReferenceText() string {
	return "Concept: " + c.name
}

var (
	Foreseeability1 = (&Concept{
		name:      "foreseeability (basic)",
		shortName: "foreseeability1",
		level:     1,
	}).add()
	CauseInFact1 = (&Concept{
		name: "cause in fact (basic)",
		shortName: "causeinfact1",
		level: 1,
		explanation: &Explanation{
			Text:       []string{
				"Cause-in-fact causation requires a plaintiff to show that he or she would not have been " +
					"injured but for the defendant's actions. The essential question in determining the " +
					"cause-in-fact is whether the plaintiff's injuries would have resulted regardless of "+
					"the defendant's negligence.",
			},
		},
	}).add()
	ComparativeNegligence1 = (&Concept{
		name:      "comparative negligence",
		shortName: "compneg1",
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
	}).add()
	ModifiedComparativeNegligence1 = (&Concept{
		name:      "modified comparative negligence",
		shortName: "modcompneg1",
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
	}).add()
	PureComparativeNegligence1 = (&Concept{
		name:      "pure comparative negligence",
		shortName: "purecompneg1",
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
	}).add()
	ContributoryNegligence1 = (&Concept{
		name:      "contributory negligence",
		shortName: "contribneg1",
		level:     1,
	}).add()
	PreponderanceOfTheElements1 = (&Concept{
		name:      "preponderance of the elements",
		shortName: "prepond1",
		level:     1,
	}).add()
	AssumptionOfRisk1 = (&Concept{
		name:      "assumption of risk",
		shortName: "assumprisk1",
		level:     1,
	}).add()
	NegligencePerSe1 = (&Concept{
		name:      "negligence per se",
		shortName: "negperse1",
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
				&Restatement{"288"},
			},
		},
	}).add()
	ResIpsaLoquitur1 = (&Concept{
			name:      "res ipsa loquitur (basic)",
			shortName: "resipsa1",
			level:     1,
			explanation: &Explanation{
				Text: []string{
					"Res ipsa Loquitur: The thing speaks for itself.",
					"The doctrine of res ipsa loquitur can be called in when, even though it is not exactly " +
						"known what happened it is obvious that something negligent happened. The doctrinal " +
						"example is a sack of flour falling from above and landing on someone's head.",
				},
			},
		}).add()
	VicariousLiability1 = (&Concept{
		name:        "vicarious liability",
		shortName:   "vicliab1",
		level:       1,
		explanation: &Explanation{
			Text: []string{
				"Vicarious liability is a legal doctrine whereby a person who is not personally at fault is " +
					"legally required to bear the burden of another's tortious wrongdoing.",
			},
		},
	}).add()
	PunitiveDamages1 = (&Concept{
		name: "punitive damages",
		shortName: "pundam1",
		level: 1,
		explanation: &Explanation{
			Text: []string{
				"Punitive damages (as opposed to compensatory damages) are designed to prevent others from "+
					"being hurt by the same or similar actions.",
			},
		},
	}).add()
	EconomicDamages1 = (&Concept{
		name:        "economic damages",
		shortName:   "ecodam1",
		level:       1,
		explanation: &Explanation{
			Text:[]string{
				"Economic damages are compensation you receive as a result of monetary losses you suffer " +
					"because of an accident.",
			},
		},
	}).add()
	CollateralSourcePayments1 = (&Concept{
		name:        "collateral source payments",
		shortName:   "collsrcpay1",
		level:       1,
		explanation: &Explanation{
			Text:[]string{
				"Collateral source payments are payments a plaintiff might receive from for instance an "+
					"insurance company for part or all of the damages. The Collateral Source Rule states " +
					"that no evidence of collateral source payments may be introduced to the jury. Because "+
					"of this these payments do not reduce the liability for the tortfeasor. Many states have "+
					"abrogated this rule by statute.",
			},
			References: []Reference{
				&Restatement{"920A"},
				&URL{"https://www.claimsjournal.com/news/national/2018/01/11/282417.htm"},
			},
		},
	}).add()
)

func initConcepts() {
	// Links all the contributory and comparative negligence concepts to
	// one another.
	ComparativeNegligence1.related = []*Concept{
		ModifiedComparativeNegligence1,
		PureComparativeNegligence1,
		ContributoryNegligence1,
	}
	PureComparativeNegligence1.related = []*Concept{
		ComparativeNegligence1,
		ModifiedComparativeNegligence1,
		ContributoryNegligence1,
	}
	ModifiedComparativeNegligence1.related = []*Concept{
		ComparativeNegligence1,
		PureComparativeNegligence1,
		ContributoryNegligence1,
	}
	ContributoryNegligence1.related = []*Concept{
		ComparativeNegligence1,
		ModifiedComparativeNegligence1,
		PureComparativeNegligence1,
	}

	for _, c := range allConcepts {
		c.sortRelatedConcepts()
	}
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

func checkConcepts() {
	m := make(map[string]interface{})

	for _, c := range allConcepts {
		name := c.shortName
		if _, ok := m[name]; ok {
			panic(fmt.Sprintf("Duplicate concept short name: %s", name))
		}
		m[name] = nil
	}
}
