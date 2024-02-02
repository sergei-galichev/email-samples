package tools

import (
	"bufio"
	"code.sajari.com/docconv/v2"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"html/template"
	"log"
	"mail-server/internal/model"
	"net/http"
	"path/filepath"
	"strings"
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

	if file[0].Filename != "" {
		fileExtension := filepath.Ext(file[0].Filename)

		f, err := file[0].Open()
		if err != nil {
			return model.MailUpload{}, errors.New("[ReadMultiForm]: cannot open uploaded document")
		}
		defer func() {
			err = f.Close()
			if err != nil {
				log.Fatal("[ReadMultiForm]: cannot close uploaded document")
			}
		}()

		switch fileExtension {
		case ".txt":
			scanner := bufio.NewScanner(f)
			builder := strings.Builder{}

			for scanner.Scan() {
				line := fmt.Sprintf("%s<br>", scanner.Text())
				builder.WriteString(line)
			}

			mailUpload.DocxContent = builder.String()
			if err = scanner.Err(); err != nil {
				log.Fatal("[ReadMultiForm]: scanner error: cannot read uploaded document")
			}
		case ".docx", ".doc":
			res, _, err := docconv.ConvertDocx(f)
			if err != nil {
				log.Fatal("[ReadMultiForm]: docconv error: cannot convert uploaded document")
			}

			lines := strings.Split(res, "\n")

			builder := strings.Builder{}

			for _, line := range lines {
				builder.WriteString(fmt.Sprintf("%s<br>", line))
			}

			mailUpload.DocxContent = builder.String()

		default:
			return model.MailUpload{},
				errors.New(
					"[ReadMultiForm]: uploaded document not allowed: try .docx, .doc or .txt",
				)
		}

	}

	return mailUpload, nil
}

func HTMLRender(w http.ResponseWriter, r *http.Request, data interface{}) error {
	filePath := "./index.html"

	t, err := template.ParseFiles(filePath)
	if err != nil {
		return errors.New("[HTMLRender]: cannot parse template")
	}

	err = t.Execute(w, data)
	if err != nil {
		return errors.New("[HTMLRender]: cannot execute template")
	}

	return nil
}
