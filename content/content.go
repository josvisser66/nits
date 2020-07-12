package content

import "C"
import . "../nits"

func GetContent() *Content {
	return &Content{
		Questions: []Question{
			&MultipleChoiceQuestion{
				ShortName: "mc_reasonable_foreseeability1",
				Question:  []string{"The concept of reasonable foreseeability is satisfied only if:"},
				Concepts:  []*Concept{Foreseeability},
				Answers: []*Answer{
					{
						Text: "The plaintiff has proved, beyond a reasonable doubt, that he or she in fact " +
							"suffered a loss that was caused the defendant's carelessness.",
						Concepts: []*Concept{PreponderanceOfTheElements},
					},
					{
						Text:    "The defendant's behavior was virtually certain to inflict a loss on the plaintiff.",
						Correct: true,
					},
					{
						Text: "The plaintiff has proved, on a balance of probabilities, that he or she in fact " +
							"suffered a loss that was actually caused the defendant's carelessness.",
					},
					{
						Text: "The defendant's behavior was more likely than not to inflict a loss on the plaintiff.",
					},
					{
						Text:           "None of the above.",
						NoneOfTheAbove: true,
					},
				},
			},
			&MultipleChoiceQuestion{
				ShortName: "mc_pure_compneg1",
				Concepts:  []*Concept{PureComparativeNegligence},
				Question: []string{
					"Assume that the state of East Delaware has a statute under which Ellen would recover $60,000 " +
						"of her $300,000 in damages because a jury found her to be 80%% negligent in the accident " +
						"in which she was injured.",
					"East Delaware has adopted:",
				},
				Answers: []*Answer{
					{
						Text:     "The defense of pure comparative negligence.",
						Concepts: []*Concept{PureComparativeNegligence},
						Correct:  true,
					},
					{
						Text:     "The defense of modified comparative negligence",
						Concepts: []*Concept{ModifiedComparativeNegligence},
					},
					{
						Text:     "The defense of contributory negligence.",
						Concepts: []*Concept{ContributoryNegligence},
					},
					{
						Text:     "The defense of assumption of risk.",
						Concepts: []*Concept{AssumptionOfRisk},
					},
					{
						Text:     "The defense of negligence per se",
						Concepts: []*Concept{NegligencePerSe},
					},
				},
			},
			&MultipleChoiceQuestion{
				ShortName: "mc_mrbean_contribneg1",
				Concepts:  []*Concept{ContributoryNegligence},
				Question: []string{
					"In response to a number of accidents involving pedestrians, a city enacted a statute making it " +
						"illegal to walk through the business district other than on the sidewalk. The city also " +
						"enacted a statute making it illegal for a business to obstruct the sidewalk in front of " +
						"its establishment. Mr. Bean was walking along the sidewalk when he discovered that a store " +
						"has stacked a pile of boxes such that the sidewalk was totally obstructed. Mr. Bean stepped " +
						"into the street to walk around the boxes and was struck by a negligently driven taxi. " +
						"This jurisdiction follows contributory negligence rules.",
					"If Mr. Bean asserts a claim against the taxi driver, what will be the effect of Mr. Bean's " +
						"leaving the sidewalk and walking in the street?",
				},
				Answers: []*Answer{
					{
						Text:     "It will bar his recovery as a matter of law.",
						Concepts: []*Concept{ContributoryNegligence, NegligencePerSe},
						Correct:  true,
					},
					{
						Text:     "It will reduce his recovery.",
						Concepts: []*Concept{ComparativeNegligence},
					},
					{
						Text: "It may be considered by the trier of fact on the issue of the taxi driver's liability.",
					},
					{
						Text: "It is not relevant to determining Mr. Bean's rights.",
					},
				},
			},
			&PropsQuestion{
				ShortName: "props_contribneg_compneg1",
				Propositions: []*Proposition{
					{
						Proposition: "Under contributory negligence, if you contribute to your injury, you cannot recover damages.",
						Concepts:    []*Concept{ContributoryNegligence},
						True:        true,
					},
					{
						Proposition: "Under comparative negligence, if you contribute to your injury, you cannot recover damages.",
						Concepts:    []*Concept{ComparativeNegligence},
					},
				},
			},
			&MultipleChoiceQuestion{
				ShortName: "mc_amanda_surgery",
				Question: []string{
					"Amanda needed surgery on her right knee. When the anesthesia wore off after " +
						"the operation, she noticed surgical wrapping around both knees. When she asked " +
						"the nurse why both knees were wrapped, the nurse replied that the surgeon made " +
						"an incision on her left knee, discovered the mistake, and proceeded with the " +
						"operation on the right knee. What modification of the law of negligence will " +
						"Amanda probably be able to invoke to recover damages from the surgeon?",
				},
				Concepts: nil,
				Answers: []*Answer{
					{
						Text:     "Contributory negligence",
						Concepts: []*Concept{ContributoryNegligence},
					},
					{
						Text:     "Negligence per se",
						Concepts: []*Concept{NegligencePerSe},
					},
					{
						Text:     "Res ipsa loquitur",
						Concepts: []*Concept{ResIpsaLoquitur},
						Correct:  true,
					},
					{
						Text:     "Comparative negligence",
						Concepts: []*Concept{ComparativeNegligence},
					},
				},
			},
			&MultipleChoiceQuestion{
				ShortName: "mc_defense_liability_claims",
				Question:  []string{"All of the following are legal defenses to liability claims EXCEPT:"},
				Answers: []*Answer{
					{
						Text:     "Contributory negligence",
						Concepts: []*Concept{ContributoryNegligence},
					},
					{
						Text:     "Assumption of the risk",
						Concepts: []*Concept{AssumptionOfRisk},
					},
					{
						Text:     "Vicarious liability",
						Concepts: []*Concept{VicariousLiability},
						Correct:  true,
					},
					{
						Text:     "Comparative negligence",
						Concepts: []*Concept{ComparativeNegligence},
					},
				},
			},
			&MultipleChoiceQuestion{
				ShortName: "mc_managed_care",
				Question: []string{
					"A managed care company was found liable for denying valid claims for health " +
						"insurance coverage. The company was ordered to pay compensatory damages to " +
						"a group of plaintiffs. To make an example of the insurer, the court also " +
						"ordered the insurer to pay an additional $10 million to deter other insurers from " +
						"engaging in the same wrongful acts. The $10 million award is an example of:",
				},
				Answers: []*Answer{
					{
						Text:     "Punitive damages",
						Concepts: []*Concept{PunitiveDamages},
						Correct:  true,
					},
					{
						Text:     "Economic damages",
						Concepts: []*Concept{EconomicDamages},
					},
					{
						Text: "Non-economic damages",
					},
					{
						Text:     "Collateral source payments",
						Concepts: []*Concept{CollateralSourcePayments},
					},
				},
			},
			&MultipleChoiceQuestion{
				ShortName: "mc_pure_compneg2",
				Question: []string{
					"Bruce was involved in an accident in a state that uses a pure comparative " +
						"negligence rule. Bruce was found to be 75 percent responsible for the accident. " +
						"His actual damages were $20,000. How much will Bruce be able to recover from the " +
						"defendant?",
				},
				Concepts: []*Concept{PureComparativeNegligence},
				Answers: []*Answer{
					{
						Text: "$0",
					},
					{
						Text:"$5,000",
						Correct: true,
					},
					{
						Text:"$15,000",
					},
					{
						Text: "$20,000",
					},
				},
			},
		},
	}
}
