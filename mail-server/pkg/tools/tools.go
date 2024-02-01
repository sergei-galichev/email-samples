package tools

import (
	"encoding/json"
	"github.com/pkg/errors"
	"mail-server/internal/model"
	"net/http"
	"time"
)

func ReadForm(r *http.Request, subscriber model.Subscriber) (model.Subscriber, error) {
	if err := r.ParseForm(); err != nil {
		return model.Subscriber{}, errors.New("[ReadForm]: cannot parse form")
	}
	subscriber = model.Subscriber{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Interest:  r.Form.Get("interest"),
	}

	return subscriber, nil
}

func JSONWriter(w http.ResponseWriter, message string, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	bytes, err := json.Marshal(message)
	if err != nil {
		return errors.New("[JSONWriter]: cannot marshal json")
	}

	_, err = w.Write(bytes)
	if err != nil {
		return errors.New("[JSONWriter]: cannot write json")
	}

	return nil
}

func ReadMultiForm(w http.ResponseWriter, r *http.Request, mailUpload model.MailUpload) (model.MailUpload, error) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return model.MailUpload{}, errors.New("[ReadMultiForm]: cannot parse form")
	}

	form := r.MultipartForm

	mailUpload.DocxName = form.Value["docx_name"][0]

	mailUpload.Date = time.Now()

	file, ok := form.File["docx"]
	if !ok {
		return model.MailUpload{}, errors.New("[ReadMultiForm]: cannot get uploaded document")
	}

	_ = file

	return mailUpload, nil
}
