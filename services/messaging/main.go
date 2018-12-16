package messaging

import (
	"bonex-middleware/config"
	"bonex-middleware/dao"
	"bonex-middleware/dao/models"
	"bonex-middleware/log"
	"bonex-middleware/types"
	"context"
	"database/sql"
	"github.com/jasonlvhit/gocron"
	"github.com/stellar/go/clients/horizon"
	"strings"
	"time"
)

const cleanMessagesAfterDays = 7

type (
	Messaging struct {
		dao    dao.DAO
		config *config.Config
	}
)

func New(d dao.DAO, cfg *config.Config) *Messaging {
	return &Messaging{
		dao:    d,
		config: cfg,
	}
}

func (this *Messaging) Title() string {
	return "Messaging"
}

func (this *Messaging) GracefulStop(ctx context.Context) error {
	return nil
}

func (this *Messaging) Run() error {
	gocron.Every(this.config.Messaging.RunNotificatorEverySeconds).Seconds().Do(this.Notificate)
	gocron.Every(this.config.Messaging.RunCleanerEveryHours).Hours().Do(this.Clean)
	gocron.Every(this.config.Messaging.RunRemoverEveryMinutes).Minutes().Do(this.Remove)

	//run scheduler
	<-gocron.Start()

	return nil
}

func (this *Messaging) setMessageStatus(msg *models.Message, status models.MessageStatus) error {
	if msg.Status != status {
		dbTx, err := this.dao.BeginTx()
		if err != nil {
			return err
		}
		defer dbTx.RollbackTx()

		//load for update
		m, err := this.dao.GetMessageById(msg.Id, dbTx)
		if err != nil {
			return err
		}

		m.Status = status

		err = this.dao.UpdateMessage(m, dbTx)
		if err != nil {
			return err
		}

		err = dbTx.CommitTx()
		if err != nil {
			return err
		}

		//replace the value
		*msg = *m
	}

	return nil
}

func (this *Messaging) Notificate() {
	ms, err := this.dao.GetMessages(&types.MessageFilters{
		Statuses: []models.MessageStatus{models.MessageStatusCreated, models.MessageStatusWaitingForTx},
		//TODO do not load old and not removed messages
	})
	if err != nil {
		if err == sql.ErrNoRows {
			//no messages -- nothing to do
			return
		}
		log.Errorf("cannot load messages to Notificate users: %s")
		return
	}

	var tx *horizon.Transaction

	for _, m := range ms {
		tx, err = this.dao.GetTransaction(m.TxHash)
		if err != nil {
			log.Warnf("cannot load %s transaction %s: %s", this.dao.BlockchainName(), m.TxHash, err.Error())
			err = this.setMessageStatus(m, models.MessageStatusWaitingForTx)
			if err != nil {
				log.Errorf("cannot setMessageStatus: %s", err.Error())
				continue
			}

			continue
		}
		if strings.EqualFold(tx.Hash, m.TxHash) { //means tx was loaded successfully
			log.Infof("sending message notification for tx hash %s", tx.Hash)
			err = this.dao.SendMessage(m)
			if err != nil {
				log.Warnf("cannot send message for transaction %s: %s", tx.Hash, err.Error())
				continue
			}

			err = this.setMessageStatus(m, models.MessageStatusNotificationSent)
			if err != nil {
				log.Errorf("cannot setMessageStatus: %s", err.Error())
				continue
			}

			//TODO what if firebase lost notification? May be resend some time after
		}
	}
}

func (this *Messaging) Clean() {
	now := time.Now()
	weekBefore := now.AddDate(0, 0, -cleanMessagesAfterDays)
	ms, err := this.dao.GetMessages(&types.MessageFilters{
		//TODO: may be clean all possible messages
		Statuses: []models.MessageStatus{models.MessageStatusCreated, models.MessageStatusWaitingForTx},
		ToTime:   &weekBefore,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			//no messages -- nothing to do
			return
		}
		log.Errorf("cannot load messages to Clean: %s", err.Error())
		return
	}

	this.deleteMessages(ms)
}

func (this *Messaging) Remove() {
	ms, err := this.dao.GetMessages(&types.MessageFilters{
		Statuses: []models.MessageStatus{models.MessageStatusReceived},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			//no messages -- nothing to do
			return
		}
		log.Errorf("cannot load messages to Remove: %s", err.Error())
		return
	}

	this.deleteMessages(ms)
}

func (this *Messaging) deleteMessages(ms []*models.Message) {
	dbTx, err := this.dao.BeginTx()
	if err != nil {
		log.Errorf("cannot BeginTx: %s", err.Error())
		return
	}
	defer dbTx.RollbackTx()

	for _, m := range ms {
		err = this.dao.DeleteMessage(m, dbTx)
		if err != nil {
			log.Errorf("cannot DeleteMessage %d: %s", m.Id, err.Error())
			continue
		}
	}

	err = dbTx.CommitTx()
	if err != nil {
		log.Errorf("cannot CommitTx: %s", err.Error())
		return
	}

	return
}
