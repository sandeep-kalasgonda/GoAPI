// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/appointments": {
            "post": {
                "description": "Create a new appointment with customer details",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Create a new appointment",
                "parameters": [
                    {
                        "description": "Appointment Details",
                        "name": "appointment",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.Appointment"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/main.Appointment"
                        }
                    },
                    "400": {
                        "description": "Invalid JSON payload",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Failed to create appointment",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/appointments/all": {
            "get": {
                "description": "Get a list of all appointments",
                "produces": [
                    "application/json"
                ],
                "summary": "Get all appointments",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/main.Appointment"
                            }
                        }
                    },
                    "500": {
                        "description": "Failed to retrieve appointments",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/appointments/get": {
            "get": {
                "description": "Get an appointment by its ID",
                "produces": [
                    "application/json"
                ],
                "summary": "Get an appointment by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Appointment ID",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.Appointment"
                        }
                    },
                    "400": {
                        "description": "Invalid appointment ID",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Appointment not found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/appointments/update": {
            "put": {
                "description": "Update an existing appointment",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Update an appointment",
                "parameters": [
                    {
                        "description": "Appointment Details",
                        "name": "appointment",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.Appointment"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.Appointment"
                        }
                    },
                    "400": {
                        "description": "Invalid JSON payload",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Appointment not found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Failed to update appointment",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.Appointment": {
            "description": "This model is used to store appointment information",
            "type": "object",
            "properties": {
                "date_time": {
                    "type": "string"
                },
                "doctor": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "This model is used to store appointment information",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
