package email

type stakingIncomeEmailInterface struct {}
var stakingIncomeEmail stakingIncomeEmailInterface

func (*stakingIncomeEmailInterface) sendEmail(user, token string, language EmailLanguage) error {

	return nil
}
