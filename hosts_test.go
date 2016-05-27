package hostsfile_test

import (
	"testing"

	"github.com/jaytaylor/go-hostsfile"
)

func TestHostsReverseLookup(t *testing.T) {
	res, err := hostsfile.ReverseLookup("127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	if len(res) == 0 {
		t.Errorf("Expected len(res) > 0 but actaul=%v res=%+v", len(res), res)
	}
}

func TestParseHosts(t *testing.T) {
	testCases := []struct {
		hostsFileContent string
		expectedEntries  map[string]int
		forbiddenEntries []string
	}{
		{
			hostsFileContent: `##
# Host Database
#
# localhost is used to configure the loopback interface
# when the system is booting.  Do not change this entry.
##
127.0.0.1	localhost localhost.local localhost.localdomain 	jays-computer jays-computer.local jays-computer.localdomain
  	  127.0.0.1   	 	 	 	 	wow.this.line.is.Ka.r.az.y
255.255.255.255	broadcasthost
::1             localhost
fe80::1%lo0	localhost

garbage

#172.1.2.12 mesos-primary1a
#172.1.2.180 mesos-primary2a
#172.1.2.182 mesos-primary3a
#172.1.2.63 mesos-worker1a
#172.1.2.115 mesos-worker2a

127.0.0.2 foo.bar
192.168.1.34 hello-app.lan talksbythebay-lan
; 192.168.1.240 should-not.resolve really-it-should.nt
`,
			expectedEntries: map[string]int{
				"127.0.0.1":    7,
				"127.0.0.2":    1,
				"192.168.1.34": 2,
			},
			forbiddenEntries: []string{
				"172.1.2.12",
				"192.168.1.240",
			},
		},
		{
			hostsFileContent: `;;
; Host Database
;
; localhost is used to configure the loopback interface
; when the system is booting.  Do not change this entry.
;;
127.0.0.1	localhost localhost.local       localhost.localdomain jays-computer jays-computer.local jays-computer.localdomain
255.255.255.255	broadcasthost
::1             localhost
fe80::1%lo0	localhost

	192.168.1.34 hello-app.lan talksbythebay-lan
	# 192.168.1.240 should-not.resolve really-it-should.nt`,
			expectedEntries: map[string]int{
				"127.0.0.1":    6,
				"192.168.1.34": 2,
			},
			forbiddenEntries: []string{
				"192.168.1.240",
			},
		},
	}

	for i, testCase := range testCases {
		res, err := hostsfile.ParseHosts([]byte(testCase.hostsFileContent), nil)
		if err != nil {
			t.Fatalf("[i=%v] Error parsing hosts content: %s", i, err)
		}
		for entry, expectedCount := range testCase.expectedEntries {
			if reverses, ok := res[entry]; ok {
				if expected, actual := expectedCount, len(reverses); actual != expected {
					t.Errorf("[i=%v] Expected len(res['%v'])=%v but actual=%v; reverses=%+v", i, entry, expected, actual, reverses)
				}
			} else {
				t.Errorf("[i=%v] Expected to find entries for ip=%v but none were found", i, entry)
			}
		}
		for _, forbiddenEntry := range testCase.forbiddenEntries {
			if reverses, ok := res[forbiddenEntry]; ok {
				t.Errorf("[i=%v] Expected '%v' to be absent from res but actual=%+v", i, forbiddenEntry, reverses)
			}
		}
	}
}
