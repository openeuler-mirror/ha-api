/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * ha-api is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: liqiuyu
 * Date: 2022-04-19 16:49:51
 * LastEditTime: 2024-03-25 17:19:21
 * Description: 删除重复数据，获取数字单位
 ******************************************************************************/
package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
)

func IsInSlice(str string, sli []string) bool {
	//TODO

	for _, item := range sli {
		if item == str {
			return true
		}
	}
	return false
}

// RemoveDupl remove duplicates in string array
func RemoveDupl(strs []string) []string {
	strSet := map[string]bool{}
	for _, v := range strs {
		strSet[v] = true
	}
	strsDupl := []string{}
	for k := range strSet {
		strsDupl = append(strsDupl, k)
	}
	return strsDupl
}

// GetNumAndUnitFromStr gets the first number and the unit after this number
// like "20.5min" ==> ["20.5", "min"]
func GetNumAndUnitFromStr(s string) (string, string) {
	r := regexp.MustCompile("[0-9](.*)[0-9]")
	index := r.FindStringIndex(s)
	if len(index) == 0 {
		return s[:1], s[1:]
	}
	return s[:index[1]], s[index[1]:]
}

// Send http request
func SendRequest(url string, method string, data interface{}) (resp *http.Response, err error) {
	var httpResp *http.Response
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	switch method {
	case "POST":
		httpResp, err = client.Post(url, "application/json", bytes.NewReader(jsonData))
	case "GET":
		httpResp, _ = client.Get(url)
	case "DELETE":
		req, _ := http.NewRequest("DELETE", url, bytes.NewReader(jsonData))
		httpResp, err = client.Do(req)
	case "PUT":
		req, _ := http.NewRequest("PUT", url, bytes.NewReader(jsonData))
		httpResp, err = client.Do(req)
	default:
		return nil, errors.New("unsupported method")
	}

	if err != nil {
		return nil, err
	}

	return httpResp, nil
}
