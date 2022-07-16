// Package payqr helps to create QR codes for payments, mainly for the Swedish
// market but may be applicable in other places.
//
// Further documentation and basis of this library can be found at:
// * https://www.qrkod.info/
// * https://www.qrkod.info/specification.pdf
package payqr

import (
	"encoding/json"
	"time"

	"github.com/skip2/go-qrcode"
)

// Type defines the type of QR transfer.
type Type int

const (
	InvoiceType         Type = 1
	CreditInvoiceType   Type = 2
	CashPaidInvoiceType Type = 3
)

// PaymentType selects over which type of bank transfer system the payment
// should go over. Default is PaymentTypeBG.
type PaymentType string

const (
	PaymentTypeIBAN PaymentType = "IBAN"
	PaymentTypeBBAN PaymentType = "BBAN"
	PaymentTypeBG   PaymentType = "BG"
	PaymentTypePG   PaymentType = "PG"
)

// Payment is the structure for storing the data about a payment. Shoulc not
// be used directly but can be used as you see fit.
type Payment struct {
	UsingQRVersion         int         `json:"uqr"`
	Type                   Type        `json:"tp"`
	AccountName            string      `json:"nme"` // Company name of the sending party.
	CompanyID              string      `json:"cid"`
	Reference              string      `json:"iref"`
	CreditInvoiceReference string      `json:"cref,omitempty"`
	Currency               string      `json:"cur,omitempty"`
	VAT                    int         `json:"vat,omitempty"`
	HighVAT                int         `json:"vh,omitempty"`
	MediumVAT              int         `json:"vm,omitempty"`
	LowVAT                 int         `json:"vl,omitempty"`
	CreatedDate            string      `json:"idt"`
	DueDate                string      `json:"ddt"`
	DueAmount              float64     `json:"due"`
	PaymentType            PaymentType `json:"pt"`
	AccountNumber          string      `json:"acc"`
	BankCode               string      `json:"bc,omitempty"`
	CountryCode            string      `json:"cc,omitempty"`
	Address                string      `json:"adr,omitempty"`

	swishEditableFields byte
}

// Option is a modifyier for a Payment to add more data to it.
type Option func(*Payment)

// WithCreationDate sets the creation date. If not specified, default is
// today.
func WithCreationDate(t time.Time) Option {
	return func(p *Payment) {
		p.CreatedDate = t.Format("20060102")
	}
}

// WithType sets the QR transfer type.
func WithType(typ Type) Option {
	return func(p *Payment) {
		p.Type = typ
	}
}

// WithPaymentType sets the payment type. For Swedish domestic, PG and BG are
// the most common.
func WithPaymentType(typ PaymentType) Option {
	return func(p *Payment) {
		p.PaymentType = typ
	}
}

// WithCurrency sets the currency for foreign payments.
func WithCurrency(currency string) Option {
	return func(p *Payment) {
		p.Currency = currency
	}
}

// WithAddress sets the adr field in the payment.
func WithAddress(address string) Option {
	return func(p *Payment) {
		p.Address = address
	}
}

// WithCountryCode sets the country code, the format seems to be ISO 3166-1
// alpha-2.
func WithCountryCode(countryCode string) Option {
	return func(p *Payment) {
		p.CountryCode = countryCode
	}
}

// WithBankCode sets the bank code to the payment. This differs between the
// PaymentTypes and may contain BIC/SWIFT.
func WithBankCode(bankCode string) Option {
	return func(p *Payment) {
		p.BankCode = bankCode
	}
}

// New creates a new payment with the defined options. The input should give a
// fair default but may be modified with options.
func New(accountNumber, accountName, companyID, reference string, dueAmount float64, dueDate time.Time, options ...Option) *Payment {
	p := &Payment{
		UsingQRVersion: 1,
		Type:           InvoiceType,
		CreatedDate:    time.Now().Format("20060102"),
		AccountNumber:  accountNumber,
		AccountName:    accountName,
		Reference:      reference,
		DueAmount:      dueAmount,
		DueDate:        dueDate.Format("20060102"),
		PaymentType:    PaymentTypeBG,
		CompanyID:      companyID,
	}

	for _, opt := range options {
		opt(p)
	}

	return p
}

// HasRequiredFields checks if the payment has the required fields set per
// Type.
func (d *Payment) HasRequiredFields() bool {
	// These fields are required for all types.
	if d.UsingQRVersion < 1 || d.Type < 1 || d.Type > 3 || d.AccountName == "" || d.CompanyID == "" {
		return false
	}

	// Specific fields per type.
	switch d.Type {
	case CreditInvoiceType:
		return d.Reference != ""
	case InvoiceType:
		// ddt, due, pt, acc
		return d.DueDate != "" || d.DueAmount == 0 || d.AccountNumber != ""
	}

	return true
}

// QR returns a QR code that can be used to communicate how to send transfers.
func (d *Payment) QR() (*qrcode.QRCode, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	return qrcode.New(string(b), qrcode.High)
}
