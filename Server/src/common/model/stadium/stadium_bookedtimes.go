package stadium

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"common/connection"

	"gorm.io/gorm"
)

type StadiumBooked struct {
	gorm.Model
	StadiumId   uint       `gorm:"not null;index:idx_stadium_date"`
	Date        string     `gorm:"not null;index:idx_stadium_date"`
	BookedTimes BookedTime `gorm:"not null;"`
}

func InitBooked() {
	err := connection.GetDB().AutoMigrate(&StadiumBooked{})
	if err != nil {
		panic(err)
	}
}

type BookedTime map[string]string

func (idle *BookedTime) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := BookedTime{}
	err := json.Unmarshal(bytes, &result)
	*idle = result
	return err
}

// 实现 driver.Valuer 接口，Value 返回 json value
func (booked BookedTime) Value() (driver.Value, error) {

	return json.Marshal(booked)
}
