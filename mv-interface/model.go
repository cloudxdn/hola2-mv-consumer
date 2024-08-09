package mvinterface

type MvInterfaceMsg struct {
	Ctime             string  `json:"ctime"`
	EquipID           int     `json:"equip_id"`
	EquipAddr         string  `json:"equip_addr"`
	SysName           string  `json:"sys_name"`
	Pkey              int     `json:"pkey"`
	OctetsIn          int     `json:"octets_in"`
	OctetsOut         int     `json:"octets_out"`
	UtilIn            float64 `json:"util_in"`
	UtilOut           float64 `json:"util_out"`
	UpktsIn           int     `json:"upkts_in"`
	UpktsOut          int     `json:"upkts_out"`
	NupktsIn          int     `json:"nupkts_in"`
	NupktsOut         int     `json:"nupkts_out"`
	OctetsBpsIn       int     `json:"octets_bps_in"`
	OctetsBpsOut      int     `json:"octets_bps_out"`
	OctetsPpsIn       int     `json:"octets_pps_in"`
	OctetsPpsOut      int     `json:"octets_pps_out"`
	ErrorsIn          int     `json:"errors_in"`
	ErrorsOut         int     `json:"errors_out"`
	ErrorIn           float64 `json:"error_in"`
	ErrorOut          float64 `json:"error_out"`
	DiscardsIn        int     `json:"discards_in"`
	DiscardsOut       int     `json:"discards_out"`
	DiscardIn         float64 `json:"discard_in"`
	DiscardOut        float64 `json:"discard_out"`
	Crc               int     `json:"crc"`
	Collision         int     `json:"collision"`
	IfUnknownProtosIn int     `json:"if_unknown_protos_in"`
	McastPktsIn       int     `json:"mcast_pkts_in"`
	McastPktsOut      int     `json:"mcast_pkts_out"`
	QdropsIn          int     `json:"qdrops_in"`
	QdropsOut         int     `json:"qdrops_out"`
	RxPower           float64 `json:"rx_power"`
	TxPower           float64 `json:"tx_power"`
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
