package email

type withdrawSuccessEmailInterface struct {}
var withdrawSuccessEmail withdrawSuccessEmailInterface

func (*withdrawSuccessEmailInterface) sendEmail(user, token string, language EmailLanguage) error {

	return nil
}
