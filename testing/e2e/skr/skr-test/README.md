# SKR test

This test covers creation and updates Kyma runtimes using existing KCP infrastructure.

## Usage modes

You can use the SKR test in two modes - with and without provisioning.

### With provisioning

In this mode, the test executes the following steps:

1. Provision an SKR cluster.
2. Run the OIDC test.
3. Deprovision the SKR instance and clean up the resources.

### Without Provisioning.

In this mode the test additionally needs the following environment variables:
- `SKIP_PROVISIONING`, set to `true`
- `INSTANCE_ID` the uuid of the provisioned SKR instance

In this mode, the test executes the following steps:
1. Ensure the SKR exists.
2. Run the OIDC test.
3. Clean up the resources.

## Manually test execution

1. Before you run the test, prepare the `.env` file based on the following `.env.template`

2. To set up the environment variables in your system, run:

```bash
export $(xargs < .env)
```

3. Choose whether you want to run the test with or without provisioning.
    - To run the test **with** provisioning, call the following target:

    ```bash
    npm run skr-test
    #or
    make skr-test
    ```
    - To run the SKR test **without** provisioning, use the following command:

    ```bash
    make skr-test SKIP_PROVISIONING=true
    #or
    npm run skr-test # when all env vars are exported
    ```
