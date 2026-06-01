package examplesuite

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// Run this example: go test example_test.go -v

type ExampleSuite struct {
	suite.Suite
}

// TestExampleSuite is the thing that calls all the actual suite tests to be run. It is required by the testify suite package.
func TestExampleSuite(t *testing.T) {
	suite.Run(t, &ExampleSuite{})
}

// These 2 functions run before and after the entire suite
func (ts *ExampleSuite) SetupSuite() {
	ts.T().Log("SetupSuite")
}

func (ts *ExampleSuite) TearDownSuite() {
	ts.T().Log("TearDownSuite")
}
//////////////////////////////////////////////////////////////////////////////////////////

// These 2 functions run before and after each test, and also happen before/after the BeforeTest and AfterTest functions.
// Unsure why both sets actually exist but there you go
func (ts *ExampleSuite) SetupTest() {
	ts.T().Log("SetupTest")
}

func (ts *ExampleSuite) TearDownTest() {
	ts.T().Log("TearDownTest")
}
////////////////////////////////////////////////////////////////////////////////////////////

// These 2 functions run before and after each test, inside the Setup/Teardown ones
// Unsure why both sets actually exist but there you go
func (ts *ExampleSuite) AfterTest(suiteName, testName string) {
	ts.T().Logf("AfterTest: %s - %s", suiteName, testName)
}

func (ts *ExampleSuite) BeforeTest(suiteName, testName string) {
	ts.T().Logf("BeforeTest: %s - %s", suiteName, testName)
}
////////////////////////////////////////////////////////////////////////////////////////////

// Needed so there is an actual test to run, otherwise the suite won't run at all
func (ts *ExampleSuite) Test1() {
	ts.T().Log("ACTUAL TEST RUNS")
}
////////////////////////////////////////////////////////////////////////////////////////////


// Sub test example
// SUB tests just dont seem to work with testify

func (ss *ExampleSuite) SetupSubTest() {
	ss.T().Log("WONT BE SHOWN SetupSubTest")
}

func (ss *ExampleSuite) TearDownSubTest() {
	ss.T().Log("WONT BE SHOWN TearDownSubTest")
}

func (ss *ExampleSuite) tESTSubTest1() {
	ss.T().Run("FIRST", func(t *testing.T) {
		ss.T().Log("FIRST SUBTEST RAN")
	})
	ss.T().Run("SECOND", func(t *testing.T) {
		ss.T().Log("SECOND SUBTEST RAN")
	})
}

