// Package concurrency provides support for distributed concurrency capabilities
//
// The goal is to simplify the development of concurrency capabilities.
package concurrency

import (
	"fmt"

	"github.com/QOSGroup/cassini/types"
)

// Mutex distributed lock interface
type Mutex interface {
	// Lock get lock through chain id and sequence
	//
	// If it returned an error, indicates that the call failed.
	// Whether successful or not,
	// the current sequence saved in the distributed lock is returned.
	// Negative sequence(<0) are returned unless there are some unknown exceptions.
	Lock(sequence int64) (int64, error)

	// Update update the sequence saved in the distributed lock.
	Update(sequence int64) error

	// Unlock after successfully acquiring the lock, the lock needs to be unlocked.
	//
	// If it returned an error, indicates that the call failed.
	Unlock(success bool) error

	// Close close the lock
	Close() error
}

// NewMutex new mutex based on configuration.
func NewMutex(name, address string) (Mutex, error) {
	protocol, _ := types.ParseAddrs(address)

	switch protocol {
	case "etcd":
		return nil, fmt.Errorf("etcd not supported, %s / %s", name, address)
	}

	return NewStandaloneMutex(name), nil
}
