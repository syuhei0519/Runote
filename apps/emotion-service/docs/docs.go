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
        "/emotions": {
            "get": {
                "description": "すべての登録感情を取得",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "emotions"
                ],
                "summary": "感情一覧を取得",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Emotion"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "ユーザー定義の感情を登録する（重複不可）",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "emotions"
                ],
                "summary": "感情を新規登録する",
                "parameters": [
                    {
                        "description": "感情登録リクエスト",
                        "name": "emotion",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.EmotionCreateRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.Emotion"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/emotions/unused": {
            "get": {
                "description": "投稿に紐づいていない感情のみを返す",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "emotions"
                ],
                "summary": "未使用感情一覧を取得",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Emotion"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/emotions/{id}": {
            "put": {
                "description": "感情の名前を変更する",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "emotions"
                ],
                "summary": "感情名の更新",
                "parameters": [
                    {
                        "type": "string",
                        "description": "感情ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "感情名更新内容",
                        "name": "emotion",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.EmotionNameUpdateRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Emotion"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "投稿で使用されていない感情を削除する",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "emotions"
                ],
                "summary": "感情をIDで削除",
                "parameters": [
                    {
                        "type": "string",
                        "description": "感情ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.SuccessResponse"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/post-emotions/{post_id}/{user_id}": {
            "get": {
                "description": "投稿ID・ユーザーIDに紐づく感情をRedisまたはDBから取得",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "post-emotions"
                ],
                "summary": "投稿に紐づく感情を取得",
                "parameters": [
                    {
                        "type": "string",
                        "description": "投稿ID",
                        "name": "post_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ユーザーID",
                        "name": "user_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.PostEmotion"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "投稿ID・ユーザーIDを指定して感情を更新、キャッシュも更新",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "post-emotions"
                ],
                "summary": "投稿に紐づく感情を更新",
                "parameters": [
                    {
                        "type": "string",
                        "description": "投稿ID",
                        "name": "post_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ユーザーID",
                        "name": "user_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "感情更新内容",
                        "name": "emotion",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.PostEmotionUpdateRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.PostEmotion"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "指定された投稿ID・ユーザーIDに紐づく感情を削除、Redisキャッシュも削除",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "post-emotions"
                ],
                "summary": "投稿に紐づく感情を削除",
                "parameters": [
                    {
                        "type": "string",
                        "description": "投稿ID",
                        "name": "post_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ユーザーID",
                        "name": "user_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handlers.SuccessResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handlers.EmotionCreateRequest": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "handlers.EmotionNameUpdateRequest": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "handlers.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "不正な入力です"
                }
            }
        },
        "handlers.PostEmotionUpdateRequest": {
            "type": "object",
            "required": [
                "emotion_id",
                "intensity"
            ],
            "properties": {
                "emotion_id": {
                    "type": "integer"
                },
                "intensity": {
                    "type": "integer"
                }
            }
        },
        "handlers.SuccessResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "操作が成功しました"
                }
            }
        },
        "models.Emotion": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "isPreset": {
                    "type": "boolean"
                },
                "name": {
                    "type": "string"
                },
                "posts": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.PostEmotion"
                    }
                }
            }
        },
        "models.Post": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "emotions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.PostEmotion"
                    }
                },
                "id": {
                    "type": "integer"
                }
            }
        },
        "models.PostEmotion": {
            "type": "object",
            "properties": {
                "custom": {
                    "type": "string"
                },
                "emotion": {
                    "$ref": "#/definitions/models.Emotion"
                },
                "emotionID": {
                    "type": "integer"
                },
                "intensity": {
                    "type": "integer"
                },
                "post": {
                    "$ref": "#/definitions/models.Post"
                },
                "postID": {
                    "type": "integer"
                },
                "userID": {
                    "type": "integer"
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
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
