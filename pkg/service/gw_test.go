package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
	"gotest.tools/assert"
)

func TestGateway(t *testing.T) {

	id, err := model.AppManager().GetMaxAppIDFromDB()

	fmt.Println("id", id, "err", err)

	gw := Instance()

	{
		input := io.UpsertAddonConfigInput{
			AddonID:  "addon_id",
			TenantID: "tenant_id",
			Config:   "config",
		}
		output := gw.UpsertAddonConfig(input)
		assert.Assert(t, output.Code == 0)
	}

	{
		input := io.DeleteAddonConfigInput{
			AddonID:  "addon_id",
			TenantID: "tenant_id",
		}
		output := gw.DeleteAddonConfig(input)
		assert.Assert(t, output.Code == 0)
	}

	{
		input := io.ShellInput{
			Shell: "echo -n 43",
		}
		output := gw.Shell(input)
		assert.Assert(t, output.Out == "43", output)
	}

	{
		txHash := "txh123"
		{

			input := io.EnqueTxInput{TxHash: txHash, Admin: true}
			output := gw.EnqueTx(input)
			assert.Assert(t, output.Code == 0)

		}

		{
			input := io.DequeTxInput{TxHash: txHash, Admin: true}
			output := gw.DequeTx(input)
			assert.Assert(t, output.Code == 0 && output.Exists)
		}
	}

	{
		test = true

		paymentConfigID := "test_PaymentConfigID"
		paymentID := "test_PaymentID"
		orderID := "test_OrderID"
		{
			input := io.CreatePaymentConfigInput{
				PaymentConfigID: paymentConfigID,
				AmountOptions:   []int{100, 200},
				PeriodOptions:   []model.PayPeriod{model.PayPeriodMonthly, model.PayPeriodSeasonly},
				CoinType:        model.CoinTypeONG}
			output := gw.CreatePaymentConfig(input)
			assert.Assert(t, output.Code == 0)
		}
		{
			input := io.CreatePaymentOrderInput{
				PaymentConfigID: paymentConfigID,
				PaymentID:       paymentID,
				OrderID:         orderID,
				PayPeriod:       model.PayPeriodMonthly,
				PayMethod:       model.PayMethodBeforeUse,
				Amount:          150,
				CoinType:        model.CoinTypeONG,
			}
			output := gw.CreatePaymentOrder(input)
			assert.Assert(t, output.Code == 0 && output.Balance == 50)
		}

		{
			input := io.GetPaymentInfoInput{PaymentID: paymentID}
			output := gw.GetPaymentInfo(input)
			assert.Assert(t,
				output.Code == 0 &&
					output.Payment.Balance == 50 && output.Payment.BalanceExpireTime.After(time.Now()) &&
					output.PaymentConfig != nil &&
					len(output.PaymentOrders) == 1)
		}

		// {
		// 	n, err := model.PaymentOrderManager().DeletePaymentOrders(0, paymentID)
		// 	assert.Assert(t, n == 1 && err == nil)
		// 	exists, err := model.PaymentManager().DeleteOne(0, paymentID)
		// 	assert.Assert(t, exists && err == nil)
		// 	exists, err = model.PaymentConfigManager().DeleteOne(0, paymentConfigID)
		// 	assert.Assert(t, exists && err == nil)
		// }

	}

}
