package models

import (
	"bytes"
	"encoding/json"
	"fli-gateway-api/lib"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"gopkg.in/mgo.v2/bson"
)

var url_UploadFile = os.Getenv("DOMAIN_WEBHOOK")

type FreshdeskModel struct{}
type ContactFreshdaskModel struct {
	Address          string `json:"address"`
	UniqueExternalID string `json:"unique_external_id"`
}
type Text struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
type SendMessageModel struct {
	To       string `json:"to"`
	Messages []Text `json:"messages"`
}
type File struct {
	Type               string `json:"type"`
	OriginalContentUrl string `json:"originalContentUrl"`
	PreviewImageUrl    string `json:"previewImageUrl"`
}
type SendFileModel struct {
	To       string `json:"to"`
	Messages []File `json:"messages"`
}
type Attachments struct {
	AttachmentURL string `json:"attachment_url"`
	Name          string `json:"name"`
	File_size     int    `json:"file_size"`
}

func (FreshdeskModel) getUser(domain string) (*User, error) {
	connection := lib.DBConnection()

	Domain := "https://" + domain + "/"

	detail_user := &User{}
	err := connection.Collection(USER_DB).FindOne(bson.M{"sub_domain": Domain}, detail_user)

	return detail_user, err
}

func (FreshdeskModel) CheckContact(domain string, userID int) (string, int, []byte, error) {
	detail := new(FreshdeskModel)
	user_Db, errDB := detail.getUser(domain)
	if errDB != nil {
		return "", 0, nil, errDB
	}

	Url := user_Db.SubDomain + "api/v2/contacts/" + strconv.Itoa(userID)

	req, _ := http.NewRequest("GET", Url, nil)
	req.SetBasicAuth(user_Db.ApiKey, "X")

	res, _ := http.DefaultClient.Do(req)
	body, _ := io.ReadAll(res.Body)

	if res.StatusCode < 400 {
		if len(body) > 2 {
			contact_detail := ContactFreshdaskModel{}
			err := json.Unmarshal(body, &contact_detail)
			return contact_detail.UniqueExternalID, res.StatusCode, body, err
		} else {
			fmt.Println("You do not have a contact in freshdesk.")
			return "", res.StatusCode, body, nil
		}
	} else {
		return "", res.StatusCode, body, nil
	}
}

func (FreshdeskModel) SendMessage(clientID, uniqueExternalID, message string) (int, []byte, error) {

	text := Text{
		Type: "text",
		Text: message,
	}

	detailEvent := &SendMessageModel{
		To: uniqueExternalID,
		Messages: []Text{
			text,
		},
	}

	data, err := json.Marshal(detailEvent)
	if err != nil {
		fmt.Println("err")
	}

	renderData := bytes.NewReader(data)

	detail := new(LineModels)
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

func (FreshdeskModel) SendFile(clientID, uniqueExternalID string, file Attachments) (int, []byte, string, string, error) {

	detail := new(LineModels)
	fileNameOri := ""
	fileNamePre := ""
	var detailEvent *SendFileModel
	sprint := strings.Split(fmt.Sprintf("%v", strings.ReplaceAll(file.Name, " ", "_")), ".")

	if sprint[1] == "jpg" || sprint[1] == "png" || sprint[1] == "jpeg" || sprint[1] == "gif" || sprint[1] == "bmp" {
		fileNameOri = detail.LoadFileUrl(file.AttachmentURL)
		sprint := strings.Split(fmt.Sprintf("%v", fileNameOri), ".")
		fileNamePre = sprint[0] + ".jpg"

		src, _ := imaging.Open("assets/" + fileNameOri)
		src = imaging.Resize(src, 1600, 1600, imaging.Lanczos)
		imaging.Save(src, "assets/resize/"+fileNamePre)

		file := File{
			Type:               "image",
			OriginalContentUrl: file.AttachmentURL,
			PreviewImageUrl:    url_UploadFile + "/freshdesk_webhook/assets/resize/" + fileNamePre,
		}

		detailEvent = &SendFileModel{
			To: uniqueExternalID,
			Messages: []File{
				file,
			},
		}
	} else if sprint[1] == "mp4" || sprint[1] == "avi" || sprint[1] == "mkv" || sprint[1] == "mpg" || sprint[1] == "mpeg" || sprint[1] == "mov" || sprint[1] == "wmv" || sprint[1] == "m4v" {
		fileNameOri = detail.LoadFileUrl(file.AttachmentURL)
		sprint := strings.Split(fmt.Sprintf("%v", fileNameOri), ".")
		fileNamePre = sprint[0] + ".jpg"

		snapshotVideo := new(FreshdeskModel)
		reader := snapshotVideo.ExtractImage("./assets/"+fileNameOri, 5)
		img, err := imaging.Decode(reader)
		if err != nil {
			fmt.Println(err)
		}
		err = imaging.Save(img, "./assets/resize/"+fileNamePre)
		if err != nil {
			fmt.Println(err)
		}

		file := File{
			Type:               "video",
			OriginalContentUrl: file.AttachmentURL,
			PreviewImageUrl:    url_UploadFile + "/freshdesk_webhook/assets/resize/" + fileNamePre,
		}

		detailEvent = &SendFileModel{
			To: uniqueExternalID,
			Messages: []File{
				file,
			},
		}
	}
	data, err := json.Marshal(detailEvent)
	if err != nil {
		fmt.Println("err")
	}

	renderData := bytes.NewReader(data)
	accessToken, errDB := detail.GenerateAccessToken(clientID)
	if errDB != nil {
		return 0, nil, "", "", errDB
	}

	Url := "https://api.line.me/v2/bot/message/push"

	req, _ := http.NewRequest("POST", Url, renderData)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, _ := http.DefaultClient.Do(req)
	body, _ := io.ReadAll(res.Body)

	detail.RevokeAccessToken(accessToken)

	return res.StatusCode, body, fileNameOri, fileNamePre, nil
}

func (FreshdeskModel) ExtractImage(inFileName string, frameNum int) io.Reader {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		fmt.Println(err)
	}
	return buf
}

func (FreshdeskModel) AutoReplyNewTicket(clientID string, requesterID, ticketID int) (int, []byte, error) {
	detail := new(LineModels)
	Line_Db, errDB := detail.GetLINE(clientID)
	user_Db, errDB := detail.GetUser(clientID)
	if errDB != nil {
		return 0, nil, errDB
	}

	data := &bytes.Buffer{}
	writer := lib.NewWriter(data)
	BodyMessage, _ := writer.CreateFormField("body")
	io.Copy(BodyMessage, strings.NewReader(Line_Db.AutoMessageReply))
	writer.Close()
	encodedData := bytes.NewReader(data.Bytes())

	Url := user_Db.SubDomain + "/api/v2/tickets/" + strconv.Itoa(ticketID) + "/reply"

	req, _ := http.NewRequest("POST", Url, encodedData)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.SetBasicAuth(user_Db.ApiKey, "X")

	res, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)

	return res.StatusCode, body, nil
}

func (FreshdeskModel) ReplyTickettoMail(email string) bool {
	sprint := strings.Split(fmt.Sprintf("%v", email), "@")

	if sprint[1] == "freshintegratio.co" || sprint[1] == "freshintegratio.ne" {
		return false
	} else {
		return true
	}
}

func (FreshdeskModel) AppInstall(userID string, set bool) error {
	connection := lib.DBConnection()

	detail_activate_key := &ActivateKey{}
	errfindDB := connection.Collection(ActivateKEY_DB).FindOne(bson.M{"user_id": userID}, detail_activate_key)

	detail_activate_key.Install = set

	connection.Collection(ActivateKEY_DB).Save(detail_activate_key)

	return errfindDB
}
