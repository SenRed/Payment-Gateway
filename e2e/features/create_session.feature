Feature: CreateSession Endpoint
  As a user of the payment gateway system,
  I want to be able to create a new session for a transaction,
  So that I can initiate payment processing.

  Background:

 Scenario: Successfully create a new session
   Given the database is empty
   When I send a "POST" request to "/v1/session/" with the following body:
     """
     {
       "sessionID": "session123",
       "merchantID": "exampleMerchant",
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
   Then the response status code should be 200
   And no error message should be returned

  Scenario: Attempt to create a session missing information
    When I send a "POST" request to "/v1/session/" with the following body:
    """
    {
      "sessionID": "session123",
      "customerCardInfo": {
        "cardNumber": "1234567890123456",
        "expiryYear": "25",
        "cvv": "123"
      }
    }
    """
    Then the response status code should be 400
    And an the response body should be:
      """
      {
        "reason":"MerchantID is required"
      }
      """


  Scenario: Attempt to create a session with a duplicated session ID
    Given the following session information is provided:
      | SessionID  | MerchantID | Amount.Currency | Amount.Value | CustomerCardInfo.CardNumber | CustomerCardInfo.ExpiryMonth | CustomerCardInfo.ExpiryYear | CustomerCardInfo.SecurityCode |
      | session123 | AmazonID   | EUR             | "100"        | "1234567890123456"          | "12"                         | "25"                        | "123"                         |

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
    Then the response status code should be 400
    And an the response body should be:
      """
      {
        "reason":"A session already exist with this id"
      }
      """