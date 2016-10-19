package main

import (
	"bytes"
	"encoding/json"
	"github.com/antonholmquist/jason"
	"net/http"
	"strings"
	"time"
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
	req, err := http.NewRequest("GET", Config.ServerAPI+"/api/server", nil)
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
		log.Error("Error parsing JSON from: ", Config.ServerAPI)
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

	data := Payload{
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

	data := Payload{
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

//Wrapper function to clean things up a little
func seekAndDestroyBadHosts() {
	deadcows := checkForDeadHosts()
	if len(deadcows) != 0 {
		findCowsToKill(deadcows)
	} else {
		log.Info("No dead rancher hosts found")
	}
}

//Checks for hosts that are disconnected/reconnecting in Rancher
func checkForDeadHosts() map[string]string {
	DeadCows := make(map[string]string)

	hosts := getRancherHosts(Config.RancherURL + "/v1/hosts")

	data, err := hosts.GetObjectArray("data")
	if err != nil {
		log.Error("Error getting hostame from parsed JSON")
	}

	for _, hosts := range data {
		state, err := hosts.GetString("agentState")
		if err != nil {
			//log.Error("Error getting agentState from parsed JSON")
		}
		if isHostDead(state) {
			hostname, err := hosts.GetString("hostname")
			if err != nil {
				log.Error("Error getting hostame from parsed JSON")
			}
			DeadHostID, _ := hosts.GetString("id")
			DeadHostIP := convertHostnameToIP(hostname)
			log.Info("Found a dead host:", DeadHostIP)
			DeadCows[DeadHostID] = DeadHostIP
		}
	}
	return DeadCows
}

//Compares servers in EC2 with disconnected rancher hosts
func findCowsToKill(deadcows map[string]string) {
	servers := getEC2Servers()
	for id, ip := range deadcows {
		if serverExistsInEC2(ip, servers) {
			log.Info("cow to kill:", id)
			if deactivateRancherHost(id) {
				log.Info("Succesfully deactivated rancher host:", id)
			}
			time.Sleep(3 * time.Second)
			if removeRancherHost(id) {
				log.Info("Succesfully removed rancher host:", id)
			}
		}
	}
}

//Simple string conversion of hostname to IP address
func convertHostnameToIP(hostname string) string {
	step1 := strings.Replace(hostname, "ip-", "", 1)
	step2 := strings.Replace(step1, "-", ".", -1)
	return step2
}
