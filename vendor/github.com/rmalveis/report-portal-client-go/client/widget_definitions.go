package client

var WidgetTypes = map[string]string{
	// TODO: allow types below to be managed by  this client
	//"casesTrend",
	//"notPassed",
	//"investigatedTrend",
	//"launchStatistics",
	//"uniqueBugTable",
	//"activityStream",
	//"launchesComparisonChart",
	//"launchesTable",
	//"passingRateSummary",
	//"passingRatePerLaunch",
	//"productStatus",
	//"mostTimeConsuming",
	//"cumulative",
	//"componentHealthCheck",
	//"topPatternTempaltes",
	//"componentHealthCheckTable",
	"oldLineChart":                          "oldLineChart", // fallback: when a type is not recognized, this value is forced
	"Launch statistics chart":               "statisticTrend",
	"Launch duration chart":                 "launchesDurationChart",
	"Failed cases trend chart":              "bugTrend",
	"Overall statistics":                    "overallStatistics",
	"Most failed test-cases table (TOP-20)": "topTestCases",
	"Flaky test cases table (TOP-20)":       "flakyTestCases",
}
var WidgetCriteria = map[string]string{
	"Total":                "statistics$executions$total",
	"Passed":               "statistics$executions$passed",
	"Failed":               "statistics$executions$failed",
	"Skipped":              "statistics$executions$skipped",
	"Product Bug":          "statistics$defects$product_bug$pb001",
	"Automation Bug":       "statistics$defects$automation_bug$ab001",
	"System Issue":         "statistics$defects$system_issue$si001",
	"No Defect":            "statistics$defects$no_defect$nd001",
	"To Investigate":       "statistics$defects$to_investigate$ti001",
	"Product Bug Total":    "statistics$defects$product_bug$total",
	"Automation Bug Total": "statistics$defects$automation_bug$total",
	"System Issue Total":   "statistics$defects$system_issue$total",
	"No Defect Total":      "statistics$defects$no_defect$total",
	"To Investigate Total": "statistics$defects$to_investigate$total",
	"Start time":           "startTime",
	"End Time":             "endTime",
	"Name":                 "name",
	"Number":               "number",
	"Status":               "status",
}
var WidgetVisualizationOptions = map[string]string{
	"Area": "area-spline",
	"Bars": "bars",
}
var WidgetModes = map[string]string{
	"Launch":   "launch",
	"timeline": "day",
}
