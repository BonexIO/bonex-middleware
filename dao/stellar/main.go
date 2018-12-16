package stellar

import (
	"bonex-middleware/config"
	"bonex-middleware/dao"
	"bonex-middleware/log"
	"fmt"
	"github.com/fatih/color"
	b "github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
	"github.com/wedancedalot/decimal"
	"net/http"
)

type Stellar struct {
	privateKey string
	client     *horizon.Client
	network    b.Network
}

func NewStellar(cfg *config.Config) dao.BlockchainDAO {
	return &Stellar{
		client: &horizon.Client{
			URL:  cfg.HorizonClientURL,
			HTTP: http.DefaultClient,
		},
		network: b.Network{cfg.NetworkPassphrase},
	}
}

func (this *Stellar) BlockchainName() string {
	return "Stellar"
}

func (this *Stellar) SendMoney(toAddress string, amount decimal.Decimal) error {
	if this.privateKey == "" {
		return fmt.Errorf("private key is not set")
	}
	log.Infof("Sending %s to %s", amount.String(), toAddress)

	tx, err := b.Transaction(
		b.SourceAccount{AddressOrSeed: this.privateKey},
		this.network,
		b.AutoSequence{SequenceProvider: this.client},
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

	resp, err := this.client.SubmitTransaction(xdr)
	if err != nil {
		return err
	}

	log.Infof("Transaction %s successfully submitted to network", resp.Hash)

	return nil
}

func (this *Stellar) SetPrivateKey(pk string) error {
	if this.privateKey != "" {
		return fmt.Errorf("private key is set already")
	}

	kp, err := keypair.Parse(pk)
	if err != nil {
		return err
	}

	if _, ok := kp.(*keypair.Full); !ok {
		return fmt.Errorf("provided key is not a private key")
	}

	fmt.Printf("Emission address derrived from priv key is: %s", color.New(color.Bold, color.FgGreen).SprintFunc()(kp.Address()))

	this.privateKey = pk

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

func (this *Stellar) GetTransaction(txHash string) (*horizon.Transaction, error) {
	tx, err := this.client.LoadTransaction(txHash)
	if err != nil {
		return nil, err
	}

	return &tx, nil
}
