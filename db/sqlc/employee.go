package db

import "context"

func SelectEmployee(id int64) (Employee, error) {
	emp, err := q.SelectEmployeeById(context.Background(), id)
	if err != nil {
		return Employee{}, nil
	}
	return emp, nil
}

func AddEmployee(etype, mail, socialSecurityNumber, standard_tax_deductions, other_deductions, phone, rate string) (int64, error) {
	res, err := q.AddEmployee(context.Background(), AddEmployeeParams{
		Type:                  EmployeesType(etype),
		Mail:                  mail,
		SocialSecurityNumber:  socialSecurityNumber,
		StandardTaxDeductions: standard_tax_deductions,
		OtherDeductions:       other_deductions,
		PhoneNumber:           phone,
		SalaryRate:            rate,
	})
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return id, nil
}

func UpdateEmployee(id int64, etype, mail, socialSecurityNumber, standard_tax_deductions, other_deductions, phone, rate string) error {
	return q.UpdateEmployee(context.Background(), UpdateEmployeeParams{
		ID:                    id,
		Type:                  EmployeesType(etype),
		Mail:                  mail,
		SocialSecurityNumber:  socialSecurityNumber,
		StandardTaxDeductions: standard_tax_deductions,
		OtherDeductions:       other_deductions,
		PhoneNumber:           phone,
		SalaryRate:            rate,
	})
}

func DeleteEmployee(id int64) error {
	return q.DeleteEmployee(context.Background(), id)
}
