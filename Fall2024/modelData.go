/*
 * This package includes code inspired by the work of:
 *
 * Author: Brett Duersch
 * Title: Merck Process Hierarchy Modeling
 * Date: February 16, 2024 
 *
 * Probabilites modified to concise data
 */
package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
	"unicode/utf8"
)

type Batch struct {
	BatchID     string `json:"batchID"`
	Site        string `json:"site"`
	MaterialNum string `json:"materialNum"`
	DOM         string `json:"DOM"`
}

// Reverse string function from:
// https://stackoverflow.com/questions/1752414/how-to-reverse-a-string-in-go
func StringReverse(s string) string {
	size := len(s)
	buf := make([]byte, size)
	for start := 0; start < size; {
		r, n := utf8.DecodeRuneInString(s[start:])
		start += n
		utf8.EncodeRune(buf[size-start:], r)
	}
	return string(buf)
}

func defineP() map[string]map[string]float64 {

	// setup constants that direct formation of hierarchy
	var p map[string]map[string]float64 = make(map[string]map[string]float64)

	// "minMeasure" = minimum # of measures for this level
	// "measure" = probability of additional measures will be added
	// "minLevel" = minimum number of levels that will be added to that layer
	// "level" = probability that additional levels will be added to that layer

	p["prc"] = make(map[string]float64)
	p["prc"]["minMeasure"] = 0
	p["prc"]["measure"] = 0
	p["prc"]["minLevel"] = 3
	p["prc"]["level"] = 0.7

	p["stg"] = make(map[string]float64)
	p["stg"]["minMeasure"] = 0
	p["stg"]["measure"] = 0.25
	p["stg"]["minLevel"] = 1
	p["stg"]["level"] = 0.6

	p["op"] = make(map[string]float64)
	p["op"]["minMeasure"] = 0
	p["op"]["measure"] = 0.25
	p["op"]["minLevel"] = 1
	p["op"]["level"] = 0.6

	p["act"] = make(map[string]float64)
	p["act"]["minMeasure"] = 1
	p["act"]["measure"] = 0.7
	p["act"]["minLevel"] = 0
	p["act"]["level"] = 0

	// "minXpath" = minimum number of xpaths per measure
	// "xpath" = probability of adding extra xpaths
	// "minMeta" = minimum number of metaData elements
	// "meta" = probability of adding extra metadata elements
	p["mes"] = make(map[string]float64)
	p["mes"]["minXpath"] = 1
	p["mes"]["xpath"] = 0.3
	p["mes"]["minMeta"] = 2
	p["mes"]["meta"] = 0.5

	// "minBatch" = minimum number of batches
	// "batch" = probability of adding extra batches
	// "startDOM" = starting DOM
	// "dateStep" = base time between batches
	// "mean" = target mean of batch results
	// "stdDev" = standard devation for batch results
	// "offset" = max offset from mean for batch results
	p["res"] = make(map[string]float64)
	p["res"]["minBatch"] = 20
	p["res"]["batch"] = 0.95
	layout := "01/02/2006 3:04:05 PM"
	t, _ := time.Parse(layout, "01/01/2022 9:03:46 AM")
	p["res"]["startDOM"] = float64(t.Unix())
	p["res"]["dateStep"] = 15
	p["res"]["mean"] = 100
	p["res"]["stdDev"] = 2.5
	p["res"]["offset"] = 2.5

	// "minRawMat" = minimum number of raw materials
	p["rm"] = make(map[string]float64)
	p["rm"]["minRawMat"] = 5
	return p
}

// define the measures for a given level
// baseName - string used as the basis for new measureNames
// levelType - string used to access appropriate output in probabilty map
// pMap - probability map to define likelihood of a measure/xpath/metadata creation
func MakeMeasure(
	baseName, //
	levelType string,
	pMap map[string]map[string]float64,
	numSites int,
	siteBatches map[string][]Batch,
) (
	mesOut []Measure,
	xPaths []Xpath,
	metas []Metadata,
	results []Results,
) {

	// keep track of levels added
	var numMes float64 = 0
	var numXp float64 = 0
	var numMeta float64 = 0

	// initiate loop variables
	mes := true
	xp := true
	meta := true

	for mes {

		numMes++
		// evaluate creation of new measures
		if numMes > pMap[levelType]["minMeasure"] && rand.Float64() > pMap[levelType]["measure"] {
			break
		}

		// assign values
		measID := fmt.Sprintf("%v-M%v", baseName, numMes)
		measName := fmt.Sprintf("Measure: %v", measID)
		thisMes := Measure{}
		thisMes.Measure = measName
		thisMes.MeasureID = measID

		mesOut = append(mesOut, thisMes)

		numXp = 0
		numMeta = 0

		// independently calculate xpaths on a per-site basis
		for i := 0; i < numSites; i++ {

			siteXpaths := []Xpath{}
			siteText := fmt.Sprintf("Site-%v", i+1)
			// create the xpath components
			for xp {
				numXp++
				if numXp > pMap["mes"]["minXpath"] && rand.Float64() > pMap["mes"]["xpath"] {
					break
				}

				tempXpath := Xpath{
					Xpath:     fmt.Sprintf("//%v//site%v//x%v", measID, i+1, numXp),
					MeasureID: measID,
					Site:      siteText,
				}

				siteXpaths = append(siteXpaths, tempXpath)
			}

			offset := pMap["res"]["offset"] * rand.NormFloat64()
			stdDev := pMap["res"]["stdDev"] * rand.Float64()
			mean := pMap["res"]["mean"]

			// setup splits between different xpaths
			xPaths = append(xPaths, siteXpaths...)
			numBatches := len(siteBatches[siteText])
			bucketSize := math.Ceil(float64(numBatches) / float64(numXp-1))

			// create results linked to each xpath
			for ndx, row := range siteBatches[siteText] {
				thisResult := Results{}

				// get the index of the xpath result to use
				xpIndex := int(math.Floor(float64(ndx) / float64(bucketSize)))

				// now generate result table entries
				thisResult.Xpath = siteXpaths[xpIndex].Xpath
				thisResult.Site = siteXpaths[xpIndex].Site
				thisResult.BatchID = row.BatchID
				thisResult.DOM = row.DOM
				thisResult.ResultName = measID
				thisResult.Result = rand.NormFloat64()*stdDev + mean + offset

				results = append(results, thisResult)
			}

			numXp = 0
		}

		// create the metadata components on a per-measure basis
		for meta {
			numMeta++
			if numMeta > pMap["mes"]["minMeta"] && rand.Float64() > pMap["mes"]["meta"] {
				break
			}

			tempMeta := Metadata{
				MeasureID: measID,
				Key:       fmt.Sprintf("Key%v", numMeta),
				Value:     fmt.Sprintf("Value%v", numMeta),
			}
			metas = append(metas, tempMeta)
		}

	}
	return
}

func MakeBatches(numSites int, p map[string]map[string]float64) (
	map[string][]Batch, []RawMaterials) {

	siteBatch := make(map[string][]Batch)
	rawMat := []RawMaterials{}

	// independently calculate xpaths on a per-site basis
	for i := 0; i < numSites; i++ {

		theseBatches := []Batch{}
		siteText := fmt.Sprintf("Site-%v", i+1)

		numBatch := float64(0)
		// generate random matnr for each site
		matnr := StringReverse(fmt.Sprintf("%v", time.Now().Unix()))
		DOM := p["res"]["startDOM"]

		for true {

			numBatch++
			// evaluate continuation of process
			if numBatch > p["res"]["minBatch"] && rand.Float64() > p["res"]["batch"] {
				break
			}

			// increment DOM according to dateStep
			DOM = DOM + 86400*(rand.Float64())*p["res"]["dateStep"]

			thisBatch := Batch{}
			thisBatch.BatchID = fmt.Sprintf("%v", math.Ceil(DOM/4321))
			thisBatch.Site = siteText
			thisBatch.MaterialNum = matnr
			thisBatch.DOM = fmt.Sprintf("%v", time.Unix(int64(DOM), 0).Format("2006-01-02"))

			theseBatches = append(theseBatches, thisBatch)

			// define raw materials for the batch
			for j := 0; j < int(p["rm"]["minRawMat"]); j++ {
				thisRM := RawMaterials{}
				thisRM.ParentBatchID = thisBatch.BatchID
				thisRM.ParentMaterialNum = thisBatch.MaterialNum
				thisRM.ChildMaterialName = fmt.Sprintf("RawMat-%v", j)
				thisRM.ChildBatchID = fmt.Sprintf("%v", math.Ceil(rand.Float64()*300000))
				thisRM.ChildMaterialNum = fmt.Sprintf("%v", (j+1)*6021023) // matnr in moles -- remember chem 101 ;)

				rawMat = append(rawMat, thisRM)
			}
		}
		siteBatch[siteText] = theseBatches
	}

	return siteBatch, rawMat
}

func ModelData(name string) (ModelOutput, error) {

	outPut := ModelOutput{}

	// storage buffers
	xPaths := []Xpath{}
	metaData := []Metadata{}
	results := []Results{}

	// counters
	var numStg float64 = 0
	var numOp float64 = 0
	var numAct float64 = 0

	// for loop flags
	prc := true
	stg := true
	op := true

	p := defineP()

	hier := Process{}
	hier.Process = name

	// determine number of sites - will be used by measures / batches
	numSites := int(math.Abs(math.Ceil(rand.NormFloat64()))) + 2

	// generate the batch list and raw materails table
	siteBatches, rawMat := MakeBatches(numSites, p)

	for prc {

		numStg++
		// evaluate continuation of process
		if numStg > p["prc"]["minLevel"] && rand.Float64() > p["prc"]["level"] {
			break
		}

		stageID := fmt.Sprintf("%v-%v", name, numStg)
		stgName := fmt.Sprintf("Stage: %v", stageID)
		thisStage := Stage{}
		thisStage.Stage = stgName

		numOp = 0
		for stg {

			numOp++
			// evaluate continuation of process
			if numOp > p["stg"]["minLevel"] && rand.Float64() > p["stg"]["level"] {
				break
			}

			opID := fmt.Sprintf("%v-%v", stageID, numOp)
			opName := fmt.Sprintf("Operation: %v", opID)
			thisOp := Operation{}
			thisOp.Operation = opName

			numAct = 0
			for op {

				numAct++
				// evaluate continuation of process
				if numAct > p["op"]["minLevel"] && rand.Float64() > p["op"]["level"] {
					break
				}

				actID := fmt.Sprintf("%v-%v", opID, numAct)
				actName := fmt.Sprintf("Action: %v", actID)
				thisAct := Action{}
				thisAct.Action = actName

				mes, x, m, r := MakeMeasure(actID, "act", p, numSites, siteBatches)
				thisAct.Measures = append(thisAct.Measures, mes...)
				xPaths = append(xPaths, x...)
				metaData = append(metaData, m...)
				results = append(results, r...)

				thisOp.Actions = append(thisOp.Actions, thisAct)

			}

			mes, x, m, r := MakeMeasure(opID, "op", p, numSites, siteBatches)
			thisOp.Measures = append(thisOp.Measures, mes...)
			xPaths = append(xPaths, x...)
			metaData = append(metaData, m...)
			results = append(results, r...)

			thisStage.Operations = append(thisStage.Operations, thisOp)

		}
		mes, x, m, r := MakeMeasure(stageID, "stg", p, numSites, siteBatches)
		thisStage.Measures = append(thisStage.Measures, mes...)
		xPaths = append(xPaths, x...)
		metaData = append(metaData, m...)
		results = append(results, r...)

		hier.Stages = append(hier.Stages, thisStage)
	}

	// assign output variables
	outPut.Hierarchy = hier
	outPut.Xpath = xPaths
	outPut.Metadata = metaData
	outPut.Results = results
	outPut.RawMaterials = rawMat

	return outPut, nil

}
