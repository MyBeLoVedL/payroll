package db

func GetHours(id int64) (int, error) {
	var hours int
	err := dbIns.QueryRow("SELECT sum(hours) FROM timecard_record  WHERE id = (SELECT id FROM timecard WHERE emp_id = ?) ", id).Scan(&hours)
	if err != nil {
		return 0, err
	}
	return hours, nil
}

func GetAmountByID(id int64) (float64, error) {
	var total float64
	err := dbIns.QueryRow("select sum(amount) from purchase_order as p join order_info as o on p.id = order_id where emp_id = ? ", id).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}
