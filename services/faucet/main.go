package faucet

import (
	"bonex-middleware/config"
	"bonex-middleware/dao"
	"context"
	"fmt"
	"syscall"
	"golang.org/x/crypto/ssh/terminal"
	"github.com/jasonlvhit/gocron"
	"bonex-middleware/log"
	"github.com/wedancedalot/decimal"
	"bonex-middleware/types"
)

type (
	Faucet struct {
		dao    dao.DAO
		config *config.Config
		sender Blockchain
	}

	Blockchain interface {
		Title() string
		SendMoney(string, decimal.Decimal) error
		SetPrivateKey(string) error
		ValidateAddress(string) error
	}
)

const GreyListTTL = 24 * 60 * 60 //seconds

func New(d dao.DAO, cfg *config.Config, s Blockchain) *Faucet {
	return &Faucet{
		dao: d,
		config: cfg,
		sender: s,
	}
}

func (this *Faucet) Title() string {
	return "Faucet"
}

func (this *Faucet) GracefulStop(ctx context.Context) error {
	return nil
}

func (this *Faucet) PromtKey() error {
	//color.New(color.Bold, color.FgGreen).SprintFunc()(this.GetName()), color.New(color.FgBlue).SprintFunc()(msg)
	fmt.Printf("Enter emission account private key: ")
	privKeyBytes, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return err
	}

	err = this.sender.SetPrivateKey(string(privKeyBytes))
	fmt.Println(" ")

	return nil
}

func (this *Faucet) Run() error {
	gocron.Every(this.config.Faucet.RunEverySeconds).Seconds().Do(this.Do)

	//run scheduler
	<-gocron.Start()

	return nil
}

func (this *Faucet) Do() {
	req, err := this.dao.PopFromQueue()
	if err != nil {
		log.Errorf("Cannot PopFromQueue: %s", err.Error())
		return
	}
	if req == nil { //nothing to do
		return
	}

	err = this.sender.ValidateAddress(req.Address)
	if err != nil {
		log.Errorf("Invalid address provided: %s", err.Error())
		return
	}

	err = this.sender.SendMoney(req.Address, req.Amount)
	if err != nil {
		log.Errorf("Cannot SendMoney: %s", err.Error())
		return
	}

	return
}

func (this *Faucet) AddToQueue(qi *types.QueueItem, ipAddress string) error {
	//basic qi validations
	err := qi.Validate()
	if err != nil {
		return err
	}

	//more specific validations
	if qi.Amount.GreaterThan(this.config.Faucet.MaxAllowed24HoursValue) {
		log.Errorf("too big amount %s (max is %s)", qi.Amount.String(), this.config.Faucet.MaxAllowed24HoursValue.String())
		return types.NewError(types.ErrBadParam, "amount")
	}

	if err := this.sender.ValidateAddress(qi.Address); err != nil {
		log.Errorf("bad address format: %s", err.Error())
		return types.NewError(types.ErrBadParam, "address")
	}

	volume, err := this.dao.GetAccountVolume(ipAddress)
	if err != nil {
		log.Errorf("cannot GetAccountVolume: %s", err.Error())
		return err
	}

	volume = volume.Add(qi.Amount)

	if volume.GreaterThan(this.config.Faucet.MaxAllowed24HoursValue) {
		log.Errorf("daily limit %s reached on %s", this.config.Faucet.MaxAllowed24HoursValue.String(), volume.String())
		return types.NewError(types.ErrLimitReached, volume.String())
	}

	err = this.dao.SetAccountVolume(ipAddress, volume, GreyListTTL)
	if err != nil {
		log.Errorf("cannot SetAccountVolume: %s", err.Error())
		return err
	}

	err = this.dao.AddToQueue(qi)
	if err != nil {
		log.Errorf("cannot AddToQueue: %s", err.Error())
		return err
	}

	return nil
}