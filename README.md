# Wavelength

An opinionated build tool for developing and deploying serverless functions.

Currently only supports AWS Lambda functions using a Node.js runtime, deployed using Terraform.

## Pre-requisites
You must have the following installed and on your `$PATH`:
- [Lerna](https://github.com/lerna/lerna)
- [Terraform](https://www.terraform.io/)
- Go (I'm using version 1.15.5. I haven't tested any other versions)

## Installation
Run `go get github.com/lewis-od/wavelength@v1.0.0`

## How it works

Wavelength assumes the application code for all your lambda functions, as well as the Terraform code to deploy them, all
lives in the same repo.

Your directory structure should look something like:
```text
.
├── lambdas
│   ├── lambda-one
│   │   ├── dist
│   │   ├── package.json
│   │   └── src
│   └── lambda-two
│       ├── dist
│       ├── package.json
│       └── src
├── lerna.json
├── package.json
├── common
│   ├── domain-model
│   │   ├── package.json
│   │   └── src
│   └── logging
│       ├── package.json
│       └── src
├── terraform
│   ├── deployments
│   │   ├── app
│   │   └── artifact-storage
│   └── modules
│       ├── lambda-function
│       └── rest-endpoint
├── yarn.lock
└── .wavelength.yml
```

There's a few things to unpack here, so lets break it down:

### Function code

You must use [lerna] to build your code, with each function being contained in it's own package. Packages names should
follow the convention `@my-app/lambda-one`.

Each package's build script should create an artifact at `dist/<lambda-name>.zip`. A handy utility for doing this
is [lerna-to-lambda].

[lerna]: https://github.com/lerna/lerna
[lerna-to-lambda]: https://github.com/lafiosca/lerna-to-lambda

### Terraform

You should have 2 separate Terraform deployments; one that creates an s3 bucket in which to store your code, and another
that creates the actual lambda functions, along with any other required infrastructure (i.e. an API Gateway, IAM roles,
etc).

The module that creates the bucket must export the bucket name as an output. Wavelength will read this to figure out
where to upload your code.

The names of the lambda functions should match the name of the directory/node package they live in, and be prefixed with
your app name, e.g. `my-app-lambda-one`.

### Config

Wavelength requires a `.wavelength.yml` or `.wavelength.json` config file in the root of your repo. It supports the
following properties:

| Name                            | Description                                                                                     | Required? |
| ------------------------------- | ----------------------------------------------------------------------------------------------- | --------- |
| `artifactStorage.outputName`    | Name of the Terraform output used to export the artifact bucket name                            | [x]       |
| `artifactStorage.terraformDir`  | Path to the directory containing the Terraform to deploy your artifact storage bucket           | [x]       |
| `lambdas`                       | Name of the directory containing your lambda packages. Defaults to `"lambdas"` if not provided. | [ ]       |
| `projectName`                   | Name of your project/app. Assumed to be the package scope of all your lambdas                   | [x]       |

For the example shown above, this would look like:
```yaml
projectName: my-app
artifactStorage:
  terraformDir: terraform/deployments/artifact-storage
  outputName: bucket_name
```

## Dev Workflow
The intial deployment of your app will look something like

```shell
cd terraform/deployments/artifact-storage
terraform apply # Create the artifact bucket
cd ../../..
wavelength upload dev # Build and upload code to S3
cd terraform/deployments/app
terraform apply # Create lambdas
cd ../../..
```

After this, if you make some changes to your code and want to update the deployed lambdas:
```shell
# Build and upload edited code to S3
wavelength upload dev
# Update lambdas with new code
wavelength update dev
```
The last command will call the lambda [UpdateFunctionCode] endpoint, pointing the lambda at the code that was just
uploaded.

[UpdateFunctionCode]: https://docs.aws.amazon.com/lambda/latest/dg/API_UpdateFunctionCode.html

# CI/CD Workflow
Wavelength is intended to be used as part of a CI/CD pipeline as:
```shell
wavelength upload <commit-sha>
```
Your deployment step will then probably look something like:
```shell
terraform apply -ver version="<commit-sha>" -auto-approve
```
Where the `version` variable is used to set the `s3_key` for the lambda.

Running `wavelength update` is not required in this case as the commit SHA will be different each time, so Terraform
will detect a change to the `s3_key` of the function, and handle updating the lambda for us.

