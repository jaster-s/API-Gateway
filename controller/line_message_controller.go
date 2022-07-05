package controller

import (
	"fli-gateway-api/lib"
	"fli-gateway-api/models"
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type LineMessage struct {
	Destination string `json:"destination"`
	Events      []struct {
		Type      string `json:"type"`
		Timestamp int64  `json:"timestamp"`
		Source    struct {
			Type   string `json:"type"`
			UserID string `json:"userId"`
		} `json:"source"`
		Message struct {
			ID   string `json:"id"`
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"message"`
	} `json:"events"`
}

func (LineMessage) BuildMessage(c *fiber.Ctx) error {
	_logger := new(lib.Logger)

	LineData := new(LineMessage)
	if err := c.BodyParser(LineData); err != nil {
		return err
	}

	if len(LineData.Events) == 0 {
		_logger.InfoLine(fmt.Sprintf("%v", "event LINE is null"))
		return nil
	}

	message := LineData.Events[0].Message.Text
	id := LineData.Events[0].Message.ID
	type_message := LineData.Events[0].Message.Type
	UserID := LineData.Events[0].Source.UserID
	ClientID := c.Params("*")

	SendEmail := false
	Automessage := false
	CreateTicket := false

	activate_key := new(models.ActivateKeyModel)
	message_Line := new(models.LineModels)
	detail_activate_key, errDB_Line, errDB_ActivateKey := activate_key.LoadDataActivateKeyformDB(ClientID)
	if errDB_Line != nil {
		_logger.Error(fmt.Sprintf("%v", errDB_Line), false)
		fmt.Println("LINE office account is not in the database.")
	} else if errDB_ActivateKey != nil {
		_logger.Error(fmt.Sprintf("%v", errDB_ActivateKey), false)
		fmt.Println("Activate Key is not in the database.")
	} else if detail_activate_key.Status == true {
		getContactID, statusCode, respondContact, errDB := message_Line.CheckContact(ClientID, UserID)
		if errDB != nil {
			_logger.Error(fmt.Sprintf("%v", errDB), false)
			fmt.Println("LINE office account is not in the database.")
		} else if statusCode < 400 {
			if getContactID != 0 {
				getTicketID, statusCode, respondTicket, errDB := message_Line.CheckTicket(ClientID, getContactID)
				if errDB != nil {
					_logger.Error(fmt.Sprintf("%v", errDB), false)
					fmt.Println("User is not in the database.")
				} else if statusCode < 400 {
					if getTicketID == 0 {
						CreateTicket = true
					} else if getTicketID != 0 {
						statusCode, respondCreate_Public_Note, fileNameChat, errDB := message_Line.CreatePublicNote(ClientID, message, id, type_message, getTicketID, getContactID)
						if errDB != nil {
							_logger.Error(fmt.Sprintf("%v", errDB), false)
						} else if statusCode < 400 {
							fmt.Println("Create public note to freshdesk")
						} else {
							_logger.InfoFreshdesk("Code: " + strconv.Itoa(statusCode) + " Respond: Create Public Note " + string(respondCreate_Public_Note))
							fmt.Println("The API create public note cannot be contacted at this time.")
							SendEmail = true
						}
						if fileNameChat != "" {
							if errRemoveChat := os.Remove(fileNameChat); errRemoveChat != nil {
								_logger.Error(fmt.Sprintf("%v", errRemoveChat), false)
							}
						}
					}
				} else {
					_logger.InfoFreshdesk("Code: " + strconv.Itoa(statusCode) + " Respond: CheckTicket " + string(respondTicket))
					fmt.Println("The API CheckTicket cannot be contacted at this time.")
					SendEmail = true
				}
			} else {
				CreateTicket = true
			}
		} else {
			_logger.InfoFreshdesk("Code: " + strconv.Itoa(statusCode) + " Respond: CheckContact " + string(respondContact))
			fmt.Println("The API CheckContact cannot be contacted at this time.")
			Automessage = true
			SendEmail = true
		}

		if CreateTicket == true {
			statusCode, respondTicket, fileNameChat, errDB := message_Line.CreateTickets(ClientID, UserID, message, id, type_message)
			if errDB != nil {
				_logger.Error(fmt.Sprintf("%v", errDB), false)
			} else if statusCode < 400 {
				_, statusCode, respondContact, errDB := message_Line.CheckContact(ClientID, UserID)
				statusCode, respondUpdate, fileNameImg, _ := message_Line.UpdateContact(ClientID, UserID, respondContact)
				if errDB != nil {
					_logger.Error(fmt.Sprintf("%v", errDB), false)
					fmt.Println("LINE office account is not in the database.")
				} else if statusCode >= 400 {
					_logger.InfoFreshdesk("Code: " + strconv.Itoa(statusCode) + " Respond: UpdateContact " + string(respondUpdate))
				}
				if errRemoveImg := os.Remove(fileNameImg); errRemoveImg != nil {
					_logger.Error(fmt.Sprintf("%v", errRemoveImg), false)
				}
				if fileNameChat != "" {
					if errRemoveChat := os.Remove(fileNameChat); errRemoveChat != nil {
						_logger.Error(fmt.Sprintf("%v", errRemoveChat), false)
					}
				}
				fmt.Println("create Ticket in freshdesk")
			} else {
				_logger.InfoFreshdesk("Code: " + strconv.Itoa(statusCode) + " Respond: Create " + string(respondTicket))
				fmt.Println("The API CreateTickets cannot be contacted at this time.")
				SendEmail = true
				Automessage = true
			}
		}

	} else {
		Automessage = true
		SendEmail = true
	}

	if SendEmail == true {
		detailMail, errDB := message_Line.DetailtoSendMail(ClientID, UserID)
		if errDB != nil {
			_logger.Error(fmt.Sprintf("%v", errDB), false)
		}
		mail := new(models.SendMail)
		if detail_activate_key.Status == true {
			fmt.Println("LineToFreshdesk Error sending mail.")
			mail.API_true(detailMail.Email, detailMail.NameLINE_OA, detailMail.NameContact, detailMail.Email_contact, message, detailMail.Unique_Extanal_ID)
		} else {
			fmt.Println("Activation key has expired. sending mail.")
			mail.API_false(detailMail.Email, detailMail.NameLINE_OA, detailMail.NameContact, detailMail.Email_contact, message, detailMail.Unique_Extanal_ID)
		}
	}

	if Automessage == true { // comment เป็นการ Auto reply message ในขณะที่ ติดต่อ freshdesk ไม่ได้หรือ API ทาง freshdesk มีปัญหา
		statusCode, respond, errDB := message_Line.AutoMessageLine(ClientID, UserID)
		if errDB != nil {
			_logger.Error(fmt.Sprintf("%v", errDB), false)
		} else if statusCode < 400 {
			fmt.Println("Auto Message")
		} else {
			_logger.InfoLine("Code: " + strconv.Itoa(statusCode) + " Respond: AutoMessageToLine " + string(respond))
		}
	}

	return nil
}
