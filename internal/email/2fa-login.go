package email

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type twoFAJSON struct {
	FromText  string `json:"from"`
	Subject   string `json:"subject"`
	PlainText string `json:"plainText"`
	Title     string `json:"title"`
	Body1     string `json:"body1"`
	Body2     string `json:"body2"`
	Body3     string `json:"body3"`
	Body4     string `json:"body4"`
	Body5     string `json:"body5"`
}

type twoFAParam struct {
	// common
	FromText           string
	From               string
	Host               string
	To                 string
	Subject            string
	MsgID              string
	PlainText          string
	Title              string
	OperatorLogo       string
	DownloadAppStore   string
	DownloadAPK        string
	DownloadGoogle     string
	DownloadTestFlight string
	OperatorLegal      string
	OperatorAddress    string
	OperatorContact    string
	// body
	B1, B2, B3, B4, Token, B5 string
	// footer
	Str1, Str2, Str3, Str4, Str5, Str6 string
}

type twoFALoginEmailInterface struct {
	JSON twoFAJSON
}

var twoFALoginEmail twoFALoginEmailInterface

func (s *twoFALoginEmailInterface) getEmailParam(user string, param Param, jsonData []byte) (interface{}, error) {
	err := json.Unmarshal(jsonData, &s.JSON)
	if err != nil {
		log.WithError(err).Errorf("Parse json data error")
		return nil, err
	}

	jsonStruct := twoFAJSON{
		FromText:  fmt.Sprintf(s.JSON.FromText, email.operator.operatorName),
		Subject:   s.JSON.Subject,
		PlainText: fmt.Sprintf(s.JSON.PlainText, email.operator.operatorName, param.Token),
		Title:     fmt.Sprintf(s.JSON.Title, email.operator.operatorName),
		Body1:     s.JSON.Body1,
		Body2:     fmt.Sprintf(s.JSON.Body2, email.operator.operatorName),
		Body3:     s.JSON.Body3,
		Body4:     s.JSON.Body4,
		Body5:     s.JSON.Body5,
	}

	emailData := twoFAParam{
		FromText:           jsonStruct.FromText,
		From:               email.from,
		Host:               email.host,
		To:                 user,
		Subject:            jsonStruct.Subject,
		MsgID:              param.messageID,
		PlainText:          jsonStruct.PlainText,
		Title:              jsonStruct.Title,
		OperatorLogo:       email.operator.operatorLogo,
		DownloadAppStore:   email.operator.downloadAppStore,
		DownloadGoogle:     email.operator.downloadGoogle,
		DownloadTestFlight: email.operator.downloadTestFlight,
		DownloadAPK:        email.operator.downloadAPK,
		OperatorLegal:      email.operator.operatorLegal,
		OperatorAddress:    email.operator.operatorAddress,
		OperatorContact:    email.operator.operatorContact,
		B1:                 jsonStruct.Body1,
		B2:                 jsonStruct.Body2,
		B3:                 jsonStruct.Body3,
		B4:                 jsonStruct.Body4,
		Token:              param.Token,
		B5:                 jsonStruct.Body5,
		Str1:               param.commonJSON.Str1,
		Str2:               param.commonJSON.Str2,
		Str3:               param.commonJSON.Str3,
		Str4:               param.commonJSON.Str4,
		Str5:               param.commonJSON.Str5,
		Str6:               param.commonJSON.Str6,
	}

	return emailData, nil
}
