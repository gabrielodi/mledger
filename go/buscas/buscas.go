package buscas

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	
	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
)

// Busca veiculo cadastrado sob certa entidade
func (s *SmartContract) buscaVeiculoPorEntidade(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	// Verifica que apenas um argumento foi passado
	if len(args) != 1 {
		return shim.Error("Numero incorreto de argumentos. Inclua apenas um documento (CPF ou CNPJ, sem pontos)")
	}

	// Recebe o argumento como byte[], transforma em string e envia a query pelo documento do proprietario
	//argByte, _ := APIstub.GetState(args[0])
	//argString := BytesToString(argByte)
	queryIntString := "{\"selector\":{\"vCnpjCpf\":" + args[0] + "}}"
	queryString := fmt.Sprintf(queryIntString)
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)

	// Analisa o resultado e retorna se busca sucedida.
	fmt.Println(BytesToString(queryResults))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte(queryResults))
}

// Busca veiculo por chassis
func (s *SmartContract) buscaVeiculoPorChassis(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	// Verifica que apenas um argumento foi passado
	if len(args) != 1 {
		return shim.Error("Numero incorreto de argumentos. Inclua apenas um chassis.")
	}

	// Recebe o argumento como byte[], transforma em string e envia a query pelo documento do proprietario
	//argByte, _ := APIstub.GetState(args[0])
	//argString := BytesToString(argByte)
	queryIntString := "{\"selector\":{\"vChassis\":\"" + args[0] + "\"}}"
	queryString := fmt.Sprintf(queryIntString)
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)

	// Analisa o resultado e retorna se busca sucedida.
	fmt.Println(BytesToString(queryResults))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// Busca Manutencoes
func (s *SmartContract) buscaManutencoes(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	// Verifica que apenas um argumento foi passado
	if len(args) != 1 {
		return shim.Error("Numero incorreto de argumentos. Inclua apenas um chassis")
	}

	// Recebe o argumento como byte[], transforma em string e envia a query pelo documento do proprietario
	//argByte, _ := APIstub.GetState(args[0])
	//argString := BytesToString(argByte)
	queryIntString := "{\"selector\":{\"mChassis\":\"" + args[0] + "\"}}"
	queryString := fmt.Sprintf(queryIntString)
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)

	// Analisa o resultado e retorna se busca sucedida.
	fmt.Println(BytesToString(queryResults))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// Busca entidade
func (s *SmartContract) buscaEntidade(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	// Verifica que apenas um argumento foi passado
	if len(args) != 1 {
		return shim.Error("Numero incorreto de argumentos. Inclua apenas um documento (CPF ou CNPJ, sem pontos)")
	}

	// Recebe o argumento como byte[], transforma em string e envia a query pelo documento do proprietario
	//argByte, _ := APIstub.GetState(args[0])
	//argString := BytesToString(argByte)
	queryIntString := "{\"selector\":{\"eCnpjCpf\":" + args[0] + "}}"
	queryString := fmt.Sprintf(queryIntString)
	queryResults, err := getQueryResultForQueryString(APIstub, queryString)

	// Analisa o resultado e retorna se busca sucedida.
	fmt.Println(BytesToString(queryResults))
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

