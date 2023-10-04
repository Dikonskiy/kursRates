// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/currency/save/{date}": {
            "post": {
                "description": "Save currency data for a specific date.",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "currency"
                ],
                "summary": "Save currency data",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Date in DD.MM.YYYY format",
                        "name": "date",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/currency/{date}": {
            "get": {
                "description": "Get currency data for a specific date.\nGet currency data for a specific date and currency code.",
                "consumes": [
                    "application/json",
                    "application/json"
                ],
                "tags": [
                    "currency",
                    "currency"
                ],
                "summary": "Get currency data by date and code",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Date in DD.MM.YYYY format",
                        "name": "date",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/currency/{date}/{code}": {
            "get": {
                "description": "Get currency data for a specific date.\nGet currency data for a specific date and currency code.",
                "consumes": [
                    "application/json",
                    "application/json"
                ],
                "tags": [
                    "currency",
                    "currency"
                ],
                "summary": "Get currency data by date and code",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Date in DD.MM.YYYY format",
                        "name": "date",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Currency code (e.g., USD)",
                        "name": "code",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/health": {
            "get": {
                "description": "Returns the health status of the application, including the database availability.",
                "produces": [
                    "application/json"
                ],
                "summary": "Check the health status of the application",
                "operationId": "health-check",
                "responses": {
                    "200": {
                        "description": "Status: Available",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "503": {
                        "description": "Status: Not available",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.1",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Swagger kursRates API",
	Description:      "A web service that, upon request, collects data from the public API of the national bank and saves the data to the local TEST database",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}