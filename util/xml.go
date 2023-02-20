/*
 * @Description:
 * @Author: yuanshisan
 * @Date: 2021-12-02 17:31:51
 * @LastEditTime: 2023-02-13 13:30:30
 * @LastEditors: yuanshisan
 */
package util

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type xmlProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

type xmlGroup struct {
	GroupName string        `xml:"name,attr"`
	Property  []xmlProperty `xml:"property"`
}

type xmlGroups struct {
	GroupsName string     `xml:"name,attr"`
	Group      []xmlGroup `xml:"group"`
}

type xmlConfig struct {
	Property []xmlProperty `xml:"property"`
	Groups   []xmlGroups   `xml:"groups"`
	Group    []xmlGroup    `xml:"group"`
}

func LoadXmlConfig(filename string) map[string]interface{} {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "LoadXmlConfig: Error: Could not open %q for reading: %s\n", filename, err)
		os.Exit(1)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "LoadXmlConfig: Error: Could not read %q: %s\n", filename, err)
		os.Exit(1)
	}
	xc := new(xmlConfig)
	if err := xml.Unmarshal(data, xc); err != nil {
		fmt.Fprintf(os.Stderr, "LoadXmlConfig: Error: Could not parse XML configuration in %q: %s\n", filename, err)
		os.Exit(1)
	}
	configMap := make(map[string]interface{})
	formatXmlProp(xc.Property, configMap)
	formatXmlGroup(xc.Group, configMap)
	formatXmlGroups(xc.Groups, configMap)
	return configMap
}

func formatXmlProp(propList []xmlProperty, originMap map[string]interface{}) {
	for _, xcProp := range propList {
		originMap[xcProp.Name] = xcProp.Value
	}
}

func formatXmlGroup(groupList []xmlGroup, originMap map[string]interface{}) {
	for idx, xcGroup := range groupList {
		var key string
		if len(xcGroup.GroupName) == 0 {
			key = strconv.Itoa(idx)
		} else {
			key = xcGroup.GroupName
		}
		tempMap := make(map[string]interface{})
		formatXmlProp(xcGroup.Property, tempMap)
		originMap[key] = tempMap
	}
}

func formatXmlGroups(groupsList []xmlGroups, originMap map[string]interface{}) {
	for _, xcGroups := range groupsList {
		tempMap := make(map[string]interface{})
		formatXmlGroup(xcGroups.Group, tempMap)
		originMap[xcGroups.GroupsName] = tempMap
	}
}
