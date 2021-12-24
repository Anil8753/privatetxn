package privatetxn

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type PrivateHistoryResult struct {
	Timestamp int64
	TxnId     string
	Value     string
}

// GetHistoryForKey gets the history of the key from the public ledger and retrives the private data records from the implicit private the data collection.
// It validates the private the data records' integrity by comparing the hash of the private data with the hashes available on the ledger.
func GetHistoryForKey(ctx contractapi.TransactionContextInterface, orgMsp string, key string) ([]PrivateHistoryResult, error) {

	itr, err := ctx.GetStub().GetHistoryForKey(key)
	if err != nil {
		return nil, err
	}

	var entries []PrivateHistoryResult

	for itr.HasNext() {
		r, err := itr.Next()
		if err != nil {
			return nil, err
		}

		if r.GetIsDelete() {
			entry := PrivateHistoryResult{Timestamp: r.Timestamp.GetSeconds(), TxnId: r.TxId, Value: ""}
			entries = append(entries, entry)
			continue
		}

		privateStateKey := string(r.GetValue())
		dataBytes, err := getPrivateData(ctx, orgsMsp, privateStateKey)
		if err != nil {
			return nil, err
		}

		entry := PrivateHistoryResult{Timestamp: r.Timestamp.GetSeconds(), TxnId: r.TxId, Value: string(dataBytes)}
		entries = append(entries, entry)
	}

	return entries, nil
}
