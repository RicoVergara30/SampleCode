package transfer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sample/db"
	"sample/middleware/encryptDecrypt"
	envrouting "sample/middleware/envRouting"
	"sample/middleware/util"
	igatemodel "sample/model/igateModel"
	"sample/model/payload"
	"sample/model/webtool"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type (
	CallBack struct {
		ReferenceId string  `json:"referenceId"` // instruction ID
		Status      string  `json:"status"`      // FAILED or SUCCESS
		Amount      float64 `json:"amount"`      // additional request, this is for janus to know the amount needed for reversal
	}

	CTRequest struct {
		SenderBIC     string `json:"senderBIC"`
		ReceivingBIC  string `json:"receivingBIC"`
		InstructionId string `json:"instructionId"`
	}
)

// @Tags			IPS
// @Summary		Get Credit Transfer Transactions
// @Description	Get a list of credit transfer transactions
// @Produce		json
// @Success		200	{array}	payload.CreditTransferJSON
// @Router			/credit-transfer [get]

// Geromme
// CALLBACKS
func GetCreditTransferTransaction(c *fiber.Ctx) error {
	m := []payload.ResponseCreditTransfer{}
	db.Database.Raw("SELECT * FROM rbi_instapay.ct_transaction").Find(&m)
	return c.Status(200).JSON(&fiber.Map{
		"message": "data successfully fetch",
		"data":    m,
	})
}

func CallbackEndpoint(c *fiber.Ctx) error {
	CallbackResponseBody := CallBack{}
	if parsErr := c.BodyParser(&CallbackResponseBody); parsErr != nil {
		return c.JSON(fiber.Map{
			"retCode": 101,
			"error":   parsErr.Error(),
		})
	}

	jsonData, err := json.Marshal(CallbackResponseBody)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return err
	}

	ServiceEP := util.GetServiceEP("CreditCallback", strings.ToLower(envrouting.Environment))
	req, err := http.NewRequest(http.MethodPut, ServiceEP, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Merchant-D", "QVBJMDAwMDU=")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return err
	}

	defer resp.Body.Close()

	// Read the response body

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	result := &map[string]interface{}{}
	unmErr := json.Unmarshal(body, &result)
	if unmErr != nil {
		fmt.Println("Callback Error:", unmErr.Error())
	}

	fmt.Println("RESPONSE:", body)
	// loggers.CallbackLogs(c.Path(), CallbackResponseBody.ReferenceId, CallbackResponseBody.ReferenceId, jsonData, result)
	log.Println("End Callback Function")

	// Send the response from the PUT request back to the client
	return c.JSON(fiber.Map{
		"response": result,
	})
}

// For Janus
func CallbackFunction(c *fiber.Ctx, status, instructionId, process, rsnCode string, amount float64) string {
	fmt.Println("Start Callback Function")
	log.Println("Start Callback Function")
	var reference CallBack

	db.Database.Raw("SELECT reference_id FROM public.transactions WHERE instruction_id = ? ", instructionId).Scan(&reference.ReferenceId)

	if status == "RJCT" || rsnCode != "" {
		reference.Status = "FAILED"
	} else {
		reference.Status = "SUCCESS"
	}

	jsonData, err := json.Marshal(reference)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return err.Error()
	}

	fmt.Println("Credit Callback:\n", string(jsonData))
	// Create the request

	ServiceEP := util.GetServiceEP("CreditCallback", strings.ToLower(envrouting.Environment))
	req, err := http.NewRequest(http.MethodPut, ServiceEP, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err.Error()
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return err.Error()
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err.Error()
	}

	result := &map[string]interface{}{}
	unmErr := json.Unmarshal(body, &result)
	if unmErr != nil {
		fmt.Println("Callback Error:", unmErr.Error())
	}

	fmt.Println("Response:", resp.Body)
	// loggers.CallbackLogs(c.Path(), reference.ReferenceId, instructionId, jsonData, result)
	log.Println("End Callback Function")
	return reference.ReferenceId
}

// For Credit Transfer as Receiver
func CompleteRequestTransaction(c *fiber.Ctx, instructionID string) (bool, error) {
	fmt.Println("Start Transfer Credit - Receving")
	transactCredit := &payload.TransactCredit{}
	db.Database.Raw("SELECT * FROM public.transactions WHERE instruction_id = ?", instructionID).Scan(transactCredit)

	// Fetch settlement account and decrypt the data
	settlementAccount := &webtool.SettlementAccount{}
	db.Database.Debug().Raw("SELECT account_number FROM public.settlements WHERE event = 'receiving'").Scan(settlementAccount)
	decryptedAccountNumber, _ := encryptDecrypt.Decrypt(settlementAccount.AccountNumber, envrouting.Environment)
	amount, _ := strconv.ParseFloat(transactCredit.Amount, 64)

	transferCredit := &igatemodel.TransferCredit{
		ReferenceNumber: transactCredit.ReferenceId,
		CreditAccount:   transactCredit.ReceivingAccount,
		DebitAccount:    decryptedAccountNumber,
		Amount:          amount,
		Description:     fmt.Sprintf("%v %v", transactCredit.ReferenceId, "Instapay Receiving Fund Transfer"),
	}

	transferCreditRequirements, err := json.Marshal(transferCredit)
	if err != nil {
		fmt.Println("Error in JSON marshal:", err)
		return false, err
	}

	fmt.Println("Transfer Credit:", transferCreditRequirements)
	// This will get the endpoint from DB
	ServiceEP := util.GetServiceEP("CreditTransfer_igate", strings.ToLower(envrouting.Environment))

	client := &http.Client{}
	req, reqErr := http.NewRequest(http.MethodPost, ServiceEP, bytes.NewBuffer(transferCreditRequirements))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Merchant-ID", "QVBJMDAwMDU=")

	if reqErr != nil {
		fmt.Println("Error requesting:", err)
		return false, reqErr
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	resp, respErr := client.Do(req)
	if respErr != nil {
		fmt.Println("Getting error request:", err)
		return false, respErr
	}

	resultTransactCredit := &igatemodel.TransferCreditResponse{}
	decodErr := json.NewDecoder(resp.Body).Decode(resultTransactCredit)
	if decodErr != nil {
		return false, decodErr
	}

	defer resp.Body.Close()

	fmt.Println("---------------------------------------")
	fmt.Println("SERVICE EP:", ServiceEP)
	fmt.Println("TRANSFER CREDIT:", transferCredit)
	fmt.Println("CORE REFERENCE ID:", resultTransactCredit.CoreReference)
	fmt.Println("TRANSACTION REFERENCE ID:", resultTransactCredit.TransReference)
	fmt.Println("RECEIVING BIC:", transactCredit.ReceivingBIC)
	fmt.Println("RECEIVING NAME:", transactCredit.ReceivingName)
	fmt.Println("RECEIVING ACCOUNT:", transactCredit.ReceivingAccount)
	fmt.Println("SENDER BIC:", transactCredit.SenderBIC)
	fmt.Println("SENDER NAME:", transactCredit.SenderName)
	fmt.Println("SENDER ACCOUNT:", transactCredit.SenderAccount)
	fmt.Println("AMOUNT:", transactCredit.Amount)
	fmt.Println("INSTRUCTION ID:", transactCredit.InstructionId)
	fmt.Println("REFERENCE ID:", resultTransactCredit.ReferenceNumber)
	fmt.Println("RESPONSE:", resp.Body)
	fmt.Println("---------------------------------------")

	// loggers.TransactCredit(c.Path(), "igate", "Transfer_Credit_Receiving", ServiceEP, instructionID, transactCredit.ReferenceId, resultTransactCredit.CoreReference, transferCreditRequirements, resultTransactCredit, req)
	log.Println("End Transfer Credit")
	fmt.Println("End Transfer Credit")
	return true, nil
}

func GetInstructionID(c *fiber.Ctx) error {
	request := &CallBack{}
	if parsErr := c.BodyParser(request); parsErr != nil {
		return c.SendString(parsErr.Error())
	}

	instructionID := &CTRequest{}
	if dbErr := db.Database.Debug().Raw("SELECT * FROM public.transactions where reference_id=? ", request.ReferenceId).Scan(instructionID).Error; dbErr != nil {
		return c.SendString(dbErr.Error())
	}
	return c.JSON(fiber.Map{
		"data": instructionID,
	})
}
