package fmevent

type FmEventMsg struct {
	Command          string `json:"command"`
	AlarmID          uint   `json:"alarmid"`
	Status           string `json:"status"`
	Organization     string `json:"organization"`
	EquipID          string `json:"equipid"`
	EquipName        string `json:"equipname"`
	Location         string `json:"location"`
	LocationName     string `json:"locationname"`
	LocationAlias    string `json:"locationalias"`
	EventType        string `json:"eventtype"`
	AlarmCode        string `json:"alarmcode"`
	AlarmCategory    string `json:"alarmcategory"`
	Severity         string `json:"severity"`
	User             string `json:"user"`
	OccurredTime     string `json:"occurredtime"`
	RecvTime         string `json:"recvtime"`
	ClearTime        string `json:"cleartime"`
	ClearUser        string `json:"clearuser"`
	AckTime          string `json:"acktime"`
	AckUser          string `json:"ackuser"`
	Times            string `json:"times"`
	UpdateTime       string `json:"updatetime"`
	RecvUpdateTime   string `json:"recvupdatetime"`
	ServerKey        string `json:"serverkey"`
	Message          string `json:"message"`
	Description      string `json:"description"`
	NetCD            string `json:"netcd"`
	NetTypeName      string `json:"nettypename"`
	NetSubTypeName   string `json:"netsubtypename"`
	EquipTypeName    string `json:"equiptypename"`
	EquipSubTypeName string `json:"equipsubtypename"`
	PKey             string `json:"pkey"`
	PortDescr        string `json:"portdescr"`
	IfIP             string `json:"ifip"`
	RingName         string `json:"ringname"`
	UpperLink        string `json:"upperlink"`
	RootLink         string `json:"rootlink"`
	IsWorking        string `json:"isworking"`
	OtherContents    string `json:"othercontents"`
	Timestamp        string `json:"@timestamp"`
}
