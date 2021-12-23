package privatetxn

import (
	"bytes"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func putPDC(ctx contractapi.TransactionContextInterface, org string, key string, data []byte) error {

	if err := ctx.GetStub().PutPrivateData(implicitCollection(org), key, data); err != nil {
		return fmt.Errorf("func: putPDC, agrs: [ org: %s, key: %s ]. %w", org, key, err)
	}

	return nil
}

func implicitCollection(msp string) string {
	return "_implicit_org_" + msp
}

func getPrivateCollectionKey(ctx contractapi.TransactionContextInterface, dataHash []byte) (string, error) {

	ts, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x_%d", dataHash, ts.GetSeconds()), nil
}

func getPrivateData(ctx contractapi.TransactionContextInterface, orgsMsp string, pvtCollectionKey string) ([]byte, error) {

	col := implicitCollection(orgsMsp)

	data, err := ctx.GetStub().GetPrivateData(col, pvtCollectionKey)
	if err != nil {
		return nil, fmt.Errorf("GetPrivateData is failed. key: %s, collection: %s. %w", pvtCollectionKey, col, err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("private data not available. key: %s, collection: %s", pvtCollectionKey, col)
	}

	// private collection data hash
	dataHash := getHash(data)

	// retrive ledger hash
	ledgerHash, err := ctx.GetStub().GetPrivateDataHash(col, pvtCollectionKey)
	if err != nil {
		return nil, fmt.Errorf("GetPrivateDataHash failed. key: %s, collection: %s. %w", pvtCollectionKey, col, err)
	}

	if len(ledgerHash) == 0 {
		return nil, fmt.Errorf("ledger hash is empty. key: %s, collection: %s", pvtCollectionKey, col)
	}

	// validate private data integrity
	if !bytes.Equal(ledgerHash, dataHash) {

		errMsg := fmt.Sprintf("data has been be tempered in private collection. key: %s, collection: %s", pvtCollectionKey, col)

		return nil, fmt.Errorf(
			"error: %s \n hash %x for passed immutable properties %s does not match on-chain hash %x",
			errMsg,
			dataHash,
			string(data),
			ledgerHash,
		)
	}

	return data, nil
}
