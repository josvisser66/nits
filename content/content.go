package content

import "C"
import . "../nits"

// A simple case courtesy of https://www.lawnerds.com/testyourself/torts_exam.html
func case2() *Case {
	david := &Person{Name: "David"}
	peter := &Person{Name: "Peter"}
	teleco := &Person{Name: "Teleco"}
	kevin := &Person{Name: "Kevin"}

	buildSafePoles := &Duty{
		Description: "Utilities building public infrastructure should do so safely",
		OwedFrom:    []*Person{teleco},
		OwedTo:      []*Person{peter},
	}
	lookBeforeSwitchingLanes := &Duty{
		Description: "One should always look before changing lanes",
		OwedFrom:    []*Person{david},
		OwedTo:      []*Person{peter},
	}

	driveCarefully := &Duty{
		Description: "Drive carefully, especially in the presence of playing children",
		OwedFrom:    []*Person{david},
		OwedTo:      []*Person{kevin},
	}

	petersCar1 := &PropertyDamage{
		Description: "David's car hits Peter's car",
		Persons:     []*Person{peter},
	}
	petersCar2 := &PropertyDamage{
		Description: "Peters car crashes into the telephone pole",
		Persons:     []*Person{peter},
	}
	poleBroken := &PropertyDamage{
		Description: "The telephone pole snaps in two",
		Persons:     []*Person{teleco},
	}
	kevinsInjury := &BodilyInjury{
		Description: "Kevin sustains bodily injuries",
		Persons:     []*Person{kevin},
	}

	telephonePoleHitsKevin := &PassiveEvent{
		Description:       "A piece of the broken telephone pole hits Kevin",
		InjuriesOrDamages: []InjuryOrDamage{kevinsInjury},
		Claims:            nil,
	}
	telephonePoleSnapsInTwo := &PassiveEvent{
		Description:       "The telephone pole snaps in two",
		Consequences:      []Event{telephonePoleHitsKevin},
		Duty:              buildSafePoles,
		InjuriesOrDamages: []InjuryOrDamage{poleBroken},
	}
	petersCarHitsTelephonePole := &PassiveEvent{
		Description:       "Peter's car hits a telephone pole",
		Consequences:      []Event{telephonePoleSnapsInTwo},
		InjuriesOrDamages: []InjuryOrDamage{petersCar2},
		Claims:            nil,
	}
	peterLosesControlOfTheCar := &Act{
		Description:  "Peter loses control of his car",
		Person:       peter,
		Consequences: []Event{petersCarHitsTelephonePole},
	}
	davidsCarHitsPetersCar := &PassiveEvent{
		Description:       "David's car hits Peter's car",
		Consequences:      []Event{peterLosesControlOfTheCar},
		InjuriesOrDamages: []InjuryOrDamage{petersCar1},
	}
	speeding := &BrokenLegalRequirement{
		Description:  "Peter was speeding and overtaking David",
		Persons:      []*Person{peter},
		Consequences: []Event{davidsCarHitsPetersCar},
	}
	peterIsOvertakingAndSpeeding := &Act{
		Description:  "Peter is speeding and overtaking David in the lefthand lane",
		Person:       peter,
		Consequences: []Event{},
		NegPerSe:     speeding,
	}
	davidSwervesIntoTheOtherLane := &Act{
		Description:  "David, without looking, swerves into the lane left of him",
		Person:       david,
		Consequences: []Event{davidsCarHitsPetersCar},
		Duty:         lookBeforeSwitchingLanes,
	}
	kevinRunsIntoTheStreet := &Act{
		Description:  "Kevin runs into the street without looking",
		Person:       kevin,
		Consequences: []Event{davidSwervesIntoTheOtherLane},
	}
	peterDriving := &Act{
		Description:  "Peter is driving 25 MPH in a 25 MPH street where there are children playing",
		Person:       peter,
		Consequences: []Event{davidSwervesIntoTheOtherLane},
		Duty:         driveCarefully,
	}

	return &Case{
		ShortName: "case_teleco",
		RootEvents: []Event{
			peterDriving,
			kevinRunsIntoTheStreet,
			peterIsOvertakingAndSpeeding,
		},
		Text: []string{
			"David is driving 25 MPH in 25 MPH zone down a four lane street where there are children playing. " +
				"One nine-year-old child, Kevin, runs into the street chasing a soccer ball. David, without " +
				"looking over his shoulder, swerves into the other lane to avoid Kevin and in the process he hits " +
				"a car, driven by Peter, that was speeding past him in the left-hand lane going in the same direction.",
			"Peter loses control of his car, hits a telephone pole and is seriously and permanently injured. The " +
				"telephone pole, owned by the local phone company TeleCo, easily snaps into two pieces and hits Kevin, " +
				"who is still in the street, knocking him unconscious and resulting in permanent injuries.",
			"TeleCo never did any testing of its poles to establish how easily the poles broke. " +
				"The only factor used in manufacturing the poles was cost. The poles were made of low quality trees" +
				"and were not treated in any significant manner except for a coating of tar. No reinforcement was " +
				"used on the poles.",
		},
	}
}

// GetContent returns the content that NITS operates on.
// This is test content containing a few multiple choice / proposition
// questions and two cases.
func GetContent() *Content {
	return &Content{
		Questions: []Question{
			&MultipleChoiceQuestion{
				ShortName: "mc_reasonable_foreseeability1",
				Question:  []string{"The concept of reasonable foreseeability is satisfied only if:"},
				Concepts:  []*Concept{Foreseeability1},
				Answers: []*Answer{
					{
						Text: "The plaintiff has proved, beyond a reasonable doubt, that he or she in fact " +
							"suffered a loss that was caused the defendant's carelessness.",
						Concepts: []*Concept{PreponderanceOfTheElements1},
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
				Concepts:  []*Concept{PureComparativeNegligence1},
				Question: []string{
					"Assume that the state of East Delaware has a statute under which Ellen would recover $60,000 " +
						"of her $300,000 in damages because a jury found her to be 80%% negligent in the accident " +
						"in which she was injured.",
					"East Delaware has adopted:",
				},
				Answers: []*Answer{
					{
						Text:     "The defense of pure comparative negligence.",
						Concepts: []*Concept{PureComparativeNegligence1},
						Correct:  true,
					},
					{
						Text:     "The defense of modified comparative negligence",
						Concepts: []*Concept{ModifiedComparativeNegligence1},
					},
					{
						Text:     "The defense of contributory negligence.",
						Concepts: []*Concept{ContributoryNegligence1},
					},
					{
						Text:     "The defense of assumption of risk.",
						Concepts: []*Concept{AssumptionOfRisk1},
					},
					{
						Text:     "The defense of negligence per se",
						Concepts: []*Concept{NegligencePerSe1},
					},
				},
			},
			&MultipleChoiceQuestion{
				ShortName: "mc_mrbean_contribneg1",
				Concepts:  []*Concept{ContributoryNegligence1},
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
						Concepts: []*Concept{ContributoryNegligence1, NegligencePerSe1},
						Correct:  true,
					},
					{
						Text:     "It will reduce his recovery.",
						Concepts: []*Concept{ComparativeNegligence1},
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
						Concepts:    []*Concept{ContributoryNegligence1},
						True:        true,
					},
					{
						Proposition: "Under comparative negligence, if you contribute to your injury, you cannot recover damages.",
						Concepts:    []*Concept{ComparativeNegligence1},
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
						Concepts: []*Concept{ContributoryNegligence1},
					},
					{
						Text:     "Negligence per se",
						Concepts: []*Concept{NegligencePerSe1},
					},
					{
						Text:     "Res ipsa loquitur",
						Concepts: []*Concept{ResIpsaLoquitur1},
						Correct:  true,
					},
					{
						Text:     "Comparative negligence",
						Concepts: []*Concept{ComparativeNegligence1},
					},
				},
			},
			&MultipleChoiceQuestion{
				ShortName: "mc_defense_liability_claims",
				Question:  []string{"All of the following are legal defenses to liability claims EXCEPT:"},
				Answers: []*Answer{
					{
						Text:     "Contributory negligence",
						Concepts: []*Concept{ContributoryNegligence1},
					},
					{
						Text:     "Assumption of the risk",
						Concepts: []*Concept{AssumptionOfRisk1},
					},
					{
						Text:     "Vicarious liability",
						Concepts: []*Concept{VicariousLiability1},
						Correct:  true,
					},
					{
						Text:     "Comparative negligence",
						Concepts: []*Concept{ComparativeNegligence1},
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
						Concepts: []*Concept{PunitiveDamages1},
						Correct:  true,
					},
					{
						Text:     "Economic damages",
						Concepts: []*Concept{EconomicDamages1},
					},
					{
						Text: "Non-economic damages",
					},
					{
						Text:     "Collateral source payments",
						Concepts: []*Concept{CollateralSourcePayments1},
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
				Concepts: []*Concept{PureComparativeNegligence1},
				Answers: []*Answer{
					{
						Text: "$0",
					},
					{
						Text:    "$5,000",
						Correct: true,
					},
					{
						Text: "$15,000",
					},
					{
						Text: "$20,000",
					},
				},
			},
			DefaultCase(),
			case2(),
		},
	}
}
