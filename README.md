# Terraform Provider for Zesty

- [Terraform Provider for Zesty](#terrform-provider-for-zesty)
    - [Prerequisites](#prerequisites)
    - [Usage](#usage)
    - [Local Development](#local-development)

## Prerequisites

- [Terraform](https://developer.hashicorp.com/terraform/install) 0.13+
- [Go](https://go.dev/doc/install) 1.24 (for local development)

## Usage

```terraform
terraform {
  required_providers {
    zesty = {
      source  = "zesty-co/zesty"
      version = "0.1.3"
    }
  }
  required_version = ">= 0.13"
}

provider "zesty" {
  token = "{{zesty-token}}"
}
```

To avoid committing your token to git, you can pass the token via environment variables:

```bash
$ ZESTY_API_TOKEN={{zesty-token}} terraform plan
```

To see additional information, you can increase the terraform log level, e.g:

```bash
$ TF_LOG=DEBUG terraform plan
```

For more examples, review the [examples](examples/) directory.

## Local Development

Before making any contribution, please make sure to review the [contribution guidelines](./CONTRIBUTING.md)

Steps to install locally:

```bash
git clone https://github.com/zesty-co/terraform-provider-zesty.git
cd terraform-provider-zesty
make terraformrc
make install
```

Steps to run unit tests:

```bash
make test
```
