package models

import (
	"fli-gateway-api/lib"
	"fmt"
	"log"
	"os"
	"strconv"

	mail "github.com/xhit/go-simple-mail/v2"
)

type SendMail struct{}

func (SendMail) API_true(email_LINE_OA, nameLINE_OA, CustomerName, email_contact, messageReq, unique_external_id string) error {
	_logger := new(lib.Logger)

	html_template := new(lib.HTML_Template)
	htmlBody := html_template.APIisTrue(nameLINE_OA, CustomerName, email_contact, messageReq, unique_external_id)

	server := mail.NewSMTPClient()

	server.Host = os.Getenv("MAIL_HOST")
	server.Port, _ = strconv.Atoi(os.Getenv("MAIL_PORT"))
	server.Username = os.Getenv("MAIL_USERNAME")
	server.Password = os.Getenv("MAIL_PASSWORD")
	server.Encryption = mail.EncryptionTLS

	smtpClient, err := server.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// Create email
	email := mail.NewMSG()
	email.SetFrom(os.Getenv("MAIL_USERNAME")).
		AddTo(email_LINE_OA).
		SetSubject("Message from " + CustomerName + " (" + unique_external_id + ")")

	email.SetBody(mail.TextHTML, htmlBody)

	// Send email
	err = email.Send(smtpClient)
	if err != nil {
		_logger.Error(fmt.Sprintf("Cannot Send Mail: %v", err), false)
	}
	return err
}

func (SendMail) API_false(email_LINE_OA, nameLINE_OA, CustomerName, email_contact, messageReq, unique_external_id string) error {
	_logger := new(lib.Logger)

	html_template := new(lib.HTML_Template)
	htmlBody := html_template.KeyExpiration(nameLINE_OA, CustomerName, messageReq)

	server := mail.NewSMTPClient()

	server.Host = os.Getenv("MAIL_HOST")
	server.Port, _ = strconv.Atoi(os.Getenv("MAIL_PORT"))
	server.Username = os.Getenv("MAIL_USERNAME")
	server.Password = os.Getenv("MAIL_PASSWORD")
	server.Encryption = mail.EncryptionTLS

	smtpClient, err := server.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// Create email
	email := mail.NewMSG()
	email.SetFrom(os.Getenv("MAIL_USERNAME")).
		AddTo(email_LINE_OA).
		SetSubject("Message from " + CustomerName + " (" + unique_external_id + ")")

	email.SetBody(mail.TextHTML, htmlBody)

	// Send email
	err = email.Send(smtpClient)
	if err != nil {
		_logger.Error(fmt.Sprintf("Cannot Send Mail: %v", err), false)
	}
	return err
}

func (SendMail) Add_recoveryEmail(email_LINE_OA, nameLINE_OA string) error {
	_logger := new(lib.Logger)

	html_template := new(lib.HTML_Template)
	htmlBody := html_template.RecoveryEmail(nameLINE_OA, email_LINE_OA)

	server := mail.NewSMTPClient()

	server.Host = os.Getenv("MAIL_HOST")
	server.Port, _ = strconv.Atoi(os.Getenv("MAIL_PORT"))
	server.Username = os.Getenv("MAIL_USERNAME")
	server.Password = os.Getenv("MAIL_PASSWORD")
	server.Encryption = mail.EncryptionTLS

	smtpClient, err := server.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// Create email
	email := mail.NewMSG()
	email.SetFrom(os.Getenv("MAIL_USERNAME")).
		AddTo(email_LINE_OA).
		SetSubject("Add recovery Email ")

	email.SetBody(mail.TextHTML, htmlBody)

	// Send email
	err = email.Send(smtpClient)
	if err != nil {
		_logger.Error(fmt.Sprintf("Cannot Send Mail: %v", err), false)
	}
	return err
}

func (SendMail) Reply_MessageInTicketToMail(message, CustomerMail string, attach []string, ticket_Id int) error {
	_logger := new(lib.Logger)

	html_template := new(lib.HTML_Template)
	htmlBody := html_template.ReMessagetoMail(message, ticket_Id)

	server := mail.NewSMTPClient()

	server.Host = os.Getenv("MAIL_HOST")
	server.Port, _ = strconv.Atoi(os.Getenv("MAIL_PORT"))
	server.Username = os.Getenv("MAIL_USERNAME")
	server.Password = os.Getenv("MAIL_PASSWORD")
	server.Encryption = mail.EncryptionTLS

	smtpClient, err := server.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// Create email
	email := mail.NewMSG()
	email.SetFrom(os.Getenv("MAIL_USERNAME")).
		AddTo(CustomerMail).
		SetSubject("Add recovery Email ")

	email.SetBody(mail.TextHTML, htmlBody)

	for i := range attach {
		email.Attach(&mail.File{FilePath: "assets/" + attach[i]})
	}

	// Send email
	err = email.Send(smtpClient)
	if err != nil {
		_logger.Error(fmt.Sprintf("Cannot Send Mail: %v", err), false)
	}
	return err
}

func (SendMail) SendActivateKeyTrialVersion(name, activateKey, CustomerMail string) error {
	_logger := new(lib.Logger)

	html_template := new(lib.HTML_Template)
	htmlBody := html_template.SendActivateKeyTrialVersion(name, activateKey)

	server := mail.NewSMTPClient()

	server.Host = os.Getenv("MAIL_HOST")
	server.Port, _ = strconv.Atoi(os.Getenv("MAIL_PORT"))
	server.Username = os.Getenv("MAIL_USERNAME")
	server.Password = os.Getenv("MAIL_PASSWORD")
	server.Encryption = mail.EncryptionTLS

	smtpClient, err := server.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// Create email
	email := mail.NewMSG()
	email.SetFrom(os.Getenv("MAIL_USERNAME")).
		AddTo(CustomerMail).
		SetSubject("Add recovery Email ")

	email.SetBody(mail.TextHTML, htmlBody)

	// Send email
	err = email.Send(smtpClient)
	if err != nil {
		_logger.Error(fmt.Sprintf("Cannot Send Mail: %v", err), false)
	}
	return err
}
