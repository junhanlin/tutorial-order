package entity

import (
	"database/sql/driver"
	"encoding/json"
)

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *JSONB) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}
	return nil
}

type Channel struct {
	ID                  int64 `gorm:"primaryKey"`
	Namespace           string
	BotProviderName     string
	CustomChannelId     string
	ListenWorkflowName  *string
	ListenProcessorName *string
	CurrContext         JSONB
}

type Blob struct {
	ChannelId int64 `gorm:"primaryKey"`
	BlobId    int64 `gorm:"primaryKey"`
	FileType  string
	FileName  *string
	Bucket    string
	FileKey   string
	Size      int64
	Mime      string
}
