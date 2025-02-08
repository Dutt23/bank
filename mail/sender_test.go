package mail

import (
	"fmt"
	"github/dutt23/bank/util"
	"testing"

	"github.com/stretchr/testify/require"
)


func TestSendEmail(t *testing.T) {
  config, err := util.LoadConfig("..")
  require.NoError(t, err)
  fmt.Println(config.EmailSenderPassword)
  sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAdr, config.EmailSenderPassword)
  subject := "test email"
  content := `
  <h1>Hello world<h1>
  <p>This is a test message from <a href="http://techschool.guru">Tech School</a></p>
  `
  to := []string{"fill whatever required"}
  attachFiles := []string{"../README.MD"}
  err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
  require.NoError(t, err)
}