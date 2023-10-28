package api

import (
	"net/http"

	"github.com/aegaxs/terraform-provider-spinnaker/spinnaker/api/errors"
	"github.com/mitchellh/mapstructure"
	gate "github.com/spinnaker/spin/cmd/gateclient"
)

func CreatePipeline(client *gate.GatewayClient, pipeline interface{}) error {
	_, resp, err := retry(func() (map[string]interface{}, *http.Response, error) {
		resp, err := client.PipelineControllerApi.SavePipelineUsingPOST(client.Context, pipeline, nil)

		return nil, resp, err
	})
	if err != nil || (resp.StatusCode != http.StatusOK) {
		return errors.NewResponseError(resp, err)
	}

	return nil
}

func GetPipeline(client *gate.GatewayClient, applicationName, pipelineName string, dest interface{}) (map[string]interface{}, error) {
	payload, resp, err := retry(func() (map[string]interface{}, *http.Response, error) {
		return client.ApplicationControllerApi.GetPipelineConfigUsingGET(
			client.Context,
			applicationName,
			pipelineName,
		)
	})
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, errors.NewResponseError(resp, err)
	}

	if err := mapstructure.Decode(payload, dest); err != nil {
		return nil, err
	}

	return payload, nil
}

func UpdatePipeline(client *gate.GatewayClient, pipelineID string, pipeline interface{}) error {
	_, resp, err := retry(func() (map[string]interface{}, *http.Response, error) {
		return client.PipelineControllerApi.UpdatePipelineUsingPUT(client.Context, pipelineID, pipeline)
	})
	if err != nil || resp.StatusCode != http.StatusOK {
		return errors.NewResponseError(resp, err)
	}

	return nil
}

func DeletePipeline(client *gate.GatewayClient, applicationName, pipelineName string) error {
	_, resp, err := retry(func() (map[string]interface{}, *http.Response, error) {
		resp, err := client.PipelineControllerApi.DeletePipelineUsingDELETE(
			client.Context,
			applicationName,
			pipelineName,
		)
		return nil, resp, err
	})
	if err != nil || resp.StatusCode != http.StatusOK {
		return errors.NewResponseError(resp, err)
	}

	return nil
}

// RecreatePipeline is a convenience function for deleting and subsequently
// recreating a pipeline. It will return an error if either of the delete and
// create operations fails.
func RecreatePipeline(client *gate.GatewayClient, applicationName, pipelineName string, pipeline interface{}) error {
	err := DeletePipeline(client, applicationName, pipelineName)
	if err != nil {
		return err
	}

	return CreatePipeline(client, pipeline)
}
