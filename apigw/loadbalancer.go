package apigw

import (
    "encoding/json"

    "github.com/hashicorp/terraform-plugin-sdk/helper/resource"
    "github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func lbMembersDiffFunc(k, old, new string, d *schema.ResourceData) bool {
    oldMembers, newMembers := d.GetChange("members")
    if oldMembers == nil && newMembers != nil{
       return false
    }

    if oldMembers != nil && newMembers == nil{
        return false
    }

    oldArray := oldMembers.([]interface{})
    newArray := newMembers.([]interface{})
    if len(oldArray) != len(newArray) {
        return false
    }

    equalCount := 0
    for _, x := range newArray{
        newObject := x.(map[string]interface{})
        newIP := newObject["ip"].(string)
        newPort := newObject["port"].(int)
        newWeight := newObject["weight"].(int)
        for _, y := range oldArray {
            oldObject := y.(map[string]interface{})
            oldIP := oldObject["ip"].(string)
            oldPort := oldObject["port"].(int)
            oldWeight := oldObject["weight"].(int)
            if newIP == oldIP && newPort == oldPort && newWeight == oldWeight{
                equalCount += 1
                break
            }
        }
    }
    return equalCount == len(oldArray)
}

func flattenLBMonitorInfo(v map[string]interface{}) []interface{} {
    monitorInfo := make([]interface{}, 1)
    info := make(map[string]interface{})
    info["delay"] = int(v["delay"].(float64))
    if expectedCodes, ok := v["expected_codes"].(string); ok {
        info["expected_codes"] = expectedCodes
    } else {
        info["expected_codes"] = ""
    }

    if httpMethod, ok := v["http_method"].(string); ok {
        info["http_method"] = httpMethod
    } else {
        info["http_method"] = ""
    }

    info["max_retries"] = int(v["max_retries"].(float64))
    info["monitor_type"] = v["monitor_type"].(string)
    info["timeout"] = int(v["timeout"].(float64))
    if urlPath, ok := v["url_path"].(string); ok {
        info["url_path"] = urlPath
    } else {
        info["url_path"] = ""
    }

    monitorInfo[0] = info
    return monitorInfo
}

func lbStateRefreshFunc(
        config *PConfig,
        host string,
        resourcePath string) resource.StateRefreshFunc {
    return func() (interface{}, string, error) {
        response, err := config.doNormalRequest(host, resourcePath, "GET", nil)
        if err != nil {
            return nil, "", err
        }

        var data map[string]interface{}
        err = json.Unmarshal([]byte(response), &data)

        if err != nil {
            return nil, "", err
        }

        return data, data["status"].(string), nil
    }
}

func lbStateRefreshForDeletedFunc(
        config *PConfig,
        host string, 
        resourcePath string) resource.StateRefreshFunc {
    return func() (interface{}, string, error) {
        response, err := config.doNormalRequest(host, resourcePath, "GET", nil)

        if err != nil {
            if _, ok := err.(ErrDefault404); ok {
                return response, "DELETED", nil
            }
            return response, "", err
        }

        var data map[string]interface{}
        err = json.Unmarshal([]byte(response), &data)

        if err != nil {
            return nil, "", err
        }

        return data, data["status"].(string), nil
    }
}
