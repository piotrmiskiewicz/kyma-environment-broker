# SQL Statements to Verify the Archiver Job's Work and Check DB for Deprovisioned Instances

## DB Statistics to Be Run Before Starting the Archiver Job

1. Number of all operations:
```sql
select count(*) from operations;
```

2. Number of operations that the archiver job will delete (belongs to deprovisioned instances):
```sql
select count(*) from operations where instance_id not in (select instance_id from instances);
```

3. Number of runtime states:
```sql
select count(*) from runtime_states;
```

## Statements to Verify the Archiver Job's Work
1. Number of all operations:
```sql
select count(*) from operations;
```

2. Number of operations that the archiver job will delete (belonging to deprovisioned instances):
```sql
select count(*) from operations where instance_id not in (select instance_id from instances);
```

3. Number of runtime states:
```sql
select count(*) from runtime_states;
```

4. Number of archived instances:
```sql
select count(*) from instances_archived;
```