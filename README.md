# privatetxn
Provide functions for implicit data collection similar to the regular stub functions for public ledger. <br>

## How to use
```
go get github.com/Anil8753/privatetxn
```

## Exported functions

```
func PutState(ctx contractapi.TransactionContextInterface, orgsMsp []string, key string, data []byte) error 
```
PutState stores the private data in the orgs' implicit private collection and corresponding private data key on the public ledger.<br>
Note: Transaction must be endorsed by all the orgs passed in the `PutState` function. It is requied to write in the implicit data collection of the organizations.

```
func GetState(ctx contractapi.TransactionContextInterface, orgsMsp string, key string) ([]byte, error)
```
GetState fetchs the private data key from the public ledger and retrive the data from the private data collection of the org. 
At the same time it checks the integrity of the private data by validating the hash of the data with public hash on the ledger.


```
func GetHistoryForKey(ctx contractapi.TransactionContextInterface, orgsMsp string, key string) ([]PrivateHistoryResult, error)
```

GetHistoryForKey gets the history of the key from the public ledger and retrives the private data records from the implicit private the data collection.
It validates the private the data records' integrity by comparing the hash of the private data with the hashes available on the ledger.

