package http

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmhttp"
	"go.elastic.co/apm/module/apmlogrus"
)

func Do(req *http.Request, ctx context.Context, log *logrus.Logger) (string, error) {
	// logger
	vars := mux.Vars(req)
	logAPM := log.WithFields(apmlogrus.TraceContext(req.Context()))
	logAPM.WithField("vars", vars).Info("handling call to service b request")

	client := apmhttp.WrapClient(http.DefaultClient)
	resp, err := client.Do(req.WithContext(ctx))

	if err != nil {
		log.WithError(err).Error("Fail to call service b")
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Error("fail to decode res of service b")
		return "", err
	}

	if resp.StatusCode != 200 {
		log.WithError(err).Error("failed to request servie b")
		return "", fmt.Errorf("StatusCode: %d, Body: %s", resp.StatusCode, body)
	}
	return string(body), nil
}
