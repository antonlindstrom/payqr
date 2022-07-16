package payqr

import (
	"fmt"

	"github.com/skip2/go-qrcode"
)

// SwishOption defines options for Swish QRs.
type SwishOption func(*Payment)

// SwishEditableField is a way to set/lock different editable fields.
type SwishEditableField byte

const (
	SwishPhoneEditable   SwishEditableField = 0b001
	SwishAmountEditable  SwishEditableField = 0b010
	SwishMessageEditable SwishEditableField = 0b100
)

// WithEditableFields sets/locks different fields for editing when the user
// opens the QR code in their app. Setting fields as editable will
// automatically lock the other ones.
func WithEditableFields(fields SwishEditableField) SwishOption {
	return func(p *Payment) {
		p.swishEditableFields = byte(fields)
	}
}

// swishEncode encodes a payment to the format used by Swish in QR codes.
func (d *Payment) swishEncode(phoneNumber string, options ...SwishOption) string {
	for _, opt := range options {
		opt(d)
	}

	return fmt.Sprintf("C%s;%.2f;%s;%d", phoneNumber, d.DueAmount, d.Reference, int(d.swishEditableFields))
}

// SwishQR returns a QR code that can be used for Swish payments.
func (d *Payment) SwishQR(phoneNumber string, options ...SwishOption) (*qrcode.QRCode, error) {
	return qrcode.New(d.swishEncode(phoneNumber, options...), qrcode.High)
}
