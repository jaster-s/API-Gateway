package lib

import (
	"strconv"
)

type HTML_Template struct{}

func (HTML_Template) MainTemplate(message string) string {
	var htmlBody = `<html>

	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
	</head>
	
	<body>
		<center>
			<table width="100%" style="width:100%;max-width:600px" align="center">
				<tr>
					<td style="padding:0px 15px 0px 15px;color:#000000;text-align:left" bgcolor="#ffffff" width="100%">
						` + message + `
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
	return htmlBody
}

func (HTML_Template) APIisTrue(nameLINE_OA, CustomerName, email_contact, messageReq, unique_external_id string) string {
	var htmlBody = `
    <table border="0" width="100%" style="table-layout:fixed">
        <tr>
            <td style="padding:18px 0px 6px 0px;line-height:22px;text-align:inherit" height="100%" valign="top">
                <div>Dear ` + nameLINE_OA + `,</div>
                <br>
                <div>We have received the message from ` + CustomerName + ` via ` + nameLINE_OA + `</div>
            </td>
        </tr>
    </table>

    <table width="100%" style="table-layout:fixed">
        <tr>
            <td style="padding:6px 0px 6px px;line-height:30px;text-align:inherit;background-color:#ccf1ff" height="100%">
                <div style="text-align:left">
                    <div style="padding:18px 50px 18px 50px">
                        Customer name: ` + nameLINE_OA + ` - ` + CustomerName + `,
                        <br>
                        <span style="padding:0px 0px 0px 50px">` + messageReq + `
                            <br>
                        </span>
                    </div>
                </div>
            </td>
        </tr>
    </table>

    <table width="100%" style="table-layout:fixed">
        <tr>
            <td style="padding:18px 0px 6px 0px;line-height:22px;text-align:inherit" height="100%" valign="top">
                Unfortunately, the system is temporary unavailable to automatically create the ticket. Please manually create the
                ticket by the following steps:<br>
                <div style="padding:0px 0px 0px 30px">
                    1. For a new customer, please create the new contact before creating the ticket.
                    Go to the top right pane, click New, and select New contact.<br>
                    Enter the following information: <br>
                    <div style="padding:0px 0px 0px 30px">
                        Contact name: <strong>` + nameLINE_OA + ` - ` + CustomerName + `</strong><br>
                        Email: <strong>` + email_contact + `</strong><br>
                        Unique external ID: <strong>` + unique_external_id + `</strong><br>
                        Then, click Create.
                        If this is an existing customer, go to section 2 to create the New ticket.
                    </div>
                    2. Create a new ticket by clicking the New, and select New ticket. 
                    Enter the contact using the unique external ID and fill the necessary information such as subject, description and then click Create. <br>
                </div>
                <div>&nbsp;</div>
                <div>Or, directly response your customer via the Line OA.
                We are working to resolve the problem as quickly as we could. Sorry for the inconvenience.</div>
            </td>
        </tr>
    </table>
`
	template := new(HTML_Template)
	fullTemplate := template.MainTemplate(htmlBody)
	return fullTemplate
}

func (HTML_Template) KeyExpiration(nameLINE_OA, CustomerName, messageReq string) string {
	var htmlBody = `
    <table border="0" width="100%" style="table-layout:fixed">
        <tr>
            <td style="padding:18px 0px 6px 0px;line-height:22px;text-align:inherit" height="100%" valign="top">
                <div>Dear ` + nameLINE_OA + `,</div>
                <br>
                <div>We have received the message from ` + CustomerName + ` via ` + nameLINE_OA + `</div>
            </td>
        </tr>
    </table>

    <table width="100%" style="table-layout:fixed">
        <tr>
            <td style="padding:6px 0px 6px px;line-height:30px;text-align:inherit;background-color:#ccf1ff" height="100%">
                <div style="text-align:left">
                    <div style="padding:18px 50px 18px 50px">
                        Customer name: ` + nameLINE_OA + ` - ` + CustomerName + `,
                        <br>
                        <span style="padding:0px 0px 0px 50px;">` + messageReq + `
                            <br>
                        </span>
                    </div>
                </div>
            </td>
        </tr>
    </table>

    <table width="100%" style="table-layout:fixed">
        <tr>
            <td style="padding:18px 0px 6px 0px;line-height:22px;text-align:inherit" height="100%" valign="top">
                Unfortunately, our system cannot create the ticket due to key expiration. Please reactivate the key and create the ticket for this issue by yourself. After the account has been reactivated, the ticket will be automatically created when there are some messages from your customer
            </td>
        </tr>
    </table>
`
	template := new(HTML_Template)
	fullTemplate := template.MainTemplate(htmlBody)
	return fullTemplate
}

func (HTML_Template) RecoveryEmail(nameLINE_OA, email_LINE_OA string) string {
	var htmlBody = `
    <table border="0" width="100%" style="table-layout:fixed">
        <tr>
            <td style="padding:18px 0px 6px 0px;line-height:22px;text-align:inherit" height="100%" valign="top">
                <div>Dear ` + nameLINE_OA + `,</div>
                <br>
            </td>
        </tr>
    </table>
    <table width="100%" style="table-layout:fixed">
        <tr>
            <td style="padding:6px 0px 6px px;line-height:30px;text-align:inherit;background-color:#ccf1ff" height="100%">
                <div style="text-align:left">
                    <div style="padding:18px 50px 18px 50px">
                    The contact email of your customer ` + nameLINE_OA + ` has been changed to ` + email_LINE_OA + `
                    </div>
                </div>
            </td>
        </tr>
    </table>
    `
	template := new(HTML_Template)
	fullTemplate := template.MainTemplate(htmlBody)
	return fullTemplate
}
func (HTML_Template) ReMessagetoMail(message string, ticket_Id int) string {
	var htmlBody = `
    <table border="0" width="100%" style="table-layout:fixed">
        <tr>
            <td style="padding:18px 0px 6px 0px;line-height:22px;text-align:inherit" height="100%" valign="top">
                <div>Dear ` + "Customer" + `,</div>
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
                        ` + message + `
                    </div>
                </div>
            </td>
        </tr>
    </table>
	`
	template := new(HTML_Template)
	fullTemplate := template.MainTemplate(htmlBody)
	return fullTemplate
}

func (HTML_Template) SendActivateKeyTrialVersion(name, activateKey string) string {
	var htmlBody = `
    <table border="0" width="100%" style="table-layout:fixed">
        <tr>
            <td style="padding:18px 0px 6px 0px;line-height:22px;text-align:inherit" height="100%" valign="top">
                <div>Dear ` + name + `,</div>
                <br>
            </td>
        </tr>
    </table>
    <table width="100%" style="table-layout:fixed">
        <tr>
            <td style="padding:6px 0px 6px px;line-height:30px;text-align:inherit;background-color:#ccf1ff" height="100%">
                <div style="text-align:left">
                    <div style="padding:18px 50px 18px 50px">
                        ขอขอบคุณที่สนใจแอปพลิเคชันของเรา ท่านสามารถนำคีย์ดังกล่าวไปใส่เพื่อเข้าสู่เวอร์ชั่นทดลองใช้ได้เลย โดยระยะเวลาทดสอลใช้ 14 วัน สามารถเพิ่ม LINE Official ได้ 1 Account  <br>
                        ` + activateKey + `
                    </div>
                </div>
            </td>
        </tr>
    </table>
	`
	template := new(HTML_Template)
	fullTemplate := template.MainTemplate(htmlBody)
	return fullTemplate
}
