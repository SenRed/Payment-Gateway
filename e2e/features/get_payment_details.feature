Feature: GetPaymentDetails Endpoint
  As a user of the payment gateway system,
  I want to be able to retrieve payment details for a session by providing the session ID,
  So that I can view the transaction history.

  Background:

  Scenario: Successfully retrieve payment details
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
    And I send a "GET" request to "/v1/payment/session123" with the following body:
      """
      """
    Then the response status code should be 200
    And an the response body should be:
      """
        [
          {
            "status": "NOT_STARTED",
            "code": "not-started",
            "message": "Transactions created, not started yet",
            "date": "2024-02-14T19:07:04.038886+01:00"
          },
          {
            "status": "SUCCESS",
            "code": "successful-transaction",
            "message": "Transaction executed successfully",
            "date": "2024-02-14T19:07:04.038886+01:00"
          }
        ]
      """