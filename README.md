# Secrets-Validator

[![Go Report Card](https://goreportcard.com/badge/github.com/Kilemonn/Secrets-Validator)](https://goreportcard.com/report/github.com/Kilemonn/Secrets-Validator)

## Overview

An commandline application that is used to verify the value/format of your stored secrets across multiple managers.

The current credential providers that are supported are:
- Google Cloud Secrets Manager
- AWS Secrets Manager
- Kubernetes Secrets
- The local machine environment variables

### Please refer to the [Wiki](https://github.com/Kilemonn/Secrets-Validator/wiki) for more details!

## Quick Start

Installation of the commandline tool can be done with the following command:

> go install github.com/Kilemonn/Secrets-Validator@latest

## Usage

The application requires a `.yaml` configuration file that defines the credential providers along with the constraints that you want to perform on each credential.

The application can be run using the following command (using -f to specify the file path):
> Secrets-Validator.exe -f path/to/file.yaml

**You can also pass in `-d` to enable debug to log all constraint and pattern matching output.**

Using the environment as an example we can define the following `yaml` configuration file to check that the database properties are set correctly (this is an example to demonstrate what kind of validation is available).

``` yaml
credential-providers:
    - Env: # Registers the environment as a credential provider
constraints:
    - database-connection-string-prefix-is-development: # Create a new constraint with arbitrary name
        pattern: db-host-name # This is the regex that will be matched against the credential name (in this case, environment variable name) - if this matches successfully then the following condition will be evaluated against this secret's value (environment variable value)
        condition: HasPrefix(jdbc://path-to-development-db)
    - database-port-number:
        pattern: db-port-number
        condition: IsNumeric
    - all-properties-are-unqiue:
        pattern: ALL
        condition: Unique
    - is-email: # Checks any property/secret that contains "email address" and matches the following regex
        pattern: *.email-address*.
        condition: Matches(.+@\..+)
```

The above configuration defines the "environment" as the only credential provider.
In this case, all environment variables are loaded and will be run against each of the defined constraints.
In the `constraints` definition, the `pattern` is a **Regular expression (Regex)** pattern that is run against the secret name (in this case the environment variable name). If the pattern matches then the application will attempt to perform the condition against the secret's value (in this case the environment varriable's value).

There is an "ALL" `pattern` keyword that will force match against all entries. In this case, the `all-properties-are-unique` will most likely fail, as generally environment variables do have duplicated values etc.

## Constraints

The constraints that are supported are:

- **Unique** - That the value of this property is unique across all matching credentials.
- **HasPrefix(\<prefix-string\>)** - Check it has the supplied prefix.
- **HasSuffix(\<suffix-string\>)** - Check it has the supplied suffix.
- **IsNumber** - Check value is numeric.
- **IsBoolean** - Check value is a boolean.
- **Matches(<regex-pattern>)** - Check that the value matches the provided regex.

## Further Documentation in the Wiki

Please refer to the [Wiki](https://github.com/Kilemonn/Secrets-Validator/wiki) for more documentation about how to configure different credential providers and their different usages.
