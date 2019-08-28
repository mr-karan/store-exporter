package main

import (
	"fmt"
)

// constructMetricData takes a map of column names with corresponding values and returns in
// format of Prometheus metric value and lables
func constructMetricData(data map[string]interface{}, valueName string, labels []string) (float64, []string, error) {
	var (
		value float64
		err   error
	)
	// if column name is in the result set, get the value
	if i, ok := data[valueName]; ok {
		value = getFloatValue(i)
	}
	labelValues := []string{}
	// iterate over labels and extract the value as k/v pair to construct labels
	for _, label := range labels {
		// fallback empty value in case column name doesn't match from config with the result set
		// This is done to ensure label cardinality
		lv := ""
		if i, ok := data[label]; ok {
			lv = getStringValue(i)
			// in case column name matches but type is incorrect, return error to let the user know
			// and not silently fail.
			if lv == "" {
				err = fmt.Errorf("Column: %s must be type text/string", label)
			}
		}
		labelValues = append(labelValues, lv)
	}
	return value, labelValues, err
}
