package main

import (
	"strings"
)

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
    if serverExistsInEC2(ip, servers) != true {
      log.Info("cow to kill:", id)
      if deactivateRancherHost(id) {
        log.Info("Succesfully deactivated rancher host:", id)
      }
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
