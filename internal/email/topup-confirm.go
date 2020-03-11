package email

type topupConfirmEmailInterface struct{}

var topupConfirmEmail topupConfirmEmailInterface

func (*topupConfirmEmailInterface) sendEmail(user, token string, language EmailLanguage) error {

	return nil
}
