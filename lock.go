package resque_status

import (
	"fmt"
	"sort"
	"strings"
)

// Lock - struct to define locks for jobs
type Lock struct {
	LockKeyPrefix string
	KeyParamNames []string
}

// NewLock - craetes new lock from job configurations
func NewLock(lockKeyPrefix string, keyParamNames []string) *Lock {
	return &Lock{lockKeyPrefix, keyParamNames}
}

// Key - generates lock key
func (l *Lock) Key(params map[string]interface{}) string {
	lockKeyParts := []string{}
	sort.Strings(l.KeyParamNames)
	for _, key := range l.KeyParamNames {
		lockKeyParts = append(lockKeyParts, fmt.Sprintf("%s=%v", key, params[key]))
	}
	return fmt.Sprintf("resque:lock:%s-%s", l.LockKeyPrefix, strings.Join(lockKeyParts, "|"))

}
