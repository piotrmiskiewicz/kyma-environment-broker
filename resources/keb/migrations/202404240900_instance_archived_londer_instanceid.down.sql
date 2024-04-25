ALTER TABLE instances_archived
    ALTER COLUMN instance_id TYPE varchar(64) USING substring(provider, 1, 64);
