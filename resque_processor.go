package resque_status

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/Shop2market/goworker"
)

type ResqueProcessor struct {
	JobName       string
	LockKeyPrefix string
	KeyParamNames []string
	Handler
}

type Handler func(map[string]interface{}) error

func NewResqueProcessor(jobName, lockKeyPrefix string, keyParamNames []string, handler Handler) *ResqueProcessor {
	return &ResqueProcessor{JobName: jobName, LockKeyPrefix: lockKeyPrefix, KeyParamNames: keyParamNames, Handler: handler}
}

func (rp *ResqueProcessor) Process(queue string, args ...interface{}) error {
	jobUuid := args[0].(string)

	err := rp.updateStatus(jobUuid, "working")
	if err != nil {
		return err
	}
	defer rp.updateStatus(jobUuid, "completed")

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
	lockKey := rp.LockKey(params)
	conn.Do("DEL", lockKey)
	return nil
}

func (rp *ResqueProcessor) LockKey(params map[string]interface{}) string {

	lockKeyParts := []string{}
	sort.Strings(rp.KeyParamNames)
	for _, key := range rp.KeyParamNames {
		lockKeyParts = append(lockKeyParts, fmt.Sprintf("%s=%v", key, params[key]))
	}
	return fmt.Sprintf("resque:lock:%s-%s", rp.LockKeyPrefix, strings.Join(lockKeyParts, "|"))
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
	return nil
}

func (rp *ResqueProcessor) readJobStatus(conn *goworker.RedisConn, uuid string) (status, error) {
	jobStatus, err := conn.Do("GET", "resque:status:"+uuid)
	if err != nil {
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
	Uuid    string                 `json:"uuid"`
	Options map[string]interface{} `json:"options"`
}
