package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

const (
	ccID      = "mledger"
	channelID = "mychannel"
	orgName   = "org1.example.com"
	orgAdmin  = "Admin"
)

func useGateway() {
	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile("./connection.json")),
		gateway.WithUser("Admin"),
	)
	if err != nil {
		fmt.Printf("Failed to connect: %v", err)
	}

	if gw == nil {
		fmt.Println("Failed to create gateway")
	}

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		fmt.Printf("Failed to get network: %v", err)
	}

	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	contract := network.GetContract("mledger")
	uuid.SetRand(nil)

	var wg sync.WaitGroup
	start := time.Now()
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			seededRand.Intn(20)
			result, err := contract.SubmitTransaction("invoke", "put", uuid.New().String(),
				strconv.Itoa(seededRand.Intn(20)))
			if err != nil {
				fmt.Printf("Failed to commit transaction: %v", err)
			} else {
				fmt.Println("Commit is successful")
			}

			fmt.Println(reflect.TypeOf(result))
			fmt.Printf("The results is %v", result)
		}()
	}
	wg.Wait()
	fmt.Println("The time took is ", time.Now().Sub(start))
}

func serve() {
	http.HandleFunc("/run", handlerRun) //deletar
	http.HandleFunc("/novoUsuario", novoUsuario)
	http.HandleFunc("/buscaUsuario", buscaUsuario)
	http.HandleFunc("/mudancaSenha", mudancaSenha)

	http.HandleFunc("/cadastraVeiculo", cadastraVeiculo)
	http.HandleFunc("/buscarVChassis", buscarVChassis)
	http.HandleFunc("/buscarVProprietario", buscarVProprietario)
	http.HandleFunc("/mudancaProprietario", mudancaProprietario)

	http.HandleFunc("/buscarHManutencao", buscarHManutencao)
	http.HandleFunc("/novaManutencao", novaManutencao)
	fmt.Println("Listening at /localhost:3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		fmt.Println(err)
	}
}

type Entidade struct {
	Name string
	Doc  string
}



//USUARIOS

func novoUsuario(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)

	result := invoke("novaEntidade",
		r.FormValue("nome"), r.FormValue("eCnpjCpf"), r.FormValue("senha"),
		r.FormValue("tipo"), r.FormValue("eDocumentoResponsavelRegistro"),
		r.FormValue("senhaResp"))

	fmt.Println("RESULTADO:")
	fmt.Println(string(result))
	fmt.Println("FIM")
	w.Write(result)

}

func buscaUsuario(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)

	result := query("buscaEntidade", r.FormValue("eCnpjCpf"))
	w.Write(result)
}

func mudancaSenha(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	fmt.Println(r.Form)

	result := invoke("mudancaSenha", r.FormValue("eCnpjCpf"), r.FormValue("senhaAntiga"), r.FormValue("senhaNova"))
	w.Write(result)
}

//VEICULOS

func cadastraVeiculo(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	fmt.Println(r.Form)

	result := invoke("novoVeiculo",
		r.FormValue("vChassis"), r.FormValue("renavam"), r.FormValue("vCnpjCpf"),
		r.FormValue("especificacao"), r.FormValue("dataFabricacao"), r.FormValue("dataUtilizacao"),
		r.FormValue("dataRevenda"), r.FormValue("vDocumentoResponsavelRegistro"), r.FormValue("vSenhaResponsavelRegistro"))
	w.Write(result)
}

func buscarVChassis(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	fmt.Println(r.Form)

	result := query("buscaVeiculoPorChassis", r.FormValue("vChassis"))
	w.Write(result)
}

func buscarVProprietario(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	fmt.Println(r.Form)

	result := query("buscaVeiculoPorEntidade", r.FormValue("vCnpjCpf"))
	w.Write(result)
}

func mudancaProprietario(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	fmt.Println(r.Form)
	result := invoke("mudancaProprietario", r.FormValue("vChassis"), r.FormValue("vCnpjCpf"), r.FormValue("dataRevenda"),
		r.FormValue("senhaNovo"), r.FormValue("senhaAntigo"),
		r.FormValue("vDocumentoResponsavelRegistro"), r.FormValue("vSenhaResponsavelRegistro"))
	w.Write(result)
}

//MANUTENCAO

func buscarHManutencao(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	fmt.Println(r.Form)

	result := query("buscaManutencoes", r.FormValue("mChassis"))
	w.Write(result)
}

func novaManutencao(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)

	result := invoke("novaManutencao", r.FormValue("mChassis"), r.FormValue("pnFalha"),
		r.FormValue("tipoFalha"), r.FormValue("sintomaFalha"), r.FormValue("investigacaoFalha"),
		r.FormValue("codigoFalha"), r.FormValue("manutencaoRealizada"), r.FormValue("pnTrocado"),
		r.FormValue("nSerieTrocado"), r.FormValue("notaFiscal"), r.FormValue("custoManutencao"),
		r.FormValue("dataReparo"), r.FormValue("kmAtual"), r.FormValue("consumoCombustivel"),
		r.FormValue("cargaCarregada"), r.FormValue("mDocumentoResponsavelRegistro"), r.FormValue("senhaResponsavelRegistro"))

	w.Write(result)
}

//ROTINAS
func query(funcBlockchain string, parametros ...string) []byte {
	wallet, err := gateway.NewFileSystemWallet("./wallets")
	if err != nil {
		fmt.Printf("Failed to create wallet: %s\n", err)
		os.Exit(1)
	}

	if !wallet.Exists("Admin") {
		fmt.Println("Failed to get Admin from wallet")
		os.Exit(1)
	}

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile("./connection.json")),
		gateway.WithUser("Admin"),
	)
	if err != nil {
		fmt.Printf("Failed to connect: %v", err)
	}

	if gw == nil {
		fmt.Println("Failed to create gateway")
	}

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		fmt.Printf("Failed to get network: %v", err)
	}

	contract := network.GetContract("mledger")

	var wg sync.WaitGroup
	start := time.Now()
	wg.Add(1)
	var result []byte
	go func() {
		defer wg.Done()

		if len(parametros) == 1 {
			result, err = contract.EvaluateTransaction(funcBlockchain, parametros[0])
		}
		if len(parametros) == 2 {
			result, err = contract.EvaluateTransaction(funcBlockchain, parametros[0], parametros[1])
		}
		if len(parametros) == 3 {
			result, err = contract.EvaluateTransaction(funcBlockchain, parametros[0], parametros[1], parametros[2])
		}
		if len(parametros) == 4 {
			result, err = contract.EvaluateTransaction(funcBlockchain, parametros[0], parametros[1], parametros[2], parametros[3])
		}
		if len(parametros) == 5 {
			result, err = contract.EvaluateTransaction(funcBlockchain, parametros[0], parametros[1], parametros[2], parametros[3], parametros[4])
		}
		if len(parametros) == 6 {
			result, err = contract.EvaluateTransaction(funcBlockchain, parametros[0], parametros[1], parametros[2], parametros[3], parametros[4], parametros[5])
		}
		if len(parametros) == 7 {
			result, err = contract.EvaluateTransaction(funcBlockchain, parametros[0], parametros[1], parametros[2], parametros[3], parametros[4], parametros[5], parametros[6])
		}
		if len(parametros) == 9 {
			result, err = contract.EvaluateTransaction(funcBlockchain, parametros[0], parametros[1], parametros[2], parametros[3], parametros[4], parametros[5], parametros[6], parametros[7], parametros[8])
		}
		if len(parametros) == 17 {
			result, err = contract.EvaluateTransaction(funcBlockchain, parametros[0], parametros[1], parametros[2], parametros[3], parametros[4], parametros[5], parametros[6], parametros[7], parametros[8],
				parametros[9], parametros[10], parametros[11], parametros[12], parametros[13], parametros[14], parametros[15], parametros[16])
		}
		if err != nil {
			fmt.Printf("Failed to commit transaction: %v", err)
		} else {
			fmt.Println("Commit is successful")
		}

	}()
	wg.Wait()
	fmt.Println("The time took is ", time.Now().Sub(start))
	result = []byte(strings.Trim(string(result), "[]"))
	fmt.Println(string(result))
	return result
}

func invoke(funcBlockchain string, parametros ...string) []byte {
	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile("./connection.json")),
		gateway.WithUser("Admin"),
	)
	if err != nil {
		fmt.Printf("Failed to connect: %v", err)
	}

	if gw == nil {
		fmt.Println("Failed to create gateway")
	}

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		fmt.Printf("Failed to get network: %v", err)
	}

	contract := network.GetContract("mledger")

	var wg sync.WaitGroup
	start := time.Now()
	wg.Add(1)
	var result []byte
	go func() {
		defer wg.Done()

		if len(parametros) == 1 {
			result, err = contract.SubmitTransaction(funcBlockchain, parametros[0])
		}
		if len(parametros) == 2 {
			result, err = contract.SubmitTransaction(funcBlockchain, parametros[0], parametros[1])
		}
		if len(parametros) == 3 {
			result, err = contract.SubmitTransaction(funcBlockchain, parametros[0], parametros[1], parametros[2])
		}
		if len(parametros) == 4 {
			result, err = contract.SubmitTransaction(funcBlockchain, parametros[0], parametros[1], parametros[2], parametros[3])
		}
		if len(parametros) == 5 {
			result, err = contract.SubmitTransaction(funcBlockchain, parametros[0], parametros[1], parametros[2], parametros[3], parametros[4])
		}
		if len(parametros) == 6 {
			result, err = contract.SubmitTransaction(funcBlockchain, parametros[0], parametros[1], parametros[2], parametros[3], parametros[4], parametros[5])
		}
		if len(parametros) == 7 {
			result, err = contract.SubmitTransaction(funcBlockchain, parametros[0], parametros[1], parametros[2], parametros[3], parametros[4], parametros[5], parametros[6])
		}
		if len(parametros) == 9 {
			result, err = contract.SubmitTransaction(funcBlockchain, parametros[0], parametros[1], parametros[2], parametros[3], parametros[4], parametros[5], parametros[6], parametros[7], parametros[8])
		}
		if len(parametros) == 17 {
			result, err = contract.SubmitTransaction(funcBlockchain, parametros[0], parametros[1], parametros[2], parametros[3], parametros[4], parametros[5], parametros[6], parametros[7], parametros[8],
				parametros[9], parametros[10], parametros[11], parametros[12], parametros[13], parametros[14], parametros[15], parametros[16])
		}
		if err != nil {
			fmt.Printf("Failed to commit transaction: %v", err)
		} else {
			fmt.Println("Commit is successful")
		}


	}()

	wg.Wait()
	fmt.Println("The time took is ", time.Now().Sub(start))
	result = []byte(strings.Trim(string(result), "[]"))
	fmt.Println(string(result))
	return result
}

func main() {
	serve()
}
