package txs

// --------------------------------------------------------
//      Parser
//
//      Parses structs containing event information into
//      unsigned transactions for validators to sign, then
//      relays the data packets as transactions on the
//      Cosmos Bridge.
// --------------------------------------------------------

import (
  "log"
  "encoding/hex"
  "strings"
  "strconv"

  "github.com/swishlabsco/cosmos-ethereum-bridge/cmd/ebrelayer/events"
  sdk "github.com/cosmos/cosmos-sdk/types"
)

// Witness claim builds a Cosmos transaction
type WitnessClaim struct {
  Nonce          int            `json:"nonce"`
  EthereumSender string         `json:"ethereum_sender"`
  CosmosReceiver sdk.AccAddress `json:"cosmos_receiver"`
  Validator      sdk.AccAddress `json:"validator"`
  Amount         sdk.Coins      `json:"amount"`
}

func ParsePayload(validator sdk.AccAddress, event *events.LockEvent) (WitnessClaim, error) {
  
  witnessClaim := WitnessClaim{}

  // Nonce type casting (*big.Int -> int)
  nonce, nonceErr := strconv.Atoi(event.Nonce.String())
  if nonceErr != nil {
    log.Fatal(nonceErr)
  }
  witnessClaim.Nonce = nonce

  // EthereumSender type casting (address.common -> string)
  witnessClaim.EthereumSender = event.From.Hex()

  // CosmosReceiver type casting (bytes[] -> sdk.AccAddress)
  recipient, recipientErr := sdk.AccAddressFromHex(hex.EncodeToString(event.Id[:]))
  if recipientErr != nil {
    log.Fatal(recipientErr)
  }
  witnessClaim.CosmosReceiver = recipient

  // Validator is already the correct type (sdk.AccAddress)
  witnessClaim.Validator = validator

  // Amount type casting (*big.Int -> sdk.Coins)
  ethereumCoin := []string {event.Value.String(),"ethereum"}
  weiAmount, coinErr := sdk.ParseCoins(strings.Join(ethereumCoin, ""))
  if coinErr != nil {
    log.Fatal(coinErr)
  }
  witnessClaim.Amount = weiAmount

  return witnessClaim, nil
}