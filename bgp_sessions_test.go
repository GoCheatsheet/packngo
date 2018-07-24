package packngo

import (
	"testing"
)

func TestAccBGPSession(t *testing.T) {
	skipUnlessAcceptanceTestsAllowed(t)

	c, projectID, teardown := setupWithProject(t)
	defer teardown()

	hn := randString8()
	cr := DeviceCreateRequest{
		Hostname:     hn,
		Facility:     testFacility(),
		Plan:         "baremetal_0",
		ProjectID:    projectID,
		BillingCycle: "hourly",
		OS:           "coreos_stable",
	}

	d, _, err := c.Devices.Create(&cr)
	if err != nil {
		t.Fatal(err)
	}

	d, err = waitDeviceActive(d.ID, c)
	if err != nil {
		t.Fatal(err)
	}

	bgpSession, _, err := c.BGPSessions.Create(d.ID, CreateBGPSessionRequest{AddressFamily: "ipv4"})
	if err != nil {
		t.Fatal(err)
	}

	sessionID := bgpSession.ID

	sessions, _, err := c.Devices.ListBGPSessions(d.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	var check *BGPSession
	for _, s := range sessions {
		if s.ID == sessionID {
			check = &s
			break
		}
	}

	if check == nil {
		t.Fatal("BGP Session not returned.")
	}

	sessions, _, err = c.Projects.ListBGPSessions(projectID, nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, s := range sessions {
		if s.ID == sessionID {
			check = &s
			break
		}
	}

	if check == nil {
		t.Fatal("BGP Session not returned.")
	}

	session, _, err := c.BGPSessions.Get(sessionID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if session == nil {
		t.Fatal("Session not retrieved")
	}

	_, err = c.BGPSessions.Delete(sessionID)
	if err != nil {
		t.Fatal(err)
	}
	session, _, err = c.BGPSessions.Get(sessionID, nil)
	if session != nil {
		t.Fatal("Session not deleted")
	}
	if err == nil {
		t.Fatal("Session not deleted")
	}

	c.Devices.Delete(d.ID)
}
