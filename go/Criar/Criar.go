package Criar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	
	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
)

func (s *SmartContract) novaEntidade(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 6 {
		return shim.Error("Numero incorreto de argumentos (6 sao esperados):\nNome, Documento (CPF/CNPJ), Senha, Tipo (Montadora, Concessionaria, Proprietario), documento do responsavel pelo registro, senha do responsavel pelo registro.")
	}

	//Transformando documentos em numeros
	NECnpjCpf, _ := strconv.Atoi(args[1])
	NEDoc, _ := strconv.Atoi(args[4])

	//Verifica se o usuario responsavel pelo registro existe e se a senha confere
	teste := make([]string, 1)
	teste[0] = args[4]
	respEntidadeByte := s.buscaEntidade(APIstub, teste)
	var respEntidade = Entidade{}
	json.Unmarshal([]byte(string(respEntidadeByte.Payload)), &respEntidade)
	if strconv.Itoa(respEntidade.ECnpjCpf) != strconv.Itoa(NEDoc) {
		return shim.Error(string(respEntidadeByte.Payload) + strconv.Itoa(respEntidade.ECnpjCpf))
	}
	if respEntidade.Senha != args[5] {
		return shim.Error("Senha do responsavel invalida.")
	}

	//Verifica se a nova entidade ja nao existe
	teste[0] = args[1]
	novaEntidadeByte := s.buscaEntidade(APIstub, teste)
	var novaEntidade = Entidade{}
	json.Unmarshal(novaEntidadeByte.Payload, &novaEntidade)
	if novaEntidade.ECnpjCpf == NECnpjCpf {
		return shim.Error("Usuario já existe.")
	}

	//Populando a variavel entidade com os dados do registro
	var entidade = Entidade{Nome: args[0], ECnpjCpf: NECnpjCpf, Senha: args[2], Tipo: args[3], EDocumentoResponsavelRegistro: NEDoc}

	//Transformando a entidade em um registro JSON
	entidadeAsBytes, _ := json.Marshal(entidade)

	//Encontra o ultimo ID da blockchain para ser o ID+1 (nkey)
	resultsIterator, err := APIstub.GetStateByRange("", "")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	i := 0
	for resultsIterator.HasNext() {
		_, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		i = i + 1
	}

	nkey := strconv.Itoa(i + 1)

	//Envia o registro da entidade em JSON para a blockchain sob o ID nkey.
	APIstub.PutState(nkey, entidadeAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) novoVeiculo(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 8 {
		return shim.Error("Numero incorreto de argumentos (7 sao esperados):\nChassis, Renavam, Documento (CPF/CNPJ), Especificacao, data de fabricacao, data de inicio de utilizacao e data de revenda e documento do responsavel pelo registro")
	}

	nVRenavam, _ := strconv.Atoi(args[1])
	nVCnpjCpf, _ := strconv.Atoi(args[2])
	nVDoc, _ := strconv.Atoi(args[7])

	var veiculo = Veiculo{VChassis: args[0], Renavam: nVRenavam, VCnpjCpf: nVCnpjCpf, Especificacao: args[3], DataFabricacao: args[4], DataUtilizacao: args[5], DataRevenda: args[6], VDocumentoResponsavelRegistro: nVDoc}

	veiculoAsBytes, _ := json.Marshal(veiculo)

	resultsIterator, err := APIstub.GetStateByRange("", "")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	i := 0
	for resultsIterator.HasNext() {
		_, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		i = i + 1
	}

	nkey := strconv.Itoa(i + 1)
	APIstub.PutState(nkey, veiculoAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) novaManutencao(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 17 {
		return shim.Error("Numero incorreto de argumentos (17 sao esperados):\nChassis, Numero da Peça falhada, tipo de falha, sintoma de falha, investigacao da falha, codigo de fala, manutencao realizada, numero da peca trocada, numero de serie da peca trocada, nota fiscal, custo de manutencao, data de reparo, km atual, consumo de combustivel atual, carga carrega, documento do responsavel pelo registro, senha do responsavel pelo registro.")
	}

	//PRECISA BUSCAR O NOME DO RESPONSAVEL PELO REGISTRO
	nomeResponsavelRegistro := args[15] //BUSCAR NOME E CORRIGIR AQUI!!!!!!!!!
	//VERIFICAR A SENHA DO RESPONSAVEL PELO REGISTRO args[16]

	nMPnFalha, _ := strconv.Atoi(args[1])
	nMCodFalha, _ := strconv.Atoi(args[5])
	nMPnTrocado, _ := strconv.Atoi(args[7])
	nMNSerieTrocado, _ := strconv.Atoi(args[8])
	nMCustoManutencao, _ := strconv.ParseFloat(args[10], 64)
	nMKmAtual, _ := strconv.Atoi(args[12])
	nMConsumo, _ := strconv.ParseFloat(args[13], 64)
	nMCarga, _ := strconv.Atoi(args[14])
	nMDoc, _ := strconv.Atoi(args[15])

	var registroManutencao = RegistroManutencao{MChassis: args[0], PnFalha: nMPnFalha, TipoFalha: args[2], SintomaFalha: args[3], InvestigacaoFalha: args[4], CodigoFalha: nMCodFalha, ManutencaoRealizada: args[6], PnTrocado: nMPnTrocado, NSerieTrocado: nMNSerieTrocado, NotaFiscal: args[9], CustoManutencao: nMCustoManutencao, DataReparo: args[11], KmAtual: nMKmAtual, ConsumoCombustivel: nMConsumo, CargaCarregada: nMCarga, MDocumentoResponsavelRegistro: nMDoc, ResponsavelRegistro: nomeResponsavelRegistro}

	registroManutencaoAsBytes, _ := json.Marshal(registroManutencao)

	resultsIterator, err := APIstub.GetStateByRange("", "")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	i := 0
	for resultsIterator.HasNext() {
		_, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		i = i + 1
	}

	nkey := strconv.Itoa(i + 1)
	APIstub.PutState(nkey, registroManutencaoAsBytes)

	return shim.Success(nil)
}

//dataRevenda"22/06/2019", vDocumentoResponsavelRegistro: "9999999999"
//0 - chassis
//1 - novo proprietario
//2 - data transacao
//3 - senha novo proprietario
//4 - senha velho proprietario
//5 - dealer
//6 - senha dealer
func (s *SmartContract) mudancaProprietario(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	chassis := args[0]
	novoProprietario, _ := strconv.Atoi(args[1])
	dataTransacao := args[2]
	senhaNovoProprietario := args[3]
	senhaVelhoProprietario := args[4]
	intermediador, _ := strconv.Atoi(args[5])
	senhaIntermediador := args[6]
	
	//Verifica se chassis, novo proprietario e intermediador existem.
	exChassis := checaExistencia(APIstub,"vChassis",chassis)
	exNovoProprietario := checaExistencia(APIstub,"eCnpjCpf",args[1])
	exIntermediador := checaExistencia(APIstub,"eCnpjCpf",args[5])
	if (exChassis == false ) || (exNovoProprietario == false ) || (exIntermediador == false ){
		return shim.Error(err.Error())
	}
	
	//Verifica se o intermediador é concessionaria ou montadora
	tipoIntermediadorCons := checaEntidade(APIstub,args[5],"concessionaria","tipo")
	tipoIntermediadorMont := checaEntidade(APIstub,args[5],"montadora","tipo")
	if (tipoIntermediadorCons == false ) && (tipoIntermediadorMont == false ){
		return shim.Error(err.Error())
	}
	
	//Verifica se senhas sao corretas
	checkSenhaNovo := checaEntidade(APIstub,args[5],senhaNovoProprietario,"senha")
	checkSenhaVelho := checaEntidade(APIstub,args[5],senhaVelhoProprietario,"senha")
	checkSenhaIntermediario := checaEntidade(APIstub,args[5],senhaIntermediador,"senha")
	if (checkSenhaNovo == false ) || (checkSenhaVelho == false ) || (checkSenhaIntermediario == false ){
		return shim.Error(err.Error())	
	}

	// Busca pelo veiculo existente
	queryIntString := "{\"selector\":{\"vChassis\":\"" + chassis + "\"}}"
	queryString := fmt.Sprintf(queryIntString)
	queryVeiculo, _ := getQueryResultForQueryString(APIstub, queryString)
	var veiculoChain = Veiculo{}
	json.Unmarshal(queryVeiculo, &veiculoChain)

	var veiculo = Veiculo{VChassis: veiculoChain.VChassis, Renavam: veiculoChain.Renavam, VCnpjCpf: novoProprietario, Especificacao: veiculoChain.Especificacao, DataFabricacao: veiculoChain.DataFabricacao, DataUtilizacao: veiculoChain.DataUtilizacao, DataRevenda: dataTransacao, VDocumentoResponsavelRegistro: intermediador}

	veiculoAsBytes, _ := json.Marshal(veiculo)
	//PRECISO PEGAR A nkey EXISTENTE PARA ATUALZIAR NA LEDGER
	APIstub.PutState("1111111111", veiculoAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) iniciarLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	entidades := []Entidade{
		Entidade{Nome: "Montadora de Caminhoes LTDA", ECnpjCpf: 9999999999, Senha: "69bdd58e17ab101986d8cf7a7f9279db", Tipo: "montadora", EDocumentoResponsavelRegistro: 9999999999}, //senha: montadora, hashed MD5
		Entidade{Nome: "Concessionaria A", ECnpjCpf: 8888888888, Senha: "9aa4898a07dc811fddd4b9d8655fce8f", Tipo: "concessionaria", EDocumentoResponsavelRegistro: 9999999999},       //senha: concessionariaa, hashed MD5
		Entidade{Nome: "Caio", ECnpjCpf: 7777777777, Senha: "af1dbc5648a563e9a5bd97d0eb68f41b", Tipo: "proprietario", EDocumentoResponsavelRegistro: 8888888888},                     //senha: papas, hashed MD5
		Entidade{Nome: "Patricia", ECnpjCpf: 6666666666, Senha: "86241b5767f022a036a93a9b55af2e71", Tipo: "proprietario", EDocumentoResponsavelRegistro: 8888888888},                 //senha: branco, hashed MD5
		Entidade{Nome: "Charlie", ECnpjCpf: 5555555555, Senha: "21591c1cb4eacdf98eab4454f9dbbd09", Tipo: "proprietario", EDocumentoResponsavelRegistro: 8888888888},                  //senha: sarayu, hashed MD5
	}

	veiculos := []Veiculo{
		Veiculo{VChassis: "E111111", Renavam: 38461734, VCnpjCpf: 9999999999, Especificacao: "Pesado 13L 6x4 500cv", DataFabricacao: "01/11/2019", DataUtilizacao: "", DataRevenda: "", VDocumentoResponsavelRegistro: 9999999999},
		Veiculo{VChassis: "E222222", Renavam: 49682922, VCnpjCpf: 9999999999, Especificacao: "Medio 7L 4x2 250cv", DataFabricacao: "05/11/2019", DataUtilizacao: "", DataRevenda: "", VDocumentoResponsavelRegistro: 9999999999},
		Veiculo{VChassis: "E333333", Renavam: 59692833, VCnpjCpf: 8888888888, Especificacao: "Medio 7L 6x4 300cv", DataFabricacao: "10/09/2019", DataUtilizacao: "", DataRevenda: "", VDocumentoResponsavelRegistro: 9999999999},
		Veiculo{VChassis: "E444444", Renavam: 66049299, VCnpjCpf: 8888888888, Especificacao: "Pesado 13L 6x2 400cv", DataFabricacao: "15/09/2019", DataUtilizacao: "", DataRevenda: "", VDocumentoResponsavelRegistro: 9999999999},
		Veiculo{VChassis: "E555555", Renavam: 94827493, VCnpjCpf: 7777777777, Especificacao: "Pesado 13L 6x4 450cv", DataFabricacao: "13/01/2019", DataUtilizacao: "22/06/2019", DataRevenda: "22/06/2019", VDocumentoResponsavelRegistro: 9999999999},
		Veiculo{VChassis: "E666666", Renavam: 20595753, VCnpjCpf: 7777777777, Especificacao: "Pesado 13L 6x4 500cv", DataFabricacao: "05/02/2019", DataUtilizacao: "11/05/2019", DataRevenda: "11/05/2019", VDocumentoResponsavelRegistro: 9999999999},
		Veiculo{VChassis: "E777777", Renavam: 28102409, VCnpjCpf: 7777777777, Especificacao: "Pesado 13L 6x2 400cv", DataFabricacao: "23/02/2019", DataUtilizacao: "14/06/2019", DataRevenda: "14/06/2019", VDocumentoResponsavelRegistro: 9999999999},
		Veiculo{VChassis: "E888888", Renavam: 21052708, VCnpjCpf: 6666666666, Especificacao: "Medio 7L 8x2 350cv", DataFabricacao: "01/03/2019", DataUtilizacao: "10/04/2019", DataRevenda: "10/04/2019", VDocumentoResponsavelRegistro: 9999999999},
		Veiculo{VChassis: "E999999", Renavam: 21052711, VCnpjCpf: 6666666666, Especificacao: "Medio 7L 8x2 350cv", DataFabricacao: "01/03/2019", DataUtilizacao: "10/04/2019", DataRevenda: "10/04/2019", VDocumentoResponsavelRegistro: 9999999999},
	}

	registroManutencao := []RegistroManutencao{
		//RegistroManutencao{PnFalha:"",TipoFalha:"",SintomaFalha:"",InvestigacaoFalha:"",CodigoFalha:"",ManutencaoRealizada:"",PnTrocado:"",NSerieTrocado:"",NotaFiscal:"",CustoManutencao:"",DataReparo:"",KmAtual:"",ConsumoCombustivel:"",CargaCarregada:"",ResponsavelRegistro :""},
		RegistroManutencao{MChassis: "E555555", PnFalha: 7769, TipoFalha: "Quebra", SintomaFalha: "Perda de potencia.", InvestigacaoFalha: "Codigos de falha escaneados", CodigoFalha: 238, ManutencaoRealizada: "Troca de peça quebrada.", PnTrocado: 5306, NSerieTrocado: 936861, NotaFiscal: "410520888888888888885500110000000011193203701", CustoManutencao: 2696, DataReparo: "01/03/2020", KmAtual: 35155, ConsumoCombustivel: 2.34, CargaCarregada: 54, MDocumentoResponsavelRegistro: 8888888888, ResponsavelRegistro: "Concessionaria A"},
		RegistroManutencao{MChassis: "E555555", PnFalha: 5727, TipoFalha: "Quebra", SintomaFalha: "Fumaça branca.", InvestigacaoFalha: "Vazamento de oleo identificado.", CodigoFalha: 222, ManutencaoRealizada: "Troca de peça quebrada.", PnTrocado: 8097, NSerieTrocado: 348430, NotaFiscal: "410520888888888888885500110000000021640017072", CustoManutencao: 2947, DataReparo: "11/03/2020", KmAtual: 35988, ConsumoCombustivel: 2.22, CargaCarregada: 54, MDocumentoResponsavelRegistro: 8888888888, ResponsavelRegistro: "Concessionaria A"},
		RegistroManutencao{MChassis: "E555555", PnFalha: 7148, TipoFalha: "Software", SintomaFalha: "Freios não funcionando.", InvestigacaoFalha: "Verificacao de codigos de falha.", CodigoFalha: 607, ManutencaoRealizada: "Download de novo software na unidade de controle.", PnTrocado: 0, NSerieTrocado: 0, NotaFiscal: "411020888888888888885500110000000031180774293", CustoManutencao: 121, DataReparo: "01/08/2020", KmAtual: 81230, ConsumoCombustivel: 2.61, CargaCarregada: 54, MDocumentoResponsavelRegistro: 8888888888, ResponsavelRegistro: "Concessionaria A"},
		RegistroManutencao{MChassis: "E666666", PnFalha: 1290, TipoFalha: "Quebra", SintomaFalha: "Perda de potência de frenagem.", InvestigacaoFalha: "Vazamento no ar comprimido.", CodigoFalha: 456, ManutencaoRealizada: "Troca de peça quebrada.", PnTrocado: 3485, NSerieTrocado: 577830, NotaFiscal: "410720888888888888885500110000000041754768734", CustoManutencao: 1949, DataReparo: "12/04/2020", KmAtual: 44033, ConsumoCombustivel: 2.8, CargaCarregada: 74, MDocumentoResponsavelRegistro: 8888888888, ResponsavelRegistro: "Concessionaria A"},
		RegistroManutencao{MChassis: "E666666", PnFalha: 1952, TipoFalha: "Quebra", SintomaFalha: "Ruido na suspensão", InvestigacaoFalha: "Desgaste na suspensao.", CodigoFalha: 515, ManutencaoRealizada: "Troca de peça quebrada.", PnTrocado: 4060, NSerieTrocado: 133920, NotaFiscal: "411020888888888888885500110000000051219054995", CustoManutencao: 4000, DataReparo: "01/08/2020", KmAtual: 56041, ConsumoCombustivel: 2.49, CargaCarregada: 74, MDocumentoResponsavelRegistro: 8888888888, ResponsavelRegistro: "Concessionaria A"},
		RegistroManutencao{MChassis: "E777777", PnFalha: 5885, TipoFalha: "Software", SintomaFalha: "Marcha não engatando.", InvestigacaoFalha: "Problema identificado no software da caixa.", CodigoFalha: 863, ManutencaoRealizada: "Download de novo software na unidade de controle.", PnTrocado: 0, NSerieTrocado: 0, NotaFiscal: "410320888888888888885500110000000061697193086", CustoManutencao: 145, DataReparo: "14/09/2020", KmAtual: 23231, ConsumoCombustivel: 2.87, CargaCarregada: 74, MDocumentoResponsavelRegistro: 8888888888, ResponsavelRegistro: "Concessionaria A"},
		RegistroManutencao{MChassis: "E777777", PnFalha: 1245, TipoFalha: "Quebra", SintomaFalha: "Ruido no motor.", InvestigacaoFalha: "Sujeira e desgaste nas correias", CodigoFalha: 795, ManutencaoRealizada: "Troca de peça quebrada.", PnTrocado: 5328, NSerieTrocado: 607783, NotaFiscal: "410720888888888888885500110000000071489081347", CustoManutencao: 2788, DataReparo: "22/11/2020", KmAtual: 44430, ConsumoCombustivel: 2.44, CargaCarregada: 74, MDocumentoResponsavelRegistro: 8888888888, ResponsavelRegistro: "Concessionaria A"},
		RegistroManutencao{MChassis: "E777777", PnFalha: 2079, TipoFalha: "Quebra", SintomaFalha: "Nivel de óleo baixo.", InvestigacaoFalha: "Consumo de oleo pela junta.", CodigoFalha: 985, ManutencaoRealizada: "Troca de peça quebrada.", PnTrocado: 6462, NSerieTrocado: 842840, NotaFiscal: "410820888888888888885500110000000081402701728", CustoManutencao: 3515, DataReparo: "14/12/2020", KmAtual: 48955, ConsumoCombustivel: 2.57, CargaCarregada: 74, MDocumentoResponsavelRegistro: 8888888888, ResponsavelRegistro: "Concessionaria A"},
		RegistroManutencao{MChassis: "E777777", PnFalha: 4512, TipoFalha: "Software", SintomaFalha: "Veículo não liga.", InvestigacaoFalha: "Falha no software do motor de arranque.", CodigoFalha: 648, ManutencaoRealizada: "Download de novo software na unidade de controle.", PnTrocado: 0, NSerieTrocado: 0, NotaFiscal: "411220888888888888885500110000000091746789739", CustoManutencao: 98, DataReparo: "01/02/2021", KmAtual: 75099, ConsumoCombustivel: 2.11, CargaCarregada: 74, MDocumentoResponsavelRegistro: 8888888888, ResponsavelRegistro: "Concessionaria A"},
		RegistroManutencao{MChassis: "E888888", PnFalha: 5841, TipoFalha: "Quebra", SintomaFalha: "Nível de água baixo.", InvestigacaoFalha: "Falha na valvula do tanque de expansao.", CodigoFalha: 627, ManutencaoRealizada: "Troca de peça quebrada.", PnTrocado: 1522, NSerieTrocado: 101289, NotaFiscal: "410220888888888888885500110000000101489270291", CustoManutencao: 2968, DataReparo: "01/04/2020", KmAtual: 26987, ConsumoCombustivel: 2.63, CargaCarregada: 22, MDocumentoResponsavelRegistro: 8888888888, ResponsavelRegistro: "Concessionaria A"},
	}

	i := 0
	for i < len(entidades) {
		fmt.Println("i is ", i)
		entidadesBytes, _ := json.Marshal(entidades[i])
		APIstub.PutState(strconv.Itoa(i), entidadesBytes)
		fmt.Println("Added", entidades[i])
		i = i + 1
	}
	j := 0
	for j < len(veiculos) {
		fmt.Println("j is ", j)
		veiculosBytes, _ := json.Marshal(veiculos[j])
		APIstub.PutState(strconv.Itoa(i+j-1), veiculosBytes)
		fmt.Println("Added", veiculos[j])
		j = j + 1
	}
	k := 0
	for k < len(registroManutencao) {
		fmt.Println("k is ", k)
		registroManutencaoBytes, _ := json.Marshal(registroManutencao[k])
		APIstub.PutState(strconv.Itoa(i+j+k-2), registroManutencaoBytes)
		fmt.Println("Added", registroManutencao[k])
		k = k + 1
	}

	return shim.Success(nil)
}