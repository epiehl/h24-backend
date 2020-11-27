// Package web H24-Notifier API.
//
// This is the h24-notifier api
// It is written in golang
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//Schemes: http, https
//Host: localhost:3000
//BasePath: /api/v1
//Version: 0.0.2
//License: none
//Contact: Erik Piehl<erik.piehl93@gmail.com>
//
//Consumes:
//- application/json
//
//Produces:
//- application/json
//
//SecurityDefinitions:
//Bearer:
//  type: apiKey
//  in: header
//  name: Authorization
//Security:
//  - Bearer: []
//
// swagger:meta
package web

// APIVersion shows this API's current version
var APIVersion string
