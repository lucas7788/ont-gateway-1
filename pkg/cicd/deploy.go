package cicd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/zhiqiangxu/ont-gateway/pkg/config"
	"github.com/zhiqiangxu/ont-gateway/pkg/forward"
	"github.com/zhiqiangxu/ont-gateway/pkg/instance"
	"github.com/zhiqiangxu/util"
	"go.uber.org/zap"
)

type (
	// DeployInput for input
	DeployInput struct {
		ID        string `json:"id"`
		Image     string `json:"image"`
		Path      string `json:"path"`
		Config    string `json:"config"`
		StatePath string `json:"state_path"`
	}
	// DeployOutput for output
	DeployOutput struct {
		URL string `json:"URL"`
	}
)

var (
	errDeployURLEmpty = errors.New("deploy url empty")
)

// Deploy specified docker image
// returns service ip
func Deploy(deploymentID, img, spec, path, conf, statePath string) (sip string, err error) {

	input := DeployInput{ID: deploymentID, Image: img, Path: path, Config: conf, StatePath: statePath}
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return
	}
	instance.Logger().Debug("Deploy", zap.String("url", config.Load().CICDConfig.AddonDeployAPI), zap.String("input", util.String(inputBytes)))
	code, _, outputBytes, err := forward.PostJSONRequest(config.Load().CICDConfig.AddonDeployAPI, inputBytes)
	if err != nil {
		instance.Logger().Error("PostJSONRequest", zap.String("json", util.String(inputBytes)), zap.Error(err))
		return
	}
	if code != 200 {
		instance.Logger().Error("deploy error", zap.Int("code", code), zap.String("body", string(outputBytes)), zap.Error(err))
		err = fmt.Errorf("PostJSONRequest not 200: %d", code)
		return
	}
	if err != nil {
		instance.Logger().Error("PostJSONRequest error", zap.Any("input", input))
		return
	}

	var output DeployOutput
	err = json.Unmarshal(outputBytes, &output)
	if err != nil {
		instance.Logger().Error("deploy error body", zap.String("body", string(outputBytes)), zap.Error(err))
		return
	}

	if output.URL == "" {
		instance.Logger().Error("deploy error empty url")
		err = errDeployURLEmpty
		return
	}

	sip = output.URL

	return
}
