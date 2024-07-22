// Package fiber provides utility functions for gofiber v2, jwt-go
// With additional validation functions, sending JSON response and parsing request bodies, getting JWT claims
package util

import (
	"bytes"
	"compress/zlib"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"encoding/xml"
	"fmt"
	"io"
	"sample/db"
	envrouting "sample/middleware/envRouting"
	"sample/model/bah"
	"sample/model/payload"

	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

// Context GoFiber Context
type Context struct {
	c *fiber.Ctx
}

// JWTConfig configuration for JWT
type JWTConfig struct {
	Duration     time.Duration
	CookieMaxAge int
	SetCookies   bool
	SecretKey    []byte
}

// Ctx Context to be initiated by the New function
var Ctx Context
var jwtConfig JWTConfig

// New Copies GoFiber context as new current context
func (ctx *Context) New(c *fiber.Ctx) {
	Ctx = Context{
		c: c,
	}
}

// Message struct for GoFiber context response
type Message struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// ParseBody Parses the request body from the copied current context
func ParseBody(in interface{}) error {
	err := Ctx.c.BodyParser(in)

	if err != nil {
		LogError(err)
		return Ctx.c.Status(503).SendString(err.Error())
	}
	return err
}

// GetParamValue Gets the parameter value from the copied current context
func GetParamValue(param string, message string) string {
	paramValue := Ctx.c.Params(param)

	if paramValue == "" {
		err := SendJSONMessage(message, false, 400)
		LogError(err)
	}

	return paramValue
}

// SendJSONMessage Sends JSON Message with HTTP Status code to current context
func SendJSONMessage(message string, isSuccess bool, httpStatusCode int) error {
	status := "failed"

	if isSuccess {
		status = "success"
	}

	return Ctx.c.Status(httpStatusCode).JSON(Message{
		Message: message,
		Status:  status,
	})
}

// SendSuccessResponse Wrapper function for SendJSONMessage of 200 Success
func SendSuccessResponse(message string) error {
	err := SendJSONMessage(message, true, 200)
	LogError(err)
	return err
}

// SendBadRequestResponse Wrapper function for SendJSONMessage of 400 Bad request
func SendBadRequestResponse(message string) error {
	err := SendJSONMessage(message, false, 400)
	LogError(err)
	return err
}

// ValidateField Validation of strings and return if valid based on specification and error message if invalid
func ValidateField(field, title string, isMandatory bool, max, min int, format string) (ok bool, message string) {
	ok = true

	if !isMandatory {
		return
	}

	if len(field) == 0 {
		message += fmt.Sprintf("'%s' cannot be empty.", title)
		_ = SendBadRequestResponse(message)
		ok = false
	} else {
		switch format {
		case "S":
			if len(field) > 2 {
				message += fmt.Sprintf("The length of '%s' cannot be greater than 2.", title)
				_ = SendBadRequestResponse(message)
				ok = false
			}
		case "N":
			if _, err := strconv.Atoi(field); err != nil {
				message += fmt.Sprintf("'%s' should only contain numbers.", title)
				_ = SendBadRequestResponse(message)
				ok = false
			}

			fallthrough
		case "ANS":
			cflOK, cflMessage := CheckFieldLength(field, title, max, min)

			if !cflOK {
				ok = false
				message += "\n" + cflMessage
			}
		}
	}

	return ok, message
}

// CheckFieldLength Validation of strings' length and return if valid based on maximum and minimum length specified and error message if invalid
func CheckFieldLength(field, title string, max, min int) (ok bool, message string) {
	ok = true

	if len(field) > max {
		message += fmt.Sprintf("The length of '%s' cannot be greater than %d.", title, max)
		_ = SendBadRequestResponse(message)
		ok = false
	}

	if len(field) < min {
		message += fmt.Sprintf("The length of '%s' cannot be less than %d.", title, min)
		_ = SendBadRequestResponse(message)
		ok = false
	}

	return
}

// GetJWTClaims Get User JWT claims of the current context
func GetJWTClaims() jwt.MapClaims {
	userToken := Ctx.c.Locals("user").(*jwt.Token)
	return userToken.Claims.(jwt.MapClaims)
}

// GetJWTClaim Wrapper function for getting a JWT claim by key
func GetJWTClaim(key string) map[string]interface{} {
	return GetJWTClaims()[key].(map[string]interface{})
}

// GetJSONFieldValues Returns a map of JSON keys and values of a struct
func GetJSONFieldValues(q interface{}) map[string]string {
	val := reflect.ValueOf(q).Elem()
	fields := make(map[string]string)

	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		fields[tag.Get("json")] = val.Field(i).String()
	}

	return fields
}

// ValidateJSONField Wrapper function for JSON field validation of a struct
func ValidateJSONField(q interface{}, tag string, isMandatory bool, max, min int, format string) (bool, string) {
	return ValidateField(GetJSONFieldValues(q)[tag], tag, isMandatory, max, min, format)
}

// LogError Logs errors
func LogError(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

// AuthenticationMiddleware ...
func AuthenticationMiddleware(j JWTConfig) func(*fiber.Ctx) error {
	jwtConfig = j
	return jwtware.New(jwtware.Config{
		SigningKey: jwtConfig.SecretKey,
	})
}

// GenerateJWTSignedString ...
func GenerateJWTSignedString(claimMaps fiber.Map) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(jwtConfig.Duration).Unix()

	for key, value := range claimMaps {
		claims[key] = value
	}

	t, err := token.SignedString(jwtConfig.SecretKey)

	if jwtConfig.SetCookies && err == nil {
		Ctx.c.Cookie(&fiber.Cookie{
			Name:   "token",
			Value:  t,
			MaxAge: jwtConfig.CookieMaxAge,
		})
	}

	return t, err
}

// GetJWTClaimOfType ...
func GetJWTClaimOfType(key string, valueType interface{}) error {
	userInfoJSON, err := json.Marshal(GetJWTClaim(key))

	if err == nil {
		err = json.Unmarshal(userInfoJSON, &valueType)
	}

	return err
}

// ---------------------------------
// initial length is 10, add value to extend the length
func GenerateIdentifier(len int) string {
	now := time.Now()
	nowTime := now.Unix()
	var identifier string

	for ctr := 0; ctr < len; ctr++ {
		rndNum := strconv.Itoa(rand.Intn(9-0) + 1)
		identifier += rndNum
	}

	identifier = strconv.Itoa(int(nowTime)) + identifier
	return identifier
}

func GenerateInstructionID(max_digits, current_count int) string {
	max_digits = max_digits - 1
	var instructionID string
	current_length := len(strconv.Itoa(current_count))

	if current_length <= max_digits {
		current_count++
		for strL := 0; strL <= max_digits-current_length; strL++ {
			instructionID += "0"
		}
	} else {
		current_count = 1
		for strL := 0; strL <= max_digits-current_length; strL++ {
			instructionID += "0"
		}
	}

	instructionID += strconv.Itoa(current_count)
	return instructionID
}

func InsertNumber(new_counter int, filePath string) error {
	// Read the current file contents
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Convert the content to a string
	data := string(content)

	// Append the new number to the string
	data += strconv.Itoa(new_counter) + "\n"

	// Write the updated content to the file
	err = ioutil.WriteFile(filePath, []byte(data), 0644)
	if err != nil {
		return err
	}

	return nil
}

func GetLastInsertedNumber(filePath string) (int, error) {
	// Read the file contents
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return 0, err
	}

	// Convert the content to a string
	data := string(content)

	// Split the content by newlines
	lines := strings.Split(data, "\n")

	// Get the last non-empty line (assumed to be the last inserted number)
	lastLine := ""
	for i := len(lines) - 1; i >= 0; i-- {
		if lines[i] != "" {
			lastLine = lines[i]
			break
		}
	}

	// Parse the last line as an integer
	number, err := strconv.Atoi(lastLine)
	if err != nil {
		return 0, err
	}
	return number, nil
}

func CreateFileText(filename string, value string) {
	rootPath, _ := os.Getwd()
	filePath := fmt.Sprintf("%s/%s", rootPath, filename)
	creatingErr := ioutil.WriteFile(filePath, []byte(value), 0644)
	if creatingErr != nil {
		fmt.Println("Error creating file:", creatingErr)
	}
	fmt.Println("File created successfully.")
}

func GenerateXMLIdentifier(filePath string, currentNumber, maxDigit int) string {
	var currentInstructionID string
	if currentNumber == 0 {
		currentInstructionID = GenerateInstructionID(maxDigit, currentNumber)
		// Insert a number
		newNumber, _ := strconv.Atoi(currentInstructionID)
		insertErr := InsertNumber(newNumber, filePath)
		if insertErr != nil {
			fmt.Println("Error inserting number:", insertErr)
		}
	}
	currentInstructionID = GenerateInstructionID(maxDigit, currentNumber)
	newNumber, _ := strconv.Atoi(currentInstructionID)
	insertErr := InsertNumber(newNumber, filePath)
	if insertErr != nil {
		fmt.Println("Error inserting number:", insertErr)
	}

	return currentInstructionID
}

func LoadCertificate(filename string) *x509.Certificate {
	// LOAD CERTIFICATE
	certFile, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading certificate file:", err.Error())
	}

	// PARSE CERTIFICATE
	block, _ := pem.Decode(certFile)
	if block == nil {
		fmt.Println("Error decoding PEM Block:", block)
	}

	cert, certErr := x509.ParseCertificate(block.Bytes)
	if certErr != nil {
		fmt.Println("Error parsing certificate:", certErr.Error())
	}

	return cert
}

func ValidateLocalInstrument(instrument string) (string, bool) {
	localInstrument := &bah.LocalInstrumentList{}
	db.Database.Raw("SELECT * FROM rbi_instapay.local_instruments_list WHERE local_instrument = ?", instrument).Scan(localInstrument)

	if localInstrument.LocalInstrument != "" {
		if localInstrument.IsEnabled {
			return localInstrument.LocalInstrument, localInstrument.IsEnabled
		}
		return localInstrument.LocalInstrument, localInstrument.IsEnabled
	}
	return "Local instrument not found", false
}

func TestLI(c *fiber.Ctx) error {
	paramModel := c.Params("li")

	localInstrument, isEnabled := ValidateLocalInstrument(paramModel)

	return c.JSON(fiber.Map{
		"LocalInstument": localInstrument,
		"isEnabled":      isEnabled,
	})
}

func GetServiceEP(service, environment string) string {
	serviceRoute := &bah.ServiceRoute{}
	db.Database.Raw("SELECT * FROM rbi_instapay.get_service(?,?)", service, environment).Scan(serviceRoute)
	return serviceRoute.ServiceUrl
}

// GetServiceEP_igate
func GetServiceEP_igate(service, environment string) string {
	serviceRoute := &bah.ServiceRoute{}
	db.Database.Raw("SELECT * FROM rbi_instapay.get_service(?,?)", service, environment).Scan(serviceRoute)
	return serviceRoute.ServiceUrl
}

func ShowServiceEP(c *fiber.Ctx) error {
	sevice := c.Params("service")
	environment := c.Params("env")
	url := GetServiceEP(sevice, environment)
	return c.SendString(url)
}

func InsertNotification(c *fiber.Ctx) error {
	fmt.Println("TEST INSERT NOTIFICATION")
	systemNotification := &payload.SystemNotificationISO20022{}
	if parsErr := c.BodyParser(systemNotification); parsErr != nil {
		return c.Status(500).SendString(parsErr.Error())
	}

	eventCode := systemNotification.Body.Body.SystemNotification.EventCode
	eventDateTime := systemNotification.Body.Body.SystemNotification.EvenTime

	eventDescription, _ := xml.Marshal(systemNotification.Body.Body.SystemNotification.EventDescription)
	notificationParams, _ := xml.Marshal(systemNotification.Body.Body.SystemNotification.EventParams)

	fmt.Println("Event Code:", string(eventCode))
	fmt.Println("Event Date Time:", string(eventDateTime))
	fmt.Println("Event Description:", string(eventDescription))
	fmt.Println("Event Description:", string(notificationParams))

	notification, marshalErr := xml.Marshal(systemNotification)
	if marshalErr != nil {
		return marshalErr
	}

	return c.Send(notification)
}

func CheckIPSStatus(c *fiber.Ctx) error {
	// Set timezone
	loc, err := time.LoadLocation(envrouting.PostgresTimeZone)
	if err != nil {
		fmt.Println("Error loading time zone:", err)
		return err
	}

	// Get current time
	currentTime := time.Now().In(loc)

	return c.JSON(bah.StatusFields{
		IsSignedOn: false,
		Remarks:    "Sever is down",
		SignedBy:   "George",
		Downtime: bah.UpDowntime{
			Date: currentTime.Format(time.DateOnly),
			Time: currentTime.Format(time.Kitchen),
		},
		Uptime: bah.UpDowntime{
			Date: currentTime.Format(time.DateOnly),
			Time: currentTime.Format(time.Kitchen),
		},
	})
}

func MarshalJson(jsonFormat interface{}) string {
	result, err := json.Marshal(jsonFormat)
	if err != nil {
		fmt.Println("Can't marshal json:", err)
	}

	return string(result)
}

func ToASCII(utf string) string {
	return base64.StdEncoding.EncodeToString([]byte(utf))
}

func ToUTF(ascii string) string {
	var in, out bytes.Buffer
	b := []byte(ascii)
	utf := []byte{}
	w := zlib.NewWriter(&in)
	w.Write(b)
	w.Close()
	r, err := zlib.NewReader(&in)

	if err == nil {
		io.Copy(&out, r)
		utf, _ = base64.StdEncoding.DecodeString(out.String())
	}

	return string(utf)
}

// Used to change special character
func CorrectString(text string) string {
	var fixedString string
	for _, char := range text {
		if !IsASCII(char) {
			fmt.Printf("SPECIAL CHARACTER FOUND! %v : %v \n", string(char), char)
			fixedString += strings.ReplaceAll(string(char), string(char), "Ñ")
		} else {
			fixedString += string(char)
		}
	}
	fmt.Println("FIXED STRING:", fixedString)
	return fixedString
}

func IsASCII(char rune) bool {
	return char >= 0 && char <= 127
}

// Postman Test
func CheckMessageSequence(c *fiber.Ctx) error {
	messageModel := &bah.MessageSequence{}

	c.BodyParser(messageModel)
	messageId := GenerateMessageIdentifier(messageModel.Prefix, messageModel.BIC, messageModel.GenerationSource, 14)

	return c.JSON(fiber.Map{
		"messageID:": messageId,
	})
}

// Generation of ISO20022 message identifier
func GenerateMessageIdentifier(prefix, bic, generationSource string, maxLength int) string {
	messageSequenceModel := &bah.MessageSequence{}
	dateNow := time.Now().Format("20060102")

	// Get the previous serial number
	if prefix == "TX" {
		db.Database.Raw("SELECT * FROM rbi_instapay.message_sequence WHERE prefix = 'TX' AND generation_source = ? LIMIT 1", generationSource).Scan(messageSequenceModel)
	} else if prefix == "HC" {
		db.Database.Raw("SELECT * FROM rbi_instapay.message_sequence WHERE prefix = 'HC' AND generation_source = ? LIMIT 1", generationSource).Scan(messageSequenceModel)
	} else if prefix == "NA" {
		db.Database.Raw("SELECT * FROM rbi_instapay.message_sequence WHERE prefix = 'NA' AND generation_source = ? LIMIT 1", generationSource).Scan(messageSequenceModel)
	} else {
		db.Database.Raw("SELECT * FROM rbi_instapay.message_sequence WHERE prefix = ? AND generation_source = ? LIMIT 1", prefix, generationSource).Scan(messageSequenceModel)
	}

	serial := GenerateSequenceNumber(prefix, messageSequenceModel.SerialNumber, maxLength)
	result := ""

	if prefix == "NA" || prefix == "TX" || prefix == "HC" {
		result = fmt.Sprintf("%s%s%s%s", dateNow, bic, generationSource, serial)
	} else {
		result = fmt.Sprintf("%s%s%s%s%s", prefix, dateNow, bic, generationSource, serial)
	}

	// Update the serial number
	if messageSequenceModel.Id == 0 {
		if prefix != "" {
			db.Database.Exec(`
			INSERT INTO rbi_instapay.message_sequence (prefix, bic, generation_source, serial_number)
			VALUES (?,?,?,?)
		`, prefix, bic, generationSource, serial)
		} else {
			db.Database.Exec(`
			INSERT INTO rbi_instapay.message_sequence (prefix, bic, generation_source, serial_number)
			VALUES ('NA',?,?,?)`, bic, generationSource, serial)
		}
	} else {
		if prefix != "" {
			db.Database.Exec(`UPDATE rbi_instapay.message_sequence SET serial_number= ?, updated_at = ? WHERE prefix = ? AND generation_source = ?`,
				serial, time.Now(), prefix, generationSource)
		} else if prefix == "TX" {
			db.Database.Exec(`UPDATE rbi_instapay.message_sequence SET serial_number= ?, updated_at = ? WHERE prefix = 'TX' AND generation_source = ?`,
				serial, time.Now(), generationSource)
		} else {
			db.Database.Exec(`UPDATE rbi_instapay.message_sequence SET serial_number= ?, updated_at = ? WHERE prefix = 'NA' AND generation_source = ?`,
				serial, time.Now(), generationSource)
		}
	}
	return result
}

// use to generate digits to complete ISO20022 xml message structure
func GenerateSequenceNumber(prefix string, previous, maxLength int) string {
	var formattedNumber = ""
	strPrevious := strconv.Itoa(previous)

	if len(strPrevious) >= maxLength || previous == 0 {
		previous = 0
	}
	series := "%0" + strconv.Itoa(maxLength) + "d"
	current := previous + 1
	if prefix == "TX" {
		if previous%2 == 0 {
			current = previous + 1
		} else {
			current = previous + 3
		}
	}

	// Loop through the numbers from start to end (inclusive)
	for i := previous; i <= current; i++ {
		// Convert the integer back to a string with leading zeros using formatting
		formattedNumber = fmt.Sprintf(series, i)
		// series = append(series, formattedNumber)

	}
	// fmt.Println("RESULT:", formattedNumber)
	return formattedNumber
}
