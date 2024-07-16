package loginpage

import (
	"log"
	"time"

	"sample/ctransaction"
	"sample/db"
	"sample/models"

	"github.com/gofiber/fiber/v2"
)

func FtransactionHandler(c *fiber.Ctx) error {
	device := c.Context().UserAgent()
	requestCreditTransfer := &models.Ftransaction{}
	if parsErr := c.BodyParser(requestCreditTransfer); parsErr != nil {
		log.Println("End Credit Transfer")
		return c.JSON(fiber.Map{
			"message": "error parsing",
			"error":   parsErr.Error(),
		})
	}

	var returnCode string
	responseError := models.ResponseError{}
	time.Sleep(time.Second * 1)
	finalResponse := &models.ResponseFtransaction{}
	finalResponse = GetFinalReponseK2C(ctransaction.ReferenceId())
	requestCreditTransfer.Application = "RICO APPS"

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
func GetFinalReponseK2C(instructionId string) *models.ResponseFtransaction {
	creditTransferResponse := &models.ResponseFtransaction{}
	db.Database.Debug().Raw("SELECT transaction_type, status, reason_code, local_instrument, instruction_id, reference_id, sender_bic, sender_name, sender_account, currency, amount, receiving_bic, receiving_name, receiving_account, application FROM public.regist WHERE instruction_id = ?", instructionId).Scan(creditTransferResponse)
	return creditTransferResponse
}
