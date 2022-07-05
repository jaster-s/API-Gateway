package controller

import (
	"encoding/json"
	"fli-gateway-api/lib"
	"fli-gateway-api/models"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type FreshdeskMessage struct {
	Data struct {
		Conversation struct {
			BodyText       string   `json:"body_text"`
			TicketID       int      `json:"ticket_id"`
			ToEmails       []string `json:"to_emails"`
			ToEmailUserIds []int    `json:"to_email_user_ids"`
			Attachments    []struct {
				AttachmentURL string `json:"attachment_url"`
				Name          string `json:"name"`
				File_size     int    `json:"file_size"`
			} `json:"attachments"`
		} `json:"conversation"`
		Ticket struct {
			Id          int `json:"id"`
			RequesterID int `json:"requester_id"`
		} `json:"ticket"`
	} `json:"data"`
	Domain string `json:"domain"`
}

type AppData struct {
	AccountID string `json:"account_id"`
	Domain    string `json:"domain"`
	Event     string `json:"event"`
	Region    string `json:"region"`
	Timestamp string `json:"timestamp"`
	Iparams   struct {
		Param1 string `json:"Param1"`
		Param2 string `json:"Param2"`
	} `json:"iparams"`
}

func (FreshdeskMessage) NewTickets(c *fiber.Ctx) error {
	_logger := new(lib.Logger)

	BodyFreshdesk := FreshdeskMessage{}
	json.Unmarshal(c.Body(), &BodyFreshdesk)

	if BodyFreshdesk.Data.Ticket.RequesterID == 0 {
		return nil
	}
	ticketID := BodyFreshdesk.Data.Ticket.Id
	userID := BodyFreshdesk.Data.Ticket.RequesterID
	domain := BodyFreshdesk.Domain

	message_freshdesk := new(models.FreshdeskModel)
	getClintID, statusCode, respond, err := message_freshdesk.CheckContact(domain, userID)
	if err != nil {
		_logger.Error(fmt.Sprintf("%v", err), false)
	} else if statusCode < 400 && getClintID != "" {
		sprint := strings.Split(getClintID, "_")
		statusCode, respond, errDB := message_freshdesk.AutoReplyNewTicket(sprint[0], userID, ticketID)
		if err != nil {
			_logger.Error(fmt.Sprintf("%v", errDB), false)
		} else if statusCode < 400 {
			fmt.Println("Tickets are created and auto-reply ticket are sent.")
		} else {
			_logger.InfoFreshdesk("Code: " + strconv.Itoa(statusCode) + " Respond: AutoReplyNewTicket " + string(respond))
			fmt.Println("The API AutoReplyNewTicket in freshdesk cannot be contacted at this time.")
		}
	} else {
		_logger.InfoFreshdesk("Code: " + strconv.Itoa(statusCode) + " Respond: CheckContact " + string(respond))
		fmt.Println("The API CheckContact in freshdesk cannot be contacted at this time.")
	}

	return nil
}

func (FreshdeskMessage) ReplytoLine(c *fiber.Ctx) error {
	_logger := new(lib.Logger)

	BodyFreshdesk := FreshdeskMessage{}
	json.Unmarshal(c.Body(), &BodyFreshdesk)

	if len(BodyFreshdesk.Data.Conversation.ToEmailUserIds) == 0 {
		return nil
	}

	SendEmail := false

	message := BodyFreshdesk.Data.Conversation.BodyText
	userID := BodyFreshdesk.Data.Conversation.ToEmailUserIds[0]
	domain := BodyFreshdesk.Domain
	customerMail := BodyFreshdesk.Data.Conversation.ToEmails[0]
	attachments := BodyFreshdesk.Data.Conversation.Attachments

	message_freshdesk := new(models.FreshdeskModel)
	getClintIDandContactExternalID, statusCode, respond, err := message_freshdesk.CheckContact(domain, userID)
	if err != nil {
		_logger.Error(fmt.Sprintf("%v", err), false)
	} else if statusCode < 400 && getClintIDandContactExternalID != "" {
		sprint := strings.Split(getClintIDandContactExternalID, "_")

		activate_key := new(models.ActivateKeyModel)
		detail_activate_key, errDB_Line, errDB_ActivateKey := activate_key.LoadDataActivateKeyformDB(sprint[0])
		if errDB_Line != nil {
			_logger.Error(fmt.Sprintf("%v", errDB_Line), false)
			fmt.Println("LINE office account is not in the database.")
		} else if errDB_ActivateKey != nil {
			_logger.Error(fmt.Sprintf("%v", errDB_ActivateKey), false)
			fmt.Println("Activate Key is not in the database.")
		} else if detail_activate_key.Status == true {
			statusCode, respond, errDB := message_freshdesk.SendMessage(sprint[0], sprint[1], message)
			if errDB != nil {
				_logger.Error(fmt.Sprintf("%v", errDB), false)
			} else if statusCode < 400 {
				fmt.Println("send message to line")
				checkMail := message_freshdesk.ReplyTickettoMail(customerMail)
				if checkMail == true {
					SendEmail = true
				}
			} else {
				_logger.InfoFreshdesk("Code: " + strconv.Itoa(statusCode) + " Respond: SendMessage to Line " + string(respond))
				fmt.Println("The API SendMessage to line cannot be contacted at this time.")
			}
			if len(BodyFreshdesk.Data.Conversation.Attachments) != 0 {
				for i := range attachments {
					attachments := BodyFreshdesk.Data.Conversation.Attachments[i]
					statusCode, respond, _, _, errDB := message_freshdesk.SendFile(sprint[0], sprint[1], attachments)
					if errDB != nil {
						_logger.Error(fmt.Sprintf("%v", errDB), false)
					} else if statusCode < 400 {
						fmt.Println("send file to line") //ลบไฟล์เวลาส่งไฟล์กลับไปที่ Line
						// time.Sleep(20 * time.Second)
						// if errRemovefile := os.Remove("./assets/" + fileNameOri); errRemovefile != nil {
						// 	_logger.Error(fmt.Sprintf("%v", fileNameOri), false)
						// }
						// if errRemovefile := os.Remove("./assets/resize/" + fileNamePre); errRemovefile != nil {
						// 	_logger.Error(fmt.Sprintf("%v", fileNamePre), false)
						// }
					} else {
						_logger.InfoFreshdesk("Code: " + strconv.Itoa(statusCode) + " Respond: SendMessage to Line " + string(respond))
						fmt.Println("The API SendFile to line cannot be contacted at this time.")
					}
				}

			}
		} else {
			fmt.Println("Activation key has expired. sending mail.")
		}
	} else {
		_logger.InfoFreshdesk("Code: " + strconv.Itoa(statusCode) + " Respond: CheckContact " + string(respond))
		fmt.Println("The API CheckContact in freshdesk cannot be contacted at this time.")
	}

	ticket_id := BodyFreshdesk.Data.Conversation.TicketID

	if SendEmail == true {
		var fileName []string
		attact_freshdesk := new(models.LineModels)
		for i := range BodyFreshdesk.Data.Conversation.Attachments {
			fileName = append(fileName, attact_freshdesk.LoadFileUrl(attachments[i].AttachmentURL))
		}
		mail := new(models.SendMail)
		fmt.Println("LineToFreshdesk auto reply sending mail.")
		mail.Reply_MessageInTicketToMail(message, customerMail, fileName, ticket_id)

		gmailA_API := new(lib.GmailAPI)
		subject := "Your ticket number is " + strconv.Itoa(ticket_id)
		gmailA_API.ConversationInTicket(customerMail, subject, message, fileName, ticket_id)
	}

	return nil
}

func (FreshdeskMessage) AppInstall(c *fiber.Ctx) error {
	_logger := new(lib.Logger)
	userID := fmt.Sprintf("%v", c.Locals("userID"))
	verify_Install := new(models.FreshdeskModel)

	errfindDB := verify_Install.AppInstall(userID, true)
	if errfindDB != nil {
		_logger.Error(fmt.Sprintf("%v", errfindDB), false)
	}
	return nil
}

func (FreshdeskMessage) AppUninstall(c *fiber.Ctx) error {
	_logger := new(lib.Logger)
	userID := fmt.Sprintf("%v", c.Locals("userID"))
	verify_Install := new(models.FreshdeskModel)

	errfindDB := verify_Install.AppInstall(userID, false)
	if errfindDB != nil {
		_logger.Error(fmt.Sprintf("%v", errfindDB), false)
	}
	return nil
}
