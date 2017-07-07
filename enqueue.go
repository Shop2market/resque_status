package resque_status

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/Shop2market/goworker"
)

type ResqueEnqueuer struct {
	JobName string
	Queue   string
	Lock
}

// UUIDGenerator for tests injection
type UUIDGenerator func() string

var GenerateUUID UUIDGenerator

func init() {
	GenerateUUID = generateUUID
}

func Enqueue(queue, jobName string, params map[string]interface{}) error {
	return goworker.Enqueue(&goworker.Job{
		Queue: queue,
		Payload: goworker.Payload{
			Class: jobName,
			Args: []interface{}{
				GenerateUUID(),
				params,
			},
		},
	})
}

func generateUUID() string {
	md5Bytes := md5.Sum([]byte(time.Now().String()))
	return hex.EncodeToString(md5Bytes[:])
}

func NewEnqueuer(jobName, queue, lockKeyPrefix string, keyParamNames []string) *ResqueEnqueuer {
	return &ResqueEnqueuer{JobName: jobName, Queue: queue, Lock: Lock{LockKeyPrefix: lockKeyPrefix, KeyParamNames: keyParamNames}}
}

func (re *ResqueEnqueuer) Enqueue(params JobParams) error {
	uuid, err := re.getUUIDByLock(re.Key(params))
	if err != nil {
		return err
	}
	if uuid != "" {
		return nil
	}
	defer re.lock(params)
	err = re.ensureQueue()
	if err != nil {
		return err
	}
	return Enqueue(re.Queue, re.JobName, params)
}

func (re *ResqueEnqueuer) ensureQueue() error {
	conn, err := goworker.GetConn()
	defer goworker.PutConn(conn)
	if err != nil {
		return err
	}
	_, err = conn.Do("SADD", "resque:queues", re.Queue)
	return err
}

func (re *ResqueEnqueuer) getUUIDByLock(lockKey string) (string, error) {
	conn, err := goworker.GetConn()
	defer goworker.PutConn(conn)
	if err != nil {
		return "", err
	}
	val, err := conn.Do("GET", lockKey)
	if err != nil {
		return "", err
	}
	if val == nil {
		return "", nil
	}
	valBytes, ok := val.([]byte)
	if !ok {
		return "", nil
	}
	return string(valBytes[:]), nil
}

func (re *ResqueEnqueuer) lock(params map[string]interface{}) error {
	conn, err := goworker.GetConn()
	defer goworker.PutConn(conn)
	if err != nil {
		return err
	}
	lockKey := re.Key(params)
	_, err = conn.Do("SET", lockKey, GenerateUUID())
	return err
}
