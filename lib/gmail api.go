package lib

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type GmailAPI struct{}

// GmailService : Gmail client for sending email
var GmailService *gmail.Service

func OAuthGmailService() {
	config := oauth2.Config{
		ClientID:     os.Getenv("Gmail_API_ClientID"),
		ClientSecret: os.Getenv("Gmail_API_ClientSecret"),
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost",
	}

	token := oauth2.Token{
		AccessToken:  os.Getenv("Gmail_API_AccessToken"),
		RefreshToken: os.Getenv("Gmail_API_RefreshToken"),
		TokenType:    "Bearer",
		Expiry:       time.Now(),
	}

	var tokenSource = config.TokenSource(context.Background(), &token)

	srv, err := gmail.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		log.Printf("Unable to retrieve Gmail client: %v", err)
	}

	GmailService = srv
	if GmailService != nil {
		fmt.Println("Email service is initialized")
	}
}

func chunkSplit(body string, limit int, end string) string {
	var charSlice []rune

	// push characters to slice
	for _, char := range body {
		charSlice = append(charSlice, char)
	}

	var result = ""

	for len(charSlice) >= 1 {
		// convert slice/array back to string
		// but insert end at specified limit
		result = result + string(charSlice[:limit]) + end

		// discard the elements that were copied over to result
		charSlice = charSlice[limit:]

		// change the limit
		// to cater for the last few words in
		if len(charSlice) < limit {
			limit = len(charSlice)
		}
	}
	return result
}

func randStr(strSize int, randType string) string {

	var dictionary string

	if randType == "alphanum" {
		dictionary = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	}

	var strBytes = make([]byte, strSize)
	_, _ = rand.Read(strBytes)
	for k, v := range strBytes {
		strBytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(strBytes)
}

func (GmailAPI) ConversationInTicket(to, subject, content string, fileName []string, ticket_Id int) (bool, error) {
	var htmlBody = `
    <html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
</head>
<body>
    <center>
        <table width="100%" style="width:100%;max-width:600px" align="center">
            <tr>
                <td style="padding:0px 15px 0px 15px;color:#000000;text-align:left" bgcolor="#ffffff" width="100%">
                    <table border="0" width="100%" style="table-layout:fixed">
                        <tr>
                            <td style="padding:18px 0px 6px 0px;line-height:22px;text-align:inherit" height="100%" valign="top">
                            <div>Dear Customer,</div>
                                <br>
                            </td>
                        </tr>
                    </table>
                    <table width="100%" style="table-layout:fixed">
                        <tr>
                            <td style="padding:6px 0px 6px px;line-height:30px;text-align:inherit;background-color:#ccf1ff" height="100%">
                                <div style="text-align:left">
                                    <div style="padding:18px 50px 18px 50px">
                                        Your ticket number is ` + strconv.Itoa(ticket_Id) + ` <br>
                                        ` + content + `
                                    </div>
                                </div>
                            </td>
                        </tr>
                    </table>
                    <table width="100%" style="table-layout:fixed">
                        <tr>
                            <td style="padding:18px 0px 6px 0px;line-height:22px;text-align:inherit" height="100%" valign="top">
                                <div>&nbsp;</div>
                                <div>Best regards,<br>
                                    Facgure tester</div>
                            </td>
                        </tr>
                    </table>
                    <table width="100%" style="table-layout:fixed">
                        <tr>
                            <td style="padding:12px 20px 14px 20px; font-size:12px; line-height:16px; font-weight:normal; color:#666666; background:#efefef;">
                                หมายเหตุ: อย่าตอบกลับอีเมลฉบับนี้ หากมีข้อสงสัยใดๆ โปรดติดต่อเราที่เว็บไซต์:<br>
                                <a href="https://www.facgure.com/contact" style="color:#1428a0; ">
                                    Facgure Automatic Digital Solution Companion
                                </a>
                            </td>
                        </tr>
                    </table>
                </td>
            </tr>
        </table>
    </center>
</body>
</html>
	`
	var message gmail.Message

	boundary := randStr(32, "alphanum")

	messageBody := "Content-Type: multipart/mixed; boundary=" + boundary + " \n" +
		"MIME-Version: 1.0\n" +
		"to: " + to + "\n" +
		"subject: " + subject + "\n\n" +

		"--" + boundary + "\n" +
		"Content-Type: text/plain; charset=" + string('"') + "UTF-8" + string('"') + "\n" +
		"MIME-Version: 1.0\n" +
		"Content-Transfer-Encoding: 7bit\n\n" +
		htmlBody + "\n\n" +
		"--" + boundary + "\n"

	if len(fileName) != 0 {
		for i := range fileName {
			fileBytes, err := ioutil.ReadFile("assets/" + fileName[i])
			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			fileMIMEType := http.DetectContentType(fileBytes)

			fileData := base64.StdEncoding.EncodeToString(fileBytes)

			attachBody := "Content-Type: " + fileMIMEType + "; name=" + string('"') + fileName[i] + string('"') + " \n" +
				"MIME-Version: 1.0\n" +
				"Content-Transfer-Encoding: base64\n" +
				"Content-Disposition: attachment; filename=" + string('"') + fileName[i] + string('"') + " \n\n" +
				chunkSplit(fileData, 76, "\n") +
				"--" + boundary + "--"
			messageBody += attachBody
		}
	}

	message.Raw = base64.URLEncoding.EncodeToString([]byte(messageBody))

	// Send the message
	_, err := GmailService.Users.Messages.Send("me", &message).Do()
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Println("Message sent!")
	}
	return true, nil
}
