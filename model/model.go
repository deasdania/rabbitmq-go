package model

type BatchUploadCSV struct {
	ID string
}

type LivetestBacktestMessageQueue struct {
	ID               string
	BatchUploadCSVID string
}

type BodyPublishTest struct {
	HandlerNameKey string
	ProductName    string
	ProductID      string
	MsgQue         LivetestBacktestMessageQueue
}
