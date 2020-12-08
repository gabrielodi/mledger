package support

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
)

func (s *SmartContract) checaExistencia(APIstub shim.ChaincodeStubInterface, jsonvar string, parametro string) sc.Response {

	// Verifica que apenas dois argumentos foram passados
	if len(args) != 2 {
		return shim.Error("Numero incorreto de argumentos buscaLedger")
	}

	queryIntString := "{\"selector\":{\"" + jsonvar + "\":\"" + parametro + "\"}}"
	queryString := fmt.Sprintf(queryIntString)
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)

	// Analisa o resultado e retorna se busca sucedida.
	fmt.Println(BytesToString(queryResults))
	if err != nil {
		return false
	}
	return true
}

func (s *SmartContract) checaEntidade(APIstub shim.ChaincodeStubInterface, parametro string, variavel string, tipoVariavel string) sc.Response {

	// Verifica que apenas dois argumentos foram passados
	if len(args) != 2 {
		return shim.Error("Numero incorreto de argumentos buscaLedger")
	}

	queryIntString := "{\"selector\":{\"eCnpjCpf\":\"" + parametro + "\"}}"
	queryString := fmt.Sprintf(queryIntString)
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)

	var ent = Entidade{}
	json.Unmarshal(queryResults, &ent)

	if tipoVariavel == "tipo" {
		if variavel != ent.Tipo {
			return false
		}
	}
	if tipoVariavel == "senha" {
		if variavel != ent.Senha {
			return false
		}
	}
	return true
}

func (s *SmartContract) getLatestKey(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := ""
	endKey := ""

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	i := 0
	for resultsIterator.HasNext() {
		_, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		i = i + 1
	}
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(i))

	fmt.Printf("--\n%s\n", buffer.String())
	//	fmt.Printf("NUMBER OF HITS:\n%s\n", strconv.Itoa(i))

	return shim.Success(buffer.Bytes())

}

// =========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

// Query per string
func (s *SmartContract) queryPerString(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	vary, _ := APIstub.GetState(args[0])
	varx := BytesToString(vary)
	stry, _ := APIstub.GetState(args[1])
	strx := BytesToString(stry)
	queryString := "{\"selector\":{\"" + varx + "\":\"" + strx + "\"}}"
	//queryString := "{\"selector\":{\"Model\":\"Prius\"}}"
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	fmt.Println(queryString)
	fmt.Println(BytesToString(queryResults))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)

}

//Byte to string
func BytesToString(data []byte) string {
	return string(data[:])
}

// ===========================================================================================
// constructQueryResponseFromIterator constructs a JSON array containing query results from
// a given result iterator
// ===========================================================================================
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return &buffer, nil
}
