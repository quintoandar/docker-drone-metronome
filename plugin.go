package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/drone/envsubst"
	"github.com/ghodss/yaml"

	log "github.com/Sirupsen/logrus"
)

// Plugin defines the parameters
type Plugin struct {
	Server        string
	Metronomefile string
	JobConfig     string
}

// Exec runs the plugin
func (p *Plugin) Exec() error {

	log.WithFields(log.Fields{
		"server":        p.Server,
		"metronomefile": p.Metronomefile,
	}).Info("attempting to start job")

	data, err := p.ReadInput()

	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to read metronomefile/job_config input data")
		return err
	}

	b, err := parseData(data)

	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to parse input data into JSON format: ", string(b))
		return err
	}

	var v map[string]interface{}

	if err := json.Unmarshal(b, &v); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Errorln("failed to unmarshal metronome data:", string(b))
		return err
	}

	if _, ok := v["id"]; !ok {
		err := errors.New("invalid data")
		log.WithFields(log.Fields{
			"err": err,
		}).Error("metronome data is missing 'id' key:", string(b))
		return err
	}

	var buff bytes.Buffer

	if err := json.Indent(&buff, b, "", "\t"); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to parse JSON: ", string(b))
		return err
	}

	log.Info("sending data to metronome server")

	u, err := url.Parse(p.Server)

	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to parser metronome url")
		return err
	}

	u.Path = fmt.Sprintf("/v1/jobs/%s", v["id"])
	log.Infoln("GET", u.String())

	resp, err := http.Get(u.String())

	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("error getting current job")
		return err
	}

	var method string

	if resp.StatusCode == 200 {
		log.Infoln("Updating job")
		method = http.MethodPut
		u.Path = fmt.Sprintf("/v0/scheduled-jobs/%s", v["id"])
	} else {
		log.Infoln("Creating new job")
		method = http.MethodPost
		u.Path = fmt.Sprintf("/v0/scheduled-jobs")
	}

	log.Infoln(method, u.String())
	req, err := http.NewRequest(method, u.String(), &buff)

	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("error creating request")
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)

	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("error creating request")
		return err
	}

	if resp.StatusCode >= 300 {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err == nil {
			err = errors.New(string(body))
		}

		log.WithFields(log.Fields{
			"status": resp.Status,
			"err":    err,
		}).Error("error updating application")
		return err
	}

	return nil
}

// ReadInput reads Metronomefile/Appconfig data
func (p Plugin) ReadInput() (data string, err error) {
	if p.Metronomefile != "" {
		log.Info("parsing metronomefile ", p.Metronomefile)

		// When 0.9 comes out, limit to secrets and other Drone variables
		b, err := ioutil.ReadFile(p.Metronomefile)

		if err != nil {
			return "", err
		}

		return envsubst.EvalEnv(string(b))
	}

	if p.JobConfig != "" {
		log.Warn("job_config is deprecated and will be removed, please use a metronomefile instead")

		return envsubst.EvalEnv(string(p.JobConfig))
	}

	err = errors.New("missing parameters")
	return
}

func parseData(data string) (b []byte, err error) {
	if isYAML(data) {
		log.Info("data is in YAML format, parsing into JSON")
		return yaml.YAMLToJSON([]byte(data))
	}

	if isJSON(data) {
		log.Info("data is in JSON format, no need to parse")
		return
	}

	err = errors.New("invalid data")
	return
}

func isJSON(s string) bool {
	var j map[string]interface{}
	return json.Unmarshal([]byte(s), &j) == nil
}

func isYAML(s string) bool {
	var y map[string]interface{}
	return yaml.Unmarshal([]byte(s), &y) == nil
}
