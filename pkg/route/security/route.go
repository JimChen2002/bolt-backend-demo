package security

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/ulule/limiter/v3"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"treehollow-v3-backend/pkg/base"
	"treehollow-v3-backend/pkg/utils"
)

func ApiListenHttp() {
	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "TOKEN")
	r.Use(cors.New(corsConfig))

	r.POST("/v3/security/login/check_email",
		checkEmailParamsCheckMiddleware,
		checkEmailRegexMiddleware,
		checkEmailRateLimitVerificationCode,
		checkEmail)
	r.POST("/v3/security/login/check_email_code",
		checkEmailCode)
	r.POST("/v3/security/login/check_phone",
		checkPhone)
	r.POST("/v3/security/login/check_phone_code",
		checkPhoneCode)

	// code below is under development

	r.POST("/v3/security/login/create_account",
		loginParamsCheckMiddleware,
		checkAccountNotRegistered,
		loginCheckIOSToken,
		createAccount)
	r.POST("/v3/security/login/create_account_invitation",
		loginParamsCheckMiddleware,
		checkAccountNotRegistered,
		loginCheckIOSToken,
		createAccountInvitation)
	r.POST("/v3/security/login/change_password",
		checkAccountIsRegistered,
		changePassword)
	r.POST("/v3/security/login/unregister",
		checkAccountIsRegistered,
		deleteAccount)
	r.GET("/v3/security/devices/list", listDevices)
	r.POST("/v3/security/devices/terminate", terminateDevice)
	r.POST("/v3/security/logout", logout)
	r.POST("/v3/security/update_ios_token", updateIOSToken)

	listenAddr := viper.GetString("security_api_listen_address")
	if strings.Contains(listenAddr, ":") {
		_ = r.Run(listenAddr)
	} else {
		_ = os.MkdirAll(filepath.Dir(listenAddr), os.ModePerm)
		_ = os.Remove(listenAddr)

		listener, err := net.Listen("unix", listenAddr)
		utils.FatalErrorHandle(&err, "bind failed")
		log.Printf("Listening and serving HTTP on unix: %s.\n"+
			"Note: 0777 is not a safe permission for the unix socket file. "+
			"It would be better if the user manually set the permission after startup\n",
			listenAddr)
		_ = os.Chmod(listenAddr, 0777)
		err = http.Serve(listener, r)
	}
}
