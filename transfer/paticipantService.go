package transfer

import (
	"encoding/json"
	"fmt"

	"sample/db"
	"sample/middleware/util"
	"sample/model/bah"
	"sample/model/payload"
	model "sample/model/payload"
	"time"

	"log"

	"github.com/gofiber/fiber/v2"
)

// Used for credit transfers and to identify the application using the instapay (KPLUS/K2C)
func ProcessCreditTransfer(c *fiber.Ctx) error {
	device := c.Context().UserAgent()
	// time.Sleep(30 * time.Second) // try DS24 error code
	log.Println("Start Credit Transfer")
	// now := time.Now().Add(time.Hour + 1)
	// nowDate := now.Format("20060102")

	// FOR KPLUS TRANSACTION
	requestCreditTransfer := &payload.FDSAPRequestCreditTransfer{}
	if parsErr := c.BodyParser(requestCreditTransfer); parsErr != nil {
		// return c.Status(503).SendString(parsErr.Error())
		log.Println("End Credit Transfer")
		return c.JSON(fiber.Map{
			"status":  101,
			"message": "error parsing",
			"error":   parsErr.Error(),
		})
	}

	requestCreditTransfer.Application = "KPLUS"
	// requestCreditTransfer.SendingBankIdentifierCode := "SRCPPHM2XXX"

	// GENERATION OF IDs FOR DIGEST VALUE AND SIGNATURE VALUE

	requestCreditTransfer.BusinessMessageId = util.GenerateMessageIdentifier("B", bah.RBIBankIdentifierCode, "B", 14)
	requestCreditTransfer.InstructionId = util.GenerateMessageIdentifier("NA", bah.RBIBankIdentifierCode, "B", 14)
	requestCreditTransfer.TransactionId = util.GenerateMessageIdentifier("TX", bah.RBIBankIdentifierCode, "B", 14)

	var returnCode string
	responseError := payload.ResponseError{}
	reqMar, _ := json.MarshalIndent(requestCreditTransfer, "", " ")
	fmt.Println("INITIAL REQUEST:", string(reqMar))

	// FOR K2C TRANSACTION PARAMETERS
	FDSIRequest := &model.FDSIRequestCreditTransfer{}
	if requestCreditTransfer.ReceivingBIC == "" && requestCreditTransfer.SenderBIC == "" {
		fmt.Println("K2C CREDIT TRANSFER SENDING")
		if parsErr := c.BodyParser(FDSIRequest); parsErr != nil {
			// return c.Status(503).SendString(parsErr.Error())
			log.Println("End Credit Transfer")
			return c.JSON(fiber.Map{
				"status":  101,
				"message": "error parsing",
				"error":   parsErr.Error(),
			})
		}

		requestCreditTransfer.Application = "K2C"
		requestCreditTransfer.SenderAccountNumber = FDSIRequest.SenderAccountNumber
		requestCreditTransfer.SenderBIC = bah.RBIBankIdentifierCode // "SRCPPHM2XXX" for workbook to get DS0H response code
		requestCreditTransfer.SenderName = FDSIRequest.SenderAccountName
		requestCreditTransfer.Amount = FDSIRequest.TransactionAmount
		requestCreditTransfer.Currency = "PHP" //"USD" for workbook to get 650 response code

		// need to get the BIC using the bank code for K2C request body
		var bic string
		db.Database.Raw("SELECT bic FROM public.ips_ips_participants WHERE bank_code = ?", FDSIRequest.RecipientBankCode).Scan(&bic)
		requestCreditTransfer.ReceivingBIC = bic
		fmt.Println("K2C - Code to BIC:", requestCreditTransfer.ReceivingBIC)
		// requestCreditTransfer.ReceivingBIC = "SRCPPHM2XXY" // for workbook to get RC04 response code
		// fmt.Println("K2C BIC:", requestCreditTransfer.ReceivingBIC)

		if bic == "" {
			returnCode = "04"
			db.Database.Raw("SELECT * FROM public.error_response WHERE code = ?", returnCode).Scan(&responseError)
			return c.JSON(payload.Response{
				Device:      string(device),
				RetCode:     returnCode,
				Description: responseError.Description,
				Response: payload.IPSResponse{
					InstructionID: requestCreditTransfer.InstructionId,
				},
			})
		}

		// requestCreditTransfer.ReceivingBIC = "SRCPPHM2XXY" FDSIRequest.RecipientAccountNumber // for workbook to get RC04 response code
		requestCreditTransfer.ReceivingAccountNumber = FDSIRequest.RecipientAccountNumber
		requestCreditTransfer.ReceivingName = FDSIRequest.RecipientAccountName

		// Get reference id from janus using the generated instruction id
		// referenceId := GetIFTJnusK2C(requestCreditTransfer.InstructionId)
		// if referenceId == "" {
		// 	fmt.Println("K2C Fund Transfer Error!")
		// 	returnCode = "99"
		// 	database.DBConn.Raw("SELECT description FROM knowledge_base.error_response WHERE code = ?", returnCode).Scan(&responseError)
		// 	return c.JSON(payload.Response{
		// 		Device:      string(device),
		// 		RetCode:     returnCode,
		// 		Description: responseError.Description,
		// 		Response: payload.IPSResponse{
		// 			InstructionID: requestCreditTransfer.InstructionId,
		// 		},
		// 	})
		// }

		// requestCreditTransfer.ReferenceID = referenceId

		if requestCreditTransfer.LocalInstrument == "" {
			requestCreditTransfer.LocalInstrument = "ICRT"
		}

		fmt.Println("----- `K2C Request Body` -----")
		fmt.Println("SENDER ACCOUNT NUMBER:", FDSIRequest.SenderAccountNumber)
		fmt.Println("SENDER ACCOUNT NAME:", FDSIRequest.SenderAccountName)
		fmt.Println("RECIPIENT BANK CODE:", FDSIRequest.RecipientBankCode)
		fmt.Println("RECIPIENT ACCOUNT NUMBER:", FDSIRequest.RecipientAccountNumber)
		fmt.Println("RECIPIENT ACCOUNT NAME:", FDSIRequest.RecipientAccountName)
		fmt.Println("TRANSACTION REFERENCE:", FDSIRequest.TransactionReference)
		fmt.Println("AMOUNT:", FDSIRequest.TransactionAmount)
		fmt.Println("CHARGE:", FDSIRequest.TransactionCharge)
		fmt.Println("----------------------------")
	}

	// Transfer this request to CreditTransfer
	transferRequestBody, err := json.MarshalIndent(requestCreditTransfer, "", " ")
	// fmt.Println("TRANSFER REQUEST BODY:", string(transferRequestBody))
	if err != nil {
		fmt.Println("Marshal: ", err)
		log.Println("End Credit Transfer")
		return c.JSON("Error")
	}

	if (FDSIRequest.RecipientBankCode == "" || FDSIRequest.RecipientAccountName == "" || FDSIRequest.RecipientAccountNumber == "") && requestCreditTransfer.ReceivingBIC == "" {
		fmt.Println("K2C Fund Transfer Error!")
		fmt.Println("K2C ERROR:", FDSIRequest)
		returnCode = "99"
		db.Database.Raw("SELECT description FROM public.error_response WHERE code = ?", returnCode).Scan(&responseError)
		return c.JSON(payload.Response{
			Device:      string(device),
			RetCode:     returnCode,
			Description: responseError.Description,
			Response: payload.IPSResponse{
				InstructionID: requestCreditTransfer.InstructionId,
			},
		})
	} else {

		c.Request().SetBody(transferRequestBody)

		// Calls credit transfer function to generate XML message that will be send to BancNet service
		// CreditTransferK2C(c)

		time.Sleep(time.Second * 1)
		finalResponse := &payload.ResponseCreditTransfer{}
		finalResponse = GetFinalReponseK2C(requestCreditTransfer.InstructionId)

		fmt.Println("RESPONSE CODE:", finalResponse.ReasonCode)
		// If PACS.008 structure is invalid
		if finalResponse.ReasonCode == "650" {
			db.Database.Debug().Raw("SELECT description FROM public.error_response WHERE code = ?", finalResponse.ReasonCode).Scan(&finalResponse.Description)
			returnCode = "04"
			db.Database.Debug().Raw("SELECT * FROM public.error_response WHERE code = ?", returnCode).Scan(&responseError)
			return c.JSON(payload.Response{
				Device:      string(device),
				RetCode:     returnCode,
				Description: responseError.Description,
				Response: payload.IPSResponse{
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

		// Will loop until the response code is not "NONE"
		if finalResponse.ReasonCode == "NONE" && (finalResponse.Status == "FAILED" || finalResponse.Status == "FAILED-RJCT") {
			// fmt.Printf("Waiting for PACS.002 in %d ...\n", ctr)
			ctr := 40
			for ctr > 0 {
				// fmt.Printf("%d - Waiting for PACS.002... \n", ctr)
				finalResponse = GetFinalReponseK2C(requestCreditTransfer.InstructionId)
				// fmt.Println("STATUS:", finalResponse.Status)
				if finalResponse.Status == "SUCCESS" || finalResponse.Status == "FAILED-RJCT" {
					break
				}
				ctr--
				time.Sleep(time.Second * 1)
			}

			// When pacs.002 waiting time to receive exceeds to 30secs
			if ctr == 0 {
				returnCode = "06"
				db.Database.Raw("SELECT * FROM public.error_response WHERE code = ?", returnCode).Scan(&responseError)
				return c.JSON(payload.Response{
					Device:      string(device),
					RetCode:     returnCode,
					Description: responseError.Description,
					Response: payload.IPSResponse{
						InstructionID: finalResponse.InstructionID,
					},
				})
			} else if finalResponse.Status == "SUCCESS" {
				returnCode = "00"
				db.Database.Raw("SELECT * FROM public.error_response WHERE code = ?", returnCode).Scan(&responseError)
			} else if finalResponse.Status == "FAILED-RJCT" {
				returnCode = "04"
				db.Database.Debug().Raw("SELECT * FROM public.error_response WHERE code = ?", returnCode).Scan(&responseError)
			} else { //exceeds to waiting time for pacs.002
				returnCode = "06"
				db.Database.Debug().Raw("SELECT * FROM public.error_response WHERE code = ?", returnCode).Scan(&responseError)
			}
		}

		// Get the description for bancnet reason code
		db.Database.Debug().Raw("SELECT description FROM public.error_response WHERE code = ?", finalResponse.ReasonCode).Scan(&finalResponse.Description)
		fmt.Println("reason code description:", finalResponse.Description)

		return c.JSON(payload.Response{
			Device:      string(device),
			RetCode:     returnCode,
			Description: responseError.Description,
			Response: payload.IPSResponse{
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
}

func GetFinalReponseK2C(instructionId string) *payload.ResponseCreditTransfer {
	creditTransferResponse := &payload.ResponseCreditTransfer{}
	db.Database.Raw("SELECT transaction_type, status, reason_code, local_instrument, instruction_id, reference_id, sender_bic, sender_name, sender_account, currency, amount, receiving_bic, receiving_name, receiving_account, application FROM public.ct_transaction WHERE instruction_id = ?", instructionId).Scan(creditTransferResponse)
	return creditTransferResponse
}
