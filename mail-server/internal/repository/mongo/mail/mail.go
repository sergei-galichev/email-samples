package mail

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"mail-server/internal/model"
	"time"
)

func (r *repository) AddMail(mailUpload model.MailUpload) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.storage.Collection("mails").InsertOne(ctx, mailUpload)
	if err != nil {
		return "", errors.New("[AddMail]: unable to add new mail")
	}

	return fmt.Sprintf("[AddMail]: new mail added"), nil
}
