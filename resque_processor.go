package resque_status

import (
	"encoding/json"

	"github.com/Shop2market/goworker"
)

// ResqueProcessor - Main Processor type
type ResqueProcessor struct {
	JobName string
	Lock
	Handler
}

// JobParams - arguments from resque jobs
type JobParams map[string]interface{}

// ExpiresIn - expiration config for jobs, default 24h
var ExpiresIn *int64

func init() {
	var defaultExpiresIn int64
	defaultExpiresIn = 24 * 60 * 60 // 24 hours in seconds
	ExpiresIn = &defaultExpiresIn
}

// Handler Job handler function
type Handler func(JobParams) error

func NewResqueProcessor(jobName, lockKeyPrefix string, keyParamNames []string, handler Handler) *ResqueProcessor {
	return &ResqueProcessor{JobName: jobName, Lock: Lock{LockKeyPrefix: lockKeyPrefix, KeyParamNames: keyParamNames}, Handler: handler}
}

func (rp *ResqueProcessor) Process(queue string, args ...interface{}) error {
	jobUUID := args[0].(string)

	err := rp.updateStatus(jobUUID, "working")
	if err != nil {
		return err
	}
	defer rp.updateStatus(jobUUID, "completed")

	params := args[1].(map[string]interface{})

	defer rp.unlock(params)

	return rp.Handler(params)
}

func (rp *ResqueProcessor) unlock(params map[string]interface{}) error {
	conn, err := goworker.GetConn()
	defer goworker.PutConn(conn)
	if err != nil {
		return err
	}
	_, err = conn.Do("DEL", rp.Key(params))
	return err
}

func (rp *ResqueProcessor) updateStatus(uuid, statusString string) error {
	conn, err := goworker.GetConn()
	defer goworker.PutConn(conn)
	if err != nil {
		return err
	}

	serializedStatus, err := rp.readJobStatus(conn, uuid)
	if err != nil {
		return err
	}

	serializedStatus.Status = statusString
	serializedStatus.Name = rp.JobName

	err = rp.saveJobStatus(conn, uuid, serializedStatus)
	if err != nil {
		return err
	}

	return nil
}

func (rp *ResqueProcessor) saveJobStatus(conn *goworker.RedisConn, uuid string, serializedStatus status) error {
	statusBytes, err := json.Marshal(serializedStatus)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", "resque:status:"+uuid, statusBytes)
	if err != nil {
		return err
	}
	if ExpiresIn != nil {
		_, err = conn.Do("EXPIRE", "resque:status:"+uuid, *ExpiresIn)
		if err != nil {
			return err
		}
	}
	return nil
}

func (rp *ResqueProcessor) readJobStatus(conn *goworker.RedisConn, uuid string) (status, error) {
	jobStatus, err := conn.Do("GET", "resque:status:"+uuid)
	if err != nil || jobStatus == nil {
		return status{}, err
	}
	serializedStatus := status{}
	json.Unmarshal(jobStatus.([]byte), &serializedStatus)
	return serializedStatus, nil
}

type status struct {
	Time    int64                  `json:"time"`
	Status  string                 `json:"status"`
	Name    string                 `json:"name"`
	UUID    string                 `json:"uuid"`
	Options map[string]interface{} `json:"options"`
}
