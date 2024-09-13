package validation

// var (
// 	passwordRegex = map[string]*regexp.Regexp{
// 		"number":  regexp.MustCompile(`[0-9]`),
// 		"upper":   regexp.MustCompile(`[A-Z]`),
// 		"special": regexp.MustCompile(`[!@#$%^&*]`),
// 	}
// )

// // ValidateRegister validates register request
// func passwordValidator(fl validator.FieldLevel) bool {
// 	password := fl.Field().String()

// 	if password == "" {
// 		return false
// 	}
// 	if len(password) < 8 {
// 		return false
// 	}
// 	if !passwordRegex["number"].MatchString(password) ||
// 		!passwordRegex["upper"].MatchString(password) ||
// 		!passwordRegex["special"].MatchString(password) {
// 		return false
// 	}
// 	return true
// }

// func init() {
// 	// Регистрируем кастомный валидатор
// 	authhandler.Validate.RegisterValidation("password", passwordValidator)
// }
