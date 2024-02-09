package postsql

import (
    "github.com/gocraft/dbr"
    "github.com/kyma-project/kyma-environment-broker/internal/storage/dberr"
)

func (ws writeSession) DeleteRuntimeStateByOperationID(operationID string) (int64, error) {
    res, err := ws.deleteFrom(RuntimeStateTableName).
        Where(dbr.Eq("operation_id", operationID)).
        Exec()

    if err != nil {
        return 0, dberr.Internal("Failed to delete record from RuntimeState table: %s", err)
    }
    return res.RowsAffected()
}