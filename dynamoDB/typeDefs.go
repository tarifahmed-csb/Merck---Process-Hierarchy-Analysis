/*
 * This package includes code inspired by the work of:
 *
 * Author: Brett Duersch
 * Title: Merck Process Hierarchy Modeling
 * Date: February 16, 2024 
 *
 * Nothing Modified
 */

package main

type ModelOutput struct {
	Hierarchy    Process       `json:"hierarchy"`
	Xpath        []Xpath       `json:"xpath"`
	Metadata     []Metadata    `json:"metadata"`
	Results      []Results     `json:"results"`
	RawMaterials []RawMaterials `json:"rawMaterials"`
}

type Process struct {
	Process string  `json:"process"`
	Stages  []Stage `json:"stages"`
}

type Stage struct {
	Stage       string       `json:"stage"`
	Operations  []Operation  `json:"operations"`
	Measures    []Measure    `json:"measures"`
}

type Operation struct {
	Operation   string      `json:"operation"`
	Actions     []Action    `json:"actions"`
	Measures    []Measure   `json:"measures"`
}

type Action struct {
	Action      string      `json:"action"`
	Measures    []Measure   `json:"measures"`
}

type Measure struct {
	Measure   string `json:"measure"`
	MeasureID string `json:"measureID"`
}

type Xpath struct {
	Xpath     string `json:"xpath"`
	MeasureID string `json:"measureID"`
	Site      string `json:"site"`
}

type Metadata struct {
	MeasureID string `json:"measureID"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

type Results struct {
	Xpath      string  `json:"xpath"`
	Site       string  `json:"site"`
	BatchID    string  `json:"batchID"`
	DOM        string  `json:"DOM"`
	ResultName string  `json:"resultName"`
	Result     float64 `json:"result"`
}

type RawMaterials struct {
	ParentBatchID     string `json:"parentBatchID"`
	ParentMaterialNum string `json:"parentMaterialNum"`
	ChildMaterialName string `json:"childMaterialName"`
	ChildBatchID      string `json:"childBatchID"`
	ChildMaterialNum  string `json:"childMaterialNum"`
}
