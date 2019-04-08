package app

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	PortEnv        = "PORT"
	JWTSecretEnv   = "JWT_SECRET"
	JWTDurationEnv = "JWT_DURATION"
)

var (
	Port        string
	JWTSecret   string
	JWTDuration int
)

var Cmd = &cobra.Command{
	Use:   "Calendoer",
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
}
