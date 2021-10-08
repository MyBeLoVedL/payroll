package db

import "log"

func GetHours(id int64) (int, error) {
	var hours int
	err := dbIns.QueryRow("SELECT sum(hours) FROM timecard_record  WHERE id = (SELECT id FROM timecard WHERE emp_id = ?) ", 1).Scan(&hours)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	return hours, nil
}
