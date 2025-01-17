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
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/createurl": {
            "post": {
                "description": "Create Short Url from Long or Original Url",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Urls"
                ],
                "summary": "Short Url",
                "parameters": [
                    {
                        "description": "this long or original url",
                        "name": "Url",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.UrlRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/urldb.Url"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/{durl}": {
            "delete": {
                "description": "Delete Url given its shorturl",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Urls"
                ],
                "summary": "Delete Url",
                "parameters": [
                    {
                        "type": "string",
                        "description": "url to delete",
                        "name": "durl",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/{surl}": {
            "get": {
                "description": "Redirect given short Url to original or long url",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Urls"
                ],
                "summary": "Redirect Url",
                "parameters": [
                    {
                        "type": "string",
                        "description": "url to redirect",
                        "name": "surl",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "302": {
                        "description": "Redirects to the long url",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.UrlRequest": {
            "description": "Request for creating shorturl",
            "type": "object",
            "required": [
                "url"
            ],
            "properties": {
                "url": {
                    "description": "url like \"http://www.example.com\", \"https://www.google.com\"",
                    "type": "string"
                }
            }
        },
        "urldb.Url": {
            "type": "object",
            "properties": {
                "key": {
                    "type": "string"
                },
                "longurl": {
                    "type": "string"
                },
                "shorturl": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Urlshortner API",
	Description:      "This is a Url Shortner api.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
