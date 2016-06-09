/*
test program
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// ShareInfoCode example simple Chaincode implementation

//this struct will be used by consumer  allowing inquiry
type inqinfoshare struct {
	Withentity string   `json:"withentity"`
	Mydata     []string `json:"mydata"`
	ShareDate  string   `json:"sharedate"`
}

//this struct will be used by consumer entity to read from for inquiry requests
type inquiry struct {
	EntityCode string `json:"entitycode"`
	About      string `json:"about"`
}

//this struct will be used to write to consumer that inquiry happened and state change
type inqinfosharedone struct {
	Withentity string `json:"withentity"`
	ShareDate  string `json:"sharedate"`
}

//this struct will be used by consumer entity after inquiry is complete and state change done
type inquirydone struct {
	About        string `json:"about"`
	DataObtained string `json:"dataobtained"`
}

type ShareInfoCode struct {
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(ShareInfoCode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *ShareInfoCode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	//arg0 = username
	if len(args) != 1 {
		return nil, errors.New("required 1 argument, the chain user")
	}
	//inq self share
	var inqs = []inqinfoshare{}
	bytestowrite, er := json.Marshal(inqs)
	if er != nil {
		return nil, errors.New("error occured marshalling")
	}
	err := stub.PutState(args[0]+"-shareinfo", bytestowrite)
	if err != nil {
		return nil, errors.New("error occured:" + err.Error())
	}

	//inq request list
	var myinqs = []inquiry{}
	bytestowrite, er = json.Marshal(myinqs)
	if er != nil {
		return nil, errors.New("error occured marshalling my inquiries")
	}
	err = stub.PutState(args[0]+"-consinq", bytestowrite)
	if err != nil {
		return nil, errors.New("error occured:" + err.Error())
	}

	//state on my inqs done list on self-share parallel
	var inqsharedone = []inqinfosharedone{}
	bytestowrite, er = json.Marshal(inqsharedone)
	if er != nil {
		return nil, errors.New("error occured marshalling my inquiries")
	}
	err = stub.PutState(args[0]+"-inqdone", bytestowrite)
	if err != nil {
		return nil, errors.New("error occured:" + err.Error())
	}

	//what inq req is complete and what was gotten
	var myinqdone = []inquirydone{}
	bytestowrite, er = json.Marshal(myinqdone)
	if er != nil {
		return nil, errors.New("error occured marshalling my inquiries")
	}
	err = stub.PutState(args[0]+"-consinqnotify", bytestowrite)
	if err != nil {
		return nil, errors.New("error occured:" + err.Error())
	}

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *ShareInfoCode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "share" {
		return t.write(stub, args)
	} else if function == "inquire" {
		return t.inquire(stub, args)
	}
	fmt.Println("invoke did not find func: " + function) //error

	return nil, errors.New("Received unknown function invocation")
}

//invoke init user
func (t *ShareInfoCode) inituser(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return nil, nil
}

func (t *ShareInfoCode) write(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	//args0=username
	//arg1=stateJSON [{"withentity":"porsche35-userid","mydata":["score3","dec1003"],"sharedate":"2016-06-02"}]

	var user string
	var err error
	fmt.Println("running write() function")

	if len(args) != 1 {
		return nil, errors.New("incorrect number of args, expecting 1. name of variable and value to set")

	}

	user = args[0]
	var res1 inqinfoshare
	err = json.Unmarshal([]byte(args[1]), &res1)
	if err != nil {
		return nil, err
	}
	fmt.Println(res1.Mydata)

	//get the current state data first
	bytesofdata, er := stub.GetState(user + "-shareinfo")

	if er != nil {
		return nil, errors.New("error reading state user-shareinfo")
	}

	var res2 []inqinfoshare
	//unmarshal into struct array from json
	err = json.Unmarshal(bytesofdata, &res2)
	if err != nil {
		return nil, errors.New("unable to unmarshall state data")
	}

	res2new := make([]inqinfoshare, len(res2)+1)
	if len(res2) > 0 {
		copy(res2new, res2[:len(res2)])
	}

	//add the new shareinfo data from user to the list
	res2new[len(res2)] = res1

	//unmarshal into new json to store in ledger
	bytestosave, er := json.Marshal(res2new)

	//save to ledger state
	err = stub.PutState(user+"-shareinfo", bytestosave)
	if err != nil {
		return nil, err
	}

	//do the same in entity ledger store where entity will read for inq requests to them
	//indicator of which business is going well? future indicator of D&A analytics/analysis
	//currently copy paste from above, refactor later
	//get the current state data first

	bytesofdata, er = stub.GetState(res1.Withentity + "-consreq")

	if er != nil {
		return nil, errors.New("error reading state user-consreq")
	}

	var entres2 []inquiry
	//unmarshal into struct array from json
	err = json.Unmarshal(bytesofdata, &entres2)
	if err != nil {
		return nil, errors.New("unable to unmarshall state data from entity state")
	}

	entres2new := make([]inquiry, len(entres2)+1)
	//slice and take all data
	if len(res2) > 0 {
		copy(entres2new, entres2[:len(entres2)])
	}

	//add the new shareinfo data from user to the list
	entres2new[len(entres2new)] = inquiry{res1.Withentity, args[0]}

	//unmarshal into new json to store in ledger
	bytestosave, er = json.Marshal(entres2new)

	//save to ledger state of entity as request queue
	err = stub.PutState(res1.Withentity+"-consreq", bytestosave)
	if err != nil {
		return nil, err
	}

	return nil, nil

}

//remove data in chaincode then post with new state record
func (t *ShareInfoCode) remove(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	//args0 = user

	var err error
	fmt.Println("running write() function")

	if len(args) == 0 {
		return nil, errors.New("incorrect number of args, expecting 1. name of variable and value to set")
	}

	err = stub.DelState(args[0] + "-shareinfo")
	if err != nil {
		return nil, err
	}

	return nil, nil

}

// Query is our entry point for queries
func (t *ShareInfoCode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "dummy_query" { //read a variable
		fmt.Println("hi there " + function) //error
		return nil, nil
	} else if function == "read" { //read a variable
		return t.read(stub, args)
	}

	fmt.Println("query did not find func: " + function) //error

	return nil, errors.New("Received unknown function query")
}

func (t *ShareInfoCode) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	var name, jsonResp, chaintype string
	var err error

	if len(args) != 2 {
		return nil, errors.New("incorrect number of args, expecting name of var to query")

	}

	name = args[0]
	chaintype = args[1]
	valAsbytes, err := stub.GetState(name + "-" + chaintype)

	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

func (t *ShareInfoCode) inquire(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	//args[0] - current user who in inquiring

	var name, jsonResp string
	var err error

	if len(args) != 2 {
		return nil, errors.New("incorrect number of args, expecting name of var to query")

	}

	name = args[0]

	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return nil, errors.New(jsonResp)
	}

	inqsbytes, er := stub.GetState(args[0] + "-consinq")

	if er != nil {
		return nil, errors.New("error getting consinq state")
	}
	var inqs []inquiry
	_ = json.Unmarshal(inqsbytes, &inqs)
	//key for the entity to know inq is done as part of inq request
	if len(inqs) == 0 {
		return nil, nil
	}

	var inquiredOnChainkey = args[0] + "-inqdone"
	var lenofinqs = len(inqs)

	//var inqsdone [lenofinqs]inqdone
	inqsdone := make([]inquirydone, 1, lenofinqs)

	for a := 0; a < len(inqs); a++ {
		var inqdone = inquirydone{inqs[a].About, time.Now().String()}
		inqsdone[a] = inqdone

		var inquiredConsumerChainkey = inqs[a].About + "-consinqnotify"
		var donemyinq = inqinfosharedone{args[0], time.Now().String()}
		bytes2write, _ := json.Marshal(donemyinq)
		stub.PutState(inquiredConsumerChainkey, bytes2write)
	}
	bytestopost, _ := json.Marshal(inqsdone)

	stub.PutState(inquiredOnChainkey, bytestopost)
	//key for the consumer to know somebody hit u

	return inqsbytes, nil

}
