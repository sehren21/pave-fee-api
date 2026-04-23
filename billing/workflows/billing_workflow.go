package workflows

import (
	"time"

	"encore.app/billing/dto"
	"github.com/shopspring/decimal"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func BillingWorkflow(ctx workflow.Context, billID string) error {
	logger := workflow.GetLogger(ctx)

	total := decimal.Zero
	closed := false

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 10,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Second * 10,
			MaximumAttempts:    5,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	billItemCh := workflow.GetSignalChannel(ctx, "add-bill-item")
	closeCh := workflow.GetSignalChannel(ctx, "close-bill")

	for {
		selector := workflow.NewSelector(ctx)

		selector.AddReceive(billItemCh, func(c workflow.ReceiveChannel, more bool) {
			if closed {
				logger.Warn("Rejecting bill item, bill already closed")
				return
			}

			var item dto.AddBillItemRequest
			c.Receive(ctx, &item)

			amount, err := decimal.NewFromString(item.Amount)
			if err != nil {
				logger.Warn("Rejecting bill item, bill amount is invalid")
				return
			}

			total = total.Add(amount)

			err = workflow.ExecuteActivity(ctx, "AddBillItemActivity", item).Get(ctx, nil)
			if err != nil {
				logger.Error("Failed to add bill item", "error", err)
				return
			}
		})

		selector.AddReceive(closeCh, func(c workflow.ReceiveChannel, more bool) {
			closed = true

			err := workflow.ExecuteActivity(ctx, "CloseBillActivity", billID, total).Get(ctx, nil)
			if err != nil {
				logger.Error("Failed to close bill item", "error", err)
				return
			}

			logger.Info("Bill closed", "total", total)
		})

		selector.Select(ctx)

		if closed {
			break
		}
	}

	return nil
}
