package parser

import (
	"strings"
)

// Given an original value, replace {account_id}, {region_name}, {uuid}
// This is
func Get_parsed_value(original_value string, account_id string, region_name string, uuid string) string {

	static_account_id := "{account_id}"
	static_region_name := "{region_name}"
	static_uuid := "{uuid}"

	result := original_value
	result = strings.Replace(result, static_account_id, account_id, -1)
	result = strings.Replace(result, static_region_name, region_name, -1)
	result = strings.Replace(result, static_uuid, uuid, -1)

	return result
}
