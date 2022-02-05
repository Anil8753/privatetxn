package privatetxn

import (
	"fmt"
)

func getImplicitPrivateCollection(mspId string) string {
	return fmt.Sprintf("_implicit_org_%s", mspId)
}
