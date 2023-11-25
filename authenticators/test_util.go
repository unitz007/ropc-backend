package authenticators

import (
	"ropc-backend/model"
)

var (
	// RightUsername test user details
	RightUsername       = "right_username"
	RightPassword       = "right_password"
	WrongPassword       = "wrong_password"
	WrongUsername       = "wrong_username"
	hashedRightPassword = "$2a$12$JadjhGXumDBw.8X9o0.EaeNkIDaeGtmHkAmxfgrqApaFT0t.ZVrm."

	// test client details
	rightClientId     = "right_clientId"
	rightClientSecret = "right_clientSecret"
	WrongClientId     = "wrong_clientId"
	WrongClientSecret = "wrong_clientSecret"

	_ = model.Application{
		ClientId:     rightClientId,
		ClientSecret: rightClientSecret,
	}

	WrongClient = model.Application{
		ClientId:     WrongClientId,
		ClientSecret: WrongClientSecret,
	}
)
