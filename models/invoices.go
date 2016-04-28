package models

import (
	"time"
)

//const msyqlDatetimeFormat = "2006-01-02 15:04:05"

type Invoices struct {
	InvoiceId uint32     `column:"invoice_id"`
	UserId    uint16     `column:"user_id"`
	Memo      string     `column:"memo"`
	DeleteFlg string     `column:"delete_flg"`
	Created   *time.Time `column:"create_datetime"`
	Updated   *time.Time `column:"update_datetime"`
}

//time.Now().Format(msyqlDatetimeFormat)
func NewInvoice(id uint32, userid uint16, memo string) Invoices {
	return Invoices{
		InvoiceId: id,
		UserId:    userid,
		Memo:      memo,
	}
}
