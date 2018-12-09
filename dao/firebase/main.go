package firebase

import (
	"bonex-middleware/config"
	"bonex-middleware/dao"
	"github.com/NaySoftware/go-fcm"
)

type firebaseDAO struct {
	client *fcm.FcmClient
}

func NewFirebase(c *config.Config) (dao.FirebaseDAO, error) {
	cl := fcm.NewFcmClient(c.Firebase.ServerKey)

	return &firebaseDAO{client: cl}, nil
}
