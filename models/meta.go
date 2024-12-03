/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * ha-api licensed under the Mulan Permissive Software License, Version 2. 
 * See LICENSE file for more details.
 * Author: 赵忠章 <jinzi120021@sina.com>
 * Date: Fri Jan 22 11:18:35 2021 +0800
 */
package models

import (
	"fmt"
	"regexp"
	"strings"

	"gitee.com/openeuler/ha-api/utils"
	"github.com/beevik/etree"
)

func GetAllResourceMetas() map[string]interface{} {
	result := map[string]interface{}{}

	out, err := utils.RunCommand(utils.CmdListResourceStandards)
	if err != nil {
		result["action"] = false
		result["error"] = err.Error()
		return result
	}
	standards := strings.Split(string(out), "\n")
	data := make(map[string]interface{})
	res := make(map[string][]string)

	for _, st := range standards {
		if st == "" {
			continue
		}
		if st == "ocf" {
			out, err := utils.RunCommand(utils.CmdListOcfProviders)
			if err != nil {
				result["action"] = false
				result["error"] = err.Error()
				return result
			}
			pvds := strings.Split(string(out), "\n")
			for _, p := range pvds {
				if p == "" {
					continue
				}
				out, err := utils.RunCommand(fmt.Sprintf(utils.CmdListOcfResourceAgent, p))
				if err != nil {
					result["action"] = false
					result["error"] = err.Error()
					return result
				}
				ag := strings.Split(string(out), "\n")
				// Eliminate Duplicates
				agMap := map[string]bool{}
				for _, agStr := range ag {
					agMap[agStr] = true
				}
				ag = []string{}
				for k := range agMap {
					ag = append(ag, k)
				} // Eliminate Duplicates Over
				res[p] = ag
			}
			data["ocf"] = res
		} else if st == "lsb" {
			continue
		} else {
			out, err := utils.RunCommand(fmt.Sprintf(utils.CmdListResourceAgent, st))
			if err != nil {
				result["action"] = false
				result["error"] = err.Error()
				return result
			}
			la := strings.Split(string(out), "\n")
			// Eliminate Duplicates
			laMap := map[string]bool{}
			for _, laStr := range la {
				if laStr != "" {
					laMap[laStr] = true
				}
			}
			la = []string{}
			for k := range laMap {
				la = append(la, k)
			} // Eliminate Duplicates Over
			laLen := len(la)
			for i := laLen - 1; i >= 0; i-- {
				if strings.HasSuffix(la[i], "@") {
					la = append(la[:i], la[(i+1):]...)
				}
			}
			for i, subla := range la {
				if ok, _ := regexp.MatchString(`.*\@`, subla); ok {
					la = append(la[:i], la[(i+1):]...)
				}
			}
			data[st] = la
		}
	}
	result["action"] = true
	result["data"] = data
	return result
}

func GetResourceMetas(rscClass, rscType, rscProvider string) map[string]interface{} {
	data := make(map[string]interface{})
	prop := []map[string]interface{}{}
	actions := []map[string]string{}

	cmd := ""
	if rscProvider == "" {
		cmd = fmt.Sprintf(utils.CmdShowMetaData, rscClass, rscType)
	} else {
		cmd = fmt.Sprintf(utils.CmdShowMetaDataWithProvider, rscClass, rscProvider, rscType)
	}
	out, err := utils.RunCommand(cmd)
	if err != nil {
		return map[string]interface{}{"action": false, "data": err.Error()}
	}
	doc := etree.NewDocument()
	if err = doc.ReadFromBytes(out); err != nil {
		return map[string]interface{}{"action": false, "data": err.Error()}
	}
	eRoot := doc.Root()
	data["name"] = eRoot.SelectAttrValue("name", "")
	parameter := eRoot.FindElements("./parameters/parameter")
	for _, i := range parameter {
		parameters := map[string]interface{}{}
		parameters["name"] = i.SelectAttrValue("name", "")
		parameters["required"] = i.SelectAttrValue("required", "")
		parameters["unique"] = i.SelectAttrValue("unique", "")
		content := i.FindElement("content")
		cnt := map[string]string{}
		cnt["default"] = content.SelectAttrValue("default", "")
		cnt["type"] = content.SelectAttrValue("type", "")
		parameters["content"] = cnt
		if i.FindElement("shortdesc") != nil {
			text := i.FindElement("shortdesc").Text()
			count := strings.Count(text, "\n")
			parameters["shortdesc"] = strings.Replace(text, "\n", " ", count)
		} else {
			parameters["shortdesc"] = ""
		}
		if i.FindElement("longdesc") != nil {
			text := i.FindElement("longdesc").Text()
			count := strings.Count(text, "\n")
			parameters["longdesc"] = strings.Replace(text, "\n", " ", count)
		} else {
			parameters["longdesc"] = ""
		}
		prop = append(prop, parameters)
		if rscClass == "stonith" {
			pcmkHostList := map[string]interface{}{}
			content := map[string]string{"default": "", "type": "string"}
			pcmkHostList["content"] = content
			pcmkHostList["longdesc"] = "A list of machines controlled by this device."
			pcmkHostList["name"] = "pcmk_host_list"
			pcmkHostList["required"] = "1"
			pcmkHostList["shortdesc"] = ""
			pcmkHostList["unique"] = ""
			pcmkHostList["value"] = ""
			prop = append(prop, pcmkHostList)
		}
	}
	actionElems := eRoot.FindElements("./actions/action")
	for _, actionElem := range actionElems {
		act := map[string]string{}
		for _, attr := range actionElem.Attr {
			act[attr.Key] = attr.Value
		}
		actions = append(actions, act)
	}
	if version := eRoot.FindElement("version"); version != nil {
		data["version"] = version.Text()
	} else {
		data["version"] = ""
	}
	if longdesc := eRoot.FindElement("longdesc"); longdesc != nil {
		text := longdesc.Text()
		count := strings.Count(text, "\n")
		data["longdesc"] = strings.Replace(text, "\n", " ", count)
	} else {
		data["longdesc"] = ""
	}
	if shortdesc := eRoot.FindElement("shortdesc"); shortdesc != nil {
		text := shortdesc.Text()
		count := strings.Count(text, "\n")
		data["shortdesc"] = strings.Replace(text, "\n", " ", count)
	} else {
		data["shortdesc"] = ""
	}
	data["parameters"] = prop
	data["actions"] = actions

	if rscType == "fence_sbd" {
		w := data["parameters"].([]map[string]interface{})
		for i, v := range w {
			ret := v["name"]
			if ret == "plug" {
				w = append(w[:i], w[(i+1):]...)
			}
		}
	}
	/* By the end of "pcs resource describe " */
	return map[string]interface{}{"action": true, "data": data}
}
