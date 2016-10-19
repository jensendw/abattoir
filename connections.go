package main

import (
	"net/http"
  "encoding/json"
	"github.com/antonholmquist/jason"
  "bytes"
)

//Decides if host is dead based on various states
func isHostDead(state string) bool {
	switch state {
	case "active":
		return false
	case "reconnecting":
		return true
	case "inactive":
		return true
	default:
		return false
	}
}

//Gets all hosts from rancher
func getRancherHosts(url string) *jason.Object {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
    log.Error("Error connecting to: ", url)
	}
	req.SetBasicAuth(Config.RancherAccessKey, Config.RancherSecretKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
    log.Error("Error connecting to: ", url)
	}
	defer resp.Body.Close()

	respJSON, err := jason.NewObjectFromReader(resp.Body)
	if err != nil {
    log.Error("Error connecting to: ", url)
	}
	return respJSON
}

//Gets all EC2 instances from ops-servers endpoint
func getEC2Servers() *jason.Object {
	req, err := http.NewRequest("GET", Config.ServerAPI + "/api/server", nil)
	if err != nil {
    log.Error("Error connecting to: ", Config.ServerAPI)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
    log.Error("Error connecting to: ", Config.ServerAPI)
	}
	defer resp.Body.Close()

	respJSON, err := jason.NewObjectFromReader(resp.Body)
	if err != nil {
    log.Error("Error parsing JSON from: ", Config.ServerAPI )
	}
	return respJSON
}

//Parses servers data to determine if an IP address is listed
func serverExistsInEC2(ip string, servers *jason.Object) bool {
  data, err := servers.GetObjectArray("Data")

	if err != nil {
    log.Error("Invalid data from server API")
	}
  for _, hosts := range data {
    ipaddress, err := hosts.GetString("PrivateIpAddress")
    if err != nil {
      //handle error
    }
    if ipaddress == ip {
      log.Info("Server found with IP: ", ipaddress)
      return true
    }
  }
  return false
}

//POST to rancher endpoint to deactivate host
func deactivateRancherHost(id string) bool {
  type Payload struct {
  }

  data := Payload {
  	// fill struct
  }
  payloadBytes, err := json.Marshal(data)
  if err != nil {
    log.Error("Something went wrong deactivating rancher host", err)
  	return false
  }
  body := bytes.NewReader(payloadBytes)

  deactivationURL := "https://rancher.origami42.com/v1/projects/" + Config.RancherProjectID + "/hosts/" + id + "/?action=deactivate"
  req, err := http.NewRequest("POST", deactivationURL, body)
  if err != nil {
    log.Error("Something went wrong deactivating rancher host", err)
  	return false
  }
  req.SetBasicAuth(Config.RancherAccessKey, Config.RancherSecretKey)
  req.Header.Set("Accept", "application/json")
  req.Header.Set("Content-Type", "application/json")

  resp, err := http.DefaultClient.Do(req)
  if err != nil {
    log.Error("Something went wrong deactivating rancher host", err)
    return false
  }
  defer resp.Body.Close()
  return true
}

//POST to rancher endpoint to remove host
func removeRancherHost(id string) bool {
  type Payload struct {
  }

  data := Payload {
  	// fill struct
  }
  payloadBytes, err := json.Marshal(data)
  if err != nil {
    log.Error("Something went wrong removing rancher host", err)
  	return false
  }
  body := bytes.NewReader(payloadBytes)

  deactivationURL := "https://rancher.origami42.com/v1/projects/" + Config.RancherProjectID + "/hosts/" + id + "/?action=remove"
  req, err := http.NewRequest("POST", deactivationURL, body)
  if err != nil {
    log.Error("Something went wrong removing rancher host", err)
  	return false
  }
  req.SetBasicAuth(Config.RancherAccessKey, Config.RancherSecretKey)
  req.Header.Set("Accept", "application/json")
  req.Header.Set("Content-Type", "application/json")

  resp, err := http.DefaultClient.Do(req)
  if err != nil {
    log.Error("Something went wrong removing rancher host", err)
    return false
  }
  defer resp.Body.Close()
  return true
}
