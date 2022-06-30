// Copyright 2022 by the TRISA contributors. All rights reserved.
// License that can be found in the LICENSE file.

/* The pkg folder contains the reference implementation code, including compiled code generated from the protocol buffer definitions in the proto folder.

	- The iso3166 folder contains language codes. Dictated by the International Organization for Standardization, these codes are an important part of TRISA's international compatibility. Language codes provide a consistent way to refer to countries across the world without the need for a lingua franca.
	- The ivms101 folder extends the generated protobuf code with JSON loading utilities, validation helpers, short constants, etc. This package enables compatibility between international entities, making it easier to transmit descriptive information about individuals and organizations in a mutually intelligible way.
	- The trisa folder contains structs and methods for a range of TRISA-related tasks, such as performing cryptography and making mTLS connections.

*/
package pkg

