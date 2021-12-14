package structures

import (
	"fmt"
	"testing"

	"github.com/USACE/go-consequences/utils"
)

// for testing with Github action
const path = "./data/DF.json"

// for testing locally
// const path = "/workspaces/go-consequences/data/DF.json"

func Test_readJson(t *testing.T) {
	var c RawDFStruct
	err := utils.ReadJson(path, &c)

	if err != nil {
		t.Errorf("Unable to parse Json from file")
	}
}

func Test_IngestDFStore(t *testing.T) {
	_, err := ingestDDFStore(path)
	if err != nil {
		t.Errorf("Unable to parse Json from file")
	}
}

func Test_InitDDFStore(t *testing.T) {
	var p DepthDFProvider
	p.Init(path)
}

func Test_GenerateDF(t *testing.T) {
	var p DepthDFProvider
	p.Init(path)

	df, _ := p.DamageFunction("COM1", "structure")

	fmt.Println(df)
}
