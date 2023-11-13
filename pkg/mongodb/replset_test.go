package mongodb

import (
	"testing"
)

func TestReplSetString(t *testing.T) {
	for i := 0; i <= 11; i++ {
		test := ReplsetStatusCode(i)
		if i == 4 || i == 11 {
			if test.String() != "No status" {
				t.Errorf("%d should have been no status but got, %s", i, test.String())
			}
		} else {
			if test.String() == "No status" {
				t.Errorf("%d returned no status when it should have been a value", i)
			}
		}
	}
}

func TestReplSetStatusString(t *testing.T) {
	rss := &ReplSetStatus{
		Members: []ReplSetMember{
			{
				Name: "t1", State: int(STATUS_PRIMARY), StateStr: STATUS_PRIMARY.String(),
			},
			{
				Name: "t2", State: int(STATUS_SECONDARY), StateStr: STATUS_SECONDARY.String(),
			},
			{
				Name: "t3", State: int(STATUS_ARBITER), StateStr: STATUS_ARBITER.String(),
			},
		},
	}

	str := rss.String()
	if str != "Member t1 is in state PRIMARY(1), Member t2 is in state SECONDARY(2), Member t3 is in state ARBITER(7)" {
		t.Errorf("replset status string incorrect, got \"%s\"", rss.String())
	}
}
