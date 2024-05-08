# Running Archiver scenario

## Gather DB Statistics

Run all the SQL statements defined in the [DB statistics which should be run before the archiver is started](db_scripts.md#db-statistics-which-should-be-run-before-the-archiver-is-started) section.

## Run the Archiver Job with Disabled Deletion

Set the following environment variables:
- **APP_PERFORM_DELETION** set to `false`
- **APP_DRY_RUN**  set to `false`

Run the archiver (`./apply.sh`).

## Verify the Archiver Job's Work

Run all the SQL statements defined in the [Statements to verify the archiver work](db_scripts.md#statements-to-verify-the-archiver-work) section.

## Delete all instances_archived

Delete all rows from instances_archived table. It will be recreated once again. Make sure that the number of operations were not decreased (deletion was not performed)
```sql
delete from instances_archived;
```

## Enable Archiving and Deletion of Operations and Instances at the End of Deprovisioning

Set the following configurations for KEB:
```
archiving:
    enabled: true
    dryRun: false
cleaning:
    enabled: true
    dryRun: false
```

Wait for KEB to restart with the new configuration.

## Run the Archiver Job with Enabled Deletion

Set the following environment variables:
- **APP_PERFORM_DELETION** set to `true`.
- **APP_DRY_RUN**  set to `false`.

## Verify the Archiver Job's Work

Run all the SQL statements defined in the [Statements to verify the archiver work](db_scripts.md#statements-to-verify-the-archiver-work) section.