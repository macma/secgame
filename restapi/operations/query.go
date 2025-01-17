// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// QueryHandlerFunc turns a function with the right signature into a query handler
type QueryHandlerFunc func(QueryParams) middleware.Responder

// Handle executing the request and returning a response
func (fn QueryHandlerFunc) Handle(params QueryParams) middleware.Responder {
	return fn(params)
}

// QueryHandler interface for that can handle valid query params
type QueryHandler interface {
	Handle(QueryParams) middleware.Responder
}

// NewQuery creates a new http.Handler for the query operation
func NewQuery(ctx *middleware.Context, handler QueryHandler) *Query {
	return &Query{Context: ctx, Handler: handler}
}

/*Query swagger:route POST /api/v1/query query

Query query API

*/
type Query struct {
	Context *middleware.Context
	Handler QueryHandler
}

func (o *Query) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewQueryParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// QueryOKBody query o k body
//
// swagger:model QueryOKBody
type QueryOKBody struct {

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this query o k body
func (o *QueryOKBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *QueryOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *QueryOKBody) UnmarshalBinary(b []byte) error {
	var res QueryOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
