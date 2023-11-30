package middlewares

//
//import (

//)
//
//
//type Security struct {
//	config         utils.Config
//	userRepository repositories.UserRepository
//}
//
//func NewSecurity(config utils.Config, repository repositories.UserRepository) *Security {
//	return &Security{config: config, userRepository: repository}
//}
//
//func (s *Security) TokenValidation(h func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
//	return func(w http.ResponseWriter, r *http.Request) {
//		accessToken := r.Header.Get(tokenHeader)
//
//		if accessToken == "" {
//			panic(errors.New(tokenHeaderErrorMsg + " for path: " + r.URL.String()))
//		}
//
//		token, err := utils.ValidateToken(accessToken, s.config.TokenSecret())
//
//		if err != nil {
//			panic(errors.New("token validation failed: " + err.Error()))
//		}
//
//		email := token["sub"].(string)
//		user, err := s.userRepository.GetUser(email)
//		if err != nil {
//			http.Error(w, "", http.StatusForbidden)
//		}
//
//		r = r.WithContext(context.WithValue(r.Context(), handlers.UserKey, user))
//
//		h(w, r)
//	}
//}
