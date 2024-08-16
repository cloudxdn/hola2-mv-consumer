package mvnode

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hola2-mv-consumer/common"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/elastic/go-elasticsearch/esapi"
	elasticsearch "github.com/elastic/go-elasticsearch/v8"
)

func ProcessTopicMvNode(consumer sarama.Consumer, es *elasticsearch.Client, topic string) {
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error subscribing to Kafka topic %s: %s", topic, err)
	}
	defer partitionConsumer.Close()

	fmt.Println("Listening for messages on Kafka topic:", topic)

	messageBuffer := make(chan MvNodeMsg, 1000)

	go processBuffer(messageBuffer, es, topic)

	for message := range partitionConsumer.Messages() {
		if strings.TrimSpace(string(message.Value)) == "completed" {
			continue
		}

		mvNodeMsg, err := parseMVNodeMessage(string(message.Value))
		if err != nil {
			log.Printf("Error parsing message: %s", err)
			continue
		}

		messageBuffer <- mvNodeMsg
	}
}

func parseMVNodeMessage(message string) (MvNodeMsg, error) {
	fields := strings.Split(message, "\x1e")
	if len(fields) < 29 { // Adjust this number based on the expected number of fields
		return MvNodeMsg{}, fmt.Errorf("invalid number of fields in message")
	}

	return MvNodeMsg{
		Ctime:               fields[0],
		EquipID:             common.ParseInt(fields[1]),
		EquipAddr:           fields[2],
		SysName:             fields[3],
		Rterrrate:           common.ParseFloat(fields[4]),
		Bufferrate:          common.ParseFloat(fields[5]),
		Cpuutil:             common.ParseFloat(fields[6]),
		Memutil:             common.ParseFloat(fields[7]),
		Usedmem:             common.ParseInt64(fields[8]),
		Totalmem:            common.ParseInt64(fields[9]),
		Icmpindestunreachs:  common.ParseFloat(fields[10]),
		Icmpinerrors:        common.ParseFloat(fields[11]),
		Icmpinmsgs:          common.ParseFloat(fields[12]),
		Icmpintimeexcds:     common.ParseFloat(fields[13]),
		Icmpoutdestunreachs: common.ParseFloat(fields[14]),
		Icmpouterrors:       common.ParseFloat(fields[15]),
		Icmpoutmsgs:         common.ParseFloat(fields[16]),
		Icmpouttimeexcds:    common.ParseFloat(fields[17]),
		Ipforwdatagrams:     common.ParseFloat(fields[18]),
		Ipinaddrerrors:      common.ParseFloat(fields[19]),
		Ipindelivers:        common.ParseFloat(fields[20]),
		Ipindiscards:        common.ParseFloat(fields[21]),
		Iphdrerrors:         common.ParseFloat(fields[22]),
		Ipinreceives:        common.ParseFloat(fields[23]),
		Ipunknownprotos:     common.ParseFloat(fields[24]),
		Ipoutdiscards:       common.ParseFloat(fields[25]),
		Ipoutnoroutes:       common.ParseFloat(fields[26]),
		Ipoutrequests:       common.ParseFloat(fields[27]),
		Iproutingdiscards:   common.ParseFloat(fields[28]),
	}, nil
}

func processBuffer(messageBuffer <-chan MvNodeMsg, es *elasticsearch.Client, topic string) {
	var buffer []MvNodeMsg
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

func bulkIndexMessages(messages []MvNodeMsg, es *elasticsearch.Client, topic string) {
	var buf bytes.Buffer

	for _, msg := range messages {
		msg.Timestamp = time.Now().UTC().Format(time.RFC3339)
		meta := []byte(fmt.Sprintf(`{ "create" : { "_id": "%s", "_index" : "%s" } }%s`, strconv.Itoa(msg.EquipID)+string(msg.Ctime), strings.ToLower(topic), "\n"))
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
		log.Printf("Successfully indexed %s batch of %d messages", topic, len(messages))
	}
}
