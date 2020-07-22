package nits

// This files contains a default case. It can be used in NITS content but is
// also used for unit testing.

func DefaultCase() *Case {
	ashton := &Person{Name: "Ashton"}
	demi := &Person{Name: "Demi"}
	bruce := &Person{Name: "Bruce"}
	rooke := &Person{Name: "Rooke"}

	doAGoodOilChange := &Duty{
		Description: "Perform a good quality oil change",
		OwedFrom:    []*Person{demi},
		OwedTo:      []*Person{ashton},
	}
	giveGoodAdvice := &Duty{
		Description: "Give good advice",
		OwedFrom:    []*Person{demi},
		OwedTo:      []*Person{ashton},
	}
	leaveCarsSafely := &Duty {
		Description: "Cars should not be left in the middle of the road.",
		OwedFrom:    []*Person{ashton},
		OwedTo:      []*Person{bruce, rooke,demi},
	}

	rookesInjury := &BodilyInjury{
		Description: "Rooke suffers serious injuries because of being thrown from the car",
		Persons:     []*Person{rooke},
	}
	brucesDamage := &PropertyDamage{
		Description: "Bruce's car is seriously damaged because of the accident",
		Persons:     []*Person{bruce},
	}

	rookeGetsThrownFromTheCar := &PassiveEvent{
		Description:       "Rooke gets thrown from the car",
		InjuriesOrDamages: []InjuryOrDamage{rookesInjury},
	}

	rookeShouldHaveWornASeatbelt := &BrokenLegalRequirement{
		Description:  "Rooke did not wear a seatbelt",
		Persons:      []*Person{rooke},
		Consequences: []Event{rookeGetsThrownFromTheCar},
		Explanation: &Explanation{
			Text: []string{
				"Explanation here",
			},
		},
	}

	rookeGetsThrownFromTheCar.NegPerSe = rookeShouldHaveWornASeatbelt

	brucesCarPlowsIntoAshtonsCar := &Act{
		shortName:         "plows",
		Person:            bruce,
		Description:       "Bruce plows his car into Ashton's car",
		Consequences:      []Event{rookeGetsThrownFromTheCar},
		InjuriesOrDamages: []InjuryOrDamage{brucesDamage},
		Claims: []*Claim{
			{
				Person:      bruce,
				Description: "Bruce claims that he did not see Ashton's car because of the truck in front of him",
				Explanation: &Explanation{
					Text: []string{
						"Explanation of why this is irrelevant",
					},
				},
			},
		},
	}

	bruceHadDrankTooMuch := &BrokenLegalRequirement{
		Description:  "Bruce had drank too much and had blood alcohol levels over the legal limit",
		Persons:      []*Person{bruce},
		Consequences: []Event{brucesCarPlowsIntoAshtonsCar},
		Explanation: &Explanation{
			Text: []string{
				"Explanation here",
			},
		},
	}

	brucesCarPlowsIntoAshtonsCar.NegPerSe = bruceHadDrankTooMuch

	ashtonDials911 := &Act{
		Person:      ashton,
		Description: "Ashton dials 911 and requests fire department and police assistance",
	}
	ashtonFleesTheCar := &Act{
		Person:       ashton,
		Description:  "Ashton abandons the car",
		Duty: leaveCarsSafely,
		Consequences: []Event{brucesCarPlowsIntoAshtonsCar},
	}
	carDies := &PassiveEvent{
		shortName:    "car_dies",
		Description:  "The engine of Ashton's car dies",
		Consequences: []Event{ashtonFleesTheCar, ashtonDials911},
	}
	smokeUnderHood := &PassiveEvent{
		Description: "Smoke comes out from under the hood of Ashton's car",
	}
	continuesDriving := &Act{
		Person:       ashton,
		Description:  "Ashton continues to drive her car",
		Consequences: []Event{smokeUnderHood, carDies},
	}
	demiGivesBadAdvice := &Act{
		Person:       demi,
		Description:  "Demi advises Ashton to continue driving and to bring the car in at his convenience",
		Consequences: []Event{continuesDriving},
		Duty:         giveGoodAdvice,
	}
	askingAdvice := &Act{
		Person:       ashton,
		Description:  "Ashton calls Demi and asks for advice",
		Consequences: []Event{demiGivesBadAdvice},
	}
	oilLightGoesOn := &PassiveEvent{
		Description:  "The low oil indicator in Ashton's car flips on",
		Consequences: []Event{askingAdvice},
	}
	lowOilPressure := &PassiveEvent{
		Description:  "Ashton's car is low on on oil",
		Consequences: []Event{oilLightGoesOn},
	}
	badOilChange := &Act{
		Description:  "Demi (or Mayko) performs a bad oil change on Ashton's car",
		Consequences: []Event{lowOilPressure},
		Duty:         doAGoodOilChange,
	}

	return &Case{
		Text: []string{
			"Ashton left his home at 5:00 p.m. on Thursday, November 12, 2010, for a doctor's " +
				"appointment. His appointment was at 5:30 p.m. and it would take him at least 25 minutes to " +
				"reach his doctor's office. As Ashton pulled into traffic, he noted that the yellow low oil " +
				"pressure light on his dashboard was on. He was concerned, because he had just taken the car " +
				"in for a routine service and oil change at Mayko the day before. He pulled over to the side of " +
				"the road, pulled his receipt from the dashboard, and used his cell phone to call Demi, the " +
				"owner of Mayko. Demi assured Ashton that the light did not really mean that the oil pressure " +
				"was low because they had just changed it the day before. Instead, Demi said, the light was " +
				"probably just the result of a failure to reset a switch when they changed the oil or some sort of " +
				"short in the wiring. Demi advised Ashton to bring the car by at his convenience, and that she " +
				"would reset or repair the light.",

			"Relieved, Ashton continued down the highway toward his doctor's office. A few minutes " +
				"later, when Ashton was less than a mile from his doctor's office, he saw smoke coming from " +
				"the hood of the car. He tried to pull over to the side of the road, but before he could make it, " +
				"his engine died completely and the volume of smoke became even greater. Ashton dashed " +
				"from the car, leaving it in the right hand lane of traffic. A small fire erupted from the sides of " +
				"the hood. Again using his cell phone, Ashton dialed 911 and requested fire department and " +
				"police assistance. ",

			"But before fire or police units arrived, a car driven by Bruce plowed into the back of Ashton's " +
				"car. Bruce was not injured, but his passenger, Rooke, was thrown from the car and suffered " +
				"serious injuries. The police determined that although Bruce had been wearing a seatbelt at the " +
				"time of the collision, Rooke was not wearing a seatbelt. The police also determined that both " +
				"Bruce and Rooke, who had been drinking together all afternoon, had blood alcohol levels over " +
				"the legal limit. Bruce claimed that he did not see Ashton's car in time to stop because the 18- " +
				"wheeler in front of him had obscured his view of what was in the lane ahead. When the truck " +
				"changed lanes just before reaching Ashton's disabled car, Bruce was suddenly able to see " +
				"Ashton's car, but not in time to stop. An investigation reveals that Ashton's car stalled " +
				"because it ran out of oil. Demi had failed to replace the oil pan properly and all the oil in the " +
				"car had drained out.",
		},
		ShortName:  "case_ashton_car_crash",
		RootEvents: []Event{badOilChange},
	}
}
