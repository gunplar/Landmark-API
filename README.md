# Landmark-API
Master thesis at Fra-UAS and HERE Technologies

## Java module
The java module is still under development, but development priority goes to
the Golang module.

## Golang module

### Prerequisites

- Go must be installed to build this module. The version used
for this demonstration is go 1.17. Please refer to
https://go.dev/doc/install

- Currently, AWS Route 53 is the only supported DNS provider
of this project. An AWS subscription is required.
Please refer to https://docs.aws.amazon.com/sdkref/latest/guide/access-users.html
for configuring AWS credentials on your machine.
- An AWS Route 53 hosted zone should also be created
beforehand, the zone name and zone id should be 
defined in the config file Landmark-API/golang/resources/appConfig.json with
content in JSON format, after cloning this repo:
```
{
  "ZoneName": "example.zone.com",
  "ZoneId": "jf89e2894hrfew8"
}
```
- The production use of this CLI will use the user's
AWS credential to encrypt the DNS content, but with
potential security vulnerability of this early
demonstration product, and ease of file reading 
considered, there should be a file named `pass`
in the resource directory, containing exactly one line
of string as the pseudo AWS credential.
```
echo "pseudo-password-instead-of-your-aws-one" > Landmark-API/golang/resources/pass
```


### Setup
Go to the repo directory and install dependencies for this
project: 
```
cd Landmark-API
go mod tidy
```

Go to the golang directory and 
install the landmark CLI:

```
cd golang
go install
```

Now the CLI is ready for use. Thanks to being built
from the Go library Cobra, autocomplete script can be 
generated. For example, running these command save
the script for bash shell in a file, and sourcing this file enables 
command autocomplete:
```
landmark completion bash > autocomplete
source autocomplete
```

### Usage

Again, for the ease of file reading, 
the CLI currently only function correctly from the 
`golang` directory.

```
cd Landmark-API/golang
```

All command usages can be displayed with the flag `-h`

An example of the work flow would be:

The end user store the data they want to share (geo-JSON,...)
in encrypted form in a DNS entry
```
landmark user publish-stored --domain usera --content "replace this with the geoJSON content"
```

The string will be encrypted with an AES key,
SHA-256 hash-derived from the pseudo AWS credential.
The ciphertext is then available at `usera.example.zone.com`.
The postal service will also generate an asymmetric
RSA-4096 key pair with the command:

```
landmark postal publish-key --domain real.dhl
```

The private key will be saved in the `resources`
directory, while the public key is available at 
`real.dhl.example.zone.com`.

When the user requires a delivery, they can encrypt the
AES key with the postal service's public key:

```
landmark user publish-shared --domain usera --postal-domain real.dhl
```

The encrypted key is available at `real.dhl.usera.example.zone.com`.
The postal service will then retrieve this key, decrypt 
it to get the AES key, and subsequently decrypt the
ciphertext on the user's domain to get the original geo-JSON,
all in one single command:

```
landmark postal show-user-data --postal-domain real.dhl --user-domain usera
```