/*
test program
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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
	//var inqs = []inqinfoshare{}
	var inqs = inqinfoshare{}
	bytestowrite, er := json.Marshal(inqs)
	if er != nil {
		return nil, errors.New("error occured marshalling")
	}
	err := stub.PutState(args[0]+"-shareinfo", bytestowrite)
	///err := stub.PutState(args[0]+"-shareinfo", []byte("INIT"))
	///if err != nil {
	///	return nil, errors.New("error occured:" + err.Error())
	///}

	//inq request list
	//var myinqs = []inquiry{}
	var myinqs = inquiry{}
	bytestowrite, er = json.Marshal(myinqs)
	if er != nil {
		return nil, errors.New("error occured marshalling my inquiries")
	}
	err = stub.PutState(args[0]+"-consinq", []byte("INIT"))
	///if err != nil {
	///	return nil, errors.New("error occured:" + err.Error())
	///}

	//state on my inqs done list on self-share parallel
	//var inqsharedone = []inqinfosharedone{}
	var inqsharedone = inqinfosharedone{}
	bytestowrite, er = json.Marshal(inqsharedone)
	if er != nil {
		return nil, errors.New("error occured marshalling my inquiries")
	}
	//err = stub.PutState(args[0]+"-inqdone", bytestowrite)
	err = stub.PutState(args[0]+"-inqdone", []byte("INIT"))
	if err != nil {
		return nil, errors.New("error occured:" + err.Error())
	}

	//what inq req is complete and what was gotten
	///var myinqdone = []inquirydone{}
	var myinqdone = inquirydone{}
	bytestowrite, er = json.Marshal(myinqdone)
	if er != nil {
		return nil, errors.New("error occured marshalling my inquiries")
	}
	err = stub.PutState(args[0]+"-consinqnotify", bytestowrite)
	///err = stub.PutState(args[0]+"-consinqnotify", []byte("INIT"))
	///if err != nil {
	///return nil, errors.New("error occured:" + err.Error())
	///}

	return nil, nil
}

func (t *ShareInfoCode) InitD(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	stub.PutState(args[0]+"-shareinfo", []byte(""))
	stub.PutState(args[0]+"-consinq", []byte(""))
	stub.PutState(args[0]+"-inqdone", []byte(""))
	stub.PutState(args[0]+"-consinqnotify", []byte(""))

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *ShareInfoCode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "initd" { //initialize the chaincode state, used as reset
		return t.InitD(stub, "initD", args)
	} else if function == "share-str" {
		return t.write(stub, args) //string concat method
	} else if function == "share-a" {
		return t.writeA(stub, args) //write same thing back
	} else if function == "share-b" {
		return t.writeB(stub, args) //write json array thru reallocation of struct array
	} else if function == "share-c" {
		return t.writeC(stub, args) //write static json array from struct array
	} else if function == "share-d" {
		return t.writeD(stub, args) //write static json array from struct array
	} else if function == "inquire" {
		return t.inquire(stub, args)
	} else if function == "shareone" {
		return t.writesingle(stub, args)
	} else if function == "inquireone" {
		return t.inquireone(stub, args)
	} else if function == "inquired" {
		return t.inquireD(stub, args)
	} else if function == "inquirydone" {
		return t.inquireone(stub, args)
	} else if function == "clearqueue" {
		return t.clearQueue(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error

	return nil, errors.New("Received unknown function invocation")
}

//invoke init user
func (t *ShareInfoCode) inituser(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return nil, nil
}

//write with string contact method
func (t *ShareInfoCode) write(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	//args0=username
	//arg1=stateJSON [{"withentity":"porsche35-userid","mydata":["score3","dec1003"],"sharedate":"2016-06-02"}]

	var err error
	fmt.Println("running write() function")

	var user = args[0]
	var sharewith = args[1]
	var mydata []string = strings.Split(args[2], "|")
	var sharedon = args[3]
	var res1 = inqinfoshare{sharewith, mydata, sharedon}
	bytestostore, _ := json.Marshal(res1)

	bytesofdata, _ := stub.GetState(user + "-shareinfo")

	var storedata = string(bytesofdata) + "^" + string(bytestostore)
	_ = stub.PutState(user+"-shareinfo", []byte(storedata))

	entbytesofdata, entErr := stub.GetState(sharewith + "-consreq")

	if entErr != nil {
		return nil, errors.New(entErr.Error() + "-error reading state user-consreq")
	}
	var inqQStored = string(entbytesofdata)
	var newinqitem = inquiry{EntityCode: sharewith, About: args[0]}
	bytesofinq, _ := json.Marshal(newinqitem)
	var newinqQstring = string(bytesofinq)

	var combinedinq string
	if len(inqQStored) != 0 {
		combinedinq = inqQStored + "^" + newinqQstring
	} else {
		combinedinq = newinqQstring
	}

	err = stub.PutState(res1.Withentity+"-consreq", []byte(combinedinq))
	if err != nil {
		return nil, errors.New(err.Error() + "-error writing to conseq state")
	}

	return nil, nil

}

//try to write samething back to state
func (t *ShareInfoCode) writeA(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	//args0=username
	//arg1=stateJSON [{"withentity":"porsche35-userid","mydata":["score3","dec1003"],"sharedate":"2016-06-02"}]

	var user = args[0]
	var sharewith = args[1]

	//var storedata = string(bytesofdata) + "^" + string(bytestostore)
	storeddata, _ := stub.GetState(user + "-shareinfo")
	_ = stub.PutState(user+"-shareinfo", storeddata)

	bytesofinq, _ := stub.GetState(sharewith + "-consreq")
	_ = stub.PutState(sharewith+"-consreq", bytesofinq)

	return nil, nil

}

//try to write struct as json array
func (t *ShareInfoCode) writeB(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	//args0=username
	//arg1=stateJSON [{"withentity":"porsche35-userid","mydata":["score3","dec1003"],"sharedate":"2016-06-02"}]

	var err error
	fmt.Println("running write() function")

	var user = args[0]
	var sharewith = args[1]
	var mydata []string = strings.Split(args[2], "|")
	var sharedon = args[3]
	var res1 = inqinfoshare{sharewith, mydata, sharedon}

	//get the current state data first
	bytesofdata, _ := stub.GetState(args[0] + "-shareinfo")
	var res2 = []inqinfoshare{}
	//unmarshal into struct array from json
	err = json.Unmarshal(bytesofdata, &res2)
	if err != nil {
		return nil, errors.New(err.Error() + "unable to unmarshall state data")
	}

	res2new := make([]inqinfoshare, len(res2)+1)
	if len(res2) > 0 {
		copy(res2new, res2[:len(res2)])
	}

	//add the new shareinfo data from user to the list
	res2new[len(res2)] = res1

	//unmarshal into new json to store in ledger
	bytestosave, _ := json.Marshal(res2new)

	//save to ledger state
	err = stub.PutState(user+"-shareinfo", bytestosave)

	///if err != nil {
	///	return nil, errors.New(er.Error() + "error storing state into shareinfo")
	///}

	//do the same in entity ledger store where entity will read for inq requests to them
	//indicator of which business is going well? future indicator of D&A analytics/analysis
	//currently copy paste from above, refactor later
	//get the current state data first

	entbytesofdata, entErr := stub.GetState(sharewith + "-consreq")

	if entErr != nil {
		return nil, errors.New(entErr.Error() + "-error reading state user-consreq")
	}
	var newinqitem = inquiry{EntityCode: sharewith, About: args[0]}

	var entres2 = []inquiry{}
	//unmarshal into struct array from json
	err = json.Unmarshal(entbytesofdata, &entres2)
	if err != nil {
		return nil, errors.New(err.Error() + "-unable to unmarshall state data from entity state")
	}

	entres2new := make([]inquiry, len(entres2)+1)
	//slice and take all data
	//should be entity result2 (entres2)
	if len(entres2) > 0 {
		copy(entres2new, entres2[:len(entres2)])
	}

	//add the new shareinfo data from user to the list
	entres2new[len(entres2new)] = newinqitem

	//unmarshal into new json to store in ledger
	bytestosave, _ = json.Marshal(entres2new)
	_ = stub.PutState(sharewith+"-consreq", bytestosave)

	//save to ledger state of entity as request queue

	if err != nil {
		return nil, errors.New(err.Error() + "-error writing to conseq state")
	}

	return nil, nil

}

//write statically struct array
func (t *ShareInfoCode) writeC(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	//args0=username
	//arg1=stateJSON [{"withentity":"porsche35-userid","mydata":["score3","dec1003"],"sharedate":"2016-06-02"}]

	var user = args[0]
	var sharewith = args[1]
	var mydata []string = strings.Split(args[2], "|")
	var sharedon = args[3]
	var r1 = inqinfoshare{sharewith, mydata, sharedon}
	var r2 = inqinfoshare{sharewith, mydata, sharedon}
	var res1 [2]inqinfoshare
	res1[0] = r1
	res1[1] = r2
	bytestostore11, _ := json.Marshal(res1)
	_ = stub.PutState(user+"-shareinfo", bytestostore11)

	//var storedata = string(bytesofdata) + "^" + string(bytestostore)
	var inqarr [2]inquiry
	inqarr[0] = inquiry{sharewith, args[0]}
	inqarr[1] = inquiry{sharewith, args[0] + "--1"}

	bytestostore22, _ := json.Marshal(inqarr)
	_ = stub.PutState(sharewith+"-consreq", bytestostore22)

	return nil, nil

}
func (t *ShareInfoCode) writeD(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	var user = args[0]
	var sharewith = args[1]
	var mydata string = args[2]
	var sharedon = args[3]
	var datastring string = sharewith + ",[" + mydata + "]," + sharedon
	storeddatastring, _ := stub.GetState(user + "-shareinfo")
	if len(string(storeddatastring)) > 0 {
		datastring = string(storeddatastring) + "^" + datastring
	}

	_ = stub.PutState(user+"-shareinfo", []byte(datastring))

	var entityqueue = sharewith + "," + user
	storedqueue, _ := stub.GetState(sharewith + "-consinq")

	if len(string(storedqueue)) > 0 {
		var newinqQ = string(storedqueue) + "^" + entityqueue
		_ = stub.PutState(sharewith+"-consinq", []byte(newinqQ))
	} else {
		_ = stub.PutState(sharewith+"-consinq", []byte(entityqueue))
	}

	return nil, nil
}

func (t *ShareInfoCode) writesingle(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	var user = args[0]
	var sharewith = args[1]
	var mydata []string = strings.Split(args[2], "|")
	var sharedon = args[3]
	var res1 = inqinfoshare{sharewith, mydata, sharedon}
	bytestostore, _ := json.Marshal(res1)
	_ = stub.PutState(user+"-shareinfo", bytestostore)

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

	if len(args) != 1 {
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

	/*
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
	*/
	var inquiredOnChainkey = args[0] + "-inqdone"
	stub.PutState(inquiredOnChainkey, inqsbytes)
	_ = stub.PutState(args[0]+"-consinq", []byte(""))
	//key for the consumer to know somebody hit u

	return inqsbytes, nil

}

func (t *ShareInfoCode) inquireone(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	//args[0] = user
	//args[1]=inquiredon
	var inquired = inqinfosharedone{args[0], time.Now().String()}
	inqDoneBytes, _ := json.Marshal(inquired)
	stub.PutState(args[1]+"-consinqnotify", inqDoneBytes)

	return inqDoneBytes, nil
}
func (t *ShareInfoCode) inquireD(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	//args[0] = user
	//args[1]=inquiredon
	var inquired string = args[0] + "," + time.Now().String()

	_ = stub.PutState(args[1]+"-consinqnotify", []byte(inquired))

	return nil, nil
}

func (t *ShareInfoCode) clearQueue(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	//args[0] = user

	//use it as state transition, not a delete for now...
	_ = stub.PutState(args[0]+"-consinq", []byte(""))

	return nil, nil
}
