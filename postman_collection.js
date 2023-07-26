{
	"info": {
		"_postman_id": "9bcc1cef-6379-4af5-ae18-538ca16fb1ea",
		"name": "Invoice Bidder",
		"schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		"_exporter_id": "21280865"
	},
	"item": [
		{
			"name": "Issuer",
			"item": [
				{
					"name": "Create issuer",
					"event": [
						{
							"listen": "test",
							"script": {
								"exec": [
									"pm.test(\"Response status code is 201\", function () {",
									"  pm.expect(pm.response.code).to.equal(201);",
									"});",
									"",
									"",
									"pm.test(\"Response has the required fields\", function () {",
									"  const responseData = pm.response.json();",
									"  ",
									"  pm.expect(responseData).to.be.an('object');",
									"  pm.expect(responseData.id).to.exist.and.to.be.a('string', \"id should be a string\");",
									"  pm.expect(responseData.fullName).to.exist.and.to.be.a('string', \"fullName should be a string\");",
									"});",
									"",
									"",
									"pm.test(\"Id is a non-empty string\", function () {",
									"    const responseData = pm.response.json();",
									"    ",
									"    pm.expect(responseData.id).to.be.a('string').and.to.have.lengthOf.at.least(1, \"Value should not be empty\");",
									"});",
									"",
									"",
									"pm.test(\"Test fullName is a non-empty string\", function () {",
									"    const responseData = pm.response.json();",
									"    ",
									"    pm.expect(responseData.fullName).to.be.a('string').and.to.have.lengthOf.at.least(1, \"Value should not be empty\");",
									"});",
									"",
									"",
									"pm.test(\"Response time is less than 500ms\", function () {",
									"  pm.expect(pm.response.responseTime).to.be.below(500);",
									"});"
								],
								"type": "text/javascript"
							}
						}
					],
					"request": {
						"method": "POST",
						"header": [],
						"body": {
							"mode": "raw",
							"raw": "{\n    \"fullName\": \"Manuel Adalid Moya\"\n}",
							"options": {
								"raw": {
									"language": "json"
								}
							}
						},
						"url": {
							"raw": "localhost:8080/ib/issuer",
							"host": [
								"localhost"
							],
							"port": "8080",
							"path": [
								"ib",
								"issuer"
							]
						}
					},
					"response": []
				},
				{
					"name": "Retrieve Issuer",
					"event": [
						{
							"listen": "test",
							"script": {
								"exec": [
									"pm.test(\"Response status code is 200\", function () {",
									"    pm.response.to.have.status(200);",
									"});",
									"",
									"",
									"pm.test(\"Response has required fields\", function () {",
									"    const responseData = pm.response.json();",
									"    ",
									"    pm.expect(responseData).to.be.an('object');",
									"    pm.expect(responseData.id).to.exist;",
									"    pm.expect(responseData.fullName).to.exist;",
									"    pm.expect(responseData.balance).to.exist;",
									"    pm.expect(responseData.invoices).to.exist;",
									"});",
									"",
									"",
									"pm.test(\"id is a non-empty string\", function () {",
									"    const responseData = pm.response.json();",
									"    ",
									"    pm.expect(responseData.id).to.be.a('string').and.to.have.lengthOf.at.least(1, \"id should not be empty\");",
									"});",
									"",
									"",
									"pm.test(\"fullName is a non-empty string\", function () {",
									"  const responseData = pm.response.json();",
									"  ",
									"  pm.expect(responseData.fullName).to.be.a('string').and.to.have.lengthOf.at.least(1, \"fullName should not be empty\");",
									"});",
									"",
									"",
									"pm.test(\"Balance is a non-negative number\", function () {",
									"    const responseData = pm.response.json();",
									"    ",
									"    pm.expect(responseData).to.be.an('object');",
									"    pm.expect(responseData.balance).to.be.a('number');",
									"    pm.expect(responseData.balance).to.be.at.least(0, \"Balance should be non-negative\");",
									"});"
								],
								"type": "text/javascript"
							}
						}
					],
					"request": {
						"method": "GET",
						"header": [],
						"url": {
							"raw": "localhost:8080/issuer/:id",
							"host": [
								"localhost"
							],
							"port": "8080",
							"path": [
								"issuer",
								":id"
							],
							"variable": [
								{
									"key": "id",
									"value": null
								}
							]
						}
					},
					"response": []
				}
			]
		},
		{
			"name": "Investor",
			"item": [
				{
					"name": "Create Investor",
					"event": [
						{
							"listen": "test",
							"script": {
								"exec": [
									"pm.test(\"Response status code is 201\", function () {",
									"    pm.response.to.have.status(201);",
									"});",
									"",
									"",
									"pm.test(\"Response has the required fields\", function () {",
									"    const responseData = pm.response.json();",
									"    ",
									"    pm.expect(responseData).to.be.an('object');",
									"    pm.expect(responseData.id).to.exist.and.to.be.a('string');",
									"    pm.expect(responseData.fullName).to.exist.and.to.be.a('string');",
									"    pm.expect(responseData.balance).to.exist.and.to.be.a('number');",
									"});",
									"",
									"",
									"pm.test(\"The id is a non-empty string\", function () {",
									"    const responseData = pm.response.json();",
									"    ",
									"    pm.expect(responseData).to.be.an('object');",
									"    pm.expect(responseData.id).to.exist.and.to.be.a('string').and.to.have.lengthOf.at.least(1, \"Value should not be empty\");",
									"});",
									"",
									"",
									"pm.test(\"fullName is a non-empty string\", function () {",
									"    const responseData = pm.response.json();",
									"    ",
									"    pm.expect(responseData).to.be.an('object');",
									"    pm.expect(responseData.fullName).to.be.a('string').and.to.have.lengthOf.at.least(1, \"Value should not be empty\");",
									"});",
									"",
									"",
									"pm.test(\"Balance is a non-negative number\", function () {",
									"    const responseData = pm.response.json();",
									"",
									"    pm.expect(responseData).to.be.an('object');",
									"    pm.expect(responseData.balance).to.exist.and.to.be.a('number').and.to.be.at.least(0);",
									"});"
								],
								"type": "text/javascript"
							}
						}
					],
					"request": {
						"method": "POST",
						"header": [],
						"body": {
							"mode": "raw",
							"raw": "{\n    \"fullName\": \"Another Investor\",\n    \"balance\": {\n        \"amount\": \"5000\",\n        \"currency\": \"USD\"\n    }\n}",
							"options": {
								"raw": {
									"language": "json"
								}
							}
						},
						"url": {
							"raw": "localhost:8080/investor/",
							"host": [
								"localhost"
							],
							"port": "8080",
							"path": [
								"investor",
								""
							]
						}
					},
					"response": []
				},
				{
					"name": "List Investors",
					"event": [
						{
							"listen": "test",
							"script": {
								"exec": [
									"pm.test(\"Response status code is 200\", function () {",
									"  pm.response.to.have.status(200);",
									"});",
									"",
									"",
									"pm.test(\"The response should be an array\", function () {",
									"    pm.expect(pm.response.json()).to.be.an('array');",
									"});",
									"",
									"",
									"pm.test(\"Each object in the array has the required fields - id, fullName, and balance\", function () {",
									"    const responseData = pm.response.json();",
									"",
									"    pm.expect(responseData).to.be.an('array').that.is.not.empty;",
									"    responseData.forEach(function (object) {",
									"        pm.expect(object.id).to.exist.and.to.be.a('string');",
									"        pm.expect(object.fullName).to.exist.and.to.be.a('string');",
									"        pm.expect(object.balance).to.exist.and.to.be.a('string');",
									"    });",
									"});",
									"",
									"",
									"pm.test(\"id is a non-empty string\", function () {",
									"    const responseData = pm.response.json();",
									"",
									"    pm.expect(responseData).to.be.an('array');",
									"    responseData.forEach(function (item) {",
									"        pm.expect(item.id).to.be.a('string').and.to.have.lengthOf.at.least(1, \"id should be a non-empty string\");",
									"    });",
									"});",
									"",
									"",
									"pm.test(\"fullName is a non-empty string\", function () {",
									"    const responseData = pm.response.json();",
									"",
									"    pm.expect(responseData).to.be.an('array');",
									"    responseData.forEach(function (investor) {",
									"        pm.expect(investor.fullName).to.be.a('string').and.to.have.lengthOf.at.least(1, \"fullName should not be empty\");",
									"    });",
									"});"
								],
								"type": "text/javascript"
							}
						}
					],
					"request": {
						"method": "GET",
						"header": [],
						"url": {
							"raw": "localhost:8080/investor",
							"host": [
								"localhost"
							],
							"port": "8080",
							"path": [
								"investor"
							],
							"query": [
								{
									"key": "ids",
									"value": "da23bee6-2ba2-11ee-ae5e-f8e43b7adc9f,c80ebae3-2ba2-11ee-ae5e-f8e43b7adc9f",
									"disabled": true
								}
							]
						}
					},
					"response": []
				}
			]
		},
		{
			"name": "Invoice",
			"item": [
				{
					"name": "Create Invoice",
					"event": [
						{
							"listen": "test",
							"script": {
								"exec": [
									"pm.test(\"Response status code is 400\", function () {",
									"    pm.expect(pm.response.code).to.equal(400);",
									"});",
									"",
									"",
									"pm.test(\"Error field is present in the response\", function () {",
									"  const responseData = pm.response.json();",
									"  ",
									"  pm.expect(responseData.error).to.exist;",
									"});",
									"",
									"",
									"pm.test(\"Error field is a non-empty string\", function () {",
									"  const responseData = pm.response.json();",
									"  ",
									"  pm.expect(responseData).to.be.an('object');",
									"  pm.expect(responseData.error).to.be.a('string').and.to.have.lengthOf.at.least(1, \"Error field should not be empty\");",
									"});",
									"",
									"",
									"pm.test(\"Response time is less than 500ms\", function () {",
									"    pm.expect(pm.response.responseTime).to.be.below(500);",
									"});",
									"",
									"",
									"pm.test(\"Validate that the request body is not empty\", function () {",
									"  const requestBody = pm.request.body;",
									"",
									"  pm.expect(requestBody).to.exist.and.to.not.be.empty;",
									"});"
								],
								"type": "text/javascript"
							}
						}
					],
					"request": {
						"method": "POST",
						"header": [],
						"body": {
							"mode": "formdata",
							"formdata": [
								{
									"key": "invoice",
									"type": "file",
									"src": "/home/manu/Downloads/Tech_challenge_-_invoice.pdf"
								},
								{
									"key": "issuer_id",
									"value": "0896579d-2ba0-11ee-a338-f8e43b7adc9f",
									"type": "text"
								},
								{
									"key": "price",
									"value": "1250",
									"type": "text"
								},
								{
									"key": "currency",
									"value": "EUR",
									"type": "text"
								}
							]
						},
						"url": {
							"raw": "localhost:8080/invoice/",
							"host": [
								"localhost"
							],
							"port": "8080",
							"path": [
								"invoice",
								""
							]
						}
					},
					"response": []
				},
				{
					"name": "Bid",
					"event": [
						{
							"listen": "test",
							"script": {
								"exec": [
									"pm.test(\"Response status code is 500\", function () {",
									"  pm.expect(pm.response.code).to.equal(500);",
									"});",
									"",
									"",
									"pm.test(\"Response has the required field - error\", function () {",
									"    const responseData = pm.response.json();",
									"    ",
									"    pm.expect(responseData).to.be.an('object');",
									"    pm.expect(responseData.error).to.exist.and.to.be.a('string');",
									"});",
									"",
									"",
									"pm.test(\"Error message is not empty\", function () {",
									"    const responseData = pm.response.json();",
									"    ",
									"    pm.expect(responseData.error).to.exist.and.to.not.be.empty;",
									"});",
									"",
									"",
									"pm.test(\"Response time is less than 500ms\", function () {",
									"  pm.expect(pm.response.responseTime).to.be.below(500);",
									"});",
									"",
									"",
									"pm.test(\"Validate request URL\", function () {",
									"    pm.expect(pm.request.url).to.equal(\"localhost:8080/invoice/ba8a44ed-2ba1-11ee-a117-f8e43b7adc9f/bid\");",
									"});"
								],
								"type": "text/javascript"
							}
						}
					],
					"request": {
						"method": "POST",
						"header": [],
						"body": {
							"mode": "raw",
							"raw": "{\n    \"investorId\": \"da23bee6-2ba2-11ee-ae5e-f8e43b7adc9f\",\n    \"amount\": {\n        \"amount\": \"3000\",\n        \"currency\": \"EUR\"\n    }\n}",
							"options": {
								"raw": {
									"language": "json"
								}
							}
						},
						"url": {
							"raw": "localhost:8080/invoice/:id/bid",
							"host": [
								"localhost"
							],
							"port": "8080",
							"path": [
								"invoice",
								":id",
								"bid"
							],
							"variable": [
								{
									"key": "id",
									"value": "ba8a44ed-2ba1-11ee-a117-f8e43b7adc9f"
								}
							]
						}
					},
					"response": []
				},
				{
					"name": "Trade",
					"event": [
						{
							"listen": "test",
							"script": {
								"exec": [
									"pm.test(\"Response status code is 500\", function () {",
									"  pm.expect(pm.response.code).to.equal(500);",
									"});",
									"",
									"",
									"pm.test(\"Response has the 'error' field\", function () {",
									"    const responseData = pm.response.json();",
									"    ",
									"    pm.expect(responseData.error).to.exist;",
									"});",
									"",
									"",
									"pm.test(\"Error field is a non-empty string\", function () {",
									"  const responseData = pm.response.json();",
									"  ",
									"  pm.expect(responseData.error).to.be.a('string').and.to.have.lengthOf.at.least(1, \"Error field should not be empty\");",
									"});",
									"",
									"",
									"pm.test(\"Response time is less than 500ms\", function () {",
									"    pm.expect(pm.response.responseTime).to.be.below(500);",
									"});",
									"",
									"",
									"pm.test(\"Verify the request URL is correct\", function () {",
									"    pm.expect(pm.request.url).to.equal(\"localhost:8080/invoice/ba8a44ed-2ba1-11ee-a117-f8e43b7adc9f/trade\");",
									"});"
								],
								"type": "text/javascript"
							}
						}
					],
					"request": {
						"method": "POST",
						"header": [],
						"body": {
							"mode": "raw",
							"raw": "{\n    \"approve\": false\n}",
							"options": {
								"raw": {
									"language": "json"
								}
							}
						},
						"url": {
							"raw": "localhost:8080/invoice/:id/trade",
							"host": [
								"localhost"
							],
							"port": "8080",
							"path": [
								"invoice",
								":id",
								"trade"
							],
							"variable": [
								{
									"key": "id",
									"value": "ba8a44ed-2ba1-11ee-a117-f8e43b7adc9f"
								}
							]
						}
					},
					"response": []
				},
				{
					"name": "Retrieve Invoice",
					"event": [
						{
							"listen": "test",
							"script": {
								"exec": [
									"pm.test(\"Response status code is 200\", function () {",
									"  pm.response.to.have.status(200);",
									"});",
									"",
									"",
									"pm.test(\"Response has the required fields\", function () {",
									"    const responseData = pm.response.json();",
									"    ",
									"    pm.expect(responseData.id).to.exist;",
									"    pm.expect(responseData.price).to.exist;",
									"    pm.expect(responseData.status).to.exist;",
									"    pm.expect(responseData.issuer).to.exist;",
									"    pm.expect(responseData.bids).to.exist;",
									"});",
									"",
									"",
									"pm.test(\"Id is a non-empty string\", function () {",
									"    const responseData = pm.response.json();",
									"",
									"    pm.expect(responseData.id).to.be.a('string').and.to.have.lengthOf.at.least(1, \"Id should not be empty\");",
									"});",
									"",
									"",
									"pm.test(\"Price is a number greater than or equal to 0\", function () {",
									"    const responseData = pm.response.json();",
									"    ",
									"    pm.expect(responseData.price).to.exist.and.to.be.a('number');",
									"    pm.expect(responseData.price).to.be.at.least(0);",
									"});",
									"",
									"",
									"pm.test(\"Status is one of the expected values\", function () {",
									"    const responseData = pm.response.json();",
									"",
									"    pm.expect(responseData.status).to.exist.and.to.be.oneOf([\"pending\", \"completed\", \"cancelled\"]);",
									"});"
								],
								"type": "text/javascript"
							}
						}
					],
					"request": {
						"method": "GET",
						"header": [],
						"url": {
							"raw": "localhost:8080/invoice/:id",
							"host": [
								"localhost"
							],
							"port": "8080",
							"path": [
								"invoice",
								":id"
							],
							"variable": [
								{
									"key": "id",
									"value": "ba8a44ed-2ba1-11ee-a117-f8e43b7adc9f"
								}
							]
						}
					},
					"response": []
				}
			]
		}
	]
}