package schemas

import "fmt"

// MqResponseBody - struct to hold MQ Response Body
type MqResponseBody struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// QueueMessageDatum - a structure to hold 1 Message from the DB Dispatcher Queue
type QueueMessageDatum struct {
	Query        string `json:"query"`
	Host         string `json:"host"`
	DatabaseName string `json:"database_name"`
	Port         string `json:"port"`
	DryRun         bool `json:"dry_run"`
}

// ToString - A method to return the QueueMessageDatum as a string
func (o QueueMessageDatum) ToString() string {
	return fmt.Sprintf("DB: %s:%s/%s\nQuery: %s\nDry Run: %t", o.Host, o.Port, o.DatabaseName, o.Query, o.DryRun)
}

// QueueMessageData - a structure to hold 1 Message from the DB Dispatcher Queue
type QueueMessageData struct {
	Query        string `json:"query"`
	Host         string `json:"host"`
	DatabaseName string `json:"database_name"`
	Port         string `json:"port"`
}

