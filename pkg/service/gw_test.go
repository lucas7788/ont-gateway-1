package service

import (
	"testing"

	"github.com/zhiqiangxu/ont-gateway/pkg/io"
	"gotest.tools/assert"
)

func TestGateway(t *testing.T) {
	gw := Instance()

	{
		input := io.PostAddonConfigInput{
			AddonID:  "addon_id",
			TenantID: "tenant_id",
			Config:   "config",
		}
		output := gw.PostAddonConfig(input)
		assert.Assert(t, output.Code == 0)
	}

}
