package authenticators

import (
	"backend-server/model"
)

var (
	// test user details
	RightUsername       = "right_username"
	RightPassword       = "right_password"
	WrongPassword       = "wrong_password"
	WrongUsername       = "wrong_username"
	hashedRightPassword = "$2a$12$JadjhGXumDBw.8X9o0.EaeNkIDaeGtmHkAmxfgrqApaFT0t.ZVrm."

	// test client details
	rightClientId           = "right_clientId"
	rightClientSecret       = "right_clientSecret"
	WrongClientId           = "wrong_clientId"
	WrongClientSecret       = "wrong_clientSecret"
	hashedRightClientSecret = "$2a$12$0TUqBhM9DdBw980nTxz1EuL3eM/jQQSABDVuO6/lrCsjuUCOCFdxy"

	_ = model.Application{
		ClientId:     rightClientId,
		ClientSecret: rightClientSecret,
	}

	WrongClient = model.Application{
		ClientId:     WrongClientId,
		ClientSecret: WrongClientSecret,
	}
)
