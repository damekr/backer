package transfer

import (
	"github.com/d8x/bftp"
	"github.com/damekr/backer/baclnt/config"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{"prefix": "transfer"})

func StartBFTPServer() {
	bftpServer := bftp.CreateBFTPServer()
	bftpServer.SetIP(config.MainConfig.DataTransferInterface)
	bftpServer.SetPort(8000)
	log.Printf("Starting server at: %s:%d", bftpServer.Params.Server, bftpServer.Params.Port)
	bftpServer.StartServer()
}
