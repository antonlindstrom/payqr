package payqr

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	_ "image/png"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSerialize(t *testing.T) {
	tests := []struct {
		name string
		have *Payment
		want string
	}{
		{
			name: "With correct json based on example (swedish domestic)",
			have: New("5536-7742", "Test AB", "1234", "1001", 50, time.Date(2022, time.August, 6, 0, 0, 0, 0, time.Local), WithCreationDate(time.Date(2022, time.July, 7, 0, 0, 0, 0, time.Local)), WithPaymentType("BG")),
			want: `{"uqr":1,"tp":1,"nme":"Test AB","cid":"1234","iref":"1001","idt":"20220707","ddt":"20220806","due":50,"pt":"BG","acc":"5536-7742"}`,
		},
		{
			name: "With different payment type (swedish domestic)",
			have: New("5536-7742", "Test AB", "1234", "1001", 50, time.Date(2022, time.August, 6, 0, 0, 0, 0, time.Local), WithCreationDate(time.Date(2022, time.July, 7, 0, 0, 0, 0, time.Local)), WithPaymentType("PG")),
			want: `{"uqr":1,"tp":1,"nme":"Test AB","cid":"1234","iref":"1001","idt":"20220707","ddt":"20220806","due":50,"pt":"PG","acc":"5536-7742"}`,
		},
		{
			name: "Foreign payment",
			have: New("DK4830004073013895", "Test company AB", "555555-5555", "934000000000159", 10.75, time.Date(2012, time.February, 15, 0, 0, 0, 0, time.Local), WithCreationDate(time.Date(2012, time.February, 15, 0, 0, 0, 0, time.Local)), WithPaymentType(PaymentTypeIBAN), WithCurrency("DKK"), WithAddress("1092 Köpenhamn"), WithCountryCode("SE"), WithBankCode("DABADKKK")),
			want: `{"uqr":1,"tp":1,"nme":"Test company AB","cid":"555555-5555","cc":"SE","iref":"934000000000159","idt":"20120215","ddt":"20120215","due": 10.75,"cur":"DKK","pt":"IBAN","acc":"DK4830004073013895","bc":"DABADKKK","adr": "1092 Köpenhamn"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := json.Marshal(test.have)
			require.NoError(t, err)
			assert.JSONEq(t, test.want, string(got))
		})
	}

}

func TestQR(t *testing.T) {
	tests := []struct {
		name string
		have *Payment
		want string
	}{
		{
			name: "With correct json based on example (swedish domestic)",
			have: New("5536-7742", "Test AB", "1234", "1001", 50, time.Date(2022, time.August, 6, 0, 0, 0, 0, time.Local), WithCreationDate(time.Date(2022, time.July, 7, 0, 0, 0, 0, time.Local)), WithPaymentType("BG")),
			want: "iVBORw0KGgoAAAANSUhEUgAAAgAAAAIAAQMAAADOtka5AAAABlBMVEX///8AAABVwtN+AAAFTUlEQVR42uydPbKkOgyF1UVAyBJYCkuDpbEUluCQgEJT1pFsmts90QSD73Hw6o2r55tEZf0dCeHh4eHh4eHh+T/PpHZ26Y5OT1ERmddpHTUNqdc4MusyreMmMuDiJKAxwJr/O+y9Ht15oW6Dpl73/KeXii75f0bdZMAVAQTcAaNqytfdISJmZuuk25hkMPt8ZUOUeZlXkXGzf0oPAloF5Gu7UPupyJAkG1Jnj48us+o6bmMigIC/Aw7Bi6T59Zk2+LA9W+KZ72ZdxVwbAY0C4JlEpNPL2YYk0u9OlVnNkOS7ayPg2QAPdfujy+ah9qDkqHbI1pUN6bQIZbZQd0xfY2UCfjcgTu/Xml8kPFO49J8awC6/HAKeDZi20QwppzwvpMnr1TPZ64MIReovCSDgHYDHx5wYXqRsc2pFGEU0k6GXUBdUAloDWLlEVY9sM5Ym5+R5E3imHPfYO+WJkODyJKAxgD0oZkiCBwXZTX47Ev6+oohilRULdfvjGuoSQAAs0ZxQ8TfwTOIvkoW6XqJT+/upFFsIaAugHtUeqKBYpFvLJYh1Uav1ULd/r6EQQACKMCjkd54eoRjn/m7Hj+yZUrW3q2RXBDQF0JIQ5xfJ28b5cjDPdIgHPovIiOwIQTEBbQGQ3eQ8prOqmcA1IdTVw0PdJWM3L7YcQgABN0C4pj6afajkuzhFUJT1V8qjIbnlTAQ0AdjCZvD0QI+Etq/01ZDMh9m/dDckAgiwYMYCF+nCFHWZ1pJJRSIlKMaZNuUeKxPQBsANyeqvltyY8sh0aTmojVu7RgvwPdQloAEAxGYJoe4rfBIq9mJitRCmWdfX3dVbu48AAmqoq4cg1A1NaymXlFAX13u/v3eNCWgCoOWnkTsHNfXZM2mYF4rzKON/qKEQ8LsB8GLoF9Y6/jpGMS7ahajuj3FJQGMAT2SkjxcpGnvDmx5p1kWgp5efQlcCng9A19cEz7X8mlPf91oaDAm/vIW6BBAQLxKqJaZH8m6fl1B2VPddnFIskYD2AFOkOailaTGkeuneyrs8un8c9iXg2QCLOkwc0vlYr9XSPIDFaB6oKM7v/X7PnQkg4GKIIT3C9JVClya122e3yYd2TgIaA6w+QRERSunHXIuyMZI1Xoc5CSDgVpTF8gD0eKCczzGOQhx9qe7njPre5SGgEYCMVjZzv2ROCBp5n6S56JGgM7i7NgJaAECDhrDWDMllJAldX/RofIJiDHd1EkDAOyA8k4+C2+OzlpwplgdgbHyLNPvn1DkBDwfYsIQMu3QuH7Byib9Ie4/RPLWZb3WdKwEE/AQgZwq5GcY5V7nkTKW6P62jbt8WrhHwfMA2pKH29SB4rpM0Z5nJ+hzqEtAEAEMzCY09Vw65Z7LUt8iJYm7zQ3GeAAJcOH8JdWNPgHsmxy6zTZLLkD7pUAh4PsAUIwmC55fX4UvYYsHIS6s22l3YfQsFAQR4AOt78zDbF88UdqjZ3jwbyYJrk7vkmoDnA8RVhxHpio/6+vCVdKdl1LE377JRgICmAFJmvn0RkQ9L1C1VnvtetysqAQR8XX9YirJrESmZprUsPfJlevu3xXsEPBdQPx3gZbP7nJbAkpZoCH3Yk0UAAfXTAZAjiXi7MBaPHB7qzlXguH/QoRDweEDZ9F1CnMkFSdULqQvnQ7pEQLsAC0Yk9ieO10EaTFDYp0Skh0aeAAI+Auq3SIqkNUnZDh8TFNvgBToCmgPEt4miySPzRVHgW2nwTZr4dEB3fPy4EQG/GlC/lobPWcUmowHz4V2ZyEJ1/9sIBwHPBvDw8PDw8PDw/MvzJwAA//9bOjn9jq1/AAAAAABJRU5ErkJggg==",
		},
		{
			name: "With different payment type (swedish domestic)",
			have: New("5536-7742", "Test AB", "1234", "1001", 50, time.Date(2022, time.August, 6, 0, 0, 0, 0, time.Local), WithCreationDate(time.Date(2022, time.July, 7, 0, 0, 0, 0, time.Local)), WithPaymentType("PG")),
			want: "iVBORw0KGgoAAAANSUhEUgAAAgAAAAIAAQMAAADOtka5AAAABlBMVEX///8AAABVwtN+AAAFS0lEQVR42uydMbK0OAyENUVAyBE4CkeDo3EUjuDQAYW2rJYMw8xstMHi1wq23rqY709UdkuWZKHRaDQajUb7f9qkZlm6vdNDVETmdVpHTUPqNUxmXaZ13EQGLBwENAZYy3+H3OveHdfFQVOvufzfS0WX8seomwxYIoCAO2BUTWW520XE3GydVhmTDOafr+KIMi/zKjJu9k/pTkCrgLJsC7rMuk66DUmKI3W2+ZRFXcdtTAQQ8O+AXbAjqS6Trn6G5eKJR1mbdRU72ghoFICTSUQ6vVhxJOmzU2VWcyT5fbQR8GyAS91+74p7FDFiUncbincVRzpMocwmdcf0UysT8LcBYb0vq+ri2xQW/VMD2OIPI+DZgGkbzZFKyPNCmLzKtNWTSV6qrlDk/JIAAt4B2HzsEMOOVHxOLQmjUDNaAb4j5W4XAloDFDVTfGa3bImamlmLI9nJVHSP7VMWOyPbUnUPAe0AbEMpjlS+LE6D6KbsHQm/VyRRLB42qdvvV6lLAAHwRD+Zds+WiP/ediSTup6iU/t9qskWAtoCqKtaW1aPpz1MDq2LXK1L3f49h0IAAS5xRAbNXYRH82L+CamLj2ybUrW9q0ZXBDQFcD/KYsLFDqHVFk0AF1UL4bOIjNi8IIoJaAuA6KZo3c6yZoKjqUhdnFcmdZeC3TzZsgsBBNwAfjT1ERAjaIriFEFSVnVBdCUDVg8hoDWAhcnq7hH1SOZJIhdHMt1SgqMPRyKAAI+Zkl0W7vJSXBtbmDwgkjqqexYxNGj+0MoEtAFA3k0jh6IWO29Wl+aLUpdxBfgudQloAIBiM7vE685sCTL2YsVqUZhmt75+XL1d9xFAwEXq7l7SGjWtNV1SpS6Wc5/fb40JaAKg9dOaS5OonEeq9oUdCcl5RNmfORQC/jgAp5jHzm87kndmRHWKZffHWCSgMYAHMmKFrjWHMnmBY02/zrqYKsbirdCVgOcDcOtrBc9n+tU6aRLaa66OhC9vUpcAAmJHQrbE6pGifAApFM/ue3FK9UQC2gNMkYL1SsTIpQ3nYkTUuOXR/LXZl4BnA0x1WHFIF229ukZNgLfmqe0ySM7nPt9jZwIIuDhilB6h+wpXPyLnbZ+tJm/aOQhoDLB6B0UolHofc03KRkvWeG3mJICAW1IWwwNQeYTKeRsegN6+S3a/RNT3Wx4CGgHIaGkzP5fsELLFImpv9Uhqt773o42ABgDQIrWHQmsZScKtL+5ovINijOPqIICAG2A9/0avr00jiaSsDw9A2/gWYfZn1zkBTwfUMVfWhYcOCpQToR0TsfNSziuvcyWAgB+AFPWr3s55jZlqdn9aR91+DFwjoAGAbT71Xk/EB1pJFClFcv671CWgCYDgTl9qXVotfEfoW8uJom/zS3KeAAI+pW7MCfCTybE2Aw09nl/qUAh4PsBUbYqCZ/XYedRU27tfetZGe+H7fQoFAQTUpnEUuloHBXo8k89Qs7l51pKFo03uJdcEPB/gJ1NVukjO1+Yr6Q6LqGNu3mWiAAFNASRm8MpZKXCdoykR+16nKyoBBPwYf1hbPD1mikGc53VfDNPL3wfvEfBkwC2Xpnrr09rdv5a4EPoyJ4sAAs6nA1COdF4Xxqgal7rzWeCYv9ShEPB4QJ30XSXO5FnZawfFEqLYSpcIaBfQoefb5ydGSat/O+uCp0Skz/09diaAgPtbJO5zFjOlc7JZdFBsgyfoCGgOEG8TRSmizJeKAp9Kgzdp4umAbv/6uBEBfxpweS3tcE+MO8QcasafPMJjafn3q3kEPBZAo9FoNBqN9l/aPwEAAP//77xJzLaPXRsAAAAASUVORK5CYII=",
		},
		{
			name: "Foreign payment",
			have: New("DK4830004073013895", "Test company AB", "555555-5555", "934000000000159", 10.75, time.Date(2012, time.February, 15, 0, 0, 0, 0, time.Local), WithCreationDate(time.Date(2012, time.February, 15, 0, 0, 0, 0, time.Local)), WithPaymentType(PaymentTypeIBAN), WithCurrency("DKK"), WithAddress("1092 Köpenhamn"), WithCountryCode("SE"), WithBankCode("DABADKKK")),
			want: "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",
		},
	}

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			q, err := test.have.QR()
			require.NoError(t, err)

			b, err := q.PNG(512)
			require.NoError(t, err)

			got := base64.StdEncoding.EncodeToString(b)
			assert.Equal(t, test.want, got)

			f, err := os.Create(fmt.Sprintf("testdata/TestQR-%d.png", i))
			require.NoError(t, err)
			defer f.Close()

			_, err = f.Write(b)
			require.NoError(t, err)
		})
	}
}

func TestSwishEncoding(t *testing.T) {
	tests := []struct {
		name               string
		have               *Payment
		haveEditableFields SwishEditableField
		want               string
	}{
		{
			name:               "With one editable field",
			have:               New("5536-7742", "Test AB", "1234", "Swish message", 50, time.Date(2022, time.August, 6, 0, 0, 0, 0, time.Local), WithCreationDate(time.Date(2022, time.July, 7, 0, 0, 0, 0, time.Local)), WithPaymentType("BG")),
			haveEditableFields: SwishAmountEditable,
			want:               "C1231111111;50.00;Swish message;2",
		},
		{
			name:               "With multiple editable fields",
			have:               New("5536-7742", "Test AB", "1234", "My message", 50, time.Date(2022, time.August, 6, 0, 0, 0, 0, time.Local), WithCreationDate(time.Date(2022, time.July, 7, 0, 0, 0, 0, time.Local)), WithPaymentType("BG")),
			haveEditableFields: SwishAmountEditable | SwishMessageEditable,
			want:               "C1231111111;50.00;My message;6",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.have.swishEncode("1231111111", WithEditableFields(test.haveEditableFields))
			assert.Equal(t, test.want, got)
		})
	}
}

func ExamplePayment_QR() {
	q, err := New("5536-7742", "Test AB", "1234", "1001", 50, time.Date(2022, time.August, 6, 0, 0, 0, 0, time.Local), WithCreationDate(time.Date(2022, time.July, 7, 0, 0, 0, 0, time.Local)), WithPaymentType("BG")).QR()
	if err != nil {
		return
	}

	b, err := q.PNG(512)
	if err != nil {
		return
	}

	fmt.Printf(`<img src="data:image/png;base64,%s" alt="QR code" />`, base64.StdEncoding.EncodeToString(b))
	// Output: <img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAgAAAAIAAQMAAADOtka5AAAABlBMVEX///8AAABVwtN+AAAFTUlEQVR42uydPbKkOgyF1UVAyBJYCkuDpbEUluCQgEJT1pFsmts90QSD73Hw6o2r55tEZf0dCeHh4eHh4eHh+T/PpHZ26Y5OT1ERmddpHTUNqdc4MusyreMmMuDiJKAxwJr/O+y9Ht15oW6Dpl73/KeXii75f0bdZMAVAQTcAaNqytfdISJmZuuk25hkMPt8ZUOUeZlXkXGzf0oPAloF5Gu7UPupyJAkG1Jnj48us+o6bmMigIC/Aw7Bi6T59Zk2+LA9W+KZ72ZdxVwbAY0C4JlEpNPL2YYk0u9OlVnNkOS7ayPg2QAPdfujy+ah9qDkqHbI1pUN6bQIZbZQd0xfY2UCfjcgTu/Xml8kPFO49J8awC6/HAKeDZi20QwppzwvpMnr1TPZ64MIReovCSDgHYDHx5wYXqRsc2pFGEU0k6GXUBdUAloDWLlEVY9sM5Ym5+R5E3imHPfYO+WJkODyJKAxgD0oZkiCBwXZTX47Ev6+oohilRULdfvjGuoSQAAs0ZxQ8TfwTOIvkoW6XqJT+/upFFsIaAugHtUeqKBYpFvLJYh1Uav1ULd/r6EQQACKMCjkd54eoRjn/m7Hj+yZUrW3q2RXBDQF0JIQ5xfJ28b5cjDPdIgHPovIiOwIQTEBbQGQ3eQ8prOqmcA1IdTVw0PdJWM3L7YcQgABN0C4pj6afajkuzhFUJT1V8qjIbnlTAQ0AdjCZvD0QI+Etq/01ZDMh9m/dDckAgiwYMYCF+nCFHWZ1pJJRSIlKMaZNuUeKxPQBsANyeqvltyY8sh0aTmojVu7RgvwPdQloAEAxGYJoe4rfBIq9mJitRCmWdfX3dVbu48AAmqoq4cg1A1NaymXlFAX13u/v3eNCWgCoOWnkTsHNfXZM2mYF4rzKON/qKEQ8LsB8GLoF9Y6/jpGMS7ahajuj3FJQGMAT2SkjxcpGnvDmx5p1kWgp5efQlcCng9A19cEz7X8mlPf91oaDAm/vIW6BBAQLxKqJaZH8m6fl1B2VPddnFIskYD2AFOkOailaTGkeuneyrs8un8c9iXg2QCLOkwc0vlYr9XSPIDFaB6oKM7v/X7PnQkg4GKIIT3C9JVClya122e3yYd2TgIaA6w+QRERSunHXIuyMZI1Xoc5CSDgVpTF8gD0eKCczzGOQhx9qe7njPre5SGgEYCMVjZzv2ROCBp5n6S56JGgM7i7NgJaAECDhrDWDMllJAldX/RofIJiDHd1EkDAOyA8k4+C2+OzlpwplgdgbHyLNPvn1DkBDwfYsIQMu3QuH7Byib9Ie4/RPLWZb3WdKwEE/AQgZwq5GcY5V7nkTKW6P62jbt8WrhHwfMA2pKH29SB4rpM0Z5nJ+hzqEtAEAEMzCY09Vw65Z7LUt8iJYm7zQ3GeAAJcOH8JdWNPgHsmxy6zTZLLkD7pUAh4PsAUIwmC55fX4UvYYsHIS6s22l3YfQsFAQR4AOt78zDbF88UdqjZ3jwbyYJrk7vkmoDnA8RVhxHpio/6+vCVdKdl1LE377JRgICmAFJmvn0RkQ9L1C1VnvtetysqAQR8XX9YirJrESmZprUsPfJlevu3xXsEPBdQPx3gZbP7nJbAkpZoCH3Yk0UAAfXTAZAjiXi7MBaPHB7qzlXguH/QoRDweEDZ9F1CnMkFSdULqQvnQ7pEQLsAC0Yk9ieO10EaTFDYp0Skh0aeAAI+Auq3SIqkNUnZDh8TFNvgBToCmgPEt4miySPzRVHgW2nwTZr4dEB3fPy4EQG/GlC/lobPWcUmowHz4V2ZyEJ1/9sIBwHPBvDw8PDw8PDw/MvzJwAA//9bOjn9jq1/AAAAAABJRU5ErkJggg==" alt="QR code" />
}
