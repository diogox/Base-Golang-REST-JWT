package app

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	PortEnv          = "PORT"
	JWTSecretEnv     = "JWT_SECRET"
	JWTDurationEnv   = "JWT_DURATION"
	EmailEnv         = "EMAIL"
	EmailHostEnv     = "EMAIL_HOST"
	EmailPortEnv     = "EMAIL_PORT"
	EmailUsernameEnv = "EMAIL_USERNAME"
	EmailPasswordEnv = "EMAIL_PASSWORD"
)

var (
	Port          string
	JWTSecret     string
	JWTDuration   int
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

	// Get port
	Cmd.PersistentFlags().StringVarP(&Port, "port", "p", viper.GetString(PortEnv), "Set the port to be used")

	// Get jwt secret
	Cmd.PersistentFlags().StringVarP(&JWTSecret, "jwt-secret", "s", viper.GetString(JWTSecretEnv), "Set the JWT secret to be used")

	// Get jwt duration
	Cmd.PersistentFlags().IntVarP(&JWTDuration, "jwt-duration", "d", viper.GetInt(JWTDurationEnv), "Set the JWT duration before it needs to be refreshed")

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
