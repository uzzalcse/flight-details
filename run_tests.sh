#!/bin/bash

# Define output files
COVERAGE_FILE="coverage.out"
FAILED_TESTS_FILE="failed_tests.log"

# Run Go tests with coverage
if go test ./... -coverprofile=$COVERAGE_FILE 2>&1 | tee test_output.log | grep -q "FAIL"; then
    echo "Tests failed. Logging errors to $FAILED_TESTS_FILE"
    grep "FAIL" test_output.log > $FAILED_TESTS_FILE
else
    echo "All tests passed. No failed test log created."
    rm -f $FAILED_TESTS_FILE
fi

# Display coverage report
go tool cover -func=$COVERAGE_FILE

# Clean up test output log
rm -f test_output.log
