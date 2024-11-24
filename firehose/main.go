package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/firehose"
)

func main() {
	// Initialize the AWS session
	sess := session.Must(session.NewSession())
	svc := firehose.New(sess, aws.NewConfig().WithRegion("ap-northeast-1"))

	// Define the Firehose stream name
	streamName := "put-record-test"

	// Sample data to be sent to Firehose
	sampleData := []map[string]interface{}{
		{"id": 1, "name": "Alice", "age": 30},
		{"id": 2, "name": "Bob", "age": 25},
		{"id": 3, "name": "Charlie", "age": 35},
	}

	// Convert the sample data to Firehose records
	var records []*firehose.Record
	for _, record := range sampleData {
		data, err := json.Marshal(record)
		if err != nil {
			log.Fatalf("Failed to marshal record: %v", err)
		}
		records = append(records, &firehose.Record{Data: append(data, '\n')})
	}

	// Use the PutRecordBatch method to send the data to Firehose
	input := &firehose.PutRecordBatchInput{
		DeliveryStreamName: aws.String(streamName),
		Records:            records,
	}

	// resolveされているendpointを表示
	fmt.Println("endpoint: ", svc.Endpoint)

	result, err := svc.PutRecordBatch(input)
	if err != nil {
		log.Fatalf("Failed to put record batch: %v", err)
	}

	// // Print the response from Firehose
	fmt.Printf("Response: %+v\n", result)
}
