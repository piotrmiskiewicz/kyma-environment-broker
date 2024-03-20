# Cleaning and Archiving

## Cleaning 

All data about deprovisioned instances is stored in the database. To keep the database clean and not store any sensitive data, Kyma Environment Broker (KEB) provides a cleanup mechanism.
This mechanism is run at the end of the deprovisioning process and removes all data about a deprovisioned instance from the database. It removes the instance from the database and all related data, such as the instance's operations and runtime states, which belong to those operations.

## Archiving

KEB provides an archiving mechanism to store the data about deprovisioned instances. 
The archiving mechanism is run at the end of the deprovisioning process (but before cleaning) and stores some part of data about the deprovisioned instance in the archive database.
Such archived instances can be used for investigations using KCP CLI.