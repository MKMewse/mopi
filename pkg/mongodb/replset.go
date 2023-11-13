package mongodb

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReplsetStatusCode int

const (
	STATUS_STARTUP    ReplsetStatusCode = 0
	STATUS_PRIMARY    ReplsetStatusCode = 1
	STATUS_SECONDARY  ReplsetStatusCode = 2
	STATUS_RECOVERING ReplsetStatusCode = 3
	STATUS_STARTUP2   ReplsetStatusCode = 5
	STATUS_UNKNOWN    ReplsetStatusCode = 6
	STATUS_ARBITER    ReplsetStatusCode = 7
	STATUS_DOWN       ReplsetStatusCode = 8
	STATUS_ROLLBACK   ReplsetStatusCode = 9
	STATUS_REMOVED    ReplsetStatusCode = 10
)

func (rsc ReplsetStatusCode) String() string {
	switch rsc {
	case STATUS_STARTUP:
		return "STARTUP"
	case STATUS_PRIMARY:
		return "PRIMARY"
	case STATUS_SECONDARY:
		return "SECONDARY"
	case STATUS_RECOVERING:
		return "RECOVERING"
	case STATUS_STARTUP2:
		return "STARTUP2"
	case STATUS_UNKNOWN:
		return "UNKNOWN"
	case STATUS_ARBITER:
		return "ARBITER"
	case STATUS_DOWN:
		return "DOWN"
	case STATUS_ROLLBACK:
		return "ROLLBACK"
	case STATUS_REMOVED:
		return "REMOVED"
	default:
		return "No status"
	}
}

type ReplSetStatus struct {
	Members []ReplSetMember `json:"members" bson:"members"`
}

type ReplSetMember struct {
	State    int    `json:"state" bson:"state"`
	StateStr string `json:"stateStr" bson:"stateStr"`
	Name     string `json:"name" bson:"name"`
}

func (rss *ReplSetStatus) String() string {
	var out []string
	for _, m := range rss.Members {
		out = append(out, m.String())
	}
	return strings.Join(out, ", ")
}

func (rsm *ReplSetMember) String() string {
	return fmt.Sprintf("Member %s is in state %s(%d)", rsm.Name, rsm.StateStr, rsm.State)
}

type MongoReplicator interface {
	Status() (*ReplSetStatus, error)
	StepDown() error
}

type MongoReplicationManager struct {
	Client  *mongo.Client
	Context context.Context
}

func (mrm *MongoReplicationManager) Status() (*ReplSetStatus, error) {
	db := mrm.Client.Database("admin", nil)
	cmd := bson.D{bson.E{Key: "replSetGetStatus", Value: 1}}

	var rss ReplSetStatus
	if err := db.RunCommand(mrm.Context, cmd).Decode(&rss); err != nil {
		return nil, err
	} else {
		return &rss, nil
	}
}

type StepDownResponse struct {
	Ok    int    `json:"ok"`
	Error string `json:"errmsg",omitempty`
}

func (mrm *MongoReplicationManager) StepDown() error {
	db := mrm.Client.Database("admin", nil)
	cmd := bson.D{bson.E{Key: "replSetStepDown", Value: 60}, bson.E{Key: "secondaryCatchUpPeriodSecs", Value: 30}}

	var sdr StepDownResponse
	if err := db.RunCommand(mrm.Context, cmd).Decode(&sdr); err != nil {
		return err
	}
	if sdr.Ok != 1 {
		return fmt.Errorf("stepdown failed. Error: %s", sdr.Error)
	}

	return nil
}
