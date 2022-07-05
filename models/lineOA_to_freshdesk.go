package models

import (
	"bytes"
	"encoding/json"
	"fli-gateway-api/lib"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"gopkg.in/mgo.v2/bson"
)

type LineModels struct{}

type ProFile struct {
	UserID        string `json:"userID"`
	DisplayName   string `json:"displayName"`
	PictireURL    string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
}

type ReplyTicket struct {
	Body      string `json:"body"`
	ContactID int    `json:"user_id"`
	Private   bool   `json:"private"`
}

type ContactFreshdask []struct {
	Email     interface{} `json:"email"`
	ID        int         `json:"id"`
	Time_Zone string      `json:"time_zone"`
	Avatar    interface{} `json:"avatar"`
}

type MailDetail struct {
	Email             string `json:"email"`
	NameLINE_OA       string `json:"nameline_oa"`
	NameContact       string `json:"namecontact"`
	Unique_Extanal_ID string `json:"unique_extanal_id"`
	Email_contact     string `json:"email_contact"`
}

type TricketFreshdask []struct {
	Ticket_id      int   `json:"ticket_id"`
	ID             int   `json:"id"`
	Status         int   `json:"status"`
	ToEmailUserIds []int `json:"to_email_user_ids"`
}

type AccessTokenLine struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type CreateTicket struct {
	Subject            string      `json:"subject"`
	Name               string      `json:"name"`
	Unique_external_id string      `json:"unique_external_id"`
	Priority           int         `json:"priority"`
	Status             int         `json:"status"`
	Source             int         `json:"source"`
	Type               string      `json:"type"`
	Group_ID           interface{} `json:"group_id"`
	Responder_ID       interface{} `json:"responder_id"`
	Description        string      `json:"description"`
}

func (LineModels) GetLINE(clientID string) (*LineOA, error) {
	connection := lib.DBConnection()

	lineOA := &LineOA{}
	err := connection.Collection(LINE_OA_DB).FindOne(bson.M{"client_id": clientID}, lineOA)
	return lineOA, err
}

func (LineModels) GetUser(clientID string) (*User, error) {
	connection := lib.DBConnection()

	detail := new(LineModels)
	lineOADetail, errDB := detail.GetLINE(clientID)
	if errDB != nil {
		return nil, errDB
	}
	user := &User{}
	err := connection.Collection(USER_DB).FindById(bson.ObjectIdHex(lineOADetail.UserID), user)

	return user, err
}

func (LineModels) GetProfileLINE(userID, clientID string) *ProFile {
	Url := "https://api.line.me/v2/bot/profile/" + userID

	detail := new(LineModels)
	accessToken, _ := detail.GenerateAccessToken(clientID)

	req, _ := http.NewRequest("GET", Url, nil)
	// req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, _ := http.DefaultClient.Do(req)
	body, _ := io.ReadAll(res.Body)

	profile := ProFile{}
	if err := json.Unmarshal(body, &profile); err != nil {
		fmt.Println("err")
	}
	detail.RevokeAccessToken(accessToken)

	return &profile
}

func (LineModels) GetProfileLINEOA(clientID string) *ProFile {
	Url := "https://api.line.me/v2/bot/info"

	detail := new(LineModels)
	accessToken, _ := detail.GenerateAccessToken(clientID)

	req, _ := http.NewRequest("GET", Url, nil)
	// req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, _ := http.DefaultClient.Do(req)
	body, _ := io.ReadAll(res.Body)

	profile := ProFile{}
	if err := json.Unmarshal(body, &profile); err != nil {
		fmt.Println("err")
	}
	detail.RevokeAccessToken(accessToken)

	return &profile
}

func (LineModels) GenerateAccessToken(clientID string) (string, error) {
	detail := new(LineModels)
	lineOA_Db, errDB := detail.GetLINE(clientID)
	if errDB != nil {
		return "err", errDB
	}

	Url := "https://api.line.me/v2/oauth/accessToken"

	data := url.Values{}
	data.Add("grant_type", "client_credentials")
	data.Add("client_id", lineOA_Db.ClientID)
	data.Add("client_secret", lineOA_Db.ClientSecret)
	encodedData := strings.NewReader(data.Encode())

	req, _ := http.NewRequest("POST", Url, encodedData)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	res, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode < 400 {
		token_detail := AccessTokenLine{}
		err := json.Unmarshal(body, &token_detail)
		return token_detail.AccessToken, err
	} else {
		fmt.Println("Your clientID is invalid.")
		return string(body), nil
	}
}

func (LineModels) RevokeAccessToken(accessToken string) []byte {

	Url := "https://api.line.me/v2/oauth/revoke"

	data := url.Values{}
	data.Add("access_token", accessToken)
	encodedData := strings.NewReader(data.Encode())

	req, _ := http.NewRequest("POST", Url, encodedData)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	return body
}

func (LineModels) LoadFileUrl(fullURL string) string {
	fileUrl, _ := url.Parse(fullURL)

	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	path := fileUrl.Path
	segments := strings.Split(path, "/")
	fileName := strings.ReplaceAll(segments[len(segments)-1], " ", "_")

	file, _ := os.Create("assets/" + fileName)

	res, err := client.Get(fullURL)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	io.Copy(file, res.Body)
	defer file.Close()

	return fileName
}

func (LineModels) LoadFileInChat(dataID, clientID, type_message string) string {
	detail := new(LineModels)
	accessToken, _ := detail.GenerateAccessToken(clientID)
	type_file := ""

	if type_message == "image" {
		type_file = ".jpg"
	} else if type_message == "video" {
		type_file = ".mp4"
	} else if type_message == "audio" {
		type_file = ".mp3"
	}

	file, _ := os.Create("assets/" + dataID + type_file)

	Url := "https://api-data.line.me/v2/bot/message/" + dataID + "/content"

	req, _ := http.NewRequest("GET", Url, nil)
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	io.Copy(file, res.Body)
	defer file.Close()

	detail.RevokeAccessToken(accessToken)
	return dataID + type_file
}

func (LineModels) CheckContact(clientID, userID string) (int, int, []byte, error) {
	detail := new(LineModels)
	user_Db, errDB := detail.GetUser(clientID)
	if errDB != nil {
		return 0, 0, nil, errDB
	}

	Url := user_Db.SubDomain + "api/v2/contacts?unique_external_id=" + clientID + "_" + userID

	req, _ := http.NewRequest("GET", Url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(user_Db.ApiKey, "X")

	res, _ := http.DefaultClient.Do(req)
	body, _ := io.ReadAll(res.Body)

	if res.StatusCode < 400 {
		if len(body) > 2 {
			contact_detail := ContactFreshdask{}
			err := json.Unmarshal(body, &contact_detail)
			return contact_detail[0].ID, res.StatusCode, body, err
		} else {
			fmt.Println("You do not have a contact in freshdesk.")
			return 0, res.StatusCode, body, nil
		}
	} else {
		return 0, res.StatusCode, body, nil
	}
}

func (LineModels) UpdateContact(clientID, userID string, response []byte) (int, []byte, string, error) {
	detail := new(LineModels)
	user_Db, errDB := detail.GetUser(clientID)
	LineDetail := detail.GetProfileLINE(userID, clientID)
	if errDB != nil {
		return 400, nil, "", errDB
	}

	contact_detail := ContactFreshdask{}
	err := json.Unmarshal(response, &contact_detail)

	data := &bytes.Buffer{}
	writer := multipart.NewWriter(data)
	addEmail := false
	fileName := ""
	addData := false

	if contact_detail[0].Email == nil {
		addEmail = true
	} else if sprint := strings.Split(fmt.Sprintf("%v", contact_detail[0].Email), "@"); sprint[1] == "freshintegratio.ne" {
		addEmail = true
	}

	if addEmail == true {
		email, _ := writer.CreateFormField("email")
		io.Copy(email, strings.NewReader(strconv.Itoa(contact_detail[0].ID)+"@freshintegratio.co"))
		addData = true
		fmt.Println("add Email")
	}
	if contact_detail[0].Time_Zone != "Bangkok" {
		time_zone, _ := writer.CreateFormField("time_zone")
		io.Copy(time_zone, strings.NewReader("Bangkok"))
		addData = true
		fmt.Println("add time Zone")
	}
	if contact_detail[0].Avatar == nil {
		fileName = "assets/" + detail.LoadFileUrl(LineDetail.PictireURL+".jpg")
		avatar, _ := writer.CreateFormFile("avatar", fileName)
		file, _ := os.Open(fileName)
		io.Copy(avatar, file)
		defer file.Close()
		addData = true
		fmt.Println("add Avatar")
	}
	writer.Close()

	if addData == true {
		encodedData := bytes.NewReader(data.Bytes())

		Url := user_Db.SubDomain + "/api/v2/contacts/" + strconv.Itoa(contact_detail[0].ID)

		req, _ := http.NewRequest("PUT", Url, encodedData)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.SetBasicAuth(user_Db.ApiKey, "X")

		res, _ := http.DefaultClient.Do(req)
		body, _ := io.ReadAll(res.Body)

		return res.StatusCode, body, fileName, err
	} else {
		return 400, nil, fileName, err
	}
}

func (LineModels) CheckTicket(clientID string, contact_ID int) (int, int, []byte, error) {
	detail := new(LineModels)
	user_Db, errDB := detail.GetUser(clientID)
	if errDB != nil {
		return 0, 0, nil, errDB
	}

	Url := user_Db.SubDomain + "/api/v2/tickets?requester_id=" + strconv.Itoa(contact_ID)

	req, _ := http.NewRequest("GET", Url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(user_Db.ApiKey, "X")

	res, _ := http.DefaultClient.Do(req)
	body, _ := io.ReadAll(res.Body)

	if res.StatusCode < 400 {
		if len(body) > 2 {
			ticket_detail := TricketFreshdask{}
			err := json.Unmarshal(body, &ticket_detail)
			if ticket_detail[0].Status == 5 {
				fmt.Println("You don't have an open ticket in Freshdesk..")
				return 0, res.StatusCode, body, err
			} else {
				return ticket_detail[0].ID, res.StatusCode, body, err
			}
		} else {
			fmt.Println("You do not have a ticket in freshdesk.")
			return 0, res.StatusCode, body, nil
		}
	} else {
		return 0, res.StatusCode, body, nil
	}
}

func (LineModels) CreateTickets(clientID, userID, message, id, type_message string) (int, []byte, string, error) {
	detail := new(LineModels)
	lineOA_Db, errLineDB := detail.GetLINE(clientID)
	user_Db, errUserDB := detail.GetUser(clientID)
	LineDetail := detail.GetProfileLINE(userID, clientID)
	if errLineDB != nil || errUserDB != nil {
		if errLineDB != nil {
			return 0, nil, "", errLineDB
		} else if errUserDB != nil {
			return 0, nil, "", errUserDB
		}
	}

	createTicket := &CreateTicket{
		Subject:            lineOA_Db.LINEname + " - " + LineDetail.DisplayName,
		Name:               lineOA_Db.LINEname + " - " + LineDetail.DisplayName,
		Unique_external_id: clientID + "_" + userID,
		Priority:           1,
		Status:             2,
		Source:             100,
		Type:               "Question",
		Description:        message,
	}

	if lineOA_Db.GroupID != 0 {
		createTicket.Group_ID = lineOA_Db.GroupID
	} else {
		createTicket.Group_ID = nil
	}
	if lineOA_Db.ResponderID != 0 {
		createTicket.Responder_ID = lineOA_Db.ResponderID
	} else {
		createTicket.Responder_ID = nil
	}

	data, err := json.Marshal(createTicket)
	if err != nil {
		fmt.Println("err")
	}

	renderData := bytes.NewReader(data)

	Url := user_Db.SubDomain + "api/v2/tickets"

	req, _ := http.NewRequest("POST", Url, renderData)
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(user_Db.ApiKey, "X")

	res, _ := http.DefaultClient.Do(req)
	body, _ := io.ReadAll(res.Body)

	if type_message != "text" {
		fileName := ""
		ticketID := 0
		contactID := 0
		if res.StatusCode < 400 {
			ticket_detail := TricketFreshdask{}
			json.Unmarshal(body, &ticket_detail)
			ticketID = ticket_detail[0].Ticket_id
			contactID = ticket_detail[0].ToEmailUserIds[0]
		} else {
			return res.StatusCode, body, "", nil
		}

		data := &bytes.Buffer{}
		writer := lib.NewWriter(data)

		User_id, _ := writer.CreateFormField("user_id")
		io.Copy(User_id, strings.NewReader(strconv.Itoa(contactID)))

		Private, _ := writer.CreateFormField("private")
		io.Copy(Private, strings.NewReader("false"))

		Body, _ := writer.CreateFormField("body")
		io.Copy(Body, strings.NewReader("attachments file"))
		fileName = "assets/" + detail.LoadFileInChat(id, clientID, type_message)
		if type_message == "image" {
			attachments, _ := writer.CreateFormImg("attachments[]", fileName)
			file, _ := os.Open(fileName)
			io.Copy(attachments, file)
			defer file.Close()
		} else if type_message == "video" {
			attachments, _ := writer.CreateFormVideo("attachments[]", fileName)
			file, _ := os.Open(fileName)
			io.Copy(attachments, file)
			defer file.Close()
		} else if type_message == "audio" {
			attachments, _ := writer.CreateFormAudio("attachments[]", fileName)
			file, _ := os.Open(fileName)
			io.Copy(attachments, file)
			defer file.Close()
		}

		writer.Close()
		encodedData := bytes.NewReader(data.Bytes())

		Url := user_Db.SubDomain + "/api/v2/tickets/" + strconv.Itoa(ticketID) + "/notes"

		req, _ := http.NewRequest("POST", Url, encodedData)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.SetBasicAuth(user_Db.ApiKey, "X")

		res, _ := http.DefaultClient.Do(req)
		body, _ := ioutil.ReadAll(res.Body)

		return res.StatusCode, body, fileName, nil
	}

	return res.StatusCode, body, "", nil
}

func (LineModels) CreatePublicNote(clientID, message, id, type_message string, ticketID, contactID int) (int, []byte, string, error) {
	detail := new(LineModels)
	user_Db, errDB := detail.GetUser(clientID)
	if errDB != nil {
		return 0, nil, "", errDB
	}

	data := &bytes.Buffer{}
	writer := lib.NewWriter(data)

	User_id, _ := writer.CreateFormField("user_id")
	io.Copy(User_id, strings.NewReader(strconv.Itoa(contactID)))

	Private, _ := writer.CreateFormField("private")
	io.Copy(Private, strings.NewReader("false"))

	Body, _ := writer.CreateFormField("body")
	fileName := ""
	if type_message == "text" {
		io.Copy(Body, strings.NewReader(message))
	} else {
		io.Copy(Body, strings.NewReader("attachments file"))
		fileName = "assets/" + detail.LoadFileInChat(id, clientID, type_message)
		if type_message == "image" {
			attachments, _ := writer.CreateFormImg("attachments[]", fileName)
			file, _ := os.Open(fileName)
			io.Copy(attachments, file)
			defer file.Close()
		} else if type_message == "video" {
			attachments, _ := writer.CreateFormVideo("attachments[]", fileName)
			file, _ := os.Open(fileName)
			io.Copy(attachments, file)
			defer file.Close()
		} else if type_message == "audio" {
			attachments, _ := writer.CreateFormAudio("attachments[]", fileName)
			file, _ := os.Open(fileName)
			io.Copy(attachments, file)
			defer file.Close()
		}
	}
	writer.Close()
	encodedData := bytes.NewReader(data.Bytes())

	Url := user_Db.SubDomain + "/api/v2/tickets/" + strconv.Itoa(ticketID) + "/notes"

	req, _ := http.NewRequest("POST", Url, encodedData)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.SetBasicAuth(user_Db.ApiKey, "X")

	res, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)

	return res.StatusCode, body, fileName, nil
}

func (LineModels) DetailtoSendMail(clientID, userID string) (*MailDetail, error) {
	detail := new(LineModels)
	lineOA_Db, errLineDB := detail.GetLINE(clientID)
	LineDetail := detail.GetProfileLINE(userID, clientID)

	min := 1000000
	max := 9999999
	email := strconv.Itoa(rand.Intn(max-min)+min) + "@freshintegratio.ne"

	if errLineDB != nil {
		return nil, errLineDB
	}

	detailMail := &MailDetail{
		Email:             lineOA_Db.RecoveryMail,
		NameLINE_OA:       lineOA_Db.LINEname,
		NameContact:       LineDetail.DisplayName,
		Unique_Extanal_ID: clientID + "_" + userID,
		Email_contact:     email,
	}

	return detailMail, nil
}

func (LineModels) AutoMessageLine(clientID, userID string) (int, []byte, error) {

	detail := new(LineModels) // ยกเลิกการเรียก ข้อความจาก DB เปลียนมาเป็นการพิมพ์ข้อความใส่ชั่วคราว
	// Line_Db, errDb := detail.LoadDataLineOAformDB(clientID)
	// if errDb != nil {
	// 	return 0, nil, errDb
	// }

	text := Text{
		Type: "text",
		// Text: Line_Db.AutoMessageReply,
		Text: "ขณะนี้ระบบมีปัญหาอยู่ระหว่างทำการแก้ไข เมื่อแก้ไขเสร็จแล้ว ทางเราจะรีบตอบกลับอย่างเร็วที่สุด",
	}

	detailEvent := &SendMessageModel{
		To: userID,
		Messages: []Text{
			text,
		},
	}

	data, err := json.Marshal(detailEvent)
	if err != nil {
		fmt.Println("err")
	}

	renderData := bytes.NewReader(data)
	accessToken, errDB := detail.GenerateAccessToken(clientID)
	if errDB != nil {
		return 0, nil, errDB
	}

	Url := "https://api.line.me/v2/bot/message/push"

	req, _ := http.NewRequest("POST", Url, renderData)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, _ := http.DefaultClient.Do(req)
	body, _ := io.ReadAll(res.Body)

	detail.RevokeAccessToken(accessToken)

	return res.StatusCode, body, nil
}
