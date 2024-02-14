Feature: StartPayment Endpoint
  As a user of the payment gateway system,
  I want to be able to start a transaction by providing the session ID,
  So that I can initiate payment processing.

  Background:

  Scenario: Successfully start a transaction
    Given the database is empty
    When I send a "POST" request to "/v1/session/" with the following body:
      """
      {
        "sessionID": "session123",
        "merchantID": "AmazonID",
        "amount": {
          "currency": "EUR",
          "value": "100"
        },
        "customerCardInfo": {
          "cardNumber": "1234567890123456",
          "expiryMonth": "12",
          "expiryYear": "25",
          "cvv": "123"
        }
      }
      """
    And I send a "POST" request to "/v1/payment/session123" with the following body:
      """
      """
    Then the response status code should be 200
    And an the response body should be:
      """
      {"status":"SUCCESS","code":"successful-transaction","message":"The transaction finished successfully"}
      """

  Scenario: Attempt to start a transaction with an invalid security code
    Given the database is empty
    When I send a "POST" request to "/v1/session/" with the following body:
      """
      {
        "sessionID": "session123",
        "merchantID": "AmazonID",
        "amount": {
          "currency": "EUR",
          "value": "100"
        },
        "customerCardInfo": {
          "cardNumber": "invalid-cvc",
          "expiryMonth": "12",
          "expiryYear": "25",
          "cvv": "123"
        }
      }
      """
    And I send a "POST" request to "/v1/payment/session123" with the following body:
      """
      """
    Then the response status code should be 200
    And an the response body should be:
      """
      {"status":"FAILURE","code":"invalid-cvc","message":"Invalid CVC"}
      """