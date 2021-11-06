// Code generated by go-swagger; DO NOT EDIT.

package blockchain_info

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// NewGetBlockParams creates a new GetBlockParams object
//
// There are no default values defined in the spec.
func NewGetBlockParams() GetBlockParams {

	return GetBlockParams{}
}

// GetBlockParams contains all the bound params for the get block operation
// typically these are obtained from a http.Request
//
// swagger:parameters getBlock
type GetBlockParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*
	  Required: true
	  In: path
	*/
	Hash string
	/*
	  Required: true
	  In: path
	*/
	NetCode string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetBlockParams() beforehand.
func (o *GetBlockParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rHash, rhkHash, _ := route.Params.GetOK("hash")
	if err := o.bindHash(rHash, rhkHash, route.Formats); err != nil {
		res = append(res, err)
	}

	rNetCode, rhkNetCode, _ := route.Params.GetOK("net_code")
	if err := o.bindNetCode(rNetCode, rhkNetCode, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindHash binds and validates parameter Hash from path.
func (o *GetBlockParams) bindHash(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route
	o.Hash = raw

	return nil
}

// bindNetCode binds and validates parameter NetCode from path.
func (o *GetBlockParams) bindNetCode(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route
	o.NetCode = raw

	if err := o.validateNetCode(formats); err != nil {
		return err
	}

	return nil
}

// validateNetCode carries on validations for parameter NetCode
func (o *GetBlockParams) validateNetCode(formats strfmt.Registry) error {

	if err := validate.EnumCase("net_code", "path", o.NetCode, []interface{}{"BTC", "LTC", "DOGE"}, true); err != nil {
		return err
	}

	return nil
}