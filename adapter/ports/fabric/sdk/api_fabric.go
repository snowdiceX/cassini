package sdk

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	ethsdk "github.com/QOSGroup/cassini/adapter/ports/ethereum/sdk"
	"github.com/QOSGroup/cassini/log"
	"github.com/pkg/errors"
)

// ChaincodeInvoke invoke chaincode
func ChaincodeInvoke(chaincodeID string, argsArray []Args) (
	result string, err error) {
	log.Info("chaincode invoke...")
	if chaincodeID == "" {
		err = fmt.Errorf("must specify the chaincode ID")
		return
	}
	action, err := newChaincodeInvokeAction()
	if err != nil {
		log.Errorf("Error while initializing invokeAction: %v", err)
		return
	}

	defer action.Terminate()

	result, err = action.invoke(Config().ChannelID, chaincodeID, argsArray)
	if err != nil {
		log.Errorf("Error while calling action.invoke(): %v", err)
	}
	return
}

// ChaincodeInvokeByString call chaincode invoke of hyperledger fabric
func ChaincodeInvokeByString(chaincodeID, argsStr string) string {
	log.Infof("chaincode invoke: %s; %s; %s", chaincodeID, argsStr)
	arg, err := ParseArgs(argsStr)
	if strings.EqualFold(arg.Func, "create") {
		RegisterWalletByString(arg.Args[0], arg.Args[1], arg.Args[2])
	}
	if err == nil {
		var argsArray []Args
		argsArray = append(argsArray, *arg)
		var ret string
		ret, err = ChaincodeInvoke(chaincodeID, argsArray)
		if err == nil {
			log.Info("chaincode invoke result: ", ret)
			return ret
		}
	}
	log.Errorf("%s %v", defaultResultJSON, err)
	return defaultResultJSON
}

// ChaincodeQuery call chaincode query of hyperledger fabric
func ChaincodeQuery(chaincodeID string, argsArray []Args) (result string, err error) {
	log.Info("chaincode query...")
	if chaincodeID == "" {
		err = fmt.Errorf("must specify the chaincode ID")
		return
	}
	var action *queryAction
	action, err = newQueryAction()
	if err != nil {
		log.Errorf("Error while initializing queryAction: %v", err)
		return
	}

	defer action.Terminate()

	result, err = action.query(Config().ChannelID, chaincodeID, argsArray)
	if err != nil {
		log.Errorf("Error while running queryAction: %v", err)
	} else if result == "" {
		err = errors.New("transaction not found")
	}
	if strings.HasPrefix(result, "v") {
		if i := strings.Index(result, ":"); i > -1 {
			result = result[i+1:]
		}
	}
	return
}

// ChaincodeQueryByString call chaincode query of hyperledger fabric
func ChaincodeQueryByString(chaincodeID, argsStr string) string {
	log.Infof("chaincode query: %s; %s; %s", chaincodeID, argsStr)
	argsArray, err := ArgsArray(argsStr)
	if err == nil {
		var ret string
		ret, err = ChaincodeQuery(chaincodeID, argsArray)
		if err == nil {
			log.Info("chaincode query result: ", ret)
			return ret
		}
	}
	log.Errorf("%s %v", defaultResultJSON, err)
	return defaultResultJSON
}

// TxInfo info of Tx
type TxInfo struct {
	Contract string `json:"contract,omitempty"`
	From     string `json:"from,omitempty"`
	To       string `json:"to,omitempty"`
	Amount   string `json:"amount,omitempty"`
	GasUsed  string `json:"gasUsed,omitempty"`
	GasPrice string `json:"gasPrice,omitempty"`
	TxHash   string `json:"txHash,omitempty"`
	Height   string `json:"height,omitempty"`
	Status   string `json:"status,omitempty"`
}

// TxRegister registered Tx
type TxRegister struct {
	Key       string  `json:"key,omitempty"`
	ChainName string  `json:"chain,omitempty"`
	TokenName string  `json:"token,omitempty"`
	Addr      string  `json:"address,omitempty"`
	Amount    string  `json:"amount,omitempty"`
	GasUsed   string  `json:"gasUsed,omitempty"`
	GasPrice  string  `json:"gasPrice,omitempty"`
	Info      *TxInfo `json:"info,omitempty"`
}

// BlockRegister registered block
type BlockRegister struct {
	Height string        `json:"height,omitempty"`
	Txs    []*TxRegister `json:"transactions,omitempty"`
}

// RegisterBlock register block of source chain into hyperledger fabric
func RegisterBlock(block *BlockRegister) string {
	var err error
	var bytes []byte
	if bytes, err = json.Marshal(block); err != nil {
		log.Errorf("register block error: %v", err)
		return defaultResultJSON
	}
	a := Args{
		Func: "register",
		Args: []string{"block", string(bytes)}}
	var args []Args
	args = append(args, a)
	var ret string
	ret, err = ChaincodeInvoke("wallet", args)
	if err != nil {
		log.Errorf("register block error: %v", err)
		return defaultResultJSON
	}
	log.Info(ret)
	return ret
}

// RegisterWalletByString create a new account
func RegisterWalletByString(key, chain, token string) string {
	if strings.EqualFold(chain, "ethereum") {
		if strings.EqualFold(token, "eth") {
			return ethNewAccount(key, chain, token)
		}
	}
	return errUnsuportedToken
}

// RegisterTokenByString query token info from chain
func RegisterTokenByString(chain, tokenAddress string) string {
	token, err := ethsdk.QueryTokenInfo(chain, tokenAddress)
	if err != nil {
		log.Errorf("query token info error: %v", err)
		return errUnsuportedToken
	}
	decimals := fmt.Sprintf("0x%s",
		strconv.FormatUint(uint64(token.Decimals), 16))
	a := Args{
		Func: "register",
		Args: []string{"token", tokenAddress, chain,
			token.Symbol, token.Name, decimals}}
	var args []Args
	args = append(args, a)
	var ret string
	ret, err = ChaincodeInvoke("wallet", args)
	if err != nil {
		log.Errorf("register token error: %v", err)
		return defaultResultJSON
	}
	log.Info(ret)
	return ret
}

func ethNewAccount(key, chain, token string) string {
	account, err := ethsdk.NewAccount("", key)
	if err != nil {
		log.Errorf("new account error: %v", err)
		return errUnsuportedToken
	}
	var strs []string
	height, err := ethsdk.EthBlockNumber()
	if err != nil {
		log.Errorf("eth get block number error: %v", err)
		return errUnsuportedToken
	}
	strs = append(strs, account.WalletAddress, chain, token, height)
	a := Args{Func: "register", Args: strs}
	var args []Args
	args = append(args, a)
	ret, err := ChaincodeInvoke("wallet", args)
	if err != nil {
		log.Errorf("new account error: %v", err)
		return errUnsuportedToken
	}
	return ret
}
