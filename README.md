# privatetxn

Provide functions for implicit data collection similar to the regular stub functions for public ledger. <br>
It usage a secret key to encrypt the data this secret key is passed from the org invoking the transaction and
this secret key is stored in the implicit private data collection of the organizations participating in the transaction.
Later this secret is used to decrypt the encrypted data in various legder functions

## How to use

```
go get github.com/Anil8753/privatetxn
```

## Exported PrivateData struct

```
type PrivateData struct {
	// Secret is used to encrypt the data stored on ledger
	Secret string
	// Tags are used for couch db rich queries
	Tags interface{}
	// Endorsing Orgs, private data will be shared between the endorsing orgs
	MSPOrgs []string
}
```

## Exported functions

```
func PutState(
	ctx contractapi.TransactionContextInterface,
	key string,
	value interface{},
	pvtData *PrivateData,
) error
```

PutState saves the data secretly on the ledger.
How it works:

1.  Encrypts the data using the secret key passed via PrivateData struct
2.  Stores the key and encrypted data on public ledger
3.  Stores the private data (contains secret key) to the implicit private data collection of participating orgs

    <br>Note: Transaction must be endorsed by all the orgs passed in the `PutState` function. It is requied to write in the implicit data collection of the participating organizations.

<br>

```
func GetState(
	ctx contractapi.TransactionContextInterface,
	key string,
) ([]byte, *PrivateData, error)
```

GetState returns

1. Private data stored in the implicit data collection
2. Public data which is decypted using the secret key before return.
   This secret key is available in the implicit private data collection of the org

<br>

```
func GetHistoryForKey(
	ctx contractapi.TransactionContextInterface,
	key string,
) ([][]byte, error)
```

GetHistoryForKey returns the history of a key
It gets the encrypted history of a key and decrypts the history
using the secret key stored in the implicit private data colelction of the org

```
func GetQueryResult(
	ctx contractapi.TransactionContextInterface,
	query string,
) ([][]byte, error)
```

GetQueryResult makes the query in private data collection and
retieves the results from the encrypted data stored on public ledger
Before returning it decryts the data. Secret key is stored in the
implicit private data collection
