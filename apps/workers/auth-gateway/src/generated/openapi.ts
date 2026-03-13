// AUTO-GENERATED — do not edit. Run `pnpm generate:openapi` to regenerate.
export const openApiSpec = {
  "openapi": "3.0.3",
  "info": {
    "title": "Requiems API",
    "version": "1.0.0",
    "description": "Unified access to enterprise-grade APIs — email validation, text utilities, and more. Authenticate with the `requiems-api-key` header."
  },
  "servers": [
    {
      "url": "https://api.requiems.xyz",
      "description": "Production"
    }
  ],
  "components": {
    "securitySchemes": {
      "requiems-api-key": {
        "type": "apiKey",
        "in": "header",
        "name": "requiems-api-key",
        "description": "Your Requiems API key"
      }
    }
  },
  "security": [
    {
      "requiems-api-key": []
    }
  ],
  "tags": [
    {
      "name": "advice",
      "description": "Get random pieces of advice and wisdom for inspiration"
    },
    {
      "name": "barcode",
      "description": "Generate barcodes in multiple formats (Code 128, Code 93, Code 39, EAN-8, EAN-13), returned as a PNG image or base64-encoded JSON"
    },
    {
      "name": "bin-lookup",
      "description": "Pass the first 6–8 digits of any payment card and get back the issuing bank, card network, type, and country"
    },
    {
      "name": "counter",
      "description": "Atomic, namespace-isolated hit counter"
    },
    {
      "name": "detect-language",
      "description": "Detect the language of any text with confidence scoring"
    },
    {
      "name": "disposable_email",
      "description": "Detect disposable and temporary email addresses to prevent fraud and improve data quality. Our comprehensive blocklist is continuously updated to catch the latest disposable email providers."
    },
    {
      "name": "email-normalize",
      "description": "Normalize email addresses to their canonical form. Lowercased, trimmed, and canonicalized with provider-specific rules including alias-domain resolution. Also available as part of the full Email Validator."
    },
    {
      "name": "email-validate",
      "description": "Full email validation in one call. Syntax check, MX record lookup, disposable domain detection, normalization, and typo suggestions. Includes everything from the Disposable Domain Checker and Email Normalizer in a single request."
    },
    {
      "name": "emoji",
      "description": "Look up emoji by name, search by keyword, or get a random emoji with full Unicode metadata."
    },
    {
      "name": "holidays",
      "description": "Get a list of holidays for a specific country and year"
    },
    {
      "name": "horoscope",
      "description": "Get daily horoscope readings for all 12 zodiac signs"
    },
    {
      "name": "lorem-ipsum",
      "description": "Generate placeholder text for design mockups and prototypes"
    },
    {
      "name": "password-generator",
      "description": "Generate cryptographically secure random passwords with customizable length and character sets"
    },
    {
      "name": "phone-validation",
      "description": "Validate phone numbers globally. Detect carrier, country, type, and VOIP or virtual risk using only phone metadata."
    },
    {
      "name": "profanity",
      "description": "Detect and censor profanity in text for content moderation"
    },
    {
      "name": "qr-code",
      "description": "Generate QR codes from any text or URL, returned as a PNG image or base64-encoded JSON"
    },
    {
      "name": "quotes",
      "description": "Access a database of inspirational and famous quotes"
    },
    {
      "name": "random-user",
      "description": "Generate random fake user profiles for testing and prototyping — names, emails, phone numbers, addresses, and avatars"
    },
    {
      "name": "spell-check",
      "description": "Check spelling and get correction suggestions for misspelled words"
    },
    {
      "name": "sudoku",
      "description": "Generate Sudoku puzzles with solutions across multiple difficulty levels"
    },
    {
      "name": "text-similarity",
      "description": "Compare two texts and get a cosine similarity score between 0 and 1"
    },
    {
      "name": "thesaurus",
      "description": "Find synonyms and antonyms for any word to enhance vocabulary and writing"
    },
    {
      "name": "timezone",
      "description": "Get timezone information for any location by coordinates or city name"
    },
    {
      "name": "unit-conversion",
      "description": "Convert between units of measurement — length, weight, volume, temperature, area, and speed"
    },
    {
      "name": "useragent",
      "description": "Parse user agent strings to extract browser, OS, device type, and bot detection"
    },
    {
      "name": "random-word",
      "description": "Get random words with definitions and parts of speech. Perfect for vocabulary builders, educational apps, word games, or content inspiration."
    },
    {
      "name": "working-days",
      "description": "Calculate the number of working days between two dates with optional country-specific holidays"
    },
    {
      "name": "world-time",
      "description": "Get the current time for any IANA timezone by name"
    }
  ],
  "paths": {
    "/v1/text/advice": {
      "get": {
        "summary": "Get Random Advice",
        "tags": [
          "advice"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Returns a random piece of advice",
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "id": {
                          "type": "integer",
                          "description": "Unique identifier for the advice"
                        },
                        "advice": {
                          "type": "string",
                          "description": "A random piece of advice"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "id": 42,
                    "advice": "Don't compare yourself to others. Compare yourself to the person you were yesterday."
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          }
        }
      }
    },
    "/v1/tech/barcode": {
      "get": {
        "summary": "Generate Barcode (PNG)",
        "tags": [
          "barcode"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Returns a raw PNG image of the barcode. Ideal for direct embedding or file download.",
        "parameters": [
          {
            "name": "data",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "example": 123456789
            },
            "description": "The text or numeric string to encode in the barcode"
          },
          {
            "name": "type",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "example": "code128"
            },
            "description": "Barcode format: code128, code93, code39, ean8, ean13"
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response"
          },
          "400": {
            "description": "Missing or invalid parameters (e.g. data not provided, unsupported type)"
          },
          "422": {
            "description": "Data is invalid for the chosen barcode type (e.g. wrong digit count for EAN-8/EAN-13, non-numeric EAN data)"
          }
        }
      }
    },
    "/v1/tech/barcode/base64": {
      "get": {
        "summary": "Generate Barcode (Base64 JSON)",
        "tags": [
          "barcode"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Returns a JSON envelope containing the barcode as a base64-encoded PNG string, along with its type and dimensions.",
        "parameters": [
          {
            "name": "data",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "example": 123456789
            },
            "description": "The text or numeric string to encode in the barcode"
          },
          {
            "name": "type",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "example": "code128"
            },
            "description": "Barcode format: code128, code93, code39, ean8, ean13"
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "image": {
                          "type": "string",
                          "description": "Base64-encoded PNG image data"
                        },
                        "type": {
                          "type": "string",
                          "description": "The barcode format that was used"
                        },
                        "width": {
                          "type": "integer",
                          "description": "Width of the generated image in pixels"
                        },
                        "height": {
                          "type": "integer",
                          "description": "Height of the generated image in pixels"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "image": "<base64-encoded PNG data>",
                    "type": "code128",
                    "width": 300,
                    "height": 100
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          },
          "400": {
            "description": "Missing or invalid parameters (e.g. data not provided, unsupported type)"
          },
          "422": {
            "description": "Data is invalid for the chosen barcode type (e.g. wrong digit count for EAN-8/EAN-13, non-numeric EAN data)"
          }
        }
      }
    },
    "/v1/finance/bin/{bin}": {
      "get": {
        "summary": "BIN Lookup",
        "tags": [
          "bin-lookup"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Returns card metadata for the given 6–8 digit BIN prefix.",
        "parameters": [
          {
            "name": "bin",
            "in": "path",
            "required": true,
            "schema": {
              "type": "string",
              "example": "424242"
            },
            "description": "6–8 digit Bank Identification Number. Dashes and spaces are stripped automatically."
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "bin": {
                          "type": "string",
                          "description": "The normalised BIN prefix used for the lookup"
                        },
                        "scheme": {
                          "type": "string",
                          "description": "Card network: visa, mastercard, amex, discover, jcb, diners, unionpay, maestro, mir, rupay, private_label"
                        },
                        "card_type": {
                          "type": "string",
                          "description": "credit, debit, prepaid, or charge"
                        },
                        "card_level": {
                          "type": "string",
                          "description": "classic, gold, platinum, infinite, business, signature, or standard"
                        },
                        "issuer_name": {
                          "type": "string",
                          "description": "Name of the card-issuing bank"
                        },
                        "issuer_url": {
                          "type": "string",
                          "description": "Bank website URL"
                        },
                        "issuer_phone": {
                          "type": "string",
                          "description": "Bank customer service phone number"
                        },
                        "country_code": {
                          "type": "string",
                          "description": "ISO 3166-1 alpha-2 country code of the issuing bank (e.g. US, GB, DE)"
                        },
                        "country_name": {
                          "type": "string",
                          "description": "Full country name of the issuing bank"
                        },
                        "prepaid": {
                          "type": "boolean",
                          "description": "Whether the card is a prepaid card"
                        },
                        "luhn": {
                          "type": "boolean",
                          "description": "Whether the BIN prefix passes the Luhn algorithm check"
                        },
                        "confidence": {
                          "type": "number",
                          "description": "Data quality score (0.00–1.00). Multi-source confirmed records score higher."
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "bin": "424242",
                    "scheme": "visa",
                    "card_type": "credit",
                    "card_level": "classic",
                    "issuer_name": "Chase",
                    "issuer_url": "www.chase.com",
                    "issuer_phone": "+18002324000",
                    "country_code": "US",
                    "country_name": "United States",
                    "prepaid": false,
                    "luhn": true,
                    "confidence": 0.92
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          },
          "400": {
            "description": "BIN is not 6–8 digits or contains non-digit characters."
          },
          "404": {
            "description": "BIN prefix not found in the database."
          },
          "500": {
            "description": "Unexpected server error."
          }
        }
      }
    },
    "/v1/misc/counter/{namespace}": {
      "post": {
        "summary": "Increment Counter",
        "tags": [
          "counter"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Atomically increment a counter in the specified namespace and return the new value",
        "parameters": [
          {
            "name": "namespace",
            "in": "path",
            "required": true,
            "schema": {
              "type": "string",
              "example": "page-views"
            },
            "description": "Counter namespace (1-64 chars: alphanumeric, hyphen, underscore)"
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "namespace": {
                          "type": "string",
                          "description": "The counter namespace"
                        },
                        "value": {
                          "type": "integer",
                          "description": "The new counter value after increment"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "namespace": "page-views",
                    "value": 42
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          },
          "400": {
            "description": "Invalid namespace: must be 1–64 chars, alphanumeric, hyphen or underscore only"
          },
          "500": {
            "description": "Internal server error"
          }
        }
      },
      "get": {
        "summary": "Get Counter Value",
        "tags": [
          "counter"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Get the current value of a counter without incrementing it",
        "parameters": [
          {
            "name": "namespace",
            "in": "path",
            "required": true,
            "schema": {
              "type": "string",
              "example": "page-views"
            },
            "description": "Counter namespace (1-64 chars: alphanumeric, hyphen, underscore)"
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "namespace": {
                          "type": "string",
                          "description": "The counter namespace"
                        },
                        "value": {
                          "type": "integer",
                          "description": "The current counter value (returns 0 if counter doesn't exist)"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "namespace": "page-views",
                    "value": 42
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          },
          "400": {
            "description": "Invalid namespace: must be 1–64 chars, alphanumeric, hyphen or underscore only"
          },
          "500": {
            "description": "Internal server error"
          }
        }
      }
    },
    "/v1/ai/detect-language": {
      "post": {
        "summary": "Detect Language",
        "tags": [
          "detect-language"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Identifies the language of the provided text and returns the language name, ISO 639-1 code, and confidence score.",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "text": {
                    "type": "string",
                    "description": "The text whose language should be detected.",
                    "example": "Bonjour, comment ça va?"
                  }
                },
                "required": [
                  "text"
                ],
                "example": {
                  "text": "Bonjour, comment ça va?"
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "language": {
                          "type": "string",
                          "description": "Full name of the detected language (e.g. French, English, Spanish)"
                        },
                        "code": {
                          "type": "string",
                          "description": "ISO 639-1 two-letter language code (e.g. fr, en, es). Empty string when detection is unreliable."
                        },
                        "confidence": {
                          "type": "string",
                          "description": "Confidence score between 0.0 and 1.0. 0.0 is returned when the language cannot be reliably detected."
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "language": "French",
                    "code": "fr",
                    "confidence": 0.98
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          },
          "400": {
            "description": "The request body is missing or malformed."
          },
          "422": {
            "description": "The text field is missing or empty."
          },
          "500": {
            "description": "Unexpected server error."
          }
        }
      }
    },
    "/v1/email/disposable/check": {
      "post": {
        "summary": "Check Single Email",
        "tags": [
          "disposable_email"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Validate whether an email address uses a disposable domain",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "email": {
                    "type": "string",
                    "description": "The email address to check",
                    "example": "test@example.com"
                  }
                },
                "required": [
                  "email"
                ],
                "example": {
                  "email": "user@tempmail.com"
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "email": {
                          "type": "string",
                          "description": "The email address that was checked"
                        },
                        "is_disposable": {
                          "type": "boolean",
                          "description": "Whether the email uses a disposable domain"
                        },
                        "domain": {
                          "type": "string",
                          "description": "The domain part of the email address"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "email": "user@tempmail.com",
                    "is_disposable": true,
                    "domain": "tempmail.com"
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          },
          "400": {
            "description": "The request body is missing or malformed; The email address format is invalid"
          }
        }
      }
    },
    "/v1/email/disposable/check-batch": {
      "post": {
        "summary": "Check Batch Emails",
        "tags": [
          "disposable_email"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Validate multiple email addresses in a single request (max 100 emails)",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "emails": {
                    "type": "array",
                    "items": {},
                    "description": "Array of email addresses to check (max 100)",
                    "example": [
                      "user1@example.com",
                      "user2@tempmail.com"
                    ]
                  }
                },
                "required": [
                  "emails"
                ],
                "example": {
                  "emails": [
                    "user1@gmail.com",
                    "user2@tempmail.com",
                    "user3@guerrillamail.com"
                  ]
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "results": {
                          "type": "array",
                          "items": {},
                          "description": "Array of check results for each email"
                        },
                        "total": {
                          "type": "integer",
                          "description": "Total number of emails checked"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "results": [
                      {
                        "email": "user1@gmail.com",
                        "is_disposable": false,
                        "domain": "gmail.com"
                      },
                      {
                        "email": "user2@tempmail.com",
                        "is_disposable": true,
                        "domain": "tempmail.com"
                      },
                      {
                        "email": "user3@guerrillamail.com",
                        "is_disposable": true,
                        "domain": "guerrillamail.com"
                      }
                    ],
                    "total": 3
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          },
          "400": {
            "description": "The request body is missing or malformed; The emails field is missing; Too many emails in the request"
          }
        }
      }
    },
    "/v1/email/disposable/domain/{domain}": {
      "get": {
        "summary": "Check Domain",
        "tags": [
          "disposable_email"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Check if a specific domain is in the disposable blocklist",
        "parameters": [
          {
            "name": "domain",
            "in": "path",
            "required": true,
            "schema": {
              "type": "string",
              "example": "tempmail.com"
            },
            "description": "The domain to check"
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "domain": {
                          "type": "string",
                          "description": "The domain that was checked"
                        },
                        "is_disposable": {
                          "type": "boolean",
                          "description": "Whether the domain is in the disposable blocklist"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "domain": "tempmail.com",
                    "is_disposable": true
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          },
          "400": {
            "description": "The domain parameter is missing"
          }
        }
      }
    },
    "/v1/email/disposable/stats": {
      "get": {
        "summary": "Get Statistics",
        "tags": [
          "disposable_email"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Get statistics about the disposable email blocklist",
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "total_domains": {
                          "type": "integer",
                          "description": "Total number of disposable domains in the blocklist"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "total_domains": 10500
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          }
        }
      }
    },
    "/v1/email/disposable/domains": {
      "get": {
        "summary": "List Domains (Paginated)",
        "tags": [
          "disposable_email"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Get a paginated list of all disposable domains in the blocklist",
        "parameters": [
          {
            "name": "page",
            "in": "query",
            "required": false,
            "schema": {
              "type": "integer",
              "example": 1
            },
            "description": "Page number (default: 1)"
          },
          {
            "name": "per_page",
            "in": "query",
            "required": false,
            "schema": {
              "type": "integer",
              "example": 100
            },
            "description": "Items per page (default: 100, max: 1000)"
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "domains": {
                          "type": "array",
                          "items": {},
                          "description": "Array of domain names"
                        },
                        "total": {
                          "type": "integer",
                          "description": "Total number of domains in blocklist"
                        },
                        "page": {
                          "type": "integer",
                          "description": "Current page number"
                        },
                        "per_page": {
                          "type": "integer",
                          "description": "Number of items per page"
                        },
                        "has_more": {
                          "type": "boolean",
                          "description": "Whether there are more pages available"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "domains": [
                      "tempmail.com",
                      "guerrillamail.com",
                      "10minutemail.com"
                    ],
                    "total": 10500,
                    "page": 1,
                    "per_page": 100,
                    "has_more": true
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          }
        }
      }
    },
    "/v1/email/normalize": {
      "post": {
        "summary": "Normalize Email",
        "tags": [
          "email-normalize"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Normalizes a single email address and returns the canonical form together with a breakdown of all transformations applied.",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "email": {
                    "type": "string",
                    "description": "The email address to normalize. Must be a syntactically valid address.",
                    "example": "Te.st.User+spam@Googlemail.com"
                  }
                },
                "required": [
                  "email"
                ],
                "example": {
                  "email": "Te.st.User+spam@Googlemail.com"
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "original": {
                          "type": "string",
                          "description": "The email address exactly as supplied in the request body"
                        },
                        "normalized": {
                          "type": "string",
                          "description": "The canonical form of the address after all transformations"
                        },
                        "local": {
                          "type": "string",
                          "description": "The local part (before @) of the normalized address"
                        },
                        "domain": {
                          "type": "string",
                          "description": "The domain part (after @) of the normalized address"
                        },
                        "changes": {
                          "type": "array",
                          "items": {},
                          "description": "Ordered list of transformations applied. Possible values: lowercased, trimmed_whitespace, removed_dots, removed_plus_tag, canonicalised_domain. Empty array when no changes were needed."
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "original": "Te.st.User+spam@Googlemail.com",
                    "normalized": "testuser@gmail.com",
                    "local": "testuser",
                    "domain": "gmail.com",
                    "changes": [
                      "lowercased",
                      "removed_dots",
                      "removed_plus_tag",
                      "canonicalised_domain"
                    ]
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          },
          "400": {
            "description": "The request body is missing, not valid JSON, or contains unknown fields."
          },
          "422": {
            "description": "The email field is missing or not a valid email address format."
          },
          "500": {
            "description": "Unexpected server error."
          }
        }
      }
    },
    "/v1/email/validate": {
      "post": {
        "summary": "Validate Email",
        "tags": [
          "email-validate"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Validates a single email address and returns a full breakdown of syntax validity, MX record status, disposable domain check, normalized form, and any typo suggestion.",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "email": {
                    "type": "string",
                    "description": "The email address to validate.",
                    "example": "user@gmial.com"
                  }
                },
                "required": [
                  "email"
                ],
                "example": {
                  "email": "user@gmial.com"
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "email": {
                          "type": "string",
                          "description": "The email address exactly as supplied in the request body"
                        },
                        "valid": {
                          "type": "boolean",
                          "description": "Overall validity. True only when the address passes syntax validation and the domain has at least one MX record"
                        },
                        "syntax_valid": {
                          "type": "boolean",
                          "description": "Whether the address is syntactically valid according to RFC 5322"
                        },
                        "mx_valid": {
                          "type": "boolean",
                          "description": "Whether the domain has at least one MX record, meaning it can receive email"
                        },
                        "disposable": {
                          "type": "boolean",
                          "description": "Whether the address uses a known disposable or temporary email domain"
                        },
                        "normalized": {
                          "type": "string",
                          "description": "The canonical form of the address after normalization (lowercase, plus-tag removal, alias-domain resolution). Empty string when syntax is invalid"
                        },
                        "domain": {
                          "type": "string",
                          "description": "The domain part of the address (after @). Empty string when syntax is invalid"
                        },
                        "suggestion": {
                          "type": "string",
                          "description": "A corrected domain name when the supplied domain looks like a typo of a well-known provider (e.g. gmial.com → gmail.com). Null when no close match is found or the domain is already correct"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "email": "user@gmial.com",
                    "valid": false,
                    "syntax_valid": true,
                    "mx_valid": false,
                    "disposable": false,
                    "normalized": "user@gmial.com",
                    "domain": "gmial.com",
                    "suggestion": "gmail.com"
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          },
          "400": {
            "description": "The request body is missing, not valid JSON, or contains unknown fields."
          },
          "422": {
            "description": "The email field is missing from the request body."
          },
          "500": {
            "description": "Unexpected server error."
          }
        }
      }
    },
    "/v1/entertainment/emoji/random": {
      "get": {
        "summary": "Get Random Emoji",
        "tags": [
          "emoji"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Returns a randomly selected emoji with its full metadata.",
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "emoji": {
                          "type": "string",
                          "description": "The rendered emoji glyph"
                        },
                        "name": {
                          "type": "string",
                          "description": "CLDR short name in snake_case (e.g. grinning_face)"
                        },
                        "category": {
                          "type": "string",
                          "description": "Unicode category (e.g. Smileys & Emotion, Animals & Nature)"
                        },
                        "unicode": {
                          "type": "string",
                          "description": "Unicode code-point in U+XXXX notation (e.g. U+1F600)"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "emoji": "😀",
                    "name": "grinning_face",
                    "category": "Smileys & Emotion",
                    "unicode": "U+1F600"
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          }
        }
      }
    },
    "/v1/entertainment/emoji/search": {
      "get": {
        "summary": "Search Emoji",
        "tags": [
          "emoji"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Search for emojis whose name or category contains the given query string (case-insensitive). Returns a list of all matches.",
        "parameters": [
          {
            "name": "q",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "example": "happy"
            },
            "description": "Search term to match against emoji names and categories (e.g. happy, heart, food)"
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "items": {
                          "type": "array",
                          "items": {},
                          "description": "List of matching emoji objects"
                        },
                        "total": {
                          "type": "integer",
                          "description": "Total number of matches"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "items": [
                      {
                        "emoji": "😄",
                        "name": "grinning_face_with_smiling_eyes",
                        "category": "Smileys & Emotion",
                        "unicode": "U+1F604"
                      }
                    ],
                    "total": 1
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          },
          "400": {
            "description": "The q query parameter is missing or empty."
          }
        }
      }
    },
    "/v1/entertainment/emoji/{name}": {
      "get": {
        "summary": "Get Emoji by Name",
        "tags": [
          "emoji"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Returns a specific emoji by its CLDR snake_case name. The name is case-insensitive.",
        "parameters": [
          {
            "name": "name",
            "in": "path",
            "required": true,
            "schema": {
              "type": "string",
              "example": "grinning_face"
            },
            "description": "CLDR snake_case emoji name (e.g. grinning_face, thumbs_up)"
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "emoji": {
                          "type": "string",
                          "description": "The rendered emoji glyph"
                        },
                        "name": {
                          "type": "string",
                          "description": "CLDR short name in snake_case"
                        },
                        "category": {
                          "type": "string",
                          "description": "Unicode category"
                        },
                        "unicode": {
                          "type": "string",
                          "description": "Unicode code-point in U+XXXX notation"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "emoji": "😀",
                    "name": "grinning_face",
                    "category": "Smileys & Emotion",
                    "unicode": "U+1F600"
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          },
          "404": {
            "description": "No emoji found with the given name."
          }
        }
      }
    },
    "/v1/places/holidays": {
      "get": {
        "summary": "Get Holidays",
        "tags": [
          "holidays"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Returns a list of public holidays for the specified country and year",
        "parameters": [
          {
            "name": "country",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "example": "US"
            },
            "description": "ISO 3166-1 alpha-2 country code (e.g., \"US\", \"GB\", \"DE\")"
          },
          {
            "name": "year",
            "in": "query",
            "required": true,
            "schema": {
              "type": "integer",
              "example": 2025
            },
            "description": "Year for which to retrieve holidays (e.g., 2025)"
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "country": {
                          "type": "string",
                          "description": "ISO 3166-1 alpha-2 country code"
                        },
                        "year": {
                          "type": "integer",
                          "description": "Year for which holidays are returned"
                        },
                        "holidays": {
                          "type": "array",
                          "items": {},
                          "description": "Array of holiday objects"
                        },
                        "holidays[].date": {
                          "type": "string",
                          "description": "Holiday date in YYYY-MM-DD format"
                        },
                        "holidays[].name": {
                          "type": "string",
                          "description": "Name of the holiday"
                        },
                        "total": {
                          "type": "integer",
                          "description": "Total number of holidays for the country/year"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "country": "US",
                    "year": 2025,
                    "holidays": [
                      {
                        "date": "2025-01-01",
                        "name": "New Year's Day"
                      },
                      {
                        "date": "2025-07-04",
                        "name": "Independence Day"
                      }
                    ],
                    "total": 11
                  },
                  "metadata": {
                    "timestamp": "2025-01-15T10:30:00Z"
                  }
                }
              }
            }
          },
          "400": {
            "description": "Missing or invalid country code or year parameter"
          },
          "404": {
            "description": "No holidays found for the specified country and year"
          }
        }
      }
    },
    "/v1/entertainment/horoscope/{sign}": {
      "get": {
        "summary": "Get Daily Horoscope",
        "tags": [
          "horoscope"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Returns a daily horoscope reading for the specified zodiac sign.",
        "parameters": [
          {
            "name": "sign",
            "in": "path",
            "required": true,
            "schema": {
              "type": "string"
            },
            "description": "Zodiac sign (case-insensitive). Supported values: aries, taurus, gemini, cancer, leo, virgo, libra, scorpio, sagittarius, capricorn, aquarius, pisces"
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "sign": {
                          "type": "string",
                          "description": "Normalized zodiac sign (lowercase)"
                        },
                        "date": {
                          "type": "string",
                          "description": "Today's date in YYYY-MM-DD format (UTC)"
                        },
                        "horoscope": {
                          "type": "string",
                          "description": "Daily horoscope reading"
                        },
                        "lucky_number": {
                          "type": "integer",
                          "description": "Lucky number for the day (1-99)"
                        },
                        "mood": {
                          "type": "string",
                          "description": "Suggested mood for the day"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "sign": "aries",
                    "date": "2024-12-15",
                    "horoscope": "Today is a great day for new beginnings. Trust your instincts and take that first step toward your goals.",
                    "lucky_number": 7,
                    "mood": "energetic"
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          }
        }
      }
    },
    "/v1/text/lorem": {
      "get": {
        "summary": "Generate Lorem Ipsum",
        "tags": [
          "lorem-ipsum"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Generate Lorem Ipsum placeholder text with customizable length and format",
        "parameters": [
          {
            "name": "paragraphs",
            "in": "query",
            "required": false,
            "schema": {
              "type": "integer",
              "example": 3
            },
            "description": "Number of paragraphs to generate (1-20)"
          },
          {
            "name": "sentences",
            "in": "query",
            "required": false,
            "schema": {
              "type": "integer",
              "example": 5
            },
            "description": "Number of sentences per paragraph (1-20)"
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "text": {
                          "type": "string",
                          "description": "Generated Lorem Ipsum text"
                        },
                        "paragraphs": {
                          "type": "integer",
                          "description": "Number of paragraphs generated"
                        },
                        "wordCount": {
                          "type": "integer",
                          "description": "Total number of words in generated text"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "text": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris.",
                    "paragraphs": 1,
                    "wordCount": 45
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          },
          "400": {
            "description": "The paragraphs parameter is out of valid range; The sentences parameter is out of valid range"
          }
        }
      }
    },
    "/v1/tech/password": {
      "get": {
        "summary": "Generate Password",
        "tags": [
          "password-generator"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Generate a cryptographically secure random password with customizable character sets and length",
        "parameters": [
          {
            "name": "length",
            "in": "query",
            "required": false,
            "schema": {
              "type": "integer",
              "example": 16
            },
            "description": "Password length (8-128 characters)"
          },
          {
            "name": "uppercase",
            "in": "query",
            "required": false,
            "schema": {
              "type": "boolean",
              "example": true
            },
            "description": "Include uppercase letters (A-Z)"
          },
          {
            "name": "numbers",
            "in": "query",
            "required": false,
            "schema": {
              "type": "boolean",
              "example": true
            },
            "description": "Include numbers (0-9)"
          },
          {
            "name": "symbols",
            "in": "query",
            "required": false,
            "schema": {
              "type": "boolean",
              "example": true
            },
            "description": "Include special characters (!@#$%^&*()-_=+[]{}|;:,.<>?)"
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "password": {
                          "type": "string",
                          "description": "The generated password"
                        },
                        "length": {
                          "type": "integer",
                          "description": "Length of the generated password"
                        },
                        "strength": {
                          "type": "string",
                          "description": "Password strength assessment (weak, medium, or strong)"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "password": "aB3#cDeFgHiJkLmN",
                    "length": 16,
                    "strength": "strong"
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          },
          "400": {
            "description": "The length parameter is out of valid range (8-128)"
          },
          "500": {
            "description": "Failed to generate password (rare cryptographic failure)"
          }
        }
      }
    },
    "/v1/tech/validate/phone": {
      "get": {
        "summary": "Validate Phone Number",
        "tags": [
          "phone-validation"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Validates a single phone number and returns its country, type, formatted representation, carrier, and VOIP/virtual risk flags.",
        "parameters": [
          {
            "name": "number",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "example": "+447400123456"
            },
            "description": "The phone number to validate. Must include the country calling code (e.g. +12015551234)."
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "number": {
                          "type": "string",
                          "description": "The original number as supplied in the request"
                        },
                        "valid": {
                          "type": "boolean",
                          "description": "Whether the number is a valid, dialable phone number"
                        },
                        "country": {
                          "type": "string",
                          "description": "ISO 3166-1 alpha-2 country code (omitted when valid is false)"
                        },
                        "type": {
                          "type": "string",
                          "description": "Number type: mobile, landline, landline_or_mobile, toll_free, voip, premium_rate, shared_cost, personal_number, pager, uan, voicemail, or unknown (omitted when valid is false)"
                        },
                        "formatted": {
                          "type": "string",
                          "description": "International format of the number, e.g. +44 7400 123456 (omitted when valid is false)"
                        },
                        "carrier.name": {
                          "type": "string",
                          "description": "Carrier name from phone prefix metadata (omitted when carrier cannot be determined)"
                        },
                        "carrier.source": {
                          "type": "string",
                          "description": "How the carrier was determined. Always \"metadata\" when present"
                        },
                        "risk.is_voip": {
                          "type": "boolean",
                          "description": "true when the number type is voip"
                        },
                        "risk.is_virtual": {
                          "type": "boolean",
                          "description": "true when the number is not tied to a physical SIM or fixed line: voip, personal_number, uan, pager, or voicemail"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "number": "+447400123456",
                    "valid": true,
                    "country": "GB",
                    "type": "mobile",
                    "formatted": "+44 7400 123456",
                    "carrier": {
                      "name": "Three",
                      "source": "metadata"
                    },
                    "risk": {
                      "is_voip": false,
                      "is_virtual": false
                    }
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          },
          "400": {
            "description": "The number query parameter is missing."
          }
        }
      }
    },
    "/v1/tech/validate/phone/batch": {
      "post": {
        "summary": "Batch Validate Phone Numbers",
        "tags": [
          "phone-validation"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Validates up to 50 phone numbers in a single request. Results are returned in the same order as the input.",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "numbers": {
                    "type": "array",
                    "items": {},
                    "description": "Array of phone numbers to validate (min: 1, max: 50). Each must include the country calling code.",
                    "example": "[\"+447400123456\", \"+12015551234\"]"
                  }
                },
                "required": [
                  "numbers"
                ],
                "example": {
                  "numbers": [
                    "+447400123456",
                    "+12015551234",
                    "12345"
                  ]
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "results": {
                          "type": "array",
                          "items": {},
                          "description": "Validation result for each number in the same order as the input. Each item has the same fields as the single validate endpoint."
                        },
                        "total": {
                          "type": "integer",
                          "description": "Number of results returned. Matches the length of the input array."
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "results": [
                      {
                        "number": "+447400123456",
                        "valid": true,
                        "country": "GB",
                        "type": "mobile",
                        "formatted": "+44 7400 123456",
                        "carrier": {
                          "name": "Three",
                          "source": "metadata"
                        },
                        "risk": {
                          "is_voip": false,
                          "is_virtual": false
                        }
                      },
                      {
                        "number": "+12015551234",
                        "valid": true,
                        "country": "US",
                        "type": "landline_or_mobile",
                        "formatted": "+1 201-555-1234",
                        "risk": {
                          "is_voip": false,
                          "is_virtual": false
                        }
                      },
                      {
                        "number": "12345",
                        "valid": false
                      }
                    ],
                    "total": 3
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          },
          "422": {
            "description": "The numbers array is missing, empty, or contains more than 50 items."
          }
        }
      }
    },
    "/v1/text/profanity": {
      "post": {
        "summary": "Check Profanity",
        "tags": [
          "profanity"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Checks text for profanity, returning a censored version and the list of flagged words.",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "text": {
                    "type": "string",
                    "description": "The text to check for profanity.",
                    "example": "What the heck is going on"
                  }
                },
                "required": [
                  "text"
                ],
                "example": {
                  "text": "Some text to check"
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "has_profanity": {
                          "type": "boolean",
                          "description": "Whether any profanity was detected in the text"
                        },
                        "censored": {
                          "type": "string",
                          "description": "The input text with profane words replaced by asterisks"
                        },
                        "flagged_words": {
                          "type": "string",
                          "description": "Deduplicated list of profane words found (lowercase)"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "has_profanity": false,
                    "censored": "Some text to check",
                    "flagged_words": []
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          },
          "400": {
            "description": "The request body is missing or malformed."
          },
          "422": {
            "description": "The text field is missing or empty."
          },
          "500": {
            "description": "Unexpected server error."
          }
        }
      }
    },
    "/v1/tech/qr": {
      "get": {
        "summary": "Generate QR Code (PNG)",
        "tags": [
          "qr-code"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Returns a raw PNG image of the QR code. Ideal for direct embedding or file download.",
        "parameters": [
          {
            "name": "data",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "example": "https://example.com"
            },
            "description": "The text or URL to encode in the QR code"
          },
          {
            "name": "size",
            "in": "query",
            "required": false,
            "schema": {
              "type": "integer",
              "example": 200
            },
            "description": "Image size in pixels (default: 256, min: 50, max: 1000)"
          },
          {
            "name": "recovery",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "example": "high"
            },
            "description": "Error-correction level: low (7%), medium (15%), high (25%), highest (30%). Higher levels are more robust to physical damage but produce larger images. Default: medium"
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response"
          },
          "400": {
            "description": "Missing or invalid parameters (e.g. data not provided, size out of range, unknown recovery level)"
          },
          "500": {
            "description": "Failed to generate QR code"
          }
        }
      }
    },
    "/v1/tech/qr/base64": {
      "get": {
        "summary": "Generate QR Code (Base64 JSON)",
        "tags": [
          "qr-code"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Returns a JSON envelope containing the QR code as a base64-encoded PNG string, along with its dimensions.",
        "parameters": [
          {
            "name": "data",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "example": "https://example.com"
            },
            "description": "The text or URL to encode in the QR code"
          },
          {
            "name": "size",
            "in": "query",
            "required": false,
            "schema": {
              "type": "integer",
              "example": 200
            },
            "description": "Image size in pixels (default: 256, min: 50, max: 1000)"
          },
          {
            "name": "recovery",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "example": "highest"
            },
            "description": "Error-correction level: low (7%), medium (15%), high (25%), highest (30%). Default: medium"
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "image": {
                          "type": "string",
                          "description": "Base64-encoded PNG image data"
                        },
                        "width": {
                          "type": "integer",
                          "description": "Width of the generated image in pixels"
                        },
                        "height": {
                          "type": "integer",
                          "description": "Height of the generated image in pixels"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "image": "<base64-encoded PNG data>",
                    "width": 256,
                    "height": 256
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          },
          "400": {
            "description": "Missing or invalid parameters (e.g. data not provided, size out of range, unknown recovery level)"
          },
          "500": {
            "description": "Failed to generate QR code"
          }
        }
      }
    },
    "/v1/text/quotes/random": {
      "get": {
        "summary": "Get Random Quote",
        "tags": [
          "quotes"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Returns a random inspirational quote with author attribution",
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "id": {
                          "type": "integer",
                          "description": "Unique identifier for the quote"
                        },
                        "text": {
                          "type": "string",
                          "description": "The quote text"
                        },
                        "author": {
                          "type": "string",
                          "description": "Name of the person who said or wrote the quote"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "id": 42,
                    "text": "The only way to do great work is to love what you do.",
                    "author": "Steve Jobs"
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          },
          "503": {
            "description": "No quotes available in the database"
          }
        }
      }
    },
    "/v1/misc/random-user": {
      "get": {
        "summary": "Get Random User",
        "tags": [
          "random-user"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Returns a randomly generated fake user profile.",
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "name": {
                          "type": "string",
                          "description": "Full name of the generated user"
                        },
                        "email": {
                          "type": "string",
                          "description": "Email address of the generated user"
                        },
                        "phone": {
                          "type": "string",
                          "description": "Phone number of the generated user"
                        },
                        "address.street": {
                          "type": "string",
                          "description": "Street address"
                        },
                        "address.city": {
                          "type": "string",
                          "description": "City name"
                        },
                        "address.state": {
                          "type": "string",
                          "description": "State or region"
                        },
                        "address.zip": {
                          "type": "string",
                          "description": "Postal / ZIP code"
                        },
                        "address.country": {
                          "type": "string",
                          "description": "Country name"
                        },
                        "avatar": {
                          "type": "string",
                          "description": "URL to a unique identicon avatar for the generated user (DiceBear)"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "name": "Grace Lopez",
                    "email": "grace.lopez@example.org",
                    "phone": "555-123-4567",
                    "address": {
                      "street": "4821 Maple Avenue",
                      "city": "North Judyton",
                      "state": "California",
                      "zip": "94103",
                      "country": "United States of America"
                    },
                    "avatar": "https://api.dicebear.com/9.x/identicon/svg?seed=Grace+Lopez"
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          },
          "500": {
            "description": "Internal server error"
          }
        }
      }
    },
    "/v1/text/spellcheck": {
      "post": {
        "summary": "Check Spelling",
        "tags": [
          "spell-check"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Checks the input text for spelling mistakes and returns a corrected version along with per-word corrections.",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "text": {
                    "type": "string",
                    "description": "The text to spell-check.",
                    "example": "Ths is a smiple tset"
                  }
                },
                "required": [
                  "text"
                ],
                "example": {
                  "text": "Ths is a smiple tset"
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "corrected": {
                          "type": "string",
                          "description": "The full input text with all misspelled words replaced by their suggested corrections"
                        },
                        "corrections": {
                          "type": "string",
                          "description": "List of individual corrections. Each item contains: original (the misspelled word), suggested (the correction), and position (0-based character offset in the original text)"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "corrected": "This is a simple test",
                    "corrections": [
                      {
                        "original": "Ths",
                        "suggested": "This",
                        "position": 0
                      },
                      {
                        "original": "smiple",
                        "suggested": "simple",
                        "position": 9
                      },
                      {
                        "original": "tset",
                        "suggested": "test",
                        "position": 16
                      }
                    ]
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          },
          "400": {
            "description": "The request body is missing or malformed."
          },
          "422": {
            "description": "The text field is missing or empty."
          },
          "500": {
            "description": "Unexpected server error."
          }
        }
      }
    },
    "/v1/entertainment/sudoku": {
      "get": {
        "summary": "Get Sudoku Puzzle",
        "tags": [
          "sudoku"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Returns a randomly generated Sudoku puzzle and its solution. Difficulty defaults to medium when not specified.",
        "parameters": [
          {
            "name": "difficulty",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string"
            },
            "description": "Puzzle difficulty level. One of: easy, medium, hard. Defaults to medium."
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "difficulty": {
                          "type": "string",
                          "description": "The difficulty level of the returned puzzle (easy, medium, or hard)"
                        },
                        "puzzle": {
                          "type": "string",
                          "description": "9×9 grid representing the puzzle — 0 means an empty cell to be filled in"
                        },
                        "solution": {
                          "type": "string",
                          "description": "9×9 grid containing the complete, valid solution"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "difficulty": "hard",
                    "puzzle": [
                      [
                        5,
                        3,
                        0,
                        0,
                        7,
                        0,
                        0,
                        0,
                        0
                      ],
                      [
                        6,
                        0,
                        0,
                        1,
                        9,
                        5,
                        0,
                        0,
                        0
                      ],
                      [
                        0,
                        9,
                        8,
                        0,
                        0,
                        0,
                        0,
                        6,
                        0
                      ],
                      [
                        8,
                        0,
                        0,
                        0,
                        6,
                        0,
                        0,
                        0,
                        3
                      ],
                      [
                        4,
                        0,
                        0,
                        8,
                        0,
                        3,
                        0,
                        0,
                        1
                      ],
                      [
                        7,
                        0,
                        0,
                        0,
                        2,
                        0,
                        0,
                        0,
                        6
                      ],
                      [
                        0,
                        6,
                        0,
                        0,
                        0,
                        0,
                        2,
                        8,
                        0
                      ],
                      [
                        0,
                        0,
                        0,
                        4,
                        1,
                        9,
                        0,
                        0,
                        5
                      ],
                      [
                        0,
                        0,
                        0,
                        0,
                        8,
                        0,
                        0,
                        7,
                        9
                      ]
                    ],
                    "solution": [
                      [
                        5,
                        3,
                        4,
                        6,
                        7,
                        8,
                        9,
                        1,
                        2
                      ],
                      [
                        6,
                        7,
                        2,
                        1,
                        9,
                        5,
                        3,
                        4,
                        8
                      ],
                      [
                        1,
                        9,
                        8,
                        3,
                        4,
                        2,
                        5,
                        6,
                        7
                      ],
                      [
                        8,
                        5,
                        9,
                        7,
                        6,
                        1,
                        4,
                        2,
                        3
                      ],
                      [
                        4,
                        2,
                        6,
                        8,
                        5,
                        3,
                        7,
                        9,
                        1
                      ],
                      [
                        7,
                        1,
                        3,
                        9,
                        2,
                        4,
                        8,
                        5,
                        6
                      ],
                      [
                        9,
                        6,
                        1,
                        5,
                        3,
                        7,
                        2,
                        8,
                        4
                      ],
                      [
                        2,
                        8,
                        7,
                        4,
                        1,
                        9,
                        6,
                        3,
                        5
                      ],
                      [
                        3,
                        4,
                        5,
                        2,
                        8,
                        6,
                        1,
                        7,
                        9
                      ]
                    ]
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          },
          "400": {
            "description": "The difficulty parameter is not one of easy, medium, or hard"
          },
          "401": {
            "description": "Missing API key"
          },
          "403": {
            "description": "Invalid or revoked API key"
          }
        }
      }
    },
    "/v1/ai/similarity": {
      "post": {
        "summary": "Compare Text Similarity",
        "tags": [
          "text-similarity"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Compares two texts and returns a cosine similarity score.",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "type": "object",
                "properties": {
                  "text1": {
                    "type": "string",
                    "description": "The first text to compare.",
                    "example": "The cat sat on the mat"
                  },
                  "text2": {
                    "type": "string",
                    "description": "The second text to compare.",
                    "example": "A cat was sitting on a mat"
                  }
                },
                "required": [
                  "text1",
                  "text2"
                ],
                "example": {
                  "text1": "The cat sat on the mat",
                  "text2": "A cat was sitting on a mat"
                }
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "similarity": {
                          "type": "number",
                          "description": "Cosine similarity score between the two texts, in the range [0, 1]."
                        },
                        "method": {
                          "type": "string",
                          "description": "The algorithm used. Currently always 'cosine'."
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "similarity": 0.4364,
                    "method": "cosine"
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          },
          "400": {
            "description": "The request body is missing or malformed."
          },
          "422": {
            "description": "One or both text fields are missing or empty."
          },
          "500": {
            "description": "Unexpected server error."
          }
        }
      }
    },
    "/v1/text/thesaurus/{word}": {
      "get": {
        "summary": "Thesaurus Lookup",
        "tags": [
          "thesaurus"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Returns synonyms and antonyms for the given word.",
        "parameters": [
          {
            "name": "word",
            "in": "path",
            "required": true,
            "schema": {
              "type": "string",
              "example": "happy"
            },
            "description": "The word to look up in the thesaurus"
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "word": {
                          "type": "string",
                          "description": "The normalized (lowercased) word that was looked up"
                        },
                        "synonyms": {
                          "type": "string",
                          "description": "List of words with similar meaning"
                        },
                        "antonyms": {
                          "type": "string",
                          "description": "List of words with opposite meaning"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "word": "happy",
                    "synonyms": [
                      "joyful",
                      "cheerful",
                      "content",
                      "pleased",
                      "delighted",
                      "glad",
                      "elated",
                      "blissful"
                    ],
                    "antonyms": [
                      "sad",
                      "unhappy",
                      "miserable",
                      "sorrowful",
                      "dejected",
                      "gloomy",
                      "melancholy"
                    ]
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          },
          "400": {
            "description": "The word path parameter is missing."
          },
          "404": {
            "description": "The word was not found in the thesaurus dataset."
          }
        }
      }
    },
    "/v1/places/timezone": {
      "get": {
        "summary": "Get Timezone",
        "tags": [
          "timezone"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Returns timezone information for the given coordinates or city name. Provide either `city` or both `lat` and `lon`.",
        "parameters": [
          {
            "name": "lat",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string"
            },
            "description": "Latitude of the location (-90 to 90). Required when using coordinate-based lookup."
          },
          {
            "name": "lon",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string"
            },
            "description": "Longitude of the location (-180 to 180). Required when using coordinate-based lookup."
          },
          {
            "name": "city",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string"
            },
            "description": "City name for city-based lookup (e.g. 'Tokyo', 'London'). Required when not using coordinates."
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "timezone": {
                          "type": "string",
                          "description": "IANA timezone identifier (e.g. \"Europe/London\", \"Asia/Tokyo\")"
                        },
                        "offset": {
                          "type": "string",
                          "description": "UTC offset in +HH:MM or -HH:MM format (e.g. '+05:30', '-05:00')"
                        },
                        "current_time": {
                          "type": "string",
                          "description": "Current time in UTC, formatted as RFC 3339 (e.g. \"2024-12-15T14:30:00Z\")"
                        },
                        "is_dst": {
                          "type": "boolean",
                          "description": "Whether the location is currently observing daylight saving time"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "timezone": "Europe/London",
                  "offset": "+00:00",
                  "current_time": "2024-12-15T14:30:00Z",
                  "is_dst": false
                }
              }
            }
          }
        }
      }
    },
    "/v1/misc/convert": {
      "get": {
        "summary": "Convert Units",
        "tags": [
          "unit-conversion"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Convert a value from one unit to another",
        "parameters": [
          {
            "name": "from",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "example": "miles"
            },
            "description": "Source unit key (e.g. miles, kg, c)"
          },
          {
            "name": "to",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "example": "km"
            },
            "description": "Target unit key (e.g. km, lb, f)"
          },
          {
            "name": "value",
            "in": "query",
            "required": true,
            "schema": {
              "type": "number",
              "example": 10
            },
            "description": "Numeric value to convert"
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "from": {
                          "type": "string",
                          "description": "Source unit key"
                        },
                        "to": {
                          "type": "string",
                          "description": "Target unit key"
                        },
                        "input": {
                          "type": "number",
                          "description": "The original input value"
                        },
                        "result": {
                          "type": "number",
                          "description": "The converted value (rounded to 6 decimal places)"
                        },
                        "formula": {
                          "type": "string",
                          "description": "Human-readable conversion formula"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "from": "miles",
                    "to": "km",
                    "input": 10,
                    "result": 16.09344,
                    "formula": "miles × 1.609344"
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          }
        }
      }
    },
    "/v1/misc/convert/units": {
      "get": {
        "summary": "List Available Units",
        "tags": [
          "unit-conversion"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Returns all available unit conversion types grouped by measurement category",
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "length": {
                          "type": "array",
                          "items": {},
                          "description": "Available length units: millimeter (mm), centimeter (cm), meter (m), kilometer (km), inch (in), foot (ft), yard (yd), mile (miles), nautical mile (nmi)"
                        },
                        "weight": {
                          "type": "array",
                          "items": {},
                          "description": "Available weight units: milligram (mg), gram (g), kilogram (kg), metric ton (t), ounce (oz), pound (lb), stone (stone)"
                        },
                        "volume": {
                          "type": "array",
                          "items": {},
                          "description": "Available volume units: milliliter (ml), liter (l), teaspoon (tsp), tablespoon (tbsp), fluid ounce (fl_oz), cup (cup), pint (pt), quart (qt), gallon (gal)"
                        },
                        "temperature": {
                          "type": "array",
                          "items": {},
                          "description": "Available temperature units: celsius (c), fahrenheit (f), kelvin (k)"
                        },
                        "area": {
                          "type": "array",
                          "items": {},
                          "description": "Available area units: square millimeter (mm2), square centimeter (cm2), square meter (m2), square kilometer (km2), square inch (in2), square foot (ft2), square yard (yd2), acre (acre), hectare (ha)"
                        },
                        "speed": {
                          "type": "array",
                          "items": {},
                          "description": "Available speed units: meters per second (m_s), kilometers per hour (km_h), miles per hour (mph), knots (knots)"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "length": [
                      "cm",
                      "ft",
                      "in",
                      "km",
                      "m",
                      "miles",
                      "mm",
                      "nmi",
                      "yd"
                    ],
                    "weight": [
                      "g",
                      "kg",
                      "lb",
                      "mg",
                      "oz",
                      "stone",
                      "t"
                    ],
                    "volume": [
                      "cup",
                      "fl_oz",
                      "gal",
                      "l",
                      "ml",
                      "pt",
                      "qt",
                      "tbsp",
                      "tsp"
                    ],
                    "temperature": [
                      "c",
                      "f",
                      "k"
                    ],
                    "area": [
                      "acre",
                      "cm2",
                      "ft2",
                      "ha",
                      "in2",
                      "km2",
                      "m2",
                      "mm2",
                      "yd2"
                    ],
                    "speed": [
                      "km_h",
                      "knots",
                      "m_s",
                      "mph"
                    ]
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          }
        }
      }
    },
    "/v1/tech/useragent": {
      "get": {
        "summary": "Parse User Agent",
        "tags": [
          "useragent"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Parses a user agent string and returns structured information about the browser, OS, device, and bot status.",
        "parameters": [
          {
            "name": "ua",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "example": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
            },
            "description": "The user agent string to parse."
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "browser": {
                          "type": "string",
                          "description": "Detected browser name (e.g. Chrome, Firefox, Safari, Edge, Opera, Internet Explorer, Other)"
                        },
                        "browser_version": {
                          "type": "string",
                          "description": "Detected browser version (major.minor)"
                        },
                        "os": {
                          "type": "string",
                          "description": "Detected operating system (e.g. Windows, macOS, Linux, Android, iOS, ChromeOS, Other)"
                        },
                        "os_version": {
                          "type": "string",
                          "description": "Detected OS version (format varies by platform)"
                        },
                        "device": {
                          "type": "string",
                          "description": "Device type — one of desktop, mobile, tablet, bot, or unknown"
                        },
                        "is_bot": {
                          "type": "boolean",
                          "description": "True when the user agent matches a known bot or crawler pattern"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "browser": "Chrome",
                    "browser_version": "120.0",
                    "os": "Windows",
                    "os_version": "10/11",
                    "device": "desktop",
                    "is_bot": false
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          },
          "400": {
            "description": "The ua query parameter is missing."
          }
        }
      }
    },
    "/v1/text/words/random": {
      "get": {
        "summary": "Get Random Word",
        "tags": [
          "random-word"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Returns a random word with its definition and part of speech",
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "id": {
                          "type": "integer",
                          "description": "Unique identifier for the word"
                        },
                        "word": {
                          "type": "string",
                          "description": "The random word"
                        },
                        "definition": {
                          "type": "string",
                          "description": "Dictionary definition of the word"
                        },
                        "part_of_speech": {
                          "type": "string",
                          "description": "Grammatical classification (e.g., noun, verb, adjective, adverb)"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "id": 123,
                    "word": "ephemeral",
                    "definition": "lasting for a very short time",
                    "part_of_speech": "adjective"
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          },
          "503": {
            "description": "No words available in the database"
          }
        }
      }
    },
    "/v1/places/working-days": {
      "get": {
        "summary": "Calculate Working Days",
        "tags": [
          "working-days"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Calculate the number of working days between two dates, optionally accounting for country-specific holidays",
        "parameters": [
          {
            "name": "from",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "example": "2024-02-23"
            },
            "description": "Start date in YYYY-MM-DD format (ISO 8601)"
          },
          {
            "name": "to",
            "in": "query",
            "required": true,
            "schema": {
              "type": "string",
              "example": "2024-02-28"
            },
            "description": "End date in YYYY-MM-DD format (ISO 8601). Must be >= from date."
          },
          {
            "name": "country",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "example": "US"
            },
            "description": "ISO 3166-1 alpha-2 country code (e.g., \"US\", \"GB\", \"FR\"). When provided, country-specific holidays are excluded from working days count."
          },
          {
            "name": "subdivision",
            "in": "query",
            "required": false,
            "schema": {
              "type": "string",
              "example": "NY"
            },
            "description": "ISO 3166-2 subdivision code for state/region within the country (e.g., \"NY\" for New York, \"CA\" for California). Only used when country is provided."
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "working_days": {
                          "type": "integer",
                          "description": "Number of working days between the two dates (excluding weekends and optionally holidays)"
                        },
                        "from": {
                          "type": "string",
                          "description": "Start date (echoed from request)"
                        },
                        "to": {
                          "type": "string",
                          "description": "End date (echoed from request)"
                        },
                        "country": {
                          "type": "string",
                          "description": "Country code (echoed from request, empty string if not provided)"
                        },
                        "subdivision": {
                          "type": "string",
                          "description": "Subdivision code (echoed from request, empty string if not provided)"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "working_days": 4,
                    "from": "2024-02-23",
                    "to": "2024-02-28",
                    "country": "US",
                    "subdivision": "NY"
                  },
                  "metadata": {
                    "timestamp": "2026-01-01T00:00:00Z"
                  }
                }
              }
            }
          },
          "400": {
            "description": "The from and to parameters are required, or to date is before from date, or invalid date format"
          }
        }
      }
    },
    "/v1/places/time/{timezone}": {
      "get": {
        "summary": "Get Current Time by Timezone",
        "tags": [
          "world-time"
        ],
        "security": [
          {
            "requiems-api-key": []
          }
        ],
        "description": "Returns the current time for the given IANA timezone identifier. The timezone is supplied as a path parameter (e.g. `America/New_York`, `Europe/London`, `UTC`).",
        "parameters": [
          {
            "name": "timezone",
            "in": "path",
            "required": true,
            "schema": {
              "type": "string",
              "example": "America/New_York"
            },
            "description": "IANA timezone identifier (e.g. 'America/New_York', 'Europe/London', 'Asia/Kolkata')"
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "content": {
              "application/json": {
                "schema": {
                  "type": "object",
                  "properties": {
                    "data": {
                      "type": "object",
                      "properties": {
                        "timezone": {
                          "type": "string",
                          "description": "IANA timezone identifier (e.g. \"America/New_York\")"
                        },
                        "offset": {
                          "type": "string",
                          "description": "UTC offset in +HH:MM or -HH:MM format (e.g. '-05:00', '+05:30')"
                        },
                        "current_time": {
                          "type": "string",
                          "description": "Current time in UTC, formatted as RFC 3339 (e.g. \"2024-12-15T14:30:00Z\")"
                        },
                        "is_dst": {
                          "type": "boolean",
                          "description": "Whether the timezone is currently observing daylight saving time"
                        }
                      }
                    },
                    "metadata": {
                      "type": "object",
                      "properties": {
                        "timestamp": {
                          "type": "string",
                          "format": "date-time"
                        }
                      }
                    }
                  }
                },
                "example": {
                  "data": {
                    "timezone": "America/New_York",
                    "offset": "-05:00",
                    "current_time": "2024-12-15T14:30:00Z",
                    "is_dst": false
                  },
                  "metadata": {
                    "timestamp": "2024-12-15T14:30:00Z"
                  }
                }
              }
            }
          }
        }
      }
    }
  }
};
