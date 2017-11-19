package transfer

import (
	"github.com/d8x/bftp/common"
	"net"
	"testing"

	"github.com/d8x/bftp/server"
	"io/ioutil"
	"log"
	"os"
)

func TestSessionNegotiateMismatch(t *testing.T) {
	srv, clnt := net.Pipe()

	params := common.NewConnParameters()
	sessionClnt := NewSession(1, params, clnt)
	sessionSrv := server.NewSession(1, params, srv)
	sessionSrv.Conn = srv
	sessionClnt.Conn = clnt
	go func() {
		defer clnt.Close()
		defer srv.Close()
		err := sessionClnt.negotiate("0.2")
		if err != nil {
			t.Error(err)
		}
	}()

	err := sessionSrv.Negotiate("0.1")
	if err != common.ProtocolVersionMismatch {
		t.Error("Wrong handled mismatch protocol version")
	}

}

func TestSessionNegotiateCorrect(t *testing.T) {
	srv, clnt := net.Pipe()
	params := common.NewConnParameters()
	sessionClnt := NewSession(1, params, clnt)
	sessionSrv := server.NewSession(1, params, srv)
	sessionSrv.Conn = srv
	sessionClnt.Conn = clnt
	go func() {
		defer clnt.Close()
		defer srv.Close()
		err := sessionClnt.negotiate("0.1")
		if err != nil {
			t.Error(err)
		}

	}()

	err := sessionSrv.Negotiate("0.1")
	if err != nil {
		t.Error("Protocol version mismatch")
	}
}

func TestSessionAuthenticateWrong(t *testing.T) {
	srv, clnt := net.Pipe()
	params := common.NewConnParameters()
	sessionClnt := NewSession(1, params, clnt)
	sessionSrv := server.NewSession(1, params, srv)
	sessionSrv.Conn = srv
	sessionClnt.Conn = clnt
	go func() {
		defer clnt.Close()
		defer srv.Close()
		err := sessionClnt.authenticate("WrongPassword")
		if err != common.AuthenticationFailed {
			t.Error("Authentication should fail, error: ", err)
		}

	}()

	err := sessionSrv.Authenticate("CorrectPassword")
	if err != common.AuthenticationFailed {
		t.Error("Authentication should fail, error: ", err)
	}
}

func TestSessionAuthenticateCorrect(t *testing.T) {
	srv, clnt := net.Pipe()
	params := common.NewConnParameters()
	sessionClnt := NewSession(1, params, clnt)
	sessionSrv := server.NewSession(1, params, srv)
	sessionSrv.Conn = srv
	sessionClnt.Conn = clnt
	go func() {
		defer clnt.Close()
		defer srv.Close()
		err := sessionClnt.authenticate("Correct")
		if err != nil {
			t.Error("Authentication should fail, error: ", err)
		}

	}()

	err := sessionSrv.Authenticate("Correct")
	if err != nil {
		t.Error("Authentication should fail, error: ", err)
	}
}

func TestSessionFileUpload(t *testing.T) {
	srv, clnt := net.Pipe()
	params := common.NewConnParameters()
	transfer := new(common.Transfer)
	transfer.Buffer = 1024
	sessionClnt := NewSession(1, params, clnt)
	sessionClnt.Transfer = transfer
	sessionClnt.Conn = clnt
	sessionClntDownload := NewSession(1, params, clnt)
	sessionClntDownload.Transfer = transfer
	sessionClntDownload.Conn = srv
	size10K := int64(10 * 1024)
	tmpFile10K, err := ioutil.TempFile(os.TempDir(), "common-file10K")
	if err != nil {
		t.Log("Cannot create temporary file 10K for testing")
	}
	tmpFile10KDownload, err := ioutil.TempFile(os.TempDir(), "common-file10KDownlaod")
	if err != nil {
		t.Log("Cannot create temporary file 10K for testing")
	}
	t.Log("Created file to upload")
	if err := tmpFile10K.Truncate(size10K); err != nil {
		log.Fatal(err)
	}

	size128M := int64(128 * 1024 * 1024)
	tmpFile128M, err := ioutil.TempFile(os.TempDir(), "common-file128M")
	if err != nil {
		t.Log("Cannot create temporary file 128M for testing")
	}
	tmpFile128MDownload, err := ioutil.TempFile(os.TempDir(), "common-file128MDownload")
	if err != nil {
		t.Log("Cannot create temporary file 10K for testing")
	}
	t.Log("Created file to upload")
	if err := tmpFile128M.Truncate(size128M); err != nil {
		t.Fatal(err)
	}

	size128 := int64(128)
	tmpFile128, err := ioutil.TempFile(os.TempDir(), "common-file128")
	if err != nil {
		t.Log("Cannot create temporary file 128 for testing")
	}
	tmpFile128Download, err := ioutil.TempFile(os.TempDir(), "common-file128Download")
	if err != nil {
		t.Log("Cannot create temporary file 128 for testing")
	}
	t.Log("Created file to upload")
	if err := tmpFile128.Truncate(size128); err != nil {
		t.Fatal(err)
	}

	defer func() {
		clnt.Close()
		srv.Close()
		os.Remove(tmpFile10K.Name())
		os.Remove(tmpFile10KDownload.Name())
		os.Remove(tmpFile128M.Name())
		os.Remove(tmpFile128MDownload.Name())
		os.Remove(tmpFile128.Name())
		os.Remove(tmpFile128Download.Name())
	}()

	go func() {
		err := sessionClnt.uploadFile(tmpFile10K, size10K)
		if err != nil {
			t.Error("Cannot upload file, error: ", err)
		}
	}()

	err = sessionClntDownload.downloadFile(tmpFile10KDownload, size10K)
	if err != nil {
		t.Error("Error during downloading file, error: ", err)
	}
	t.Log("10K file test pass")

	go func() {
		err := sessionClnt.uploadFile(tmpFile128M, size128M)
		if err != nil {
			t.Error("Cannot upload file, error: ", err)
		}
	}()

	err = sessionClntDownload.downloadFile(tmpFile128MDownload, size128M)
	if err != nil {
		t.Error("Error during downloading file, error: ", err)
	}

	t.Log("128M file test pass")

	go func() {
		err := sessionClnt.uploadFile(tmpFile128, size128)
		if err != nil {
			t.Error("Cannot upload file, error: ", err)
		}
	}()

	err = sessionClntDownload.downloadFile(tmpFile128Download, size128)
	if err != nil {
		t.Error("Error during downloading file, error: ", err)
	}

	t.Log("128 bytes file test pass")
}
