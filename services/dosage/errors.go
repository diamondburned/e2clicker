package dosage

import (
	"libdb.so/e2clicker/internal/publicerrors"
	"libdb.so/xcsv"
)

func init() {
	publicerrors.MarkTypePublic[*xcsv.RecordUnmarshalingError]()
}
