package email

type twoFALoginEmailInterface struct{}

var twoFALoginEmail twoFALoginEmailInterface

func (*twoFALoginEmailInterface) sendEmail(user, token string, language EmailLanguage) error {

	return nil
}
