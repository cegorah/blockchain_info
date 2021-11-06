package restapi

import (
	"github.com/cegorah/blockchain_info/restapi/operations"
	"github.com/go-openapi/loads"
	"net/http"
)

func getApi() (*operations.BlockchainInfoAPI, error) {
	swaggerSpec, e := loads.Analyzed(SwaggerJSON, "")
	if e != nil {
		return nil, e
	}
	return operations.NewBlockchainInfoAPI(swaggerSpec), nil
}

func GetApiHandler() (http.Handler, error) {
	api, e := getApi()
	if e != nil {
		return nil, e
	}
	h := configureAPI(api)
	e = api.Validate()
	if e != nil {
		return nil, e
	}
	return h, nil
}

