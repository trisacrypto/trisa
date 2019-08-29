
# Travel Rule Information Sharing Architecture for Virtual Asset Service Providers (TRISA)


### TRISA Protocol Proposal

* Transaction envelope: https://github.com/trisacrypto/trisa/blob/master/proto/trisa/protocol/v1alpha1/trisa.proto
* Identity examples (per country/region): https://github.com/trisacrypto/trisa/tree/master/proto/trisa/identity
* Blockchain specific: https://github.com/trisacrypto/trisa/tree/master/proto/trisa/data


### Getting Started

1. Install `bazel` https://docs.bazel.build/versions/0.28.0/install.html.
2. Generate certificates for each VASP
3. Create VASP server configuration files for each
4. Run `bazel run //cmd/trisa server` to startup each VASP


### Development

* Run `make build` and `make test` for compilation and test runs.


### Coming Soon

* Dockerization and automated demo setup will be available soon.
