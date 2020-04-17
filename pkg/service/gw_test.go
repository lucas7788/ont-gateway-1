package service

import (
	"testing"

	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"gotest.tools/assert"
)

func TestGateway(t *testing.T) {
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

			input := io.EnqueTxInput{App: 1, TxHash: txHash}
			output := gw.EnqueTx(input)
			assert.Assert(t, output.Code == 0)

		}

		{
			input := io.DequeTxInput{App: 1, TxHash: txHash}
			output := gw.DequeTx(input)
			assert.Assert(t, output.Code == 0 && output.Exists)
		}

	}

}
