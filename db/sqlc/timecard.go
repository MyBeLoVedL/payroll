package db

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type project struct {
	chargeNumber int
	projectName  string
}

func getProjects() []project {
	return []project{
		{2011, "database"},
		{2022, "operating system"},
		{2023, "compiler"},
		{2024, "web application"},
	}
}

func UpdateTimecard(id int64, chargeNumber, hours int, date time.Time) error {
	var cardID int64
	card, err := q.SelectActiveTimecard(context.Background(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			result, err := q.AddTimecard(context.Background(), id)
			if err != nil {
				return err
			}
			cardID, _ = result.LastInsertId()
		} else {
			return err
		}
	} else {
		cardID = card.ID
	}
	err = q.AddTimecardRecord(context.Background(), AddTimecardRecordParams{
		ChargeNumber: int64(chargeNumber),
		CardID:       cardID,
		Hours:        int32(hours),
		Date:         date,
	})
	log.Printf("inssert ress :%v\n", err)

	return err
}
