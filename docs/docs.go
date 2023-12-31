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
        "/admin/categories": {
            "post": {
                "description": "Create a new category by providing the category details.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Create a new category",
                "operationId": "create-category",
                "parameters": [
                    {
                        "description": "Category details",
                        "name": "category",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.Category"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success\": \"Category added successfully\" entity.Category",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "error\": \"Invalid input\" entity.ErrorResponse",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/admin/categories/{id}": {
            "put": {
                "description": "Edit a category based on the provided JSON data",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "Edit a category",
                "operationId": "editCategory",
                "parameters": [
                    {
                        "type": "integer",
                        "format": "int64",
                        "description": "Category ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Category object to be edited",
                        "name": "category",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.Category"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success\": \"product edited successfully\", \"edited category\": entity.Category",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "error\": \"editing category failed",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete an existing category based on the provided ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "summary": "Delete a category",
                "operationId": "deleteCategory",
                "parameters": [
                    {
                        "type": "integer",
                        "format": "int64",
                        "description": "Category ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success\": \"Category deleted successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "error\": \"Failed to delete category",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/admin/login": {
            "post": {
                "description": "Authenticate admin using email and password and generate an authentication token.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Admin Login with Password",
                "operationId": "admin-login",
                "parameters": [
                    {
                        "description": "Admin Data",
                        "name": "admin",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entity.AdminLogin"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "message\": \"Admin logged in successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "error\": \"Empty request body",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/admin/products": {
            "get": {
                "description": "Retrieve a list of products for the admin dashboard.",
                "produces": [
                    "application/json"
                ],
                "summary": "Get a list of products for admin",
                "operationId": "get-admin-products",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.ProductWithQuantityResponse"
                            }
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/admin/users": {
            "get": {
                "description": "Get a paginated list of users.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "List Users",
                "operationId": "list-users",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Page number (default is 1)",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Number of users per page (default is 5)",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.ListUsersResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/entity.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/admin/users/toggle-permission/{id}": {
            "put": {
                "description": "Toggle the permission of a user by providing the user's ID.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Toggle User Permission",
                "operationId": "toggle-user-permission",
                "parameters": [
                    {
                        "minimum": 1,
                        "type": "integer",
                        "format": "int32",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success: User permission toggled successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "error: Invalid user ID",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "error: User not found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "entity.AdminLogin": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "entity.Category": {
            "type": "object",
            "required": [
                "description",
                "name"
            ],
            "properties": {
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "entity.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "entity.ListUsersResponse": {
            "type": "object",
            "properties": {
                "users": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.User"
                    }
                }
            }
        },
        "entity.User": {
            "type": "object",
            "required": [
                "email",
                "name",
                "password",
                "phone"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "minLength": 8
                },
                "phone": {
                    "type": "string"
                },
                "referalcode": {
                    "type": "string"
                },
                "wallet": {
                    "type": "integer"
                }
            }
        },
        "models.ProductWithQuantityResponse": {
            "type": "object",
            "properties": {
                "category": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "image_url": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "offerprice": {
                    "type": "integer"
                },
                "price": {
                    "type": "integer"
                },
                "quantity": {
                    "type": "integer"
                },
                "size": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "JWT": {
            "type": "apiKey",
            "name": "Authorise",
            "in": "cookie"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "lapify eCommerce API",
	Description:      "API for ecommerce website",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
