package app

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// Server Config
	AppUrlEnv             = "APP_URL"
	PortEnv               = "PORT"
	JWTSecretEnv          = "JWT_SECRET"
	JWTAuthDurationEnv    = "JWT_AUTH_DURATION"
	JWTRefreshDurationEnv = "JWT_REFRESH_DURATION"

	// Databases Config
	PrismaHostEnv = "PRISMA_HOST"
	RedisHostEnv  = "REDIS_HOST"

	// Email Service Config
	EmailEnv         = "EMAIL"
	EmailHostEnv     = "EMAIL_HOST"
	EmailPortEnv     = "EMAIL_PORT"
	EmailUsernameEnv = "EMAIL_USERNAME"
	EmailPasswordEnv = "EMAIL_PASSWORD"
)

var (
	// Server Config
	AppUrl             string
	Port               string
	JWTSecret          string
	JWTAuthDuration    int
	JWTRefreshDuration int

	// Databases Config
	PrismaHost string
	RedisHost  string

	// Email Service Config
	Email         string
	EmailHost     string
	EmailPort     int
	EmailUsername string
	EmailPassword string
)

var Cmd = &cobra.Command{
	Use:   "REST API",
	Short: "A REST API server.",
	Long:  "A Golang REST API server w/ auth.",
}

func init() {
	// Start viper
	viper.AutomaticEnv()

	/* Server Configuration */

	// Get url
	Cmd.PersistentFlags().StringVarP(&AppUrl, "app-url", "u", viper.GetString(AppUrlEnv), "Set the app's url to be used")

	// Get port
	Cmd.PersistentFlags().StringVarP(&Port, "port", "p", viper.GetString(PortEnv), "Set the port to be used")

	// Get jwt secret
	Cmd.PersistentFlags().StringVarP(&JWTSecret, "jwt-secret", "s", viper.GetString(JWTSecretEnv), "Set the JWT secret to be used")

	// Get jwt duration
	Cmd.PersistentFlags().IntVarP(&JWTAuthDuration, "jwt-auth-duration", "t", viper.GetInt(JWTAuthDurationEnv), "Set the JWT Auth Token duration before it needs to be refreshed")
	Cmd.PersistentFlags().IntVarP(&JWTRefreshDuration, "jwt-refresh-duration", "r", viper.GetInt(JWTRefreshDurationEnv), "Set the JWT Refresh Token duration before it needs to be refreshed")

	/* Databases Configuration */
	Cmd.PersistentFlags().StringVarP(&PrismaHost, "prisma", "", viper.GetString(PrismaHostEnv), "Set the Host name for the prisma service.")
	Cmd.PersistentFlags().StringVarP(&RedisHost, "redis", "", viper.GetString(RedisHostEnv), "Set the Host name for the redis database.")

	/* Email Service Configuration */

	// Get our email
	Cmd.PersistentFlags().StringVarP(&Email, "email", "", viper.GetString(EmailEnv), "Set the email to e used")

	// Get email service host
	Cmd.PersistentFlags().StringVarP(&EmailHost, "email-host", "", viper.GetString(EmailHostEnv), "Set the email service's host")

	// Get email service port
	Cmd.PersistentFlags().IntVarP(&EmailPort, "email-port", "", viper.GetInt(EmailPortEnv), "Set the email service's port")

	// Get email service username
	Cmd.PersistentFlags().StringVarP(&EmailUsername, "email-username", "", viper.GetString(EmailUsernameEnv), "Set the email service's username")

	// Get email service username
	Cmd.PersistentFlags().StringVarP(&EmailPassword, "email-password", "", viper.GetString(EmailPasswordEnv), "Set the email service's password")
}
