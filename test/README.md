# cloudscale tests

This directory contains integration tests, aside from the unit tests in the 
root directory. While the unit tests suite runs very quickly because they
don't make any network calls, this can take some time to run.

## integration

This folder contains tests for every type of operation in the Cloud Scale API
and runs tests against it.

Since te tests are run against live data, there is a higher chance of false
positives and test failures due to network issues, data changes, etc.

Run the tests using: 

````
CLOUDSCALE_TOKEN="HELPIMTRAPPEDINATOKENGENERATOR" go test -v -tags=integration  ./integration

````