package pkg

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"

	"google.golang.org/api/drive/v3"
)

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// Contains check string exist in slice
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Rename rename existing file and replace with new path
// same like move.
func Rename(oldPath string, path string, f os.FileInfo) error {
	return os.Rename(
		fmt.Sprintf("%s/%s", oldPath, string(f.Name())),
		fmt.Sprintf("%s/%s", path, f.Name()),
	)
}

// SendMail send mail notification along with
// file's information.
func SendMail(mail *Mail, f *drive.File) error {
	from := mail.From
	pass := mail.Password
	to := mail.To

	msg := fmt.Sprintf(
		"From: %s \n"+"To: %s \n"+"Subject: [%s]:%s \n\n"+"%s %v %v",
		from,
		to,
		"AUTOMATE SCRIPTS",
		"Automation Backup From Server",
		"This is to inform you that the backup process has successfully done. \n\n" +
		"Download Link: ",
		f.WebViewLink,
		"\nFolder: https://drive.google.com/drive/folders/" + strings.Join(f.Parents, ""),
	)

	err := smtp.SendMail(fmt.Sprintf("%s:%s", mail.Host, mail.Port),
		smtp.PlainAuth("", mail.Username, pass, mail.Host),
		from, []string{to}, []byte(msg))

	if err != nil {
		return err
	}

	return nil
}