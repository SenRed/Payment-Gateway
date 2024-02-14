package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/payment-gateway/internal/bootstrap"
	pgInfra "github.com/payment-gateway/internal/infrastructure/postgres"
	"github.com/rs/zerolog/log"
	"github.com/wI2L/jsondiff"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var tContext TestingContext

type TestingContext struct {
	app  *bootstrap.Bootstrap
	resp *httptest.ResponseRecorder
	gorm *gorm.DB
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer:  InitializeScenario,
		TestSuiteInitializer: InitializeTestSuite,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}

	CloseResources()
}

func initAPI() {
	newApp := bootstrap.Init()

	tContext = TestingContext{
		app: &newApp,
	}
	tContext.initDatabaseConnexion()
	tContext.resp = httptest.NewRecorder()
}

func cleanUp() {
	tContext.gorm.Exec("TRUNCATE TABLE transaction_entities CASCADE ")
	tContext.gorm.Exec("TRUNCATE TABLE session_entities CASCADE")
}
func (tCtx *TestingContext) initDatabaseConnexion() {
	pgConfig := tCtx.app.Config.PostgresConfig
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		pgConfig.Host, pgConfig.User, pgConfig.Password, pgConfig.Database, pgConfig.Port, pgConfig.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic(err)
	}
	tCtx.gorm = db
}

func CloseResources() {
	if tContext.gorm != nil {
		db, err := tContext.gorm.DB()
		if err != nil {
			log.Error().Err(err).Msg("Error getting database connection")
			return
		}
		err = db.Close()
		if err != nil {
			panic(err)
		}
		log.Info().Msg("Database connection closed")
	}
}
func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(initAPI)
}

func InitializeScenario(s *godog.ScenarioContext) {
	// Given steps
	s.Step(`^the database is empty$`, tContext.theDatabaseIsEmpty)
	s.Step(`^the following session information is provided:$`, tContext.theFollowingSessionInformationIsProvided)

	// When steps
	s.Step(`^I send a "(GET|POST)" request to "([^"]*)" with the following body:$`, tContext.iSendARequestToWithTheFollowingSessionInformation)

	// Assertion steps
	s.Step(`^the response status code should be (\d+)$`, tContext.theResponseStatusCodeShouldBe)
	s.Step(`^no error message should be returned$`, tContext.noErrorMessageShouldBeReturned)
	s.Step(`^an the response body should be:$`, tContext.responseBodyShouldBe)
}

func (tCtx *TestingContext) theDatabaseIsEmpty() error {
	// clean tables
	cleanUp()
	return nil
}
func (tCtx *TestingContext) theFollowingSessionInformationIsProvided(table *godog.Table) error {
	// clean tables
	cleanUp()
	// Parse the provided session information from the table
	var sessionEntity pgInfra.SessionEntity

	for _, row := range table.Rows {
		// Map each row of the table to the sessionEntity fields
		sessionEntity.ID = row.Cells[0].Value
		sessionEntity.MerchantID = row.Cells[1].Value
		sessionEntity.Amount.Currency = row.Cells[2].Value
		sessionEntity.Amount.Value = row.Cells[3].Value

		// Parse customer card info
		sessionEntity.CustomerCardInfo.CardNumber = row.Cells[4].Value
		sessionEntity.CustomerCardInfo.ExpiryMonth = row.Cells[5].Value
		sessionEntity.CustomerCardInfo.ExpiryYear = row.Cells[6].Value
		sessionEntity.CustomerCardInfo.SecurityCode = row.Cells[7].Value

		// Save sessionEntity to the database
		if err := tCtx.gorm.Create(&sessionEntity).Error; err != nil {
			return err
		}
	}

	return nil
}

func (tCtx *TestingContext) iSendARequestToWithTheFollowingSessionInformation(verb string, endpoint string, jsonString *godog.DocString) error {
	// Create a new HTTP request to the specified endpoint with the session information in the body
	req, err := http.NewRequest(verb, endpoint, bytes.NewBuffer([]byte(jsonString.Content)))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// create a httpResponse recorder
	tCtx.resp = httptest.NewRecorder()

	// execute de request and save the response
	tCtx.app.Router.ServeHTTP(tCtx.resp, req)

	return nil
}

func (tCtx *TestingContext) theResponseStatusCodeShouldBe(expectedStatusCode int) error {
	// Verify the response status code
	if tCtx.resp.Code != expectedStatusCode {

		return fmt.Errorf("expected status code %d, got %d body: %v", expectedStatusCode, tCtx.resp.Code, tCtx.resp.Body)
	}
	return nil
}

func (tCtx *TestingContext) responseBodyShouldBe(jsonString *godog.DocString) error {
	var expected, actual any

	// re-encode expected response
	if err := json.Unmarshal([]byte(jsonString.Content), &expected); err != nil {
		return err
	}

	// re-encode actual response too
	if err := json.Unmarshal(tCtx.resp.Body.Bytes(), &actual); err != nil {
		return err
	}
	patch, err := jsondiff.Compare(
		actual,
		expected,
		jsondiff.Ignores("/0/date", "/1/date"),
	)
	if err != nil {
		return err
	}
	if len(patch) > 0 {
		fmt.Println("Diff")
		for _, op := range patch {
			fmt.Printf("%s\n", op)
		}
		return fmt.Errorf("expected JSON does not match actual, %v vs. %v", jsonString.Content, tCtx.resp.Body)
	}

	return nil
}
func (tCtx *TestingContext) noErrorMessageShouldBeReturned() error {
	// Verify that no error message is returned in the response body
	respBody, err := io.ReadAll(tCtx.resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}
	if strings.Contains(string(respBody), "error") {
		return fmt.Errorf("unexpected error message in response body: %s", respBody)
	}
	return nil
}
