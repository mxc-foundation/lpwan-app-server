package email

type twoFAWithdrawEmailInterface struct{}

var twoFAWithdrawEmail twoFAWithdrawEmailInterface

func (*twoFAWithdrawEmailInterface) sendEmail(user, token string, language EmailLanguage) error {

	return nil
}
