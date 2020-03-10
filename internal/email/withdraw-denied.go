package email

type withdrawDeniedEmailInterface struct {}
var withdrawDeniedEmail withdrawDeniedEmailInterface

func (*withdrawDeniedEmailInterface) sendEmail(user, token string, language EmailLanguage) error {

	return nil
}
