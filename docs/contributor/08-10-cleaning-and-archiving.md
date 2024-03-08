# Cleaning

All data about deprovisioned instances is stored in the database. To keep the database clean and do not store any sensitive data, the Kyma Environment Broker (KEB) provides a cleanup mechanism.
This mechanism is run at the end of the deprovisioning process and removes all data about the deprovisioned instance from the database. It removes the instance from the database and all related data, such as the instance's operations and runtime states, which belongs to those operations.

# Archiving

The Kyma Environment Broker (KEB) provides an archiving mechanism to store the data about the deprovisioned instances. 
The archiving mechanism is run at the end of the deprovisioning process (but before cleaning) and stores some part of data about the deprovisioned instance in the archive database.
Such archived instance can be used for investigations using kcp cli.