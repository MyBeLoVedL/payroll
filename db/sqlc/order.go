package db

import (
	"context"
	"time"
)

//   `id` bigint PRIMARY KEY AUTO_INCREMENT,
//   `emp_id` bigint,
//   `customer_contact` varchar(255),
//   `customer_address` varchar(255),
//   `order_info_id` bigint,
//   `date` timestamp,
//   `status` smallint

//   `order_id` bigint,
//   `product_id` bigint,
//   `amount` int

func AddOrderInfo(orderID, productID int64, amount string) error {
	return q.AddOrderInfo(context.Background(), AddOrderInfoParams{orderID, productID, amount})
}

func AddPurchaseOrder(empID int64, contact, address string, date time.Time) (int64, error) {
	result, err := q.AddPurchaseOrder(context.Background(), AddPurchaseOrderParams{empID, contact, address, date})
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return id, nil
}

type AddOrderParams struct {
	EmpID     int64     `form:"empId"`
	Contact   string    `form:"contact"`
	Address   string    `form:"addres"`
	Date      time.Time `form:"date" time_format:"2006-01-02"`
	ProductID int64     `form:"productId"`
	Amount    string    `form:"amount"`
}

func (s *Store) AddOrder(ctx context.Context, arg AddOrderParams) (int64, error) {
	var err error
	var orderID int64
	err = s.execTrx(ctx, func(q *Queries) error {
		ordID, err := AddPurchaseOrder(arg.EmpID, arg.Contact, arg.Address, arg.Date)
		if err != nil {
			return err
		}

		err = AddOrderInfo(ordID, arg.ProductID, arg.Amount)
		if err != nil {
			return err
		}
		orderID = ordID
		return nil
	})
	return orderID, err
}

func UpdateOrderInfo(orderID, productID int64, amount string) error {
	return q.UpdateOrderInfo(context.Background(), UpdateOrderInfoParams{productID, amount, orderID})
}

func UpdatePurchaseOrder(orderID int64, contact, address string, date time.Time) error {
	return q.UpdatePurchaseOrder(context.Background(), UpdatePurchaseOrderParams{contact, address, date, orderID})
}

type UpdateOrderParams struct {
	ID        int64     `form:"empId"`
	Contact   string    `form:"contact"`
	Address   string    `form:"addres"`
	Date      time.Time `form:"date" time_format:"2006-01-02"`
	ProductID int64     `form:"productId"`
	Amount    string    `form:"amount"`
}

func (s *Store) UpdateOrder(ctx context.Context, arg UpdateOrderParams) error {
	var err error
	err = s.execTrx(ctx, func(q *Queries) error {
		err = UpdateOrderInfo(arg.ID, arg.ProductID, arg.Amount)
		if err != nil {
			return err
		}

		err = UpdatePurchaseOrder(arg.ID, arg.Contact, arg.Address, arg.Date)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

type Order struct {
	ID        int64     `form:"empId"`
	Contact   string    `form:"contact"`
	Address   string    `form:"addres"`
	Date      time.Time `form:"date" time_format:"2006-01-02"`
	ProductID int64     `form:"productId"`
	Amount    string    `form:"amount"`
	Closed    int
}

func SelectOrderByID(id int64) (Order, error) {
	pur, err := q.SelectOrderById(context.Background(), id)
	if err != nil {
		return Order{}, err
	}
	info, err := q.SelectOrderInfoById(context.Background(), id)
	if err != nil {
		return Order{}, err
	}

	return Order{
		ID:        pur.ID,
		Contact:   pur.CustomerContact,
		Address:   pur.CustomerAddress,
		Date:      pur.Date,
		ProductID: info.ProductID,
		Amount:    info.Amount,
		Closed:    int(pur.Closed.Int32),
	}, nil

}

func DeleteOrder(orderID int64) error {
	var err error

	err = q.DeleteOrderInfoById(context.Background(), orderID)
	if err != nil {
		return err
	}

	err = q.DeletePurchaseOrderById(context.Background(), orderID)
	if err != nil {
		return err
	}

	return nil
}
