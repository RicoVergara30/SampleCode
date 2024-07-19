package loginpage

import (
	"sample/ctransaction"
	"sample/models"

	"github.com/gofiber/fiber/v2"
)

func FtransactionHandler(c *fiber.Ctx) error {
	trans := &models.Transbody{}
	if parsErr := c.BodyParser(trans); parsErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": parsErr.Error(),
		})
	}

	// requestTrigger := time.Now().Format("2006-01-02 03:04:05")

	response := &models.Transrequest{
		InstructionID:   ctransaction.GenerateInstructionID(),
		ReferenceId:     ctransaction.Iftgenerate(),
		TransactionType: "RECEIVING",
		SenderBIC:       "CAMZPHM2XXX",
		ReceivingBIC:    "CBMFPHM1XXX",
		AmountCurrency:  "PHP",
		// LocalInstrument: "ICRT",
		Description: "Invalid transaction amount",
		ReasonCode:  "AC01",
		Status:      "FAILED-RJCT",
	}

	if trans.TransactionAmount > 100 {
		response.TransactionType = "SENDING"
		response.ReasonCode = "ACTC"
		response.Status = "SUCCESS"
		response.Description = "Transaction processed successfully"
	}

	// insertQuery := `INSERT INTO trace_alerts.credit_transfer( instruction_id, amount_currency, description,  local_instrument, reason_code, receiving_account, receiving_bic, receiving_name, reference_id, sender_account, sender_amount, sender_bic, sender_name, status, transaction_type, date_time)
	//     VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	// transacSave := db.Database.Debug().Exec(insertQuery,
	// 	instructionId, response.AmountCurrency, response.Description, response.LocalInstrument, response.ReasonCode,
	// 	trans.Recipientaccountnumber, response.ReceivingBic, trans.Recipientaccountname, referenceId,
	// 	trans.Senderaccountnumber, trans.Transactionamount, response.SenderBic, trans.Senderaccountname,
	// 	response.Status, response.TransactionType, requestTrigger).Error

	// if transacSave != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal Server Error"})
	// }

	return c.JSON(response)
}
