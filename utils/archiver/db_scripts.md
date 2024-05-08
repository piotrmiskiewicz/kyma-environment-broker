# SQL statements to verify the archiver work and check DB about deprovisioned instances

## DB statistics which should be run before the archiver is started

1. Number of all operations:
```sql
select count(*) from operations;
```

2. Number of operations, which will be deleted by the archiver (belongs to deprovisioned instances):
```sql
select count(*) from operations where instance_id not in (select instance_id from instances);
```

3. Number of runtime states:
```sql
select count(*) from runtime_states;
```

## Statements to verify the archiver work
1. Number of all operations:
```sql
select count(*) from operations;
```

2. Number of operations, which will be deleted by the archiver (belongs to deprovisioned instances):
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