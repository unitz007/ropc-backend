package kernel

import (
	ctx "context"
	"net/http"
	"ropc-backend/model"
	"ropc-backend/utils"
)

type Security struct {
	config         utils.Config
	userRepository Repository[model.User]
}

func NewSecurity(config utils.Config, userRepository Repository[model.User]) *Security {
	return &Security{
		config:         config,
		userRepository: userRepository,
	}
}

func (s *Security) Jwt(h func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get(tokenHeader)

		if accessToken == utils.Blank {
			panic(NewError(http.StatusUnauthorized, tokenHeaderErrorMsg+" for path: "+r.URL.String()))
		}

		token, err := utils.ValidateToken(accessToken, s.config.TokenSecret())

		if err != nil {
			panic(NewError(http.StatusUnauthorized, "token validation failed: "+err.Error()))
		}

		email := token["sub"].(string)
		conditions := utils.Queries[utils.WhereUsernameOrEmailIs](email)
		user, err := s.userRepository.Get(conditions)
		if err != nil {
			http.Error(w, "", http.StatusForbidden)
		}

		r = r.WithContext(ctx.WithValue(r.Context(), utils.UserKey, user))

		h(w, r)
	}
}

//func (s *Security) RegisterHandlers(handlers []func(path string, method string, handlerFunc func(w http.ResponseWriter, r *http.Request))) {
//	for _, handler := range handlers {
//		s.Jwt(handler)
//	}
//}
