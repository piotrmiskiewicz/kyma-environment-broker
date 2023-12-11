# EDP tools

This folder contains tools which allows to get information about subaccounts registered in EDP and execute registraion.

## EDP tool

EDP tool allows to connect to the EDP and execute the following commands:
 - `get` - gathers information about registered subaccount. If not found, returns a message `Not found`
 - `register` - performs registration 
 - `deregister` - removes the registration

Above commands implementation contains a copied code from existing steps implementation. Before running check if the command implementation is up-to-date.

### Build

Run the following command to build the binary:

```
go build -o edp main.go
```

### Running

#### Setting environment variables

Before using `edp` tool you must set environment variables:

1. Copy existing template file, for example: 
`cp env.dev.template env.dev`
2. Set missing secret value in the file
3. Export environment variables:
`export $(grep -v '^#' env.dev | xargs)`

#### Running a command

Get metadata from EDP:
```shell
./edp get <subaccountID>
```

Register subaccount in EDP:
```shell
./edp register <subaccount ID> <platform region> <plan>
```
for example:
```shell
./edp register 41ba3cf2-041d-4223-adfe-c5de3458acbe cf-us21 standard
```

Deregister subaccount from EDP:
```shell
./edp deregister 41ba3cf2-041d-4223-adfe-c5de3458acbe
```