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
	UserData    string `json:"userData"`
	DocData     string `json:"docData"`
	Amount      string `json:"amount"`
    SlotTime    string `json:"slotTime"`
    SlotDate    string `json:"slotDate"`
	Date        string `json:"date"`
	Cancelled   bool   `json:"cancelled"`
	Payment     bool   `json:"payment"`
	IsCompleted bool   `json:"isCompleted"`
	Action      string `json:"action"`
	Timestamp   string `json:"timestamp"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []PatientRecord{
		{ID: "3", UserID: "U234", DocID: "D567", SlotDate: "2025-04-07T09:00:00Z", UserData: "Alice Walker", DocData: "Dr. Kim", Amount: "20", Date: "2025-04-07", Cancelled: false, Payment: true, IsCompleted: true, Action: "Completed", Timestamp: "2025-04-02T08:00:00Z"},
		{ID: "4", UserID: "U345", DocID: "D678", SlotDate: "2025-04-08T11:30:00Z", UserData: "Bob Stone", DocData: "Dr. Lee", Amount: "15", Date: "2025-04-08", Cancelled: false, Payment: false, IsCompleted: false, Action: "Pending", Timestamp: "2025-04-03T10:10:00Z"},
		{ID: "5", UserID: "U456", DocID: "D789", SlotDate: "2025-04-09T14:00:00Z", UserData: "Carol Swift", DocData: "Dr. Chen", Amount: "18", Date: "2025-04-09", Cancelled: true, Payment: false, IsCompleted: false, Action: "Cancelled", Timestamp: "2025-04-04T12:20:00Z"},
		{ID: "6", UserID: "U567", DocID: "D890", SlotDate: "2025-04-10T10:45:00Z", UserData: "Dan Grey", DocData: "Dr. Kumar", Amount: "22", Date: "2025-04-10", Cancelled: false, Payment: true, IsCompleted: false, Action: "Booked", Timestamp: "2025-04-05T11:15:00Z"},
		{ID: "7", UserID: "U678", DocID: "D901", SlotDate: "2025-04-11T16:00:00Z", UserData: "Ella Frost", DocData: "Dr. Ali", Amount: "25", Date: "2025-04-11", Cancelled: false, Payment: true, IsCompleted: true, Action: "Completed", Timestamp: "2025-04-06T09:45:00Z"},
		{ID: "8", UserID: "U789", DocID: "D012", SlotDate: "2025-04-12T13:15:00Z", UserData: "Frank Ocean", DocData: "Dr. Patel", Amount: "19", Date: "2025-04-12", Cancelled: false, Payment: false, IsCompleted: false, Action: "Pending", Timestamp: "2025-04-07T14:30:00Z"},
		{ID: "9", UserID: "U890", DocID: "D123", SlotDate: "2025-04-13T08:30:00Z", UserData: "Grace Bell", DocData: "Dr. Moore", Amount: "16", Date: "2025-04-13", Cancelled: true, Payment: false, IsCompleted: false, Action: "Cancelled", Timestamp: "2025-04-08T07:20:00Z"},
		{ID: "10", UserID: "U901", DocID: "D234", SlotDate: "2025-04-14T15:00:00Z", UserData: "Henry Blaze", DocData: "Dr. Green", Amount: "23", Date: "2025-04-14", Cancelled: false, Payment: true, IsCompleted: true, Action: "Completed", Timestamp: "2025-04-09T13:50:00Z"},
		{ID: "11", UserID: "U012", DocID: "D345", SlotDate: "2025-04-15T12:30:00Z", UserData: "Isla Ray", DocData: "Dr. White", Amount: "21", Date: "2025-04-15", Cancelled: false, Payment: true, IsCompleted: false, Action: "Booked", Timestamp: "2025-04-10T12:00:00Z"},
		{ID: "12", UserID: "U123", DocID: "D456", SlotDate: "2025-04-16T09:15:00Z", UserData: "Jack Noir", DocData: "Dr. Black", Amount: "17", Date: "2025-04-16", Cancelled: true, Payment: false, IsCompleted: false, Action: "Cancelled", Timestamp: "2025-04-11T08:30:00Z"},
	}

	for _, asset := range assets {
		s.CreateRecord(ctx, asset.ID, asset.UserID, asset.DocID, asset.SlotDate, asset.SlotTime, asset.UserData, asset.DocData,
			asset.Amount, asset.Date,
			asset.Cancelled, asset.Payment, asset.IsCompleted, asset.Action)
	}

	return nil
}

func (s *SmartContract) CreateRecord(ctx contractapi.TransactionContextInterface,
	id, userID, docId, slotDate, slotTime, userData, docData, amount, date string,
	cancelled, payment, isCompleted bool,
	action string,
) error {
	record := PatientRecord{
		ID: id, UserID: userID, DocID: docId, SlotTime: slotTime, SlotDate: slotDate, UserData: userData, DocData: docData, Amount: amount, Date: date,
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
