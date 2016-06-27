// vote
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type candidates struct {
	Candidate string `json:"candidateid"`
	Votes     int    `json:"votes"`
}

type proposition struct {
	Proposal string `json:"proposal"`
	Yesvote  int    `json:"yesvote"`
	Novote   int    `json:"novote"`
}

type candidatevote struct {
}

type votingcode struct {
}

func main() {
	err := shim.Start(new(votingcode))
	if err != nil {
		fmt.Printf("Error starting voting chaincode: %s", err)
	}
}

func (t *votingcode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	if args[0] == "G" {

		_ = stub.PutState("votetype", []byte("G")) //general

	} else {
		_ = stub.PutState("votetype", []byte("P")) //proposal
	}

	err := stub.PutState("votingstarted", []byte(time.Now().String()))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *votingcode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "vote" { //initialize the chaincode state, used as reset
		return t.vote(stub, "init", args)
	} else if function == "registercandidate" {
		return t.registercandidates(stub, function, args)
	} else if function == "registerproposal" {
		return t.registerproposals(stub, function, args)
	}

	fmt.Println("invoke did not find func: " + function) //error

	return nil, errors.New("Received unknown function invocation")
}

// Query is our entry point for queries
func (t *votingcode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "dummy_query" { //read a variable
		fmt.Println("hi there " + function) //error
		return nil, nil
	} else if function == "voteresult" {
		return t.getresults(stub, args)
	} else if function == "read" {
		r, _ := stub.GetState(strings.Replace(args[0], " ", "_", -1))
		return r, nil
	}

	fmt.Println("query did not find func: " + function) //error

	return nil, errors.New("Received unknown function query")
}

func (t *votingcode) vote(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	//args0 = userid
	//arg1 = candidate

	//call efx to chain to confirm user is a voter

	//if voter is valid, register the vote

	vtype, _ := stub.GetState("votetype")

	if string(vtype) == "G" {

		votesbytes, err := stub.GetState("candidates")
		if err != nil {
			return nil, errors.New("error reading votes from state")
		}
		var votesstruct []candidates
		_ = json.Unmarshal(votesbytes, &votesstruct)
		for a := 0; a < len(votesstruct); a++ {
			if votesstruct[a].Candidate == args[1] {
				votesstruct[a].Votes = votesstruct[a].Votes + 1
				storedVal, _ := stub.GetState(strings.Replace(args[1], " ", "_", -1))
				if len(string(storedVal)) > 0 {
					val, _ := strconv.Atoi(string(storedVal))
					val++
					stub.PutState(strings.Replace(args[1], " ", "_", -1), []byte(strconv.Itoa(val)))
				}
				break
			}
		}

		barray, _ := json.Marshal(votesstruct)

		stub.PutState("votes", barray)
	} else {
		votesbytes, err := stub.GetState("proposals")
		if err != nil {
			return nil, errors.New("error reading votes from state")
		}
		var votesstruct []proposition
		_ = json.Unmarshal(votesbytes, &votesstruct)
		for a := 0; a < len(votesstruct); a++ {
			if votesstruct[a].Proposal == args[1] {
				if args[2] == "Y" {
					votesstruct[a].Yesvote += 1
				} else {
					votesstruct[a].Novote += 1
				}

				storedVal, _ := stub.GetState(args[1] + "-" + args[2])
				if len(string(storedVal)) > 0 {
					val, _ := strconv.Atoi(string(storedVal))
					val++
					stub.PutState(args[1]+"-"+args[2], []byte(strconv.Itoa(val)))
				}
				break
			}
		}

		barray, _ := json.Marshal(votesstruct)

		stub.PutState("votes", barray)
	}
	return nil, nil
}

func (t *votingcode) registercandidates(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	//args[0] = user
	//args[1]=candidateuserid

	var l = len(args)

	var candidatestruct []candidates
	//star from 0+1th element
	for a := 1; a < l; a++ {
		candidatestruct = resetCandidateArray(candidatestruct, candidates{Candidate: args[a], Votes: 0})
	}

	jsonarr, _ := json.Marshal(candidatestruct)
	_ = stub.PutState("candidates", jsonarr)

	return nil, nil

}

func resetCandidateArray(inarray []candidates, newelement candidates) []candidates {
	var newarray []candidates
	if inarray != nil {
		newarray = make([]candidates, len(inarray)+1)
		copy(newarray, inarray[:len(inarray)])
		newarray[len(inarray)] = newelement
	} else {
		newarray = make([]candidates, 1)
		newarray[len(inarray)] = newelement
	}
	return newarray
}

func (t *votingcode) registerproposals(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	//args[0] = user
	//args[1]=candidateuserid

	var l = len(args)

	var props []proposition
	//star from 0+1th element
	for a := 1; a < l; a++ {
		props = resetPropArray(props, proposition{Proposal: args[a], Yesvote: 0, Novote: 0})
	}

	jsonarr, _ := json.Marshal(props)
	_ = stub.PutState("proposals", jsonarr)

	return nil, nil

}

func resetPropArray(inarray []proposition, newelement proposition) []proposition {
	var newarray []proposition
	if inarray != nil {
		newarray = make([]proposition, len(inarray)+1)
		copy(newarray, inarray[:len(inarray)])
		newarray[len(inarray)] = newelement
	} else {
		newarray = make([]proposition, 1)
		newarray[len(inarray)] = newelement
	}
	return newarray
}

func (t *votingcode) getresults(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	barr, _ := stub.GetState("votetype")

	if string(barr) == "G" {
		votesbytes, err := stub.GetState("candidates")
		var votesstruct []candidates
		_ = json.Unmarshal(votesbytes, &votesstruct)

		for a := 0; a < len(votesstruct); a++ {
			storedVal, _ := stub.GetState(votesstruct[a].Candidate)
			votesstruct[a].Votes, _ = strconv.Atoi(string(storedVal))
			break
		}
		votesbytes, _ = json.Marshal(votesstruct)

		return votesbytes, err
	} else if string(barr) == "P" {
		votesbytes, err := stub.GetState("proposals")
		var votesstruct []proposition
		_ = json.Unmarshal(votesbytes, &votesstruct)

		for a := 0; a < len(votesstruct); a++ {
			yStoredVal, _ := stub.GetState(votesstruct[a].Proposal + "-Y")
			nStoredVal, _ := stub.GetState(votesstruct[a].Proposal + "-N")
			votesstruct[a].Yesvote, _ = strconv.Atoi(string(yStoredVal))
			votesstruct[a].Novote, _ = strconv.Atoi(string(nStoredVal))
			break
		}

		votesbytes, _ = json.Marshal(votesstruct)
		return votesbytes, err
	}
	return nil, errors.New("unable to resolve between candidate voting and proposal")
}
