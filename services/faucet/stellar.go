package faucet

import (
	"github.com/wedancedalot/decimal"
	"github.com/stellar/go/clients/horizon"
	"net/http"
	"bonex-middleware/config"
	"fmt"
	"github.com/stellar/go/keypair"
	"github.com/fatih/color"
	b "github.com/stellar/go/build"
	"bonex-middleware/log"
)

type Stellar struct {
	privateKey    string
	horizonClient *horizon.Client
}

func NewStellar(cfg *config.Config) Blockchain {
	return &Stellar{
		horizonClient: &horizon.Client{
			URL:  cfg.HorizonClientURL,
			HTTP: http.DefaultClient,
		},
	}
}

func (this *Stellar) Title() string {
	return "Stellar"
}

func (this *Stellar) SendMoney(toAddress string, amount decimal.Decimal) error {
	log.Infof("Sending %s to %s", amount.String(), toAddress)

	tx, err := b.Transaction(
		b.SourceAccount{AddressOrSeed: this.privateKey},
		b.TestNetwork,
		b.AutoSequence{SequenceProvider: horizon.DefaultTestNetClient},
		b.Payment(
			b.Destination{AddressOrSeed: toAddress},
			b.NativeAmount{Amount: amount.String()},
		),
	)
	if err != nil {
		return err
	}

	txSigned, err := tx.Sign(this.privateKey)
	if err != nil {
		return err
	}

	xdr, err := txSigned.Base64()
	if err != nil {
		return err
	}

	resp, err := this.horizonClient.SubmitTransaction(xdr)
	if err != nil {
		return err
	}

	log.Info("Transaction %s successfully submitted: %s", resp.Hash, resp.Result)

	return nil
}

func (this *Stellar) SetPrivateKey(string) error {
	//var p [32]byte
	//copy(p[:], privKeyBytes[0:32])
	//kp, err := keypair.FromRawSeed(p)

	kp, err := keypair.Parse(this.privateKey)
	if err != nil {
		return err
	}

	if _, ok := kp.(*keypair.Full); !ok {
		return fmt.Errorf("provided key is not a private key")
	}

	fmt.Printf("Emission address derrived from priv key is: %s", color.New(color.Bold, color.FgGreen).SprintFunc()(kp.Address()))

	return nil
}

func (this *Stellar) ValidateAddress(addr string) error {
	kp, err := keypair.Parse(addr)
	if err != nil {
		return err
	}

	if _, ok := kp.(*keypair.FromAddress); !ok {
		return fmt.Errorf("provided key is not a public key")
	}

	return nil
}
