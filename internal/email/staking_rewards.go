package email

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

type stakingIncomeJSON struct {
	FromText  string `json:"from"`
	Subject   string `json:"subject"`
	PlainText string `json:"plainText"`
	Title     string `json:"title"`
	Body1     string `json:"body1"`
	Body2     string `json:"body2"`
	Body3     string `json:"body3"`
	Body4     string `json:"body4"`
	Body5     string `json:"body5"`
	Body6     string `json:"body6"`
	Body7     string `json:"body7"`
	Body8     string `json:"body8"`
	Body9     string `json:"body9"`
	Body10    string `json:"body10"`
}

type stakingIncomeParam struct {
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
	B1, B2, B3, B4, B5, StakeRevenueDate, B6, StakeIncomeAmount, B7, StakeTotalIncome, B8, B9, B10, StakingInterest string
	// footer
	Str1, Str2, Str3, Str4, Str5, Str6 string
}

type stakingIncomeEmailInterface struct {
	JSON stakingIncomeJSON
}

var stakingIncomeEmail stakingIncomeEmailInterface

const (
	StakeIncomeAmount      string = "stakeIncomeAmount"
	StakeIncomeTotalAmount string = "stakeIncomeTotalAmount"
	StakeIncomeInterest    string = "stakeIncomeInterest"

	UserID         string = "userID"
	StakeID        string = "stakeID"
	StakeRevenueID string = "stakeRevenueID"

	StakeRevenueDate string = "stakeRevenueDate"
	StakeStartDate   string = "stakeStartDate"
)

func paramCheck(param Param) error {
	if param.Amount[StakeIncomeAmount] == "" {
		return errors.New("StakeIncomeAmount")
	}
	if param.Amount[StakeIncomeTotalAmount] == "" {
		return errors.New("StakeIncomeTotalAmount")
	}
	if param.Amount[StakeIncomeInterest] == "" {
		return errors.New("StakeIncomeInterest")
	}
	if param.ItemID[UserID] == "" {
		return errors.New("UserID")
	}
	if param.ItemID[StakeID] == "" {
		return errors.New("StakeID")
	}
	if param.ItemID[StakeRevenueID] == "" {
		return errors.New("StakeRevenueID")
	}
	if param.Date[StakeRevenueDate] == "" {
		return errors.New("StakeRevenueDate")
	}
	if param.Date[StakeStartDate] == "" {
		return errors.New("StakeStartDate")
	}

	return nil
}

func (s *stakingIncomeEmailInterface) getEmailParam(user string, param Param, jsonData []byte) (interface{}, error) {
	if err := paramCheck(param); err != nil {
		return nil, errors.Wrap(err, "invalid parameter for stakingIncomeEmailInterface")
	}

	err := json.Unmarshal(jsonData, &s.JSON)
	if err != nil {
		log.WithError(err).Errorf("Parse json data error")
		return nil, err
	}

	jsonStruct := stakingIncomeJSON{
		FromText: fmt.Sprintf(s.JSON.FromText, email.operator.operatorName),
		Subject:  s.JSON.Subject,
		PlainText: fmt.Sprintf(s.JSON.PlainText, param.ItemID[UserID], param.ItemID[StakeRevenueID],
			param.Date[StakeRevenueDate], param.Amount[StakeIncomeAmount], param.Amount[StakeIncomeTotalAmount]),
		Title:  fmt.Sprintf(s.JSON.Title, email.operator.operatorName),
		Body1:  s.JSON.Body1,
		Body2:  s.JSON.Body2,
		Body3:  fmt.Sprintf(s.JSON.Body3, param.ItemID[UserID]),
		Body4:  s.JSON.Body4,
		Body5:  fmt.Sprintf(s.JSON.Body5, param.ItemID[StakeRevenueID]),
		Body6:  s.JSON.Body6,
		Body7:  s.JSON.Body7,
		Body8:  fmt.Sprintf(s.JSON.Body8, param.ItemID[StakeID]),
		Body9:  fmt.Sprintf(s.JSON.Body9, param.Date[StakeStartDate]),
		Body10: s.JSON.Body10,
	}

	emailData := stakingIncomeParam{
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
		B5:                 jsonStruct.Body5,
		StakeRevenueDate:   param.Date[StakeRevenueDate],
		B6:                 jsonStruct.Body6,
		StakeIncomeAmount:  param.Amount[StakeIncomeAmount],
		B7:                 jsonStruct.Body7,
		StakeTotalIncome:   param.Amount[StakeIncomeTotalAmount],
		B8:                 jsonStruct.Body8,
		B9:                 jsonStruct.Body9,
		B10:                jsonStruct.Body10,
		StakingInterest:    param.Amount[StakeIncomeInterest],
		Str1:               param.commonJSON.Str1,
		Str2:               param.commonJSON.Str2,
		Str3:               param.commonJSON.Str3,
		Str4:               param.commonJSON.Str4,
		Str5:               param.commonJSON.Str5,
		Str6:               param.commonJSON.Str6,
	}

	return emailData, err
}
