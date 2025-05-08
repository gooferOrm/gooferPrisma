# Goofer Binaries

There are two types of binaries needed by the Go client, one being the Goofer Engine binaries and the other being the
Goofer CLI binaries.

Goofer Engine binaries are fully managed, maintained and automatically updated by the Goofer team.

Goofer CLI binaries are managed by the maintainers of the Go client.

**NOTE: This is just for documentation purposes. The Goofer CLI is automatically published.**

--------

## How to build Goofer CLI binaries

Goofer CLI binaries are automatically published to S3 by a GitHub action. You can follow the instructions below to build
these binaries yourself.

### Prerequisites

Requires NodeJS.

Install the [AWS CLI](https://aws.amazon.com/cli/) and authenticate.

### Publish the latest Goofer version

```shell script
sh publish-latest.sh
```

### Publish a specific Goofer version

```shell script
sh publish.sh <version>
# e.g.
sh publish.sh 3.0.0
```

You can check the available versions on the [Goofer releases page](https://github.com/tacherasasi/goofer/releases).

**NOTE**:

### Goofer Team

Any Goofer team member can authenticate with the Goofer Go client account.

If you want to set up Goofer CLI binaries yourself, authenticate with your own AWS account and adapt the bucket name
in `publish.sh`.
When using the client, you will need to override the URL with env vars whenever you run the Go client, specifically
`GOOFER_CLI_URL` and `GOOFER_ENGINE_URL`. You can see the shape of these values
in [binaries/version.go](https://github.com/gooferOrm/goofer/blob/main/binaries/version.go).

This will also print the query engine version which you will need in the next step.

### Bump the binaries in the Go client

Go to `binaries/version.go` and adapt
the version values.
Push to a new branch, create a PR, and merge if tests are green.

When internal breaking changes happen, adaptions may be needed.
