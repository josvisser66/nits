package nits

func test_case() *Case {
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

	rookesInjury := &BodilyInjury{
		Description: "Rooke suffers serious injuries because of being thrown from the car",
		Persons:     []*Person{rooke},
	}
	brucesDamage := &PropertyDamage{
		Description: "Bruce's car is seriously damaged because of the accident",
		Persons:     []*Person{bruce},
	}

	rookeGetsThrownFromTheCar := &Event{
		Description:       "Rooke gets thrown from the car",
		InjuriesOrDamages: []InjuryOrDamage{rookesInjury},
	}

	rookeShouldHaveWornASeatbelt := &BrokenLegalRequirement{
		Description:  "Rooke did not wear a seatbelt",
		Persons:      []*Person{rooke},
		Consequences: []*Event{rookeGetsThrownFromTheCar},
		Explanation: &Explanation{
			Text: []string{
				"Explanation here",
			},
		},
	}

	rookeGetsThrownFromTheCar.NegPerSe = rookeShouldHaveWornASeatbelt

	brucesCarPlowsIntoAshtonsCar := &Event{
		shortName: "plows",
		Description:       "Bruce's car plows into Ashton's car",
		Consequences: []*Event{rookeGetsThrownFromTheCar},
		InjuriesOrDamages: []InjuryOrDamage{brucesDamage},
		IrrelevantCause: &IrrelevantCause{
			Description: "Bruce claims that he did not see Ashton's car because of the truck in front of him",
			Explanation: &Explanation{
				Text: []string{
					"Explanation of why this is irrelevant",
				},
			},
		},
	}

	bruceHadDrankTooMuch := &BrokenLegalRequirement{
		Description:  "Bruce had drank too much and had blood alcohol levels over the legal limit",
		Persons:      []*Person{bruce},
		Consequences: []*Event{brucesCarPlowsIntoAshtonsCar},
		Explanation: &Explanation{
			Text: []string{
				"Explanation here",
			},
		},
	}

	brucesCarPlowsIntoAshtonsCar.NegPerSe = bruceHadDrankTooMuch

	ashtonDials911 := &Event{
		Description: "Ashton dials 911 and requests fire department and police assistance",
	}
	ashtonFleesTheCar := &Event{
		Description:  "Ashton abandons the car",
		Consequences: []*Event{brucesCarPlowsIntoAshtonsCar},
	}
	carDies := &Event{
		shortName: "car_dies",
		Description:  "The engine of Ashton's car dies",
		Consequences: []*Event{ashtonFleesTheCar, ashtonDials911},
	}
	smokeUnderHood := &Event{
		Description: "Smoke comes out from under the hood of Ashton's car",
	}
	continuesDriving := &Event{
		Description:  "Ashton continues to drive her car",
		Consequences: []*Event{smokeUnderHood, carDies},
	}
	demiGivesBadAdvice := &Event{
		Description:  "Demi advices Ashton to continue driving and to bring the car in at his convenience",
		Consequences: []*Event{continuesDriving},
		Duty:         giveGoodAdvice,
	}
	askingAdvice := &Event{
		Description:  "Ashton calls Demi and asks for advice",
		Consequences: []*Event{demiGivesBadAdvice},
	}
	oilLightGoesOn := &Event{
		Description:  "The low oil indicator in Ashton's car flips on",
		Consequences: []*Event{askingAdvice},
	}
	lowOilPressure := &Event{
		Description:  "Ashton's low on on oil",
		Consequences: []*Event{oilLightGoesOn},
	}
	badOilChange := &Event{
		Description:  "Demi (or Mayko) performs a bad oil change on Ashton's car",
		Consequences: []*Event{lowOilPressure},
		Duty:         doAGoodOilChange,
	}

	return &Case{
		ShortName:  "case_ashton_car_crash",
		RootEvents: []*Event{badOilChange},
	}
}

