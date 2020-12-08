package main

/* Imports
 * Bibliotecas necessarias para manipulacao de bites, arquivos JSON, strings e formatacao
 * Adicionalmente, bibliotecas para smart contracts
 */
import (
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
)

// Estrutura SmartContract
type SmartContract struct{}

// Definição de classes
// Entidade (Proprietario, Concessionaria, Montadora)
type Entidade struct {
	Nome                          string `json:"nome"`
	ECnpjCpf                      int    `json:"eCnpjCpf"`
	Senha                         string `json:"senha"`
	Tipo                          string `json:"tipo"`
	EDocumentoResponsavelRegistro int    `json:"eDocumentoResponsavelRegistro"`
}

// Veiculo
type Veiculo struct {
	VChassis                      string `json:"vChassis"`
	Renavam                       int    `json:"renavam"`
	VCnpjCpf                      int    `json:"vCnpjCpf"`
	Especificacao                 string `json:"especificacao"`
	DataFabricacao                string `json:"dataFabricacao"`
	DataUtilizacao                string `json:"dataUtilizacao"`
	DataRevenda                   string `json:"dataRevenda"`
	VDocumentoResponsavelRegistro int    `json:"vDocumentoResponsavelRegistro"`
}

// Registros de Manutencao
type RegistroManutencao struct {
	MChassis                      string  `json:"mChassis"`
	PnFalha                       int     `json:"pnFalha"`
	TipoFalha                     string  `json:"tipoFalha"`
	SintomaFalha                  string  `json:"sintomaFalha"`
	InvestigacaoFalha             string  `json:"investigacaoFalha"`
	CodigoFalha                   int     `json:"codigoFalha"`
	ManutencaoRealizada           string  `json:"manutencaoRealizada"`
	PnTrocado                     int     `json:"pnTrocado"`
	NSerieTrocado                 int     `json:"nSerieTrocado"`
	NotaFiscal                    string  `json:"notaFiscal"`
	CustoManutencao               float64 `json:"custoManutencao"`
	DataReparo                    string  `json:"dataReparo"`
	KmAtual                       int     `json:"kmAtual"`
	ConsumoCombustivel            float64 `json:"consumoCombustivel"`
	CargaCarregada                int     `json:"cargaCarregada"`
	MDocumentoResponsavelRegistro int     `json:"mDocumentoResponsavelRegistro"`
	ResponsavelRegistro           string  `json:"responsavelRegistro"`
}

func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}

/*
 * Funcao Init é chamado pelo fabric na inicializacao da blockchain.
 * Uma melhor pratica documentada é de manter essa funcao em branco,
 * uma vez que roda sempre que o chaincode é iniciado, e manter uma funcao separada (iniciarLedger())
 * com um código que pode ser rodado quando apropriado.
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * A funcao Invoke é o recurso chamado para quaisquer requisicoes para o chaincode.
 * Cada função deve ser definida aqui, para que possam ser chamadas como parametros da funcao invoke.
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Trazer funcoes e argumentos da funcao GetFunctionAndParameters das bibliotecas de chaincode
	function, args := APIstub.GetFunctionAndParameters()
	// Registrar as funcoes que podem ser chamadas por parametro dentro da funcao Invoke
	if function == "buscaVeiculoPorEntidade" { //FUNCIONANDO
		return s.buscaVeiculoPorEntidade(APIstub, args)
	} else if function == "buscaVeiculoPorChassis" { //FUNCIONANDO
		return s.buscaVeiculoPorChassis(APIstub, args)
	} else if function == "buscaManutencoes" { //FUNCIONANDO
		return s.buscaManutencoes(APIstub, args)
	} else if function == "buscaEntidade" { //FUNCIONANDO
		return s.buscaEntidade(APIstub, args)
	} else if function == "iniciarLedger" { //FUNCIONANDO
		return s.iniciarLedger(APIstub)
	} else if function == "novaEntidade" { //FEITO
		return s.novaEntidade(APIstub, args)
	} else if function == "novoVeiculo" { //FEITO
		return s.novoVeiculo(APIstub, args)
	} else if function == "novaManutencao" { //FEITO, FALTA PUXAR O NOME DO RESPONSAVEL PELO REGISTRO
		return s.novaManutencao(APIstub, args)
	} else if function == "mudancaProprietario" {
		return s.mudancaProprietario(APIstub, args)
		//	} else if function == "mudancaSenha" {
		//		return s.mudancaSenha(APIstub, args)
	} else if function == "getLatestKey" { //NAO PRECISA ESTAR AQUI, CHAMADA SEPARADAMENTE
		return s.getLatestKey(APIstub) //NAO PRECISA ESTAR AQUI, CHAMADA SEPARADAMENTE
	} else if function == "queryPerString" { //DELETAR DEPOIS
		return s.queryPerString(APIstub, args) // DELETAR DEṔOIS
	}
	return shim.Error(`Funcao Invalida.\nFuncoes Disponiveis:
	\nbuscaVeiculoPorEntidade(Documento)
	\nbuscaManutencoes(Chassis)
	\nbuscaEntidade(Documento)
	\niniciarLedger()
	\nnovaEntidade(nome,documento,senha,tipo)
	\nnovoVeiculo(chassis,renavam,documento_proprietario,especificacao,data_fabricacao,data_inicio_utilizacao,data_venda)
	\nnovaManutencao(pn_falha,tipo_falha,sintoma_falha,investigacao_falha,codigo_falha,manutencao_realizada,pn_Trocado,n_Serie_PN_Trocado,nota_fiscal,custo_manutencao,data_reparo,km_atual,consumo_combsutivel,carga_carregada,responsavel_registro)
	\nmudancaProprietario(chassis,documento_comprador,senha_comprador,senha_vendedor,documento_concessionaria,senha_concessionaria)
	\nmudancaSenha(documento,senha_antiga,senha_nova)
	`)
}
