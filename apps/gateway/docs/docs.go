package docs

import "github.com/swaggo/swag"

const docTemplate = `{
  "swagger": "2.0",
  "info": {
    "description": "API Gateway для проекта Лавка",
    "title": "Лавка Gateway API",
    "termsOfService": "http://swagger.io/terms/",
    "contact": {
      "name": "API Support",
      "url": "https://github.com/ultard/fusion"
    },
    "license": {
      "name": "MIT",
      "url": "https://opensource.org/licenses/MIT"
    },
    "version": "1.0"
  },
  "host": "localhost:8080",
  "basePath": "/",
  "paths": {
    "/auth/register": {
      "post": {
        "description": "Register a new user and send verification email",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "Auth service"
        ],
        "summary": "Register a new user",
        "parameters": [
          {
            "description": "Registration Information",
            "name": "request",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/auth.RegisterRequest"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Registration successful",
            "schema": {
              "$ref": "#/definitions/auth.RegisterResponse"
            }
          },
          "500": {
            "description": "Internal Server Error",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          },
          "503": {
            "description": "Service Unavailable",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      }
    },
    "/auth/confirmemail": {
      "get": {
        "description": "Confirm user email using verification token",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "Auth service"
        ],
        "summary": "Confirm user email",
        "parameters": [
          {
            "type": "string",
            "description": "Verification Token",
            "name": "token",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "201": {
            "description": "Email verified successfully",
            "schema": {
              "$ref": "#/definitions/auth.ConfirmEmailResponse"
            }
          },
          "500": {
            "description": "Internal Server Error",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          },
          "503": {
            "description": "Service Unavailable",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      }
    },
    "/auth/login": {
      "post": {
        "description": "Authenticate a user and return JWT tokens",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "Auth service"
        ],
        "summary": "Login user",
        "parameters": [
          {
            "description": "Login Credentials",
            "name": "request",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/auth.LoginRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Login successful",
            "schema": {
              "$ref": "#/definitions/auth.LoginResponse"
            }
          },
          "500": {
            "description": "Internal Server Error",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          },
          "503": {
            "description": "Service Unavailable",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      }
    },
    "/auth/ssologin": {
      "post": {
        "description": "Login using Single Sign-On token from trusted provider",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "Auth service"
        ],
        "summary": "Login with SSO (In progress)",
        "parameters": [
          {
            "description": "SSO Token",
            "name": "request",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/auth.SSOLoginRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/auth.LoginResponse"
            }
          },
          "500": {
            "description": "Internal Server Error",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          },
          "503": {
            "description": "Service Unavailable",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      }
    },
    "/auth/refresh": {
      "put": {
        "description": "Refresh access and refresh tokens",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "Auth service"
        ],
        "summary": "Refresh tokens",
        "parameters": [
          {
            "description": "Refresh Token",
            "name": "request",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/auth.RefreshRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "New tokens",
            "schema": {
              "$ref": "#/definitions/auth.RefreshResponse"
            }
          },
          "500": {
            "description": "Internal Server Error",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          },
          "503": {
            "description": "Service Unavailable",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      }
    },
    "/auth/forgotpassword": {
      "get": {
        "description": "Send password reset instructions to user's email",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "Auth service"
        ],
        "summary": "Request password reset",
        "parameters": [
          {
            "description": "Email Information",
            "name": "request",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/auth.ForgotPasswordRequest"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Reset instructions sent",
            "schema": {
              "$ref": "#/definitions/auth.ForgotPasswordResponse"
            }
          },
          "500": {
            "description": "Internal Server Error",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          },
          "503": {
            "description": "Service Unavailable",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      }
    },
    "/auth/logout": {
      "delete": {
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "description": "Invalidate user session and tokens",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "Auth service"
        ],
        "summary": "Log out user (In progress)",
        "parameters": [
          {
            "description": "Refresh Token",
            "name": "request",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/auth.LogoutRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Logout successful",
            "schema": {
              "$ref": "#/definitions/auth.LogoutResponse"
            }
          },
          "500": {
            "description": "Internal Server Error",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          },
          "503": {
            "description": "Service Unavailable",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      }
    },
    "/auth/resetpassword": {
      "put": {
        "description": "Reset user password with token",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "Auth service"
        ],
        "summary": "Reset password (In progress)",
        "parameters": [
          {
            "description": "Reset Password Data",
            "name": "request",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/auth.ResetPasswordRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Password reset successful",
            "schema": {
              "$ref": "#/definitions/auth.ResetPasswordResponse"
            }
          },
          "500": {
            "description": "Internal Server Error",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          },
          "503": {
            "description": "Service Unavailable",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      }
    },
    "/auth/me": {
      "get": {
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "description": "Get information about the authenticated client",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "Auth service"
        ],
        "summary": "Get client information",
        "responses": {
          "200": {
            "description": "OK",
            "schema": {
              "$ref": "#/definitions/auth.ClientInfoResponse"
            }
          },
          "500": {
            "description": "Internal Server Error",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          },
          "503": {
            "description": "Service Unavailable",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      }
    }
  },
  "securityDefinitions": {
    "BearerAuth": {
      "type": "apiKey",
      "name": "Authorization",
      "in": "header",
      "description": "Type \"Bearer\" followed by a space and JWT token."
    }
  },
  "definitions": {
	"ErrorResponse": {
		  "type": "object",
		  "properties": {
			"message": {
			  "type": "string",
			  "example": "error text"
			}
		  }
		},
    "auth.ClientInfoResponse": {
      "type": "object",
      "properties": {
        "avatar": {
          "type": "string"
        },
        "email": {
          "type": "string"
        },
        "message": {
          "type": "string"
        },
        "phone": {
          "type": "string"
        },
        "username": {
          "type": "string"
        }
      }
    },
    "auth.ConfirmEmailRequest": {
      "type": "object",
      "properties": {
        "token": {
          "type": "string"
        }
      }
    },
    "auth.ConfirmEmailResponse": {
      "type": "object",
      "properties": {
        "message": {
          "type": "string",
          "example": "Email verified successfully."
        }
      }
    },
    "auth.ForgotPasswordRequest": {
      "type": "object",
      "required": [
        "email",
        "redirectUrl"
      ],
      "properties": {
        "email": {
          "type": "string",
          "example": "user@example.com"
        },
        "redirectUrl": {
          "type": "string",
          "example": "https://fusion.com/reset-password"
        }
      }
    },
    "auth.ForgotPasswordResponse": {
      "type": "object",
      "properties": {
        "message": {
          "type": "string",
          "example": "If your email is registered, you'll receive a password reset link"
        }
      }
    },
    "auth.LoginRequest": {
      "type": "object",
      "required": [
        "login",
        "password"
      ],
      "properties": {
        "login": {
          "type": "string",
          "example": "user@example.com"
        },
        "password": {
          "type": "string",
          "minLength": 8,
          "example": "securePassword123"
        }
      }
    },
    "auth.LoginResponse": {
      "type": "object",
      "properties": {
        "accessToken": {
          "type": "string",
          "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
        },
        "message": {
          "type": "string",
          "example": "Login successful"
        },
        "refreshToken": {
          "type": "string",
          "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
        }
      }
    },
    "auth.LogoutRequest": {
      "type": "object",
      "required": [
        "refreshToken"
      ],
      "properties": {
        "refreshToken": {
          "type": "string",
          "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
        }
      }
    },
    "auth.LogoutResponse": {
      "type": "object",
      "properties": {
        "message": {
          "type": "string",
          "example": "Successfully logged out"
        }
      }
    },
    "auth.RefreshRequest": {
      "type": "object",
      "required": [
        "token"
      ],
      "properties": {
        "token": {
          "type": "string",
          "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
        }
      }
    },
    "auth.RefreshResponse": {
      "type": "object",
      "properties": {
        "accessToken": {
          "type": "string",
          "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
        },
        "message": {
          "type": "string",
          "example": "Tokens refreshed successfully"
        },
        "refreshToken": {
          "type": "string",
          "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
        }
      }
    },
    "auth.RegisterRequest": {
      "type": "object",
      "required": [
        "email",
        "password",
        "redirectUrl",
        "username"
      ],
      "properties": {
        "email": {
          "type": "string",
          "example": "user@example.com"
        },
        "password": {
          "type": "string",
          "minLength": 6,
          "example": "securePassword123"
        },
        "redirectUrl": {
          "type": "string",
          "example": "https://fusion.com/verify-email"
        },
        "username": {
          "type": "string",
          "example": "username"
        }
      }
    },
    "auth.RegisterResponse": {
      "type": "object",
      "properties": {
        "message": {
          "type": "string",
          "example": "User registered successfully, please verify your email"
        }
      }
    },
    "auth.ResetPasswordRequest": {
      "type": "object",
      "required": [
        "password",
        "resetToken"
      ],
      "properties": {
        "password": {
          "type": "string",
          "minLength": 8,
          "example": "newSecurePassword123"
        },
        "resetToken": {
          "type": "string",
          "example": "abc123def456"
        }
      }
    },
    "auth.ResetPasswordResponse": {
      "type": "object",
      "properties": {
        "message": {
          "type": "string",
          "example": "Password successfully reset"
        }
      }
    },
    "auth.SSOLoginRequest": {
      "type": "object",
      "required": [
        "ssoToken"
      ],
      "properties": {
        "ssoToken": {
          "type": "string",
          "example": "oauth2-token-from-provider"
        }
      }
    }
  },
  "securitySchemes": {
    "BearerAuth": {
      "description": "Type \"Bearer\" followed by a space and JWT token.",
      "type": "http",
      "scheme": "bearer",
      "bearerFormat": "JWT"
    }
  }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{"http", "https"},
	Title:            "Лавка Gateway API",
	Description:      "API Gateway для проекта Fusion",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
