package mvinterface

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

func ProcessTopicMvInterface(consumer sarama.Consumer, es *elasticsearch.Client, topic string) {
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error subscribing to Kafka topic %s: %s", topic, err)
	}
	defer partitionConsumer.Close()

	fmt.Println("Listening for messages on Kafka topic:", topic)

	messageBuffer := make(chan MvInterfaceMsg, 1000)

	go processBuffer(messageBuffer, es, topic)

	for message := range partitionConsumer.Messages() {
		if strings.TrimSpace(string(message.Value)) == "completed" {
			continue
		}

		mvInterfaceMsg, err := parseMvInterfaceMessage(string(message.Value))
		if err != nil {
			log.Printf("Error parsing message: %s", err)
			continue
		}

		messageBuffer <- mvInterfaceMsg
	}
}

func parseMvInterfaceMessage(message string) (MvInterfaceMsg, error) {
	fields := strings.Split(message, "\x1e")
	if len(fields) < 42 {
		return MvInterfaceMsg{}, fmt.Errorf("invalid number of fields in message")
	}

	return MvInterfaceMsg{
		Ctime:             fields[0],
		EquipID:           common.ParseInt(fields[1]),
		EquipAddr:         fields[2],
		SysName:           fields[3],
		Pkey:              common.ParseInt(fields[4]),
		OctetsIn:          common.ParseInt(fields[5]),
		OctetsOut:         common.ParseInt(fields[6]),
		UtilIn:            common.ParseFloat(fields[7]),
		UtilOut:           common.ParseFloat(fields[8]),
		UpktsIn:           common.ParseInt(fields[9]),
		UpktsOut:          common.ParseInt(fields[10]),
		NupktsIn:          common.ParseInt(fields[11]),
		NupktsOut:         common.ParseInt(fields[12]),
		OctetsBpsIn:       common.ParseInt(fields[13]),
		OctetsBpsOut:      common.ParseInt(fields[14]),
		OctetsPpsIn:       common.ParseInt(fields[15]),
		OctetsPpsOut:      common.ParseInt(fields[16]),
		ErrorsIn:          common.ParseInt(fields[17]),
		ErrorsOut:         common.ParseInt(fields[18]),
		ErrorIn:           common.ParseFloat(fields[19]),
		ErrorOut:          common.ParseFloat(fields[20]),
		DiscardsIn:        common.ParseInt(fields[21]),
		DiscardsOut:       common.ParseInt(fields[22]),
		DiscardIn:         common.ParseFloat(fields[23]),
		DiscardOut:        common.ParseFloat(fields[24]),
		Crc:               common.ParseInt(fields[25]),
		Collision:         common.ParseInt(fields[26]),
		IfUnknownProtosIn: common.ParseInt(fields[27]),
		McastPktsIn:       common.ParseInt(fields[28]),
		McastPktsOut:      common.ParseInt(fields[29]),
		QdropsIn:          common.ParseInt(fields[30]),
		QdropsOut:         common.ParseInt(fields[31]),
		RxPower:           common.ParseFloat(fields[32]),
		TxPower:           common.ParseFloat(fields[33]),
		RxLane1:           common.ParseFloat(fields[34]),
		TxLane1:           common.ParseFloat(fields[35]),
		RxLane2:           common.ParseFloat(fields[36]),
		TxLane2:           common.ParseFloat(fields[37]),
		RxLane3:           common.ParseFloat(fields[38]),
		TxLane3:           common.ParseFloat(fields[39]),
		RxLane4:           common.ParseFloat(fields[40]),
		TxLane4:           common.ParseFloat(fields[41]),
	}, nil
}

func processBuffer(messageBuffer <-chan MvInterfaceMsg, es *elasticsearch.Client, topic string) {
	var buffer []MvInterfaceMsg
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

func bulkIndexMessages(messages []MvInterfaceMsg, es *elasticsearch.Client, topic string) {
	var buf bytes.Buffer

	for _, msg := range messages {
		msg.Timestamp = time.Now().UTC().Format(time.RFC3339)
		meta := []byte(fmt.Sprintf(`{ "create" : { "_id": "%s", "_index" : "%s" } }%s`, string(rune(msg.EquipID))+string(msg.Ctime), strings.ToLower(topic), "\n"))
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
