package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Date struct {
	time.Time
}

const dateLayout = "2006-01-02"

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Format(dateLayout))
}

func (d *Date) UnmarshalJSON(data []byte) error {
	var dateStr string
	if err := json.Unmarshal(data, &dateStr); err != nil {
		return err
	}

	parsedTime, err := time.Parse(dateLayout, dateStr)
	if err != nil {
		return err
	}

	d.Time = parsedTime
	return nil
}

func (d Date) Value() (driver.Value, error) {
	if d.IsZero() {
		return nil, nil
	}
	return d.Format(dateLayout), nil
}

func (d *Date) Scan(value interface{}) error {
	if value == nil {
		d.Time = time.Time{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		d.Time = v
		return nil
	case string:
		parsedTime, err := time.Parse(dateLayout, v)
		if err != nil {
			return err
		}
		d.Time = parsedTime
		return nil
	case []byte:
		parsedTime, err := time.Parse(dateLayout, string(v))
		if err != nil {
			return err
		}
		d.Time = parsedTime
		return nil
	default:
		return errors.New("неподдерживаемый тип для Date")
	}
}
