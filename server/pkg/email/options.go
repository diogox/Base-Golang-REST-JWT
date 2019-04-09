package email

type EmailClientOptions struct {
	Host     string
	Port     int
	Username string
	Password string
}

type NewEmailOptions struct {
	Subject string
	Message string
}