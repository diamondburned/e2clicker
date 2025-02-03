package dosage

import (
	"e2clicker.app/internal/publicerrors"
	"libdb.so/xcsv"
)

func init() {
	publicerrors.MarkTypePublic[*xcsv.RecordUnmarshalingError]()
}
