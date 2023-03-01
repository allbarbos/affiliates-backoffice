package tests

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ReadFile(t *testing.T, path string) []byte {
	/* #nosec */
	expected, err := os.ReadFile(path)
	assert.NoError(t, err)
	return expected
}

func escapeQuotes(s string) string {
	return strings.NewReplacer("\\", "\\\\", `"`, "\\\"").Replace(s)
}

func BuildAttachmentFile(t *testing.T, attachment bool, filePath, contentType string) (*bytes.Buffer, string) {
	var boundary string
	body := new(bytes.Buffer)
	if attachment {
		writer := multipart.NewWriter(body)
		/* #nosec */
		file, err := os.Open(filePath)
		assert.NoError(t, err)
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes("attachment"), escapeQuotes(filePath)))
		h.Set("Content-Type", contentType)
		part, err := writer.CreatePart(h)
		assert.NoError(t, err)
		bytes, err := io.ReadAll(file)
		assert.NoError(t, err)
		_, _ = part.Write(bytes)
		boundary = writer.FormDataContentType()
		assert.NoError(t, writer.Close())
	}
	return body, boundary
}

func BuildMockDB(t *testing.T) (sqlmock.Sqlmock, *gorm.DB) {
	var err error
	conn, mock, err := sqlmock.New()
	assert.NoError(t, err)
	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 conn,
		PreferSimpleProtocol: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)
	return mock, db
}
