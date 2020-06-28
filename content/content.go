package content

import "C"
import . "../nits"

func GetContent() *Content {
	foreseeability := &Concept{
		Name:  "foreseeability",
		Level: 1,
	}
	comparativeNegligence := &Concept{
		Name:  "comparative negligence",
		Level: 1,
	}
	modifiedComparativeNegligence := &Concept{
		Name:  "modified comparative negligence",
		Level: 1,
	}
	pureComparativeNegligence := &Concept{
		Name:  "pure comparative negligence",
		Level: 1,
	}
	contributoryNegligence := &Concept{
		Name:  "contributory negligence",
		Level: 1,
	}
	// Links all the contributory and comparative negligence concepts to
	// one another.
	comparativeNegligence.Related = []*Concept{modifiedComparativeNegligence, pureComparativeNegligence, contributoryNegligence}
	pureComparativeNegligence.Related = []*Concept{comparativeNegligence, modifiedComparativeNegligence, contributoryNegligence}
	modifiedComparativeNegligence.Related = []*Concept{comparativeNegligence, pureComparativeNegligence, contributoryNegligence}
	contributoryNegligence.Related = []*Concept{comparativeNegligence, modifiedComparativeNegligence, pureComparativeNegligence}
	preponderanceOfTheElements := &Concept{
		Name:  "preponderance of the elements",
		Level: 1,
	}
	assumptionOfRisk := &Concept{
		Name:"assumption of risk",
		Level: 1,
	}
	negligencePerSe := &Concept{
		Name:  "negligence per se",
		Level: 1,
		Explanation: &Explanation{
			Text: []string{
				"In order for there to be negligence per se, the defendant must have been acting in violation of a statute or regulation.",
				"To prove negligence per se the plaintiff must prove that the defendant was in violation of a statute or regulation, the the statue or regulation was designed to prevent the kind of harm that the plaintiff suffered, and that the plaintiff is in the class of people that the statute or regulation sought to protect.",
			},
			References: []Reference{
				&Restatement{288},
			},
		},
	}
	return &Content{
		Questions: []Question{
			&MultipleChoiceQuestion{
				Question: []string{"The concept of reasonable foreseeability is satisfied only if:"},
				Concepts: []*Concept{foreseeability},
				Answers: []*Answer{
					{
						Text: "The plaintiff has proved, beyond a reasonable doubt, that he or she in fact suffered a loss that was caused the defendant’s carelessness.",
						Concepts: []*Concept{preponderanceOfTheElements},
					},
					{
						Text: "The defendant’s behavior was virtually certain to inflict a loss on the plaintiff.",
					},
					{
						Text: "The plaintiff has proved, on a balance of probabilities, that he or she in fact suffered a loss that was actually caused the defendant’s carelessness.",
					},
					{
						Text: "The defendant’s behavior was more likely than not to inflict a loss on the plaintiff.",
					},
					{
						Text: "None of the above.",
					},
				},
			},
			&MultipleChoiceQuestion{
				Question: []string{
					"Assume that the state of East Delaware has a statute under which Ellen would recover $60,000 of her $300,000 in damages because a jury found her to be 80% negligent in the accident in which she was injured.",
					"East Delaware has adopted:",
				},
				Answers: []*Answer{
					{
						Text: "The defense of pure comparative negligence.",
						Concepts: []*Concept{pureComparativeNegligence},
						Correct: true,
					},
					{
						Text: "The defense of modified comparative negligence",
						Concepts: []*Concept{modifiedComparativeNegligence},
					},
					{
						Text: "The defense of contributory negligence.",
						Concepts: []*Concept{contributoryNegligence},
					},
					{
						Text: "The defense of assumption of risk.",
						Concepts: []*Concept{assumptionOfRisk},
					},
					{
						Text: "The defense of negligence per se",
						Concepts: []*Concept{negligencePerSe},
					},
				},
			},
			&MultipleChoiceQuestion{
				Question:      []string{
					"In response to a number of accidents involving pedestrians, a city enacted a statute making it illegal to walk through the business district other than on the sidewalk. The city also enacted a statute making it illegal for a business to obstruct the sidewalk in front of its establishment. Mr. Bean was walking along the sidewalk when he discovered that a store has stacked a pile of boxes such that the sidewalk was totally obstructed. Mr. Bean stepped into the street to walk around the boxes and was struck by a negligently driven taxi. This jurisdiction follows contributory negligence rules.",
					"If Mr. Bean asserts a claim against the taxi driver, what will be the effect of Mr. Bean’s leaving the sidewalk and walking in the street?",
				},
				Answers: []*Answer{
					{
						Text: "It will bar his recovery as a matter of law.",
						Concepts: []*Concept{contributoryNegligence, negligencePerSe},
					},
					{
						Text: "It will reduce his recovery.",
						Concepts: []*Concept{comparativeNegligence},
					},
					{
						Text: "It may be considered by the trier of fact on the issue of the taxi driver’s liability.",
					},
					{
						Text: "It is not relevant to determining Mr. Bean’s rights.",
					},
				},
			},
		},
	}
}
