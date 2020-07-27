package email

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type passwordResetUnknownJSON struct {
	FromText  string `json:"from"`
	Subject   string `json:"subject"`
	PlainText string `json:"plainText"`
	Title     string `json:"title"`
	Body1     string `json:"body1"`
	Body2     string `json:"body2"`
	Body3     string `json:"body3"`
	Body4     string `json:"body4"`
}

type passwordResetUnknownParam struct {
	// common
	FromText         string
	From             string
	Host             string
	To               string
	Subject          string
	MsgID            string
	PlainText        string
	Title            string
	OperatorLogo     string
	DownloadAppStore string
	DownloadGoogle   string
	OperatorLegal    string
	OperatorAddress  string
	OperatorContact  string
	// body
	B1, B2, B3, B4 string
	// footer
	Str1, Str2, Str3, Str4 string
}

type passwordResetUnknownEmailInterface struct {
	JSON passwordResetUnknownJSON
}

var passwordResetUnknownEmail passwordResetUnknownEmailInterface

func (s *passwordResetUnknownEmailInterface) getEmailParam(user string, param Param, jsonData []byte) (interface{}, error) {
	err := json.Unmarshal(jsonData, &s.JSON)
	if err != nil {
		log.WithError(err).Errorf("Parse json data error")
		return nil, err
	}

	jsonStruct := passwordResetUnknownJSON{
		FromText:  fmt.Sprintf(s.JSON.FromText, email.operator.operatorName),
		Subject:   fmt.Sprintf(s.JSON.Subject, email.operator.operatorName),
		PlainText: fmt.Sprintf(s.JSON.PlainText, email.operator.operatorName),
		Title:     fmt.Sprintf(s.JSON.Title, email.operator.operatorName),
		Body1:     s.JSON.Body1,
		Body2:     fmt.Sprintf(s.JSON.Body2, email.operator.operatorName),
		Body3:     s.JSON.Body3,
		Body4:     s.JSON.Body4,
	}

	emailData := passwordResetUnknownParam{
		FromText:         jsonStruct.FromText,
		From:             email.from,
		Host:             email.host,
		To:               user,
		Subject:          jsonStruct.Subject,
		MsgID:            param.messageID,
		PlainText:        jsonStruct.PlainText,
		Title:            jsonStruct.Title,
		OperatorLogo:     email.operator.operatorLogo,
		DownloadAppStore: email.operator.downloadAppStore,
		DownloadGoogle:   email.operator.downloadGoogle,
		OperatorLegal:    email.operator.operatorLegal,
		OperatorAddress:  email.operator.operatorAddress,
		OperatorContact:  email.operator.operatorContact,
		B1:               jsonStruct.Body1,
		B2:               jsonStruct.Body2,
		B3:               jsonStruct.Body3,
		B4:               jsonStruct.Body4,
		Str1:             param.commonJSON.Str1,
		Str2:             param.commonJSON.Str2,
		Str3:             param.commonJSON.Str3,
		Str4:             param.commonJSON.Str4,
	}

	return emailData, nil
}
