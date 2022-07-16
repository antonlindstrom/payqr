payqr - Payments with QR
================

payqr helps to create QR codes for payments, mainly for the Swedish market but
may be applicable in other places.

The goal is to be able to generate multiple QR codes with various payment
methods with one library, applicable to the Swedish market. Further markets
may be looked into in the future.

Usage
-------------

Here follows a short example, more extensive documentation is available at the
Go docs.

	q, err := New("5536-7742", "Test AB", "1234", "My message", 50, time.Now()).QR()
	if err != nil {
		return
	}

	b, err := q.PNG(512)
	if err != nil {
		return
	}

	fmt.Printf(`<img src="data:image/png;base64,%s" alt="QR code" />`, base64.StdEncoding.EncodeToString(b))

For now, this supports:

* Bank transfers (BG, PG, IBAN and BBAN).
* Swish

Acknowledgements
-------------

Further documentation and basis of this library can be found at:
* https://www.qrkod.info/
* https://www.qrkod.info/specification.pdf

Reporting bugs
--------------

If you find any bugs or want to provide feedback, you can file bugs in the project's [GitHub Issues page](https://github.com/antonlindstrom/payqr).

Author
------

This project is maintained by [Anton Lindström](https://www.antonlindstrom.com) ([GitHub](https://github.com/antonlindstrom) | [Twitter](https://twitter.com/mycap))

License
-------

APACHE LICENSE 2.0
Copyright 2022 Anton Lindström

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
