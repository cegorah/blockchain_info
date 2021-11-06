// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "title": "blockchain_info",
    "version": "1.0.0"
  },
  "host": "localhost:8080",
  "basePath": "/v1",
  "paths": {
    "/block/{net_code}/{hash}": {
      "get": {
        "security": [
          {
            "Bearer": []
          }
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "blockchain_info"
        ],
        "operationId": "getBlock",
        "parameters": [
          {
            "type": "string",
            "name": "hash",
            "in": "path",
            "required": true
          },
          {
            "enum": [
              "BTC",
              "LTC",
              "DOGE"
            ],
            "type": "string",
            "name": "net_code",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Returns an information about a block with a hash",
            "schema": {
              "$ref": "#/definitions/BlockInfo"
            }
          },
          "400": {
            "description": "Wrong request parameters"
          },
          "404": {
            "description": "Block not found"
          },
          "500": {
            "description": "Internal server error"
          }
        }
      }
    },
    "/tx/{net_code}/{hash}": {
      "get": {
        "security": [
          {
            "Bearer": []
          }
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "blockchain_info"
        ],
        "operationId": "getTxInfo",
        "parameters": [
          {
            "type": "string",
            "name": "hash",
            "in": "path",
            "required": true
          },
          {
            "enum": [
              "BTC",
              "LTC",
              "DOGE"
            ],
            "type": "string",
            "name": "net_code",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Returns an information about a transaction",
            "schema": {
              "$ref": "#/definitions/Transactions"
            }
          },
          "400": {
            "description": "Wrong request parameters"
          },
          "404": {
            "description": "Transaction not found"
          },
          "500": {
            "description": "Internal server error"
          }
        }
      }
    }
  },
  "definitions": {
    "BlockInfo": {
      "type": "object",
      "properties": {
        "first_ten_transactions": {
          "type": "array",
          "maxLength": 10,
          "items": {
            "$ref": "#/definitions/TransactionInfo"
          }
        },
        "id": {
          "type": "integer"
        },
        "net_code": {
          "type": "string"
        },
        "next_hash": {
          "type": "string"
        },
        "prev_hash": {
          "type": "string"
        },
        "size": {
          "type": "integer"
        },
        "timestamp": {
          "type": "string",
          "format": "date"
        }
      }
    },
    "TransactionInfo": {
      "type": "object",
      "properties": {
        "fee": {
          "type": "number"
        },
        "id": {
          "type": "integer"
        },
        "sent_value": {
          "type": "integer"
        },
        "timestamp": {
          "type": "string",
          "format": "date"
        }
      }
    },
    "Transactions": {
      "type": "array",
      "items": {
        "$ref": "#/definitions/TransactionInfo"
      }
    }
  },
  "securityDefinitions": {
    "Bearer": {
      "type": "apiKey",
      "name": "Authorization",
      "in": "header"
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "title": "blockchain_info",
    "version": "1.0.0"
  },
  "host": "localhost:8080",
  "basePath": "/v1",
  "paths": {
    "/block/{net_code}/{hash}": {
      "get": {
        "security": [
          {
            "Bearer": []
          }
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "blockchain_info"
        ],
        "operationId": "getBlock",
        "parameters": [
          {
            "type": "string",
            "name": "hash",
            "in": "path",
            "required": true
          },
          {
            "enum": [
              "BTC",
              "LTC",
              "DOGE"
            ],
            "type": "string",
            "name": "net_code",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Returns an information about a block with a hash",
            "schema": {
              "$ref": "#/definitions/BlockInfo"
            }
          },
          "400": {
            "description": "Wrong request parameters"
          },
          "404": {
            "description": "Block not found"
          },
          "500": {
            "description": "Internal server error"
          }
        }
      }
    },
    "/tx/{net_code}/{hash}": {
      "get": {
        "security": [
          {
            "Bearer": []
          }
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "blockchain_info"
        ],
        "operationId": "getTxInfo",
        "parameters": [
          {
            "type": "string",
            "name": "hash",
            "in": "path",
            "required": true
          },
          {
            "enum": [
              "BTC",
              "LTC",
              "DOGE"
            ],
            "type": "string",
            "name": "net_code",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Returns an information about a transaction",
            "schema": {
              "$ref": "#/definitions/Transactions"
            }
          },
          "400": {
            "description": "Wrong request parameters"
          },
          "404": {
            "description": "Transaction not found"
          },
          "500": {
            "description": "Internal server error"
          }
        }
      }
    }
  },
  "definitions": {
    "BlockInfo": {
      "type": "object",
      "properties": {
        "first_ten_transactions": {
          "type": "array",
          "maxLength": 10,
          "items": {
            "$ref": "#/definitions/TransactionInfo"
          }
        },
        "id": {
          "type": "integer"
        },
        "net_code": {
          "type": "string"
        },
        "next_hash": {
          "type": "string"
        },
        "prev_hash": {
          "type": "string"
        },
        "size": {
          "type": "integer"
        },
        "timestamp": {
          "type": "string",
          "format": "date"
        }
      }
    },
    "TransactionInfo": {
      "type": "object",
      "properties": {
        "fee": {
          "type": "number"
        },
        "id": {
          "type": "integer"
        },
        "sent_value": {
          "type": "integer"
        },
        "timestamp": {
          "type": "string",
          "format": "date"
        }
      }
    },
    "Transactions": {
      "type": "array",
      "items": {
        "$ref": "#/definitions/TransactionInfo"
      }
    }
  },
  "securityDefinitions": {
    "Bearer": {
      "type": "apiKey",
      "name": "Authorization",
      "in": "header"
    }
  }
}`))
}