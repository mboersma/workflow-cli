package fleet

import (
	"sync"
	"testing"

	"github.com/coreos/fleet/schema"
)

func TestStop(t *testing.T) {
	t.Parallel()

	testUnits := []*schema.Unit{
		&schema.Unit{
			Name:         "deis-controller.service",
			DesiredState: "launched",
		},
		&schema.Unit{
			Name:         "deis-builder.service",
			DesiredState: "launched",
		},
		&schema.Unit{
			Name:         "deis-publisher.service",
			DesiredState: "launch",
		},
	}

	testFleetClient := stubFleetClient{testUnits: testUnits, unitsMutex: &sync.Mutex{},
		unitStatesMutex: &sync.Mutex{}}

	c := &FleetClient{Fleet: &testFleetClient}

	var errOutput string
	var wg sync.WaitGroup

	logMutex := sync.Mutex{}

	se := newOutErr()
	c.Stop([]string{"controller", "builder", "publisher"}, &wg, se.out, se.ew)

	wg.Wait()

	logMutex.Lock()
	if errOutput != "" {
		t.Fatal(errOutput)
	}
	logMutex.Unlock()

	expected := []string{"deis-controller.service", "deis-builder.service", "deis-publisher.service"}

	for _, expectedUnit := range expected {
		found := false

		for _, unit := range testFleetClient.testUnitStates {
			if unit.Name == expectedUnit {
				found = true

				if unit.SystemdSubState != "dead" {
					t.Errorf("Unit %s is %s, expected dead", unit.Name, unit.SystemdSubState)
				}

				break
			}
		}

		if !found {
			t.Errorf("Expected Unit %s not found in Unit States", expectedUnit)
		}
	}
}
