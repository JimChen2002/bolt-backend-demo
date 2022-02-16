package mail

import (
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
	"strconv"
	"net/http"
)

func SendValidationSMS(code string, recipient string) error {
	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

	// curl -X POST https://api.twilio.com/2010-04-01/Accounts/ACe87a9f8ace3bdeda0a17388ebef4a66c/Messages.json \
	//                 --data-urlencode "Body=Hello from Twilio" \
	//                 --data-urlencode "From=+19205285793" \
	//                 --data-urlencode "To=+14127582618" \
	//                 -u ACe87a9f8ace3bdeda0a17388ebef4a66c:cefed84fe008f909dc7e654a95581ef0

	req, err := http.NewRequest("POST", "https://api.twilio.com/2010-04-01/Accounts/ACe87a9f8ace3bdeda0a17388ebef4a66c/Messages.json", 
															url.Values{"Body": {"Your verification code is: "+code}, "From": {"+19205285793"}, "To": {"+1"+recipient}})
	if err != nil {
		return err
	}
	req.SetBasicAuth("ACe87a9f8ace3bdeda0a17388ebef4a66c", "cefed84fe008f909dc7e654a95581ef0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func SendValidationEmail(code string, recipient string) error {
	websiteName := viper.GetString("name")
	m := gomail.NewMessage()
	m.SetHeader("From", viper.GetString("from_domain"))
	m.SetHeader("To", recipient)
	var title string
	if(len(code) > 6) {
		title = "[" + websiteName + "] Invitation Code"
	} else {
		title = "[" + websiteName + "] Validation Code"
	}
	m.SetHeader("Subject", title)

	msg := `<!DOCTYPE html>
<html lang="cn">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>` + title + `</title>
</head>
<body>
<p>Welcome to ` + websiteName + `!</p>
<p>This is your verification code. It is valid for 12 hours.</p>
<p><strong>` + code + `</strong></p>
</body>
</html>`

	port, err := strconv.Atoi(viper.GetString("smtp_port"))
	if err != nil {
		return err
	}
	m.SetBody("text/html", msg)
	m.AddAlternative("text/plain", "Hi,\n\nWelcome to "+websiteName+"!\n\n"+code+"\nThis is your verification code. It is valid for 12 hours.\n")
	d := gomail.NewDialer(viper.GetString("smtp_host"), port, viper.GetString("smtp_username"), viper.GetString("smtp_password"))

	if err = d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func SendUnregisterValidationEmail(code string, recipient string) error {
	websiteName := viper.GetString("name")
	m := gomail.NewMessage()
	m.SetHeader("From", viper.GetString("from_domain"))
	m.SetHeader("To", recipient)
	title := "[" + websiteName + "] Verification Code"
	m.SetHeader("Subject", title)

	msg := `<!DOCTYPE html>
<html lang="cn">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>` + title + `</title>
</head>
<body>
<p>Hi, You are deleting your account on ` + websiteName + `。</p>
<p>This is your verification code. It is valid for 12 hours.</p>
<p><strong>` + code + `</strong></p>
</body>
</html>`

	port, err := strconv.Atoi(viper.GetString("smtp_port"))
	if err != nil {
		return err
	}
	m.SetBody("text/html", msg)
	m.AddAlternative("text/plain", "Hi,\n\nYou are deleting your account on "+websiteName+"。\n\n"+code+"\nThis is your verification code. It is valid for 12 hours.\n")
	d := gomail.NewDialer(viper.GetString("smtp_host"), port, viper.GetString("smtp_username"), viper.GetString("smtp_password"))

	if err = d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func SendPasswordNonceEmail(nonce string, recipient string) error {
	websiteName := viper.GetString("name")
	m := gomail.NewMessage()
	m.SetHeader("From", viper.GetString("smtp_username"))
	m.SetHeader("To", recipient)
	title := "Welcome to " + websiteName
	m.SetHeader("Subject", title)

	msg := `<!DOCTYPE html>
<html lang="cn">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>` + title + `</title>
</head>
<body>
<p>Welcome to ` + websiteName + `!</p>
<p>The string below is necessary to delete your account. Please keep it safe.</p>
<p><strong>` + nonce + `</strong></p>
</body>
</html>`

	port, err := strconv.Atoi(viper.GetString("smtp_port"))
	if err != nil {
		return err
	}
	m.SetBody("text/html", msg)
	m.AddAlternative("text/plain", "Hi,\n\nWelcome to "+websiteName+"!\nThe string below is necessary to delete your account. Please keep it safe.\n"+nonce+"\n")
	d := gomail.NewDialer(viper.GetString("smtp_host"), port, viper.GetString("smtp_username"), viper.GetString("smtp_password"))

	if err = d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
