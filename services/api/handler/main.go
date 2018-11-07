package handler

import (
	"github.com/BonexIO/go/build"
	"github.com/BonexIO/go/clients/horizon"
	"github.com/BonexIO/go/keypair"
	"github.com/BonexIO/go/network"
	"github.com/BonexIO/go/xdr"
)

func CreateAccountWithNativeTx(source, destination, amount string, accountType xdr.AccountType) error {

	passphrase := network.PublicNetworkPassphrase

	tx, err := build.Transaction(
		build.PublicNetwork,
		build.SourceAccount{source},
		build.Network{passphrase},
		build.AutoSequence{horizon.DefaultPublicNetClient},
		build.CreateAccount(
			build.Destination{destination},
			build.NativeAmount{amount},
			build.AccountType{accountType},
		),
	)

	if err != nil {
		return err
	}

	txe, err := tx.Sign(source)
	if err != nil {
		return err
	}

	txeB64, err := txe.Base64()
	if err != nil {
		return err
	}

	_, err = horizon.DefaultPublicNetClient.SubmitTransaction(txeB64)
	if err != nil {
		return err
	}

	return nil
}

func SendNativeTx(source, destination, amount string) error {

	passphrase := network.PublicNetworkPassphrase

	tx, err := build.Transaction(
		build.PublicNetwork,
		build.SourceAccount{source},
		build.Network{passphrase},
		build.AutoSequence{horizon.DefaultPublicNetClient},
		build.Payment(
			build.Destination{destination},
			build.NativeAmount{amount},
		),
	)

	if err != nil {
		return err
	}

	txe, err := tx.Sign(source)
	if err != nil {
		return err
	}

	txeB64, err := txe.Base64()
	if err != nil {
		return err
	}

	_, err = horizon.DefaultPublicNetClient.SubmitTransaction(txeB64)
	if err != nil {
		return err
	}

	return nil
}

func SendAssetTx(sourcePrivateKey, destination, asset, amount string) error {

	sourceKeyPair, _ := keypair.Parse(sourcePrivateKey)

	tx, err := build.Transaction(
		build.PublicNetwork,
		build.SourceAccount{sourcePrivateKey},
		build.AutoSequence{horizon.DefaultPublicNetClient},
		build.Payment(
			build.Destination{destination},
			build.CreditAmount{asset, sourceKeyPair.Address(), amount},
		),
	)

	if err != nil {
		return err
	}

	txe, err := tx.Sign(sourcePrivateKey)
	if err != nil {
		return err
	}

	txeB64, err := txe.Base64()
	if err != nil {
		return err
	}

	_, err = horizon.DefaultPublicNetClient.SubmitTransaction(txeB64)
	if err != nil {
		return err
	}

	return nil
}
