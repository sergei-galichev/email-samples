package model

import (
	"time"
)

type Mail struct {
	Source      string
	Destination string
	Message     string
	Subject     string
	Name        string
}

type MailUpload struct {
	DocxName    string    `bson:"docx_name" json:"docx_name"`
	DocxContent string    `bson:"docx" json:"docx"`
	Date        time.Time `bson:"date" json:"date"`
}
