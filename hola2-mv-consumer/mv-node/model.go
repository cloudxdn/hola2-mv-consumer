package mvnode

type MvNodeMsg struct {
	Ctime               string  `json:"ctime"`
	EquipID             int     `json:"equip_id"`
	EquipAddr           string  `json:"equip_addr"`
	SysName             string  `json:"sys_name"`
	Rterrrate           float64 `json:"rterrrate"`
	Bufferrate          float64 `json:"bufferrate"`
	Cpuutil             float64 `json:"cpuutil"`
	Memutil             float64 `json:"memutil"`
	Usedmem             int64   `json:"usedmem"`
	Totalmem            int64   `json:"totalmem"`
	Icmpindestunreachs  float64 `json:"icmpindestunreachs"`
	Icmpinerrors        float64 `json:"icmpinerrors"`
	Icmpinmsgs          float64 `json:"icmpinmsgs"`
	Icmpintimeexcds     float64 `json:"icmpintimeexcds"`
	Icmpoutdestunreachs float64 `json:"icmpoutdestunreachs"`
	Icmpouterrors       float64 `json:"icmpouterrors"`
	Icmpoutmsgs         float64 `json:"icmpoutmsgs"`
	Icmpouttimeexcds    float64 `json:"icmpouttimeexcds"`
	Ipforwdatagrams     float64 `json:"ipforwdatagrams"`
	Ipinaddrerrors      float64 `json:"ipinaddrerrors"`
	Ipindelivers        float64 `json:"ipindelivers"`
	Ipindiscards        float64 `json:"ipindiscards"`
	Iphdrerrors         float64 `json:"iphdrerrors"`
	Ipinreceives        float64 `json:"ipinreceives"`
	Ipunknownprotos     float64 `json:"ipunknownprotos"`
	Ipoutdiscards       float64 `json:"ipoutdiscards"`
	Ipoutnoroutes       float64 `json:"ipoutnoroutes"`
	Ipoutrequests       float64 `json:"ipoutrequests"`
	Iproutingdiscards   float64 `json:"iproutingdiscards"`
	Timestamp           string  `json:"@timestamp"`
}
