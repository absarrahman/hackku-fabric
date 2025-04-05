package chaincode

import (
    "encoding/json"
    "github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
    "time"
)

type SmartContract struct {
    contractapi.Contract
}

type PatientRecord struct {
    ID        string `json:"id"`
    Hash      string `json:"hash"`
    Metadata  string `json:"metadata"`
    Owner     string `json:"owner"`
    Timestamp string `json:"timestamp"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
    assets := []PatientRecord{
        {ID: "asset1", Hash: "blue", Metadata: "Meta1", Owner: "Tomoko", Timestamp: "1"},
        {ID: "asset2", Hash: "red", Metadata: "Meta2", Owner: "Brad", Timestamp: "2"},
        {ID: "asset3", Hash: "green", Metadata:  "Meta3", Owner: "Jin Soo", Timestamp: "3"},
        {ID: "asset4", Hash: "yellow", Metadata: "Meta4", Owner: "Max", Timestamp: "4"},
        {ID: "asset5", Hash: "black", Metadata: "Meta5", Owner: "Adriana", Timestamp: "5"},
        {ID: "asset6", Hash: "white", Metadata: "Meta6", Owner: "Michel", Timestamp: "6"},
    }

    for _, asset := range assets {
        s.CreateRecord(ctx, asset.ID, asset.Hash, asset.Metadata, asset.Owner)
    }

    return nil
}

func (s *SmartContract) CreateRecord(ctx contractapi.TransactionContextInterface, id, hash, metadata, owner string) error {
    record := PatientRecord{
        ID: id, Hash: hash, Metadata: metadata, Owner: owner,
        Timestamp: time.Now().Format(time.RFC3339),
    }
    data, _ := json.Marshal(record)
    return ctx.GetStub().PutState(id, data)
}

func (s *SmartContract) ReadRecord(ctx contractapi.TransactionContextInterface, id string) (*PatientRecord, error) {
    data, err := ctx.GetStub().GetState(id)
    if err != nil || data == nil {
        return nil, err
    }
    var record PatientRecord
    _ = json.Unmarshal(data, &record)
    return &record, nil
}

func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*PatientRecord, error) {
    // range query with empty string for startKey and endKey does an
    // open-ended query of all assets in the chaincode namespace.
    resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
    if err != nil {
        return nil, err
    }
    defer resultsIterator.Close()

    var assets []*PatientRecord
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return nil, err
        }

        var asset PatientRecord
        err = json.Unmarshal(queryResponse.Value, &asset)
        if err != nil {
            return nil, err
        }
        assets = append(assets, &asset)
    }

    return assets, nil
}
