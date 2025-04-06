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
	ID          string `json:"id"`
	UserID      string `json:"userId"`
	DocID       string `json:"docId"`
	SlotDate    string `json:"slotDate"`
	UserData    string `json:"userData"`
	DocData     string `json:"docData"`
	Date        string `json:"date"`
	Cancelled   bool   `json:"cancelled"`
	Payment     bool   `json:"payment"`
	IsCompleted bool   `json:"isCompleted"`
	Action      string `json:"action"`
	Timestamp   string `json:"timestamp"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []PatientRecord{
		{ID: "2", UserID: "U123", DocID: "D456", SlotDate: "2025-04-06T10:00:00Z", UserData: "John Doe", DocData: "Dr. Smith", Date: "2025-04-06", Cancelled: false, Payment: true, IsCompleted: false, Action: "Booked", Timestamp: "2025-04-01T09:15:00Z"},
		{ID: "3", UserID: "U789", DocID: "D321", SlotDate: "2025-04-08T14:30:00Z", UserData: "Jane Roe", DocData: "Dr. Adams", Date: "2025-04-08", Cancelled: false, Payment: true, IsCompleted: true, Action: "Completed", Timestamp: "2025-04-02T11:25:00Z"},
		{ID: "4", UserID: "U111", DocID: "D654", SlotDate: "2025-04-10T09:00:00Z", UserData: "Mike Black", DocData: "Dr. Lin", Date: "2025-04-10", Cancelled: true, Payment: false, IsCompleted: false, Action: "Cancelled", Timestamp: "2025-04-03T08:50:00Z"},
		{ID: "5", UserID: "U222", DocID: "D987", SlotDate: "2025-04-07T13:00:00Z", UserData: "Anna White", DocData: "Dr. Green", Date: "2025-04-07", Cancelled: false, Payment: true, IsCompleted: true, Action: "Completed", Timestamp: "2025-04-02T14:20:00Z"},
		{ID: "6", UserID: "U333", DocID: "D147", SlotDate: "2025-04-09T16:00:00Z", UserData: "Chris Blue", DocData: "Dr. Young", Date: "2025-04-09", Cancelled: false, Payment: false, IsCompleted: false, Action: "Pending", Timestamp: "2025-04-04T10:10:00Z"},
		{ID: "7", UserID: "U444", DocID: "D258", SlotDate: "2025-04-11T11:00:00Z", UserData: "Laura Grey", DocData: "Dr. Brown", Date: "2025-04-11", Cancelled: true, Payment: true, IsCompleted: false, Action: "Cancelled", Timestamp: "2025-04-05T12:30:00Z"},
		{ID: "8", UserID: "U555", DocID: "D369", SlotDate: "2025-04-06T15:00:00Z", UserData: "David Gold", DocData: "Dr. King", Date: "2025-04-06", Cancelled: false, Payment: true, IsCompleted: true, Action: "Completed", Timestamp: "2025-04-01T13:45:00Z"},
		{ID: "9", UserID: "U666", DocID: "D741", SlotDate: "2025-04-07T10:30:00Z", UserData: "Emily Rose", DocData: "Dr. Carter", Date: "2025-04-07", Cancelled: false, Payment: true, IsCompleted: false, Action: "Booked", Timestamp: "2025-04-03T15:00:00Z"},
		{ID: "10", UserID: "U777", DocID: "D852", SlotDate: "2025-04-08T12:00:00Z", UserData: "Steve Cyan", DocData: "Dr. Lopez", Date: "2025-04-08", Cancelled: false, Payment: false, IsCompleted: false, Action: "Pending", Timestamp: "2025-04-02T17:20:00Z"},
		{ID: "11", UserID: "U888", DocID: "D963", SlotDate: "2025-04-12T09:30:00Z", UserData: "Rachel Amber", DocData: "Dr. Patel", Date: "2025-04-12", Cancelled: false, Payment: true, IsCompleted: true, Action: "Completed", Timestamp: "2025-04-04T18:45:00Z"},
	}

	for _, asset := range assets {
		s.CreateRecord(ctx, asset.ID, asset.UserID, asset.DocID, asset.SlotDate, asset.UserData, asset.DocData, asset.Date,
            asset.Cancelled, asset.Payment, asset.IsCompleted, asset.Action, asset.Timestamp)
	}

	return nil
}

func (s *SmartContract) CreateRecord(ctx contractapi.TransactionContextInterface,
	id, userID, docId, slotDate, userData, docData, date string,
    cancelled, payment, isCompleted bool,
    action, timeStamp string,
) error {
	record := PatientRecord{
		ID: id, UserID: userID, DocID: docId, SlotDate: slotDate, UserData: userData, DocData: docData, Date: date,
        Cancelled: cancelled, Payment: payment, IsCompleted: isCompleted,
		Action: action, Timestamp: time.Now().Format(time.RFC3339),
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
