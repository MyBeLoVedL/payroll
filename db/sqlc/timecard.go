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
		if card.Committed.Int32 == 1 {

		} else {
			cardID = card.ID
		}
	}
	err = q.AddTimecardRecord(context.Background(), AddTimecardRecordParams{
		ChargeNumber: int64(chargeNumber),
		CardID:       cardID,
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
