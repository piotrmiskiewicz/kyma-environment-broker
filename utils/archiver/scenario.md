# Running Archiver scenario

## Gather DB statistics

Run all SQL statements defined in [DB statistics which should be run before the archiver is started](db_scripts.md#db-statistics-which-should-be-run-before-the-archiver-is-started) section.

## Run the archiver with disabled deletion

Set the following environment variables:
1. `APP_PERFORM_DELETION` set to `false`.
2. `APP_DRY_RUN`  set to `false`.

Run the archiver (`./apply.sh`).

## Verify the archiver work

Run all SQL statements defined in [Statements to verify the archiver work](db_scripts.md#statements-to-verify-the-archiver-work) section.

## Delete all instances_archived

Delete all rows from instances_archived table. It will be recreated once again.
```sql
delete from instances_archived;
```

## Enable archiving and deletion of operations and instances at the and of deprovisioning.

Set the following configurations for KEB:
```
archiving:
    enabled: true
    dryRun: false
cleaning:
    enabled: true
    dryRun: false
```

and wait for KEB restart with new configuration.

## Run the archiver with enabled deletion

Set the following environment variables:
1. `APP_PERFORM_DELETION` set to `true`.
2. `APP_DRY_RUN`  set to `false`.

## Verify the archiver work

Run all SQL statements defined in [Statements to verify the archiver work](db_scripts.md#statements-to-verify-the-archiver-work) section.