/*
Copyright (c) 2022 by the TRISA contributors. All rights reserved.
Use of this source code is governed by the MIT License that can be found in
the LICENSE file.

Package trisa hosts several subpackages, methods, and structs for a range of
TRISA-related tasks. The packages include:

- api

- crypto

- envelope

- data

- gds

- handler

- keys

- mtls

- peers

Package crypto describes interfaces for the various encryption and hmac
algorithms that might be used to encrypt and sign transaction envelopes
being passed securely in the TRISA network. Subpackages implement
specific algorithms such as aesgcm for symmetric encryption or rsa for
asymmetric encryption. Note that not all encryption mechanisms are
legal in different countries, these interfaces allow the use of different
algorithms and methodologies in the protocol without specifying what must
be used.

Package envelope replaces the handler Package to provide utilities for
encrypting and decrypting trisa. SecureEnvelopes as well as sealing and
unsealing them. SecureEnvelopes are the unit of transfer in a TRISA
transaction and are used to securely exchange sensitive compliance
information. Security is provided through two forms of cryptography:
symmetric cryptography to encrypt and sign the TRISA payload and
asymmetric cryptography to seal the keys and secrets of the envelopes
so that only the recipient can open it.

Package data contains the Generic Transaction message for TRISA
transaction payloads. The payload aims to provide enough information
to link Travel Rule Compliance information in the identity payload
with a transaction on the blockchain or network. All fields are optional,
this message serves as a convenience for parsing transaction payloads.

Package gds contains the TRISADirectory API subpackage and TRISADirectory
models subpackage. TheTRISADirectory API implements the TRISADirectory service
that specifies how TRISA clients should interact with a directory service to
facilitate P2P transfers. And the TRISADirectory models subpackage contains
the VASP struct definition that represents the top-level directory entry for
certificate public key exchange. A VASP entry is also the primary point of
replication between directories that implement the directory replication
service. It maintains the version information to detect changes with respect
to a specific registered directory and facilitates anti-entropy gossip
protocols.

Package handler provides the envelope struct that wraps a SecureEnvelope
containing all of the information necessary to access the payload data.
The envelope can be edited and resealed to simplify TRISA exchanges.

Package keys provides interfaces and handlers for managing public/private key
pairs that are used for sealing and unsealing secure envelopes. This package
is not intended for use with symmetric keys that are used to encrypt payloads.

Package mtls provides methods that return the standard TLS configuration
for the TRISA network, loading the certificate from the specified provider.
Using the TLS configuration ensures that all TRISA peer-to-peer connections
are handled and verified correctly.

Package peers provides structs and methods to facilitate information exchanges
to other members of the TRISA network and directory service lookups.

It is important to note that most of the subpackages in this repository are
independent and they are implemented and tested separately from other
subpackages.
*/
package trisa
