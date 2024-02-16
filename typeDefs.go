// -------------------------------------------------------------------------
//   table structure types
// -------------------------------------------------------------------------

package main

// ===== types to support the hierarchy output =====
// This is a nested structure: Process >> Stage >> Operation >> Action
// Measures can appear at any level in the hierarchy
// MeasureID is added to simplify connections to the xpath and tags tables

type Measure struct {
	Measure   string `json:"measure"`
	MeasureID string `json:"measuresID"`
}

type Action struct {
	Action   string    `json:"action"`
	Measures []Measure `json:"measures"`
}

type Operation struct {
	Operation string    `json:"operation"`
	Actions   []Action  `json:"actions"`
	Measures  []Measure `json:"measures"`
}

type Stage struct {
	Stage      string      `json:"stage"`
	Operations []Operation `json:"operations"`
	Measures   []Measure   `json:"measures"`
}

type Process struct {
	Process  string    `json:"process"`
	Stages   []Stage   `json:"stages"`
	Measures []Measure `json:"measures"`
}

// ===== Xpath Table =====
// Table connects measures to raw data and the hierarchy
// MeasureID connections to the hierarchy
// Xpath connects to raw data

type Xpath struct {
	Xpath     string `json:"xpath"`
	MeasureID string `json:"measureID"`
	Site      string `json:"site"`
}

// ===== Metadata Table =====
// Table supports adding metadata to measures in the hierarchy
// Simple Key:Value pair store
// MeasureID connections to the hierarchy

type Metadata struct {
	MeasureID string `json:"measureID"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

// ===== Results Table =====
// Stores raw data that users need to access
// Xpath connects raw data element to xpath table

type Results struct {
	Xpath       string  `json:"xpath"`
	MaterialNum string  `json:"materialNum"`
	BatchID     string  `json:"batchID"`
	Site        string  `json:"site"`
	DOM         string  `json:"DOM"` // Date of Manufacture
	ResultName  string  `json:"resultName"`
	Result      float64 `json:"result"`
}

// ===== Raw Materials Table =====
// Stores relationships between batches produced and raw materials used in the batch
// Xpath connects raw data element to xpath table

type RawMaterials struct {
	ParentBatchID     string `json:"parentBatchID"`
	ParentMaterialNum string `json:"parentMaterialNum"`
	ChildMaterialName string `json:"childMaterialName"`
	ChildBatchID      string `json:"childBatchID"`
	ChildMaterialNum  string `json:"childMaterialNum"`
}

// ===== Model Output =====
// Combined table that merges all the model component into a single output

type ModelOutput struct {
	Hierarchy    Process        `json:"hierarchy"`
	Xpath        []Xpath        `json:"xpath"`
	Metadata     []Metadata     `json:"metadata"`
	Results      []Results      `json:"results"`
	RawMaterials []RawMaterials `json:"rawMaterials"`
}
