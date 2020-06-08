package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/zhiqiangxu/ddxf"
	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"github.com/zhiqiangxu/ont-gateway/pkg/model"
	"github.com/zhiqiangxu/util"
	"gotest.tools/assert"
)

func TestGateway(t *testing.T) {

	id, err := model.AppManager().GetMaxAppIDFromDB()

	fmt.Println("id", id, "err", err)

	gw := Instance()

	// addon
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

	// shell
	{
		input := io.ShellInput{
			Shell: "echo -n 43",
		}
		output := gw.Shell(input)
		assert.Assert(t, output.Out == "43", output)
	}

	// tx
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

	// payment
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
				CoinType:        model.CoinTypeONG,
				PayMethods:      []model.PayMethod{model.PayMethodBeforeUse}}
			output := gw.CreatePaymentConfig(input)
			assert.Assert(t, output.Code == 0)
		}
		{
			input := io.CreatePaymentOrderInput{
				PaymentConfigID: paymentConfigID,
				PaymentID:       paymentID,
				OrderID:         orderID,
				PaymentInfo: &io.PaymentInfo{
					PayPeriod: model.PayPeriodMonthly,
					PayMethod: model.PayMethodBeforeUse,
				},
				Amount:   150,
				CoinType: model.CoinTypeONG,
			}
			output := gw.CreatePaymentOrder(input)
			assert.Assert(t, output.Code == 0 && output.Balance == 50)
		}

		{
			input := io.GetPaymentInfoInput{PaymentID: paymentID, WithOrders: true}
			output := gw.GetPaymentInfo(input)
			assert.Assert(t,
				output.Code == 0 &&
					output.Payment.Balance == 50 && output.Payment.BalanceExpireTime.After(time.Now()) &&
					output.PaymentConfig != nil &&
					len(output.PaymentOrders) == 1)
		}

		{
			n, err := model.PaymentOrderManager().DeletePaymentOrders(0, paymentID)
			assert.Assert(t, n == 1 && err == nil)
			exists, err := model.PaymentManager().DeleteOne(0, paymentID)
			assert.Assert(t, exists && err == nil)
			exists, err = model.PaymentConfigManager().DeleteOne(0, paymentConfigID)
			assert.Assert(t, exists && err == nil)
		}

	}

	// wallet
	walletName := "testw"
	walletContent := "text content"
	{
		input := io.ImportWalletInput{
			WalletName: walletName,
			Content:    walletContent,
		}

		output := gw.ImportWallet(input)
		assert.Assert(t, output.Error() == nil, output)
	}
	{
		input := io.GetWalletInput{
			WalletName: walletName,
		}

		output := gw.GetWallet(input)
		assert.Assert(t, output.Exists && output.Content == walletContent)

		{
			output := gw.DeleteWallet(io.DeleteWalletInput{WalletName: walletName})
			assert.Assert(t, output.Error() == nil)
		}

	}

	{
		output := gw.CreateWallet(io.CreateWalletInput{WalletName: "testxxx"})
		assert.Assert(t, output.Error() == nil && output.Content != "")
	}

	{
		id := "testID"
		desc := "test desc"
		block := uint32(1)
		h := ddxf.Sha256Bytes(util.Slice(desc))
		hash := string(h[:])
		input := io.UpdateResourceInput{
			RV: model.ResourceVersion{ID: id, Block: block, Desc: desc, Hash: hash},
		}
		output := gw.UpdateResource(input)
		assert.Assert(t, output.Error() == nil)

		output = gw.UpdateResource(input)
		assert.Assert(t, output.Error() != nil)

		input.Force = true
		output = gw.UpdateResource(input)
		assert.Assert(t, output.Error() == nil && output.Exists)

		{
			output := gw.GetResource(io.GetResourceInput{ID: id, Block: block})
			assert.Assert(t, output.Error() == nil && output.Desc == desc && output.DescHash == hash)

			output = gw.GetResource(io.GetResourceInput{ID: id, Hash: hash})
			assert.Assert(t, output.Error() == nil && output.Desc == desc && output.DescHash == hash)
		}

		n, err := model.ResourceVersionManager().DeleteResourceByID(0, id)
		assert.Assert(t, err == nil && n == 1)
	}

}
