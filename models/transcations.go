package models

type Bill struct {
	Id          string        `json:"id"`
	BillDate    string        `json:"billDate" binding:"required"`
	EntryDate   string        `json:"entryDate" binding:"required"`
	FinishDate  string        `json:"finishDate"`
	EmployeeID  string        `json:"employeeId" binding:"required"`
	CustomerID  string        `json:"customerId" binding:"required"`
	BillDetails []BillDetails `json:"billDetails" binding:"required"`
}

type BillDetails struct {
	Id           string `json:"id"`
	BillID       string `json:"billId" binding:"required"`
	ProductID    string `json:"productId" binding:"required"`
	ProductPrice int    `json:"productPrice" binding:"required"`
	Qty          int    `json:"qty" binding:"required"`
}