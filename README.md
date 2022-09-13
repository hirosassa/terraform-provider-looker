# Terraform Provider: Looker

This is a terraform provider plugin for managing [Looker](https://www.looker.com/) accounts.

## Install

You can use [Explicit Provider Source Locations](https://www.terraform.io/upgrade-guides/0-13.html#explicit-provider-source-locations).

```terraform
terraform {
  required_providers {
    looker = {
      source = "hirosassa/looker"
      version = "0.8.8"
    }
  }
}
```

## Usage

In-depth docs are available [on the Terraform registry](https://registry.terraform.io/providers/hirosassa/looker/latest).


## For developers

### How to run acceptance test

Before running acceptance test, you should set following environment variables:

```shell
export LOOKER_API_CLIENT_ID=YOUR_CLIENT_ID
export LOOKER_API_CLIENT_SECRET=YOUR_CLIENT_SECRET
export LOOKER_API_BASE_URL="https://example.com/"
```

Then run following command:

```shell
TF_ACC=1 go test ./...
```

or you can also use

```shell
make test-acceptance
```
