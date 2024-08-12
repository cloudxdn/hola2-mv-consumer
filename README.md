# Hola2 MV Data Consumer

HOLA 2.0 IPEMS의 MV(성능) 데이터를 Kafka로부터 전달 받아 Elasticsearch에 입력합니다.



## 폴더 구조

```
├── README.md
├── common
│   ├── common.go
│   └── logger.go
├── consumer_group_example
├── fm-event
│   ├── fm_event.go
│   └── model.go
├── go.mod
├── go.sum
├── main.go
├── mv-interface
│   ├── model.go
│   └── mv_interface.go
└── mv-node
    ├── model.go
    └── mv_node.go
```

