package db

import (
	"context"
	"log"
	"strconv"
)

func GetHoursByEmpID(id int64) (int, error) {
	var hours int
	err := dbIns.QueryRow("select sum(hours)  from timecard_record  where card_id =  ( select id from timecard where emp_id = ?);", id).Scan(&hours)
	if err != nil {
		return 0, err
	}
	return hours, nil
}

func GetHoursByProject(id int64, charge int64) (int, error) {
	log.Printf("empID %v charge %v\n", id, charge)
	var hours int
	err := dbIns.QueryRow("select sum(hours)  from timecard_record  where card_id =  ( select id from timecard where emp_id = ?) and charge_number = ? ;", id, charge).Scan(&hours)
	if err != nil {
		return 0, err
	}
	return hours, nil
}

func GetPayYearToDate(empID int64) (float64, error) {
	var hours string
	err := dbIns.QueryRow("select sum(amount) from paycheck where emp_id = ? and year(end_date) >= year(now());", empID).Scan(&hours)
	if err != nil {
		return 0, err
	}
	amount, err := strconv.ParseFloat(hours, 64)
	if err != nil {
		return 0, err
	}
	return amount, nil
}

func GetIDByName(name string) ([]int64, error) {
	ids, err := q.GetIDByName(context.Background(), name)
	if err != nil {
		return nil, err
	}
	return ids, nil
}
