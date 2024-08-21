package mvinterface

type MvInterfaceMsg struct {
	Ctime             string  `json:"ctime"`
	EquipID           int64   `json:"equip_id"`
	EquipAddr         string  `json:"equip_addr"`
	SysName           string  `json:"sys_name"`
	Pkey              int64   `json:"pkey"`
	OctetsIn          int64   `json:"octets_in"`
	OctetsOut         int64   `json:"octets_out"`
	UtilIn            float64 `json:"util_in"`
	UtilOut           float64 `json:"util_out"`
	UpktsIn           int64   `json:"upkts_in"`
	UpktsOut          int64   `json:"upkts_out"`
	NupktsIn          int64   `json:"nupkts_in"`
	NupktsOut         int64   `json:"nupkts_out"`
	OctetsBpsIn       int64   `json:"octetsbps_in"`
	OctetsBpsOut      int64   `json:"octetsbps_out"`
	OctetsPpsIn       int64   `json:"octetspps_in"`
	OctetsPpsOut      int64   `json:"octetspps_out"`
	ErrorsIn          int64   `json:"errors_in"`
	ErrorsOut         int64   `json:"errors_out"`
	ErrorIn           float64 `json:"error_in"`
	ErrorOut          float64 `json:"error_out"`
	DiscardsIn        int64   `json:"discards_in"`
	DiscardsOut       int64   `json:"discards_out"`
	DiscardIn         float64 `json:"discard_in"`
	DiscardOut        float64 `json:"discard_out"`
	Crc               int64   `json:"crc"`
	Collision         int64   `json:"collision"`
	IfUnknownProtosIn int64   `json:"ifunknownprotos_in"`
	McastPktsIn       int64   `json:"mcastpkts_in"`
	McastPktsOut      int64   `json:"mcastpkts_out"`
	QdropsIn          int64   `json:"qdrops_in"`
	QdropsOut         int64   `json:"qdrops_out"`
	RxPower           float64 `json:"rxpower"`
	TxPower           float64 `json:"txpower"`
	RxLane1           float64 `json:"rx_lane1"`
	TxLane1           float64 `json:"tx_lane1"`
	RxLane2           float64 `json:"rx_lane2"`
	TxLane2           float64 `json:"tx_lane2"`
	RxLane3           float64 `json:"rx_lane3"`
	TxLane3           float64 `json:"tx_lane3"`
	RxLane4           float64 `json:"rx_lane4"`
	TxLane4           float64 `json:"tx_lane4"`
	Timestamp         string  `json:"@timestamp"`
}
