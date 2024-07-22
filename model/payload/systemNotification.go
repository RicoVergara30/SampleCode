package payload

import (
	"encoding/xml"
	"sample/model/bah"
)

// SYSTEM STATUS BROADCAST MESSAGE
type (
	SystemNotificationISO20022 struct {
		XMLName   xml.Name                       `xml:"Message"`
		XMLNS     string                         `xml:"xmlns,attr"`
		XMLNsHead string                         `xml:"head,attr"`
		XMLNsNe   string                         `xml:"ne,attr"`
		Header    bah.HCRequestApplicationHeader `xml:"AppHdr"`
		Body      SystemNotificationWrapper
	}

	SystemNotificationWrapper struct {
		XMLName xml.Name               `xml:"SystemNotificationEvent"`
		Body    SystemNotificationBody `xml:"SysEvtNtfctn"`
	}

	SystemNotificationBody struct {
		XMLName            xml.Name           `xml:"SysEvtNtfctn"`
		SystemNotification SystemNotification `xml:"EvtInf"`
	}
	SystemNotification struct {
		EventCode        string   `xml:"EvtCd"`
		EventParams      []string `xml:"EvtParam"`
		EventDescription []string `xml:"EvtDesc"`
		EvenTime         string   `xml:"EvtTm"`
	}
)

type (
	NotificationReceived struct {
		Notification interface{}
	}

	CheckLetter struct {
		Input string `json:"letter"`
	}

	Notification struct {
		Id          int    `json:"id"`
		EventCode   string `json:"eventCode"`
		Description string `json:"description"`
		Parameters  string `json:"parameters"`
		XmlReceived string `json:"xml_received"`
		DateTime    string `json:"date_time"`
	}
)
