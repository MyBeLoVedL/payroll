package db

import (
	"context"
	"database/sql"
	"log"
	"time"
)

func IfCommitted(id int64) bool {
	card, err := q.SelectActiveTimecard(context.Background(), id)
	if err != nil {
		return false
	}
	log.Printf("timecard info %+v\n", card)
	return card.Committed.Int32 == 1
}

func CheckTimecard(id int64) error {
	_, err := q.SelectActiveTimecard(context.Background(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err := q.AddTimecard(context.Background(), id)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func UpdateTimecard(id int64, chargeNumber, hours int, date time.Time) error {
	card, _ := q.SelectActiveTimecard(context.Background(), id)

	err := q.AddTimecardRecord(context.Background(), AddTimecardRecordParams{
		ChargeNumber: int64(chargeNumber),
		CardID:       card.ID,
		Hours:        int32(hours),
		Date:         date,
	})

	return err
}

func SelectTimeCard(empID int64) (Timecard, error) {
	card, err := q.SelectActiveTimecard(context.Background(), empID)
	if err != nil {
		return Timecard{}, err
	}
	return card, nil
}
