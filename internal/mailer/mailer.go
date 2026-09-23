package mailer

import "embed"

const (
	FromName            = "GoSocial"
	maxRetries          = 3
	UserWelcomeTemplate = "user_invitation.tmpl"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(tempplateFile, username, email string, data any, isSandbox bool) error
}
