package loginpage

import (
	"encoding/json"
	"fmt"
	"log"
	"sample/ctransaction"
	"sample/db"
	"sample/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TransferAccount(c *fiber.Ctx) error {
	device := c.Context().UserAgent()
	trans := &models.Ftransaction{}
	Ft := &models.CreditTransfer{}

	if err := c.BodyParser(trans); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	trans.Application = "K2C"
	Ft.SenderBIC = models.RBIBankIdentifierCode

	// Generate instruction ID and validate
	trans.InstructionId = ctransaction.GenerateInstructionID()
	if trans.InstructionId == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate instruction ID",
		})
	}

	var bic string
	var returnCode string
	responseError := models.ResponseError{}
	reqMar, _ := json.MarshalIndent(trans, "", " ")
	fmt.Println("INITIAL REQUEST:", string(reqMar))

	// Retrieve BIC based on RecipientBankCode
	// db.Database.Raw("SELECT bic FROM rbi_instapay.ips_participants WHERE bank_code = ?", trans.RecipientBankCode).Scan(&bic)
	Ft.ReceivingBIC = bic

	if bic == "" {
		returnCode = "04"
		db.Database.Raw("SELECT * FROM public.error_response WHERE code = ?", returnCode).Scan(&responseError)
		return c.JSON(models.Response{
			Device:      string(device),
			RetCode:     returnCode,
			Description: responseError.Description,
			Response: models.IPSResponse{
				InstructionID: trans.InstructionId,
			},
		})
	}

	if Ft.LocalInstrument == "" {
		Ft.LocalInstrument = "ICRT"
	}

	fmt.Println("----- `K2C Request Body` -----")
	fmt.Println("SENDER ACCOUNT NUMBER:", trans.SenderAccountNumber)
	fmt.Println("SENDER ACCOUNT NAME:", trans.SenderAccountName)
	fmt.Println("RECIPIENT BANK CODE:", trans.RecipientBankCode)
	fmt.Println("RECIPIENT ACCOUNT NUMBER:", trans.RecipientAccountNumber)
	fmt.Println("RECIPIENT ACCOUNT NAME:", trans.RecipientAccountName)
	fmt.Println("TRANSACTION REFERENCE:", trans.TransactionReference)
	fmt.Println("AMOUNT:", trans.TransactionAmount)
	fmt.Println("CHARGE:", trans.TransactionCharge)
	fmt.Println("----------------------------")

	transferRequestBody, err := json.MarshalIndent(Ft, "", " ")
	if err != nil {
		fmt.Println("Marshal: ", err)
		log.Println("End Credit Transfer")
		return c.JSON("Error")
	}

	if (trans.RecipientBankCode == "" || trans.RecipientAccountName == "" || trans.RecipientAccountNumber == "") && Ft.ReceivingBIC == "" {
		fmt.Println("K2C Fund Transfer Error!")
		fmt.Println("K2C ERROR:", trans)
		returnCode = "99"
		db.Database.Raw("SELECT description FROM public.error_response WHERE code = ?", returnCode).Scan(&responseError)
		return c.JSON(models.Response{
			Device:      string(device),
			RetCode:     returnCode,
			Description: responseError.Description,
			Response: models.IPSResponse{
				InstructionID: trans.InstructionId,
			},
		})
	}

	// Simulate sending request to FtransactionHandler
	// Assuming FtransactionHandler is implemented elsewhere
	c.Request().SetBody(transferRequestBody)
	FtransactionHandler(c)
	time.Sleep(time.Second * 1)

	finalResponse, err := GetFinalResponse(trans.InstructionId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	FtransactionHandler(c)

	// Will loop until the response code is not "NONE"
	if finalResponse.ReasonCode == "NONE" && (finalResponse.Status == "FAILED" || finalResponse.Status == "FAILED-RJCT") {
		ctr := 40
		for ctr > 0 {
			finalResponse, err = GetFinalResponse(trans.InstructionId)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			if finalResponse.Status == "SUCCESS" || finalResponse.Status == "FAILED-RJCT" {
				break
			}
			ctr--
			time.Sleep(time.Second * 1)
		}

		if ctr == 0 {
			returnCode = "06"
		} else if finalResponse.Status == "SUCCESS" {
			returnCode = "00"
		} else if finalResponse.Status == "FAILED-RJCT" {
			returnCode = "04"
		} else {
			returnCode = "06"
		}

		db.Database.Raw("SELECT * FROM public.error_response WHERE code = ?", returnCode).Scan(&responseError)
	}

	// Construct and return the response
	return c.JSON(models.Response{
		Device:      string(device),
		RetCode:     returnCode,
		Description: responseError.Description,
		Response: models.IPSResponse{
			InstructionID:    finalResponse.InstructionID,
			TransactionType:  finalResponse.TransactionType,
			Status:           finalResponse.Status,
			ReasonCode:       finalResponse.ReasonCode,
			Description:      finalResponse.Description,
			ReferenceId:      finalResponse.ReferenceID,
			SenderBIC:        finalResponse.SenderBIC,
			SenderName:       finalResponse.SenderName,
			SenderAccount:    finalResponse.SenderAccount,
			ReceivingBIC:     finalResponse.ReceivingBIC,
			ReceivingName:    finalResponse.ReceivingName,
			ReceivingAccount: finalResponse.ReceivingAccount,
		},
	})
}

func GetFinalResponse(instructionId string) (*models.ResponseFtransaction, error) {
	creditTransferResponse := &models.ResponseFtransaction{}
	result := db.Database.Raw("SELECT transaction_type, status, reason_code, local_instrument, instruction_id, reference_id, sender_bic, sender_name, sender_account, currency, amount, receiving_bic, receiving_name, receiving_account, application FROM public.transactions WHERE instruction_id = ?", instructionId).Scan(creditTransferResponse)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("no transaction found with instruction ID %s", instructionId)
	}
	return creditTransferResponse, nil
}
