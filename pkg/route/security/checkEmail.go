package security

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gopkg.in/ezzarghili/recaptcha-go.v4"
	"gorm.io/gorm/clause"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
	"treehollow-v3-backend/pkg/base"
	"treehollow-v3-backend/pkg/consts"
	"treehollow-v3-backend/pkg/logger"
	"treehollow-v3-backend/pkg/mail"
	"treehollow-v3-backend/pkg/utils"
)

func checkEmailParamsCheckMiddleware(c *gin.Context) {
//	recaptchaVersion := c.PostForm("recaptcha_version")
//	recaptchaToken := c.PostForm("recaptcha_token")
	oldToken := c.PostForm("old_token")
	email := strings.ToLower(c.PostForm("email"))

	if len(email) > 100 || len(oldToken) > 32 {
		base.HttpReturnWithCodeMinusOneAndAbort(c, logger.NewSimpleError("CheckEmailParamsOutOfBound", "Wrong Parameter", logger.WARN))
		return
	}
	emailHash := utils.HashEmail(email)
	c.Set("email_hash", emailHash)
	c.Next()
}

func checkEmailRegexMiddleware(c *gin.Context) {
	email := strings.ToLower(c.PostForm("email"))
	emailCheck, err := regexp.Compile(viper.GetString("email_check_regex"))
	if err != nil {
		base.HttpReturnWithCodeMinusOneAndAbort(c, logger.NewError(err, "RegexError", "Server Error"))
		return
	}
	if !emailCheck.MatchString(email) {
		emailWhitelist := viper.GetStringSlice("email_whitelist")
		if _, ok := utils.ContainsString(emailWhitelist, email); !ok {
			base.HttpReturnWithCodeMinusOneAndAbort(c, logger.NewSimpleError("EmailRegexCheckNotPass", "Sorry, we are only open to CMU community for now", logger.INFO))
			return
		}
	}
}

func checkEmailIsRegisteredUserMiddleware(c *gin.Context) {
	emailHash := c.MustGet("email_hash").(string)
	var count int64
	//check if user is registered
	err := base.GetDb(false).Where("email_hash = ?", emailHash).Model(&base.Email{}).Count(&count).Error
	if err != nil {
		base.HttpReturnWithCodeMinusOneAndAbort(c, logger.NewError(err, "SearchEmailHashFailed", consts.DatabaseReadFailedString))
		return
	}
	if count == 1 {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
		})
		c.Abort()
		return
	}
	c.Next()
}

//compatibility settings
func checkEmailIsOldTreeholeUserMiddleware(c *gin.Context) {
	oldToken := c.PostForm("old_token")
	emailHash := c.MustGet("email_hash").(string)
	var count int64

	//check if user is old v2 version user
	err := base.GetDb(false).Where("old_email_hash = ? and old_token = ?", emailHash, oldToken).
		Model(&base.User{}).Count(&count).Error
	if err != nil {
		base.HttpReturnWithCodeMinusOneAndAbort(c, logger.NewError(err, "SearchOldEmailHashFailed", consts.DatabaseReadFailedString))
		return
	}
	if count == 1 {
		c.JSON(http.StatusOK, gin.H{
			"code": 2,
		})
		c.Abort()
		return
	}
	c.Next()
}

func checkEmailRateLimitVerificationCode(c *gin.Context) {
	emailHash := c.MustGet("email_hash").(string)

	now := utils.GetTimeStamp()
	_, timeStamp, _, _ := base.GetVerificationCode(emailHash)
	if now-timeStamp < 60 {
		base.HttpReturnWithCodeMinusOneAndAbort(c, logger.NewSimpleError("TooMuchEmailInOneMinute", "Please wait 1 minute.", logger.INFO))
		return
	}
	c.Next()
}

func checkEmail(c *gin.Context) {
	email := strings.ToLower(c.PostForm("email"))

	emailHash := c.MustGet("email_hash").(string)

	code := utils.GenCode()

	err := mail.SendValidationEmail(code, email)
	if err != nil {
		base.HttpReturnWithCodeMinusOne(c, logger.NewError(err, "SendEmailFailed"+email, "Failed to send code."))
		return
	}

	err = base.GetDb(false).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&base.VerificationCode{Code: code, EmailHash: emailHash, FailedTimes: 0, UpdatedAt: time.Now()}).Error
	if err != nil {
		base.HttpReturnWithCodeMinusOne(c, logger.NewError(err, "SaveVerificationCodeFailed", consts.DatabaseWriteFailedString))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "Code sent successfully. Please check spams and wait 1 minute to resend.",
	})
}

func checkPhone(c *gin.Context) {
	phone := strings.ToLower(c.PostForm("phone"))

	if phone != "4127582618" {
		base.HttpReturnWithCodeMinusOne(c, logger.NewSimpleError("Not registered phone. We will handle your request soon.", "Failed to send code.", logger.INFO))
		return
	}

	phoneHash := utils.HashEmail(phone)

	code := utils.GenCode()

	err := mail.SendValidationSMS(code, phone)
	if err != nil {
		base.HttpReturnWithCodeMinusOne(c, logger.NewError(err, "Send SMS Failed "+phone, "Failed to send code."))
		return
	}

	err = base.GetDb(false).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&base.VerificationCode{Code: code, EmailHash: phoneHash, FailedTimes: 0, UpdatedAt: time.Now()}).Error
	if err != nil {
		base.HttpReturnWithCodeMinusOne(c, logger.NewError(err, "SaveVerificationCodeFailed", consts.DatabaseWriteFailedString))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "Code sent successfully. Please check spams and wait 1 minute to resend.",
	})
}

func checkEmailInvitation(c *gin.Context) {
	email := strings.ToLower(c.PostForm("email"))

	code := viper.GetString("invitation_code")

	err := mail.SendValidationEmail(code, email)
	if err != nil {
		base.HttpReturnWithCodeMinusOne(c, logger.NewError(err, "SendEmailFailed"+email, "Failed to send code."))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2,
		"msg":  "Code sent successfully. Please check spams and wait 1 minute to resend.",
	})
}

func unregisterEmail(c *gin.Context) {
	email := strings.ToLower(c.PostForm("email"))

	emailHash := c.MustGet("email_hash").(string)

	code := utils.GenCode()

	err := mail.SendUnregisterValidationEmail(code, email)
	if err != nil {
		base.HttpReturnWithCodeMinusOne(c, logger.NewError(err, "SendEmailFailed"+email, "Failed to send verification code."))
		return
	}

	err = base.GetDb(false).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&base.VerificationCode{Code: code, EmailHash: emailHash, FailedTimes: 0, UpdatedAt: time.Now()}).Error
	if err != nil {
		base.HttpReturnWithCodeMinusOne(c, logger.NewError(err, "SaveVerificationCodeFailed", consts.DatabaseWriteFailedString))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "Verification code sent successfully.",
	})
}
