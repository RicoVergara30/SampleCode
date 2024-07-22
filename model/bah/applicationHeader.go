package bah

import (
	"encoding/xml"
)

// MESSAGE DEFINITION
const (
	HeadMessageDefinition                    = "head.001.001.01" // Head
	SignOnRequestMessageDefinition           = "admn.001.001.01" // Sign-On Request - sr
	SignOnResponseMessageDefinition          = "admn.002.001.01" // Sign-On Response - rs
	SignOffRequestMessageDefinition          = "admn.003.001.01" // Sign-Off Request - fr
	SignOffResponseMessageDefinition         = "admn.004.001.01" // Sign-Off Response - rf
	SystemNotificationEventMessageDefinition = "admi.004.001.02" // System Event Notification - ne
	EchoRequestMessageDefinition             = "admn.005.001.01" // Echo Request - er
	EchoResponseMessageDefinition            = "admn.006.001.01" // Echo Response - re
	CreditTransferMessageDefinition          = "pacs.008.001.08" // Credit Transfer V08 - ct
	MessageStatusReportMessageDefinition     = "pacs.002.001.10" // Message Status Report V10 (Response to Business Messages) - ps
	MessageRejectMessageDefinition           = "admi.002.001.01" // Message Reject (Admin) - mr
	PaymentCancellationMessageDefinition     = "camt.056.001.08" // Payment Cancellation V08 (Request for Return of Funds or System Time-out) - rt
	RBIBankIdentifierCode                    = "CAMZPHM2XXX"     // "CAMZPHM2XXXYY" - rc650 BIC incorrect
)

type (
	// ------------------
	// APPLICATION HEADER
	// ------------------
	ApplicationHeader struct {
		XMLName             xml.Name      `xml:"AppHdr"`
		From                string        `xml:"head:Fr>head:FIId>head:FinInstnId>head:BICFI"`
		To                  string        `xml:"head:To>head:FIId>head:FinInstnId>head:BICFI"`
		BusinessMessageID   string        `xml:"head:BizMsgIdr"`
		MessageDefinationID string        `xml:"head:MsgDefIdr"`
		CreationDateTime    string        `xml:"head:CreDt"`
		Signature           X509Signature `xml:"head:Sgntr>ds:Signature"`
		CopyDuplicate       string        `xml:"head:CpyDplct,omitempty"`
	}

	// -------------------------
	// DIGEST APPLICATION HEADER
	// -------------------------
	DigestApplicationHeader struct {
		XMLName             xml.Name `xml:"AppHdr"`
		From                string   `xml:"head:Fr>head:FIId>head:FinInstnId>head:BICFI"`
		To                  string   `xml:"head:To>head:FIId>head:FinInstnId>head:BICFI"`
		BusinessMessageID   string   `xml:"head:BizMsgIdr"`
		MessageDefinationID string   `xml:"head:MsgDefIdr"`
		CreationDateTime    string   `xml:"head:CreDt"`
		Signature           string   `xml:"head:Sgntr"`
		CopyDuplicate       string   `xml:"head:CpyDplct,omitempty"`
	}
)

type (
	DigestDuplicateApplicationHeader struct {
		XMLName             xml.Name `xml:"AppHdr"`
		From                string   `xml:"head:Fr>head:FIId>head:FinInstnId>head:BICFI"`
		To                  string   `xml:"head:To>head:FIId>head:FinInstnId>head:BICFI"`
		BusinessMessageID   string   `xml:"head:BizMsgIdr"`
		MessageDefinationID string   `xml:"head:MsgDefIdr"`
		CreationDateTime    string   `xml:"head:CreDt"`
		Signature           string   `xml:"head:Sgntr"`
		CopyDuplicate       string   `xml:"head:CpyDplct,omitempty"`
	}

	DigestBranchApplicationHeader struct {
		XMLName             xml.Name               `xml:"AppHdr"`
		From                DigestFIIdentification `xml:"head:Fr>head:FIId"`
		To                  string                 `xml:"head:To>head:FIId>head:FinInstnId>head:BICFI"`
		BusinessMessageID   string                 `xml:"head:BizMsgIdr"`
		MessageDefinationID string                 `xml:"head:MsgDefIdr"`
		CreationDateTime    string                 `xml:"head:CreDt"`
		Signature           string                 `xml:"head:Sgntr"`
	}

	// DuplicateApplicationHeader struct {
	// 	XMLName             xml.Name      `xml:"AppHdr"`
	// 	From                string        `xml:"head:Fr>head:FIId>head:FinInstnId>head:BICFI"`
	// 	To                  string        `xml:"head:To>head:FIId>head:FinInstnId>head:BICFI"`
	// 	BusinessMessageID   string        `xml:"head:BizMsgIdr"`
	// 	MessageDefinationID string        `xml:"head:MsgDefIdr"`
	// 	CreationDateTime    string        `xml:"head:CreDt"`
	// 	Signature           X509Signature `xml:"head:Sgntr>ds:Signature"`
	// 	CopyDuplicate       string        `xml:"head:CpyDplct"`
	// }

	DigestFIIdentification struct {
		BICI     string `xml:"head:FinInstnId>head:BICFI"`
		BranchID string `xml:"head:BrnchId>head:Id"`
	}

	RequestApplicationHeader struct {
		XMLName             xml.Name      `xml:"AppHdr"`
		From                string        `xml:"Fr>FIId>FinInstnId>BICFI"`
		To                  string        `xml:"To>FIId>FinInstnId>BICFI"`
		BusinessMessageID   string        `xml:"BizMsgIdr"`
		MessageDefinationID string        `xml:"MsgDefIdr"`
		CreationDateTime    string        `xml:"CreDt"`
		Signature           X509Signature `xml:"Sgntr>ds:Signature"`
	}
)

type (
	HCSignature struct {
		XMLName        xml.Name     `xml:"Signature"`
		XMLNS          string       `xml:"xmlns,attr"`
		SignedInfo     HCSignedInfo `xml:"SignedInfo"`
		SignatureValue string       `xml:"SignatureValue"`
		X509Data       HCX509Data   `xml:"KeyInfo>X509Data"`
	}

	X509Signature struct {
		XMLNS          string     `xml:"xmlns:ds,attr"`
		SignedInfo     SignedInfo `xml:"ds:SignedInfo"`
		SignatureValue string     `xml:"ds:SignatureValue"`
		X509Data       X509Data   `xml:"ds:KeyInfo>ds:X509Data"`
	}

	X509Data struct {
		SubjectName  string           `xml:"ds:X509SubjectName"`
		IssuerSerial X509IssuerSerial `xml:"ds:X509IssuerSerial"`
	}

	X509IssuerSerial struct {
		IssuerName   string `xml:"ds:X509IssuerName"`
		SerialNumber string `xml:"ds:X509SerialNumber"`
	}
)

// -----------------------------
// SIGNED INFORMATION FOR DIGEST
// -----------------------------
type (
	DigestSignedInfo struct {
		XMLName                xml.Name  `xml:"ds:SignedInfo"`
		XMLNS                  string    `xml:"xmlns,attr"`
		XMLNSct                string    `xml:"xmlns:ct,attr,omitempty"` // Credit Transfer
		XMLNSds                string    `xml:"xmlns:ds,attr,omitempty"`
		XMLNSer                string    `xml:"xmlns:er,attr,omitempty"` // Echo Request
		XMLNSfr                string    `xml:"xmlns:fr,attr,omitempty"` // SignOff Request
		XMLNShead              string    `xml:"xmlns:head,attr,omitempty"`
		XMLNSmr                string    `xml:"xmlns:mr,attr,omitempty"` // Message Reject
		XMLNSne                string    `xml:"xmlns:ne,attr,omitempty"` // System Event Notification
		XMLNSps                string    `xml:"xmlns:ps,attr,omitempty"` // Message Status
		XMLNSre                string    `xml:"xmlns:re,attr,omitempty"` // Echo Response
		XMLNSrf                string    `xml:"xmlns:rf,attr,omitempty"` // SignOff Response
		XMLNSrs                string    `xml:"xmlns:rs,attr,omitempty"` // SignOn Response
		XMLNSrt                string    `xml:"xmlns:rt,attr,omitempty"` // Payment Cancellation
		XMLNSsr                string    `xml:"xmlns:sr,attr,omitempty"` // SignOn Request
		CanonicalizationMethod Algorithm `xml:"ds:CanonicalizationMethod"`
		SignatureMethod        Algorithm `xml:"ds:SignatureMethod"`
		Reference              Reference `xml:"ds:Reference"`
	}
)

// -----------
// SIGNED INFO
// -----------
type (
	// -------------------
	// SIGNED INFO PAYLOAD
	// -------------------
	SignedInfo struct {
		XMLName                xml.Name  `xml:"ds:SignedInfo"`
		CanonicalizationMethod Algorithm `xml:"ds:CanonicalizationMethod"`
		SignatureMethod        Algorithm `xml:"ds:SignatureMethod"`
		Reference              Reference `xml:"ds:Reference"`
	}

	// ----------------------------
	// SIGNED INFO - SIGNON PAYLOAD
	// ----------------------------
	SrDigestSignedInfo struct {
		XMLName                xml.Name  `xml:"ds:SignedInfo"`
		XMLNS                  string    `xml:"xmlns,attr"`
		XMLNSds                string    `xml:"xmlns:ds,attr"`
		XMLNShead              string    `xml:"xmlns:head,attr"`
		XMLNSsr                string    `xml:"xmlns:sr,attr"`
		CanonicalizationMethod Algorithm `xml:"ds:CanonicalizationMethod"`
		SignatureMethod        Algorithm `xml:"ds:SignatureMethod"`
		Reference              Reference `xml:"ds:Reference"`
	}

	// -----------------------------
	// SIGNED INFO - SIGNOFF PAYLOAD
	// -----------------------------
	FrDigestSignedInfo struct {
		XMLName                xml.Name  `xml:"ds:SignedInfo"`
		XMLNS                  string    `xml:"xmlns,attr"`
		XMLNSds                string    `xml:"xmlns:ds,attr"`
		XMLNSFr                string    `xml:"xmlns:fr,attr"`
		XMLNSHead              string    `xml:"xmlns:head,attr"`
		CanonicalizationMethod Algorithm `xml:"ds:CanonicalizationMethod"`
		SignatureMethod        Algorithm `xml:"ds:SignatureMethod"`
		Reference              Reference `xml:"ds:Reference"`
	}

	// -------------------------------------------------------------------
	// SIGNED INFO - ECHO RESPONSE PAYLOAD (USED FOR HEALTH CHECK PAYLOAD)
	// -------------------------------------------------------------------
	ReDigestSignedInfo struct {
		XMLName                xml.Name  `xml:"ds:SignedInfo"`
		XMLNS                  string    `xml:"xmlns,attr"`
		XMLNSds                string    `xml:"xmlns:ds,attr"`
		XMLNSHead              string    `xml:"xmlns:head,attr"`
		XMLNSRe                string    `xml:"xmlns:re,attr"`
		CanonicalizationMethod Algorithm `xml:"CanonicalizationMethod"`
		SignatureMethod        Algorithm `xml:"ds:SignatureMethod"`
		Reference              Reference `xml:"ds:Reference"`
	}

	// -------------------------------------
	// SIGNED INFO - CREDIT TRANSFER PAYLOAD
	// -------------------------------------
	CtDigestSignedInfo struct {
		XMLName                xml.Name  `xml:"ds:SignedInfo"`
		XMLNS                  string    `xml:"xmlns,attr"`
		XMLNSct                string    `xml:"xmlns:ct,attr"`
		XMLNSds                string    `xml:"xmlns:ds,attr"`
		XMLNShead              string    `xml:"xmlns:head,attr"`
		CanonicalizationMethod Algorithm `xml:"ds:CanonicalizationMethod"`
		SignatureMethod        Algorithm `xml:"ds:SignatureMethod"`
		Reference              Reference `xml:"ds:Reference"`
	}

	// ------------------------------------
	// SIGNED INFO - MESSAGE STATUS PAYLOAD
	// ------------------------------------
	PsDigestSignedInfo struct {
		XMLName                xml.Name  `xml:"ds:SignedInfo"`
		XMLNS                  string    `xml:"xmlns,attr"`
		XMLNSds                string    `xml:"xmlns:ds,attr"`
		XMLNShead              string    `xml:"xmlns:head,attr"`
		XMLNSps                string    `xml:"xmlns:ps,attr"`
		CanonicalizationMethod Algorithm `xml:"ds:CanonicalizationMethod"`
		SignatureMethod        Algorithm `xml:"ds:SignatureMethod"`
		Reference              Reference `xml:"ds:Reference"`
	}

	// -----------------------------------------
	// SIGNED INFO - SYSTEM NOTIFICATION PAYLOAD
	// -----------------------------------------
	NeDigestSignedInfo struct {
		XMLName                xml.Name  `xml:"ds:SignedInfo"`
		XMLNS                  string    `xml:"xmlns,attr"`
		XMLNSds                string    `xml:"xmlns:ds,attr"`
		XMLNShead              string    `xml:"xmlns:head,attr"`
		XMLNSNe                string    `xml:"xmlns:ne,attr"`
		CanonicalizationMethod Algorithm `xml:"ds:CanonicalizationMethod"`
		SignatureMethod        Algorithm `xml:"ds:SignatureMethod"`
		Reference              Reference `xml:"ds:Reference"`
	}

	// -----------------------------------------
	// SIGNED INFO - PAYMENT CANCELATION PAYLOAD
	// -----------------------------------------
	RtDigestSignedInfo struct {
		XMLName                xml.Name  `xml:"ds:SignedInfo"`
		XMLNS                  string    `xml:"xmlns,attr"`
		XMLNSds                string    `xml:"xmlns:ds,attr"`
		XMLNShead              string    `xml:"xmlns:head,attr"`
		XMLNSrt                string    `xml:"xmlns:rt,attr"`
		CanonicalizationMethod Algorithm `xml:"ds:CanonicalizationMethod"`
		SignatureMethod        Algorithm `xml:"ds:SignatureMethod"`
		Reference              Reference `xml:"ds:Reference"`
	}

	// ------------------------------------
	// SIGNED INFO - MESSAGE REJECT PAYLOAD
	// ------------------------------------
	MrDigestSignedInfo struct {
		XMLName                xml.Name  `xml:"ds:SignedInfo"`
		XMLNS                  string    `xml:"xmlns,attr"`
		XMLNSds                string    `xml:"xmlns:ds,attr"`
		XMLNShead              string    `xml:"xmlns:head,attr"`
		XMLNSmr                string    `xml:"xmlns:mr,attr"`
		CanonicalizationMethod Algorithm `xml:"ds:CanonicalizationMethod"`
		SignatureMethod        Algorithm `xml:"ds:SignatureMethod"`
		Reference              Reference `xml:"ds:Reference"`
	}

	// -----------------------
	// REFERENCE AND ALGORITHM
	// -----------------------
	Reference struct {
		URI          string      `xml:"URI,attr"`
		Transforms   []Algorithm `xml:"ds:Transforms>ds:Transform"`
		DigestMethod Algorithm   `xml:"ds:DigestMethod"`
		DigestValue  string      `xml:"ds:DigestValue"`
	}

	Algorithm struct {
		Algorithm string `xml:"Algorithm,attr"`
	}
)

// FOR FETCHING REQUEST BODY
type (
	HCRequestApplicationHeader struct {
		XMLName             xml.Name `xml:"AppHdr"`
		From                string   `xml:"Fr>FIId>FinInstnId>BICFI"`
		To                  string   `xml:"To>FIId>FinInstnId>BICFI"`
		BusinessMessageID   string   `xml:"BizMsgIdr"`
		MessageDefinationID string   `xml:"MsgDefIdr"`
		CreationDateTime    string   `xml:"CreDt"`
		Signature           HCSgntr  `xml:"Sgntr"`
	}

	RequestBranchApplicationHeader struct {
		XMLName             xml.Name         `xml:"AppHdr"`
		From                FIIdentification `xml:"Fr>FIId"`
		To                  string           `xml:"To>FIId>FinInstnId>BICFI"`
		BusinessMessageID   string           `xml:"BizMsgIdr"`
		MessageDefinationID string           `xml:"MsgDefIdr"`
		CreationDateTime    string           `xml:"CreDt"`
		Signature           HCSgntr          `xml:"Sgntr"`
	}

	FIIdentification struct {
		BICI     string `xml:"FinInstnId>BICFI"`
		BranchID string `xml:"BrnchId>Id"`
	}

	HCSgntr struct {
		XMLName   xml.Name `xml:"Sgntr"`
		Signature HCSignature
	}
)

type (
	HCSignedInfo struct {
		XMLName                xml.Name    `xml:"SignedInfo"`
		CanonicalizationMethod Algorithm   `xml:"CanonicalizationMethod"`
		SignatureMethod        Algorithm   `xml:"SignatureMethod"`
		Reference              HCReference `xml:"Reference"`
	}

	HCReference struct {
		XMLName      xml.Name       `xml:"Reference"`
		URI          string         `xml:"URI,attr"`
		Transforms   []HCTransform  `xml:"Transforms>Transform"`
		DigestMethod HCDigestMethod `xml:"DigestMethod"`
		DigestValue  string         `xml:"DigestValue"`
	}

	HCTransform struct {
		Algorithm string `xml:"Algorithm,attr"`
	}

	HCDigestMethod struct {
		XMLName   xml.Name `xml:"DigestMethod"`
		Algorithm string   `xml:"Algorithm,attr"`
	}

	HCX509Data struct {
		SubjectName  string             `xml:"X509SubjectName"`
		IssuerSerial HCX509IssuerSerial `xml:"X509IssuerSerial"`
	}

	HCX509IssuerSerial struct {
		IssuerName   string `xml:"X509IssuerName"`
		SerialNumber string `xml:"X509SerialNumber"`
	}
)

type (
	Errors struct {
		Source      string `json:"Source"`
		ReasonCode  string `json:"ReasonCode"`
		Decription  string `json:"Description"`
		Recoverable bool   `json:"Recoverable"`
		Details     string `json:"Details"`
	}

	BICDetails struct {
		Type           string `json:"type"`
		Name           string `json:"name"`
		BankCode       string `json:"bankCode"`
		Mnemonic       string `json:"mnemonic"`
		Bic            string `json:"bic"`
		ConnectionType string `json:"connectionType"`
	}
)

type (
	MessageStructureISO20022 struct {
		XMLName   xml.Name          `xml:"Message"`
		XMLNS     string            `xml:"xmlns,attr"`
		XMLNSct   string            `xml:"xmlns:ct,attr,omitempty"` // Credit Transfer
		XMLNSds   string            `xml:"xmlns:ds,attr,omitempty"`
		XMLNSer   string            `xml:"xmlns:er,attr,omitempty"` // Echo Request
		XMLNSfr   string            `xml:"xmlns:fr,attr,omitempty"` // SignOff Request
		XMLNShead string            `xml:"xmlns:head,attr,omitempty"`
		XMLNSmr   string            `xml:"xmlns:mr,attr,omitempty"` // Message Reject
		XMLNSne   string            `xml:"xmlns:ne,attr,omitempty"` // System Event Notification
		XMLNSps   string            `xml:"xmlns:ps,attr,omitempty"` // Message Status
		XMLNSre   string            `xml:"xmlns:re,attr,omitempty"` // Echo Response
		XMLNSrf   string            `xml:"xmlns:rf,attr,omitempty"` // SignOff Response
		XMLNSrs   string            `xml:"xmlns:rs,attr,omitempty"` // SignOn Response
		XMLNSrt   string            `xml:"xmlns:rt,attr,omitempty"` // Payment Cancellation
		XMLNSsr   string            `xml:"xmlns:sr,attr,omitempty"` // SignOn Request
		Header    ApplicationHeader `xml:"AppHdr"`
		Body      interface{}
	}
)

type (
	LocalInstrumentList struct {
		LocalInstrument string `json:"localInstrument"`
		IsEnabled       bool   `json:"isEnabled"`
		Description     string `json:"description"`
	}

	ServiceRoute struct {
		ServiceUrl string `json:"url"`
	}

	StatusFields struct {
		IsSignedOn bool       `json:"isSignedOn"`
		Remarks    string     `json:"remarks"`
		SignedBy   string     `json:"signedBy"`
		Downtime   UpDowntime `json:"downtime"`
		Uptime     UpDowntime `json:"uptime"`
	}

	UpDowntime struct {
		Date string `json:"date"`
		Time string `json:"time"`
	}

	IPSStatus struct {
		SignedOn    bool   `json:"SignedOn"`
		Remarks     string `json:"remarks"`
		SignedBy    string `json:"signedBy"`
		SignonDate  string `json:"signonDate"`
		SignonTime  string `json:"signonTime"`
		SignoffDate string `json:"signoffDate"`
		SignoffTime string `json:"signoffTime"`
	}

	SetupDowntime struct {
		SetupDowntime int `json:"setupDowntime"`
	}
)

type (
	MessageSequence struct {
		Id               int    `json:"id"`
		BIC              string `json:"bic,omitempty"`
		Prefix           string `json:"prefix"`
		CreationDate     string `json:"creationDate"`
		GenerationSource string `json:"generationSource"`
		SerialNumber     int    `json:"serialNumber"`
		CreatedAt        string `json:"createdAt"`
		UpdateAt         string `json:"updatedAt"`
	}
)
