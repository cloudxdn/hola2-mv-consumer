package fmevent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hola2-mv-consumer/common"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/elastic/go-elasticsearch/esapi"
	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

func ProcessTopicFmEvent(consumer sarama.Consumer, es *elasticsearch.Client, topic string) {
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error subscribing to Kafka topic %s: %s", topic, err)
	}
	defer partitionConsumer.Close()

	fmt.Println("Listening for messages on Kafka topic:", topic)

	messageBuffer := make(chan FmEventMsg, 1000)

	go processBuffer(messageBuffer, es, topic)

	for message := range partitionConsumer.Messages() {
		if strings.TrimSpace(string(message.Value)) == "completed" {
			continue
		}

		fmEventMsg, err := parseFmEventMessage(string(message.Value))
		if err != nil {
			log.Printf("Error parsing message: %s", err)
			continue
		}

		messageBuffer <- fmEventMsg
	}
}

func parseFmEventMessage(message string) (FmEventMsg, error) {
	fields := strings.Split(message, "|^|")
	if len(fields) < 39 {
		return FmEventMsg{}, fmt.Errorf("invalid number of fields in message")
	}

	return FmEventMsg{
		Command:          fields[0],
		AlarmID:          uint(common.ParseInt(fields[1])),
		Status:           fields[2],
		Organization:     fields[3],
		EquipID:          fields[4],
		EquipName:        fields[5],
		Location:         fields[6],
		LocationName:     fields[7],
		LocationAlias:    fields[8],
		EventType:        fields[9],
		AlarmCode:        fields[10],
		AlarmCategory:    fields[11],
		Severity:         fields[12],
		User:             fields[13],
		OccurredTime:     fields[14],
		RecvTime:         fields[15],
		ClearTime:        fields[16],
		ClearUser:        fields[17],
		AckTime:          fields[18],
		AckUser:          fields[19],
		Times:            fields[20],
		UpdateTime:       fields[21],
		RecvUpdateTime:   fields[22],
		ServerKey:        fields[23],
		Message:          fields[24],
		Description:      fields[25],
		NetCD:            fields[26],
		NetTypeName:      fields[27],
		NetSubTypeName:   fields[28],
		EquipTypeName:    fields[29],
		EquipSubTypeName: fields[30],
		PKey:             fields[31],
		PortDescr:        fields[32],
		IfIP:             fields[33],
		RingName:         fields[34],
		UpperLink:        fields[35],
		RootLink:         fields[36],
		IsWorking:        fields[37],
		OtherContents:    fields[38],
	}, nil
}

func processBuffer(messageBuffer <-chan FmEventMsg, es *elasticsearch.Client, topic string) {
	var buffer []FmEventMsg
	var mu sync.Mutex
	ticker := time.NewTicker(common.BufferPeriod)

	for {
		select {
		case msg := <-messageBuffer:
			mu.Lock()
			buffer = append(buffer, msg)
			mu.Unlock()
		case <-ticker.C:
			mu.Lock()
			if len(buffer) > 0 {
				bulkIndexMessages(buffer, es, topic)
				buffer = nil // Clear the buffer
			}
			mu.Unlock()
		}
	}
}

func bulkIndexMessages(messages []FmEventMsg, es *elasticsearch.Client, topic string) {
	var buf bytes.Buffer

	for _, msg := range messages {
		msg.Timestamp = time.Now().UTC().Format(time.RFC3339)
		meta := []byte(fmt.Sprintf(`{ "create" : { "_id": "%s", "_index" : "%s" } }%s`, string(rune(msg.AlarmID))+string(msg.OccurredTime), strings.ToLower(topic), "\n"))
		data, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Error marshalling message: %s", err)
			continue
		}
		data = append(data, "\n"...)

		buf.Grow(len(meta) + len(data))
		buf.Write(meta)
		buf.Write(data)
	}

	req := esapi.BulkRequest{
		Body:    strings.NewReader(buf.String()),
		Refresh: "true",
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Printf("Failure indexing batch: %s", err)
	}
	defer res.Body.Close()

	var resBody map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		log.Printf("Error parsing the response body: %s", err)
	}

	if res.IsError() {
		log.Printf("Error indexing %s batch: %s", topic, res.String())
	} else {
		log.Printf("Successfully indexed batch %s of %d messages", topic, len(messages))
	}
}
