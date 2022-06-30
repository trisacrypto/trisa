---
title: IVMS101
date: 2022-06-29T11:25:52-04:00
lastmod: 2022-06-29T11:25:52-04:00
description: "Working with IVMS101 JSON and Protocol Buffers"
weight: 30
---

# Working with IVMS101

IVMS101 JSON data should be marshaled and unmarshaled using the `encoding/json` package to produce data mostly compatible with the [ivmsvalidator.com](https://ivmsvalidator.com/) tool produced by 21 Analytics.

This page describes how to convert between JSON and protocol buffer formatting for IVMS101 records. Unless you are specifically developing against a Go code base with marshaled protocol buffer TRISA structs, we strongly recommend that you use `encoding/json` for IVMS 101 serialization and JSON exchange. The only reason to use the protojson style is if you're specifically working with protocol buffers and need a human-readable/editable format.

- For examples of fixtures that represent the JSON style, see the files from [the `trisa` directory](https://github.com/trisacrypto/trisa/tree/main/pkg/ivms101/testdata) that end in the extension `.json`.
- Fixtures that are unmarshaled using the `protojson` package are not compatible with the validator tool. For examples of these fixtures see the files from [the `trisa` directory](https://github.com/trisacrypto/trisa/tree/main/pkg/ivms101/testdata) that end in the extension `.pb.json`.



## Marshaling and Unmarshaling

The following is a bit of code for marshaling and unmarshaling IVMS101 identity payloads to and from a file for simple reference.

### IVMS101 JSON

Use the `encoding/json` package as follows:

```go
package ivms101_test

import (
	"encoding/json"
	"io/ioutil"

	"github.com/trisacrypto/trisa/pkg/ivms101"
)

func LoadJSON(path string) (identity *ivms101.IdentityPayload, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(path); err != nil {
		return nil, err
	}

	identity = &ivms101.IdentityPayload{}
	if err = json.Unmarshal(data, identity); err != nil {
		return nil, err
	}

	return identity, nil
}

func DumpJSON(identity *ivms101.IdentityPayload, path string) (err error) {
	var data []byte
	if data, err = json.Marshal(identity); err != nil {
		return err
	}

	if err = ioutil.WriteFile(path, data, 0644); err != nil {
		return err
	}
	return nil
}
```

### Protocol Buffers JSON

Not recommended, use IVMS101 JSON as described above in most cases. For specific protocol buffer use cases, use the `protojson` file as follows:

```go
package ivms101_test

import (
	"io/ioutil"

	"github.com/trisacrypto/trisa/pkg/ivms101"
	"google.golang.org/protobuf/encoding/protojson"
)

func LoadPBJ(path string) (identity *ivms101.IdentityPayload, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(path); err != nil {
		return nil, err
	}

	jsonpb := protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	identity = &ivms101.IdentityPayload{}
	if err = jsonpb.Unmarshal(data, identity); err != nil {
		return nil, err
	}

	return identity, nil
}

func DumpPBJ(identity *ivms101.IdentityPayload, path string) (err error) {
	jsonpb := protojson.MarshalOptions{
		Multiline:       true,
		Indent:          "  ",
		AllowPartial:    true,
		UseProtoNames:   true,
		UseEnumNumbers:  false,
		EmitUnpopulated: true,
	}

	var data []byte
	if data, err = jsonpb.Marshal(identity); err != nil {
		return err
	}

	if err = ioutil.WriteFile(path, data, 0644); err != nil {
		return err
	}
	return nil
}
```