package content

import "C"
import "../nits"

func GetContent() *nits.Content {
	foreseeability := &nits.Concept{
		Name:  "foreseeability",
		Level: 1,
	}
	comparativeNegligence := &nits.Concept{
		Name:  "comparative negligence",
		Level: 1,
	}
	pureComparativeNegligence := &nits.Concept{
		Name:  "pure comparative negligence",
		Level: 1,
	}
	contributoryNegligence := &nits.Concept{
		Name:  "contributory negligence",
		Level: 1,
	}
	return &nits.Content{
		Questions: []nits.Question{
			&nits.MultipleChoiceQuestion{
				Concepts: []*nits.Concept{foreseeability},
				Question: "The concept of reasonable foreseeability is satisfied only if:",
				Answers: []string{
					/*0*/ "The plaintiff has proved, beyond a reasonable doubt, that he or she in fact suffered a loss that was caused the defendant’s carelessness.",
					/*1*/ "The defendant’s behavior was virtually certain to inflict a loss on the plaintiff.",
					/*2*/ "The plaintiff has proved, on a balance of probabilities, that he or she in fact suffered a loss that was actually caused the defendant’s carelessness.",
					/*3*/ "The defendant’s behavior was more likely than not to inflict a loss on the plaintiff.",
					/*4*/ "None of the above.",
				},
				CorrectAnswer: 3,
			},
			&nits.MultipleChoiceQuestion{
				Concepts: []*nits.Concept{
					comparativeNegligence, pureComparativeNegligence,
				},
				Question: "Assume that the state of East Delaware has a statute under which Ellen would recover $60,000 of her $300,000 in damages because a jury found her to be 80% negligent in the accident in which she was injured. ",
				Answers: []string{
					/*0*/ "The defense of pure comparative negligence.",
					/*1*/ "The defense of modified comparative negligence.",
					/*2*/ "The defense of contributory negligence.",
					/*3*/ "The defense of assumption of the risk.",
					/*4*/ "The defense of negligence per se",
				},
				CorrectAnswer: 0,
			},
			&nits.MultipleChoiceQuestion{
				Concepts:      []*nits.Concept{contributoryNegligence},
				Question:      "In response to a number of accidents involving pedestrians, a city enacted a statute making it illegal to walk through the business district other than on the sidewalk. The city also enacted a statute making it illegal for a business to obstruct the sidewalk in front of its establishment. Mr. Bean was walking along the sidewalk when he discovered that a store has stacked a pile of boxes such that the sidewalk was totally obstructed. Mr. Bean stepped into the street to walk around the boxes and was struck by a negligently driven taxi. This jurisdiction follows contributory negligence rules",
				Answers: []string{
					/*0*/ "It will bar his recovery as a matter of law.",
					/*1*/ "It will reduce his recovery.",
					/*2*/ "It may be considered by the trier of fact on the issue of the taxi driver’s liability.",
					/*3*/ "It is not relevant to determining Mr. Bean’s rights.",
				},
				CorrectAnswer: 0,
			},
		},
	}
}
