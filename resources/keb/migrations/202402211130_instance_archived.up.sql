BEGIN;

CREATE TABLE IF NOT EXISTS instances_archived (
    instance_id                      varchar(255) NOT NULL PRIMARY KEY,
    global_account_id                varchar(255) NOT NULL,
    last_runtime_id                  varchar(255),
    subscription_global_account_id   varchar(255),
    subaccount_id                    varchar(255) NOT NULL,
    plan_id                          varchar(255) NOT NULL,
    plan_name                        varchar(255),
    region                           varchar(255),
    subaccount_region                varchar(255),
    provider                         varchar(255),
    shoot_name                       varchar(255),
    internal_user                    boolean NOT NULL,


    provisioning_started_at           timestamp with time zone NOT NULL,
    provisioning_finished_at          timestamp with time zone NOT NULL,
    first_deprovisioning_started_at   timestamp with time zone NOT NULL,
    first_deprovisioning_finished_at  timestamp with time zone NOT NULL,
    last_deprovisioning_finished_at   timestamp with time zone NOT NULL

);


COMMIT;