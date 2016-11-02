package main

import (
	"reflect"
	"testing"

	"github.com/antonholmquist/jason"
	"github.com/jarcoal/httpmock"
)

func LoadTestConfig() Configuration {
	var s = Configuration{
		RancherAccessKey: "rancheraccesskey",
		RancherSecretKey: "ranchersecretkey",
		RancherURL:       "https://rancher.test.com",
		ServerAPI:        "https://serverapi.test.com",
		RancherProjectID: "e1a",
		RunInterval:      30}

	return s
}

var TestConfig = LoadTestConfig()

func TestIsHostDeadActive(t *testing.T) {
	output := isHostDead("active")
	if output {
		t.Error("Expected false, got ", output)
	}
}

func TestIsHostDeadReconnecting(t *testing.T) {
	output := isHostDead("reconnecting")
	if !output {
		t.Error("Expected true, got ", output)
	}
}

func TestIsHostDeadInactive(t *testing.T) {
	output := isHostDead("inactive")
	if !output {
		t.Error("Expected true, got ", output)
	}
}

func TestIsHostDeadOther(t *testing.T) {
	output := isHostDead("blah")
	if output {
		t.Error("Expected true, got ", output)
	}
}

func TestConvertHostnameToIP(t *testing.T) {
	output := convertHostnameToIP("ip-172-16-12-11")
	if output != "172.16.12.11" {
		t.Error("Expected true, got ", output)
	}
}

func TestGetRancherHosts(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", TestConfig.RancherURL+"/v1/hosts",
		httpmock.NewStringResponder(200, `{"id": 1, "host": "somehost"}`))
	//t.Error("Expected pointer got ", TestConfig.RancherURL)
	response := getRancherHosts(&TestConfig)

	if reflect.TypeOf(response).Kind() != reflect.Ptr {
		t.Error("Expected pointer got ", reflect.TypeOf(response).Kind())
	}
}

func TestGetEC2Servers(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", TestConfig.ServerAPI+"/api/server",
		httpmock.NewStringResponder(200, `{"id": 1, "name": "My Great Article"}`))

	response := getEC2Servers(&TestConfig)
	if reflect.TypeOf(response).Kind() != reflect.Ptr {
		t.Error("Expected pointer got ", reflect.TypeOf(response).Kind())
	}
}

func TestServerExistsInEC2(t *testing.T) {
	// Make sure that it finds an IP address when correct JSON is provided
	jsonIn := []byte(`{ "Data": [ { "PrivateIpAddress": "42.42.42.42" }, {"PrivateIpAddress": "192.168.99.1"} ] }`)
	jasonOut, _ := jason.NewObjectFromBytes(jsonIn)

	if serverExistsInEC2("42.42.42.42", jasonOut) != true {
		t.Error("Expected true got false")
	}

	if serverExistsInEC2("1.1.1.1", jasonOut) {
		t.Error("Expected false got true when using incorrect ip address")
	}

	//Test with incorrect json key for data array
	jsonIn1 := []byte(`{ "mooh": [ { "PrivateIpAddress": "42.42.42.42" }, {"PrivateIpAddress": "192.168.99.1"} ] }`)
	jsonOut1, _ := jason.NewObjectFromBytes(jsonIn1)

	if serverExistsInEC2("1.1.1.1", jsonOut1) {
		t.Error("JSON array key was not data, got true")
	}
	//Test with incorrect json key for PrivateIPAddress
	jsonIn2 := []byte(`{ "Data": [ { "IpAddress": "42.42.42.42" }, {"IpAddress": "192.168.99.1"} ] }`)
	jsonOut2, _ := jason.NewObjectFromBytes(jsonIn2)

	if serverExistsInEC2("42.42.42.42", jsonOut2) {
		t.Error("JSON key PrivateIPAddress needs to exist got true")
	}

}

func TestDeactivateRancherHost(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", TestConfig.RancherURL+"/v1/projects/"+TestConfig.RancherProjectID+"/hosts/1/?action=deactivate",
		httpmock.NewStringResponder(200, `{"id": 1, "status": "gone"}`))

	response := deactivateRancherHost("1", &TestConfig)
	if response != true {
		t.Error("Expected return true when deactivating rancher host")
	}

	response1 := deactivateRancherHost("2", &TestConfig)
	if response1 == true {
		t.Error("Expected return false when deactivating non existent host")
	}
}

func TestRemoveRancherHost(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", TestConfig.RancherURL+"/v1/projects/"+TestConfig.RancherProjectID+"/hosts/1/?action=remove",
		httpmock.NewStringResponder(200, `{"id": 1, "status": "gone"}`))

	response := removeRancherHost("1", &TestConfig)
	if response != true {
		t.Error("Expected return true when removing rancher host")
	}

	response1 := deactivateRancherHost("2", &TestConfig)
	if response1 == true {
		t.Error("Expected return false when removing non existent host")
	}
}

func TestSeekAndDestroyBadHosts(t *testing.T) {
	//test the whole thing here
}
