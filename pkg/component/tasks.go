// Copyright © 2019 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package component

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"go.thethings.network/lorawan-stack/v3/pkg/errors"
	"go.thethings.network/lorawan-stack/v3/pkg/log"
	"go.thethings.network/lorawan-stack/v3/pkg/random"
)

// TaskFunc is the task function.
type TaskFunc func(context.Context) error

// Execute executes the task function.
func (f TaskFunc) Execute(ctx context.Context, logger log.Interface) (err error) {
	defer func() {
		if p := recover(); p != nil {
			fmt.Fprintln(os.Stderr, p)
			os.Stderr.Write(debug.Stack())
			if pErr, ok := p.(error); ok {
				err = errTaskRecovered.WithCause(pErr)
			} else {
				err = errTaskRecovered.WithAttributes("panic", p)
			}
			logger.WithError(err).Error("Task panicked")
		}
	}()
	return f(ctx)
}

// TaskRestart defines a task's restart policy.
type TaskRestart uint8

const (
	// TaskRestartNever denotes a restart policy that never restarts tasks after success or failure.
	TaskRestartNever TaskRestart = iota
	// TaskRestartAlways denotes a restart policy that always restarts tasks, on success and failure.
	TaskRestartAlways
	// TaskRestartOnFailure denotes a restart policy that restarts tasks on failure.
	TaskRestartOnFailure
)

// TaskBackoffIntervalFunc is a function that decides the backoff interval based on the attempt history.
// invocation is a counter, which starts at 1 and is incremented after each task function invocation.
type TaskBackoffIntervalFunc func(ctx context.Context, executionDuration time.Duration, invocation uint, err error) time.Duration

// TaskBackoffConfig represents task backoff configuration.
type TaskBackoffConfig struct {
	Jitter       float64
	IntervalFunc TaskBackoffIntervalFunc
}

// MakeTaskBackoffIntervalFunc returns a new TaskBackoffIntervalFunc.
func MakeTaskBackoffIntervalFunc(onFailure bool, resetDuration time.Duration, intervals ...time.Duration) TaskBackoffIntervalFunc {
	return func(ctx context.Context, executionDuration time.Duration, invocation uint, err error) time.Duration {
		switch {
		case onFailure && err == nil:
			return 0
		case executionDuration > resetDuration:
			return intervals[0]
		case invocation >= uint(len(intervals)):
			return intervals[len(intervals)-1]
		default:
			return intervals[invocation-1]
		}
	}
}

// Values for DefaultTaskBackoffConfig.
const (
	DefaultTaskBackoffResetDuration = time.Minute
	DefaultTaskBackoffJitter        = 0.1
)

var (
	// DefaultTaskBackoffIntervals are the default task backoff intervals.
	DefaultTaskBackoffIntervals = [...]time.Duration{
		10 * time.Millisecond,
		50 * time.Millisecond,
		100 * time.Millisecond,
		time.Second,
	}
	// DefaultTaskBackoffIntervalFunc is the default TaskBackoffIntervalFunc.
	DefaultTaskBackoffIntervalFunc = MakeTaskBackoffIntervalFunc(false, DefaultTaskBackoffResetDuration, DefaultTaskBackoffIntervals[:]...)
	// DefaultTaskBackoffConfig is the default task backoff config.
	DefaultTaskBackoffConfig = &TaskBackoffConfig{
		Jitter:       DefaultTaskBackoffJitter,
		IntervalFunc: DefaultTaskBackoffIntervalFunc,
	}

	// DialTaskBackoffIntervals are the default task backoff intervals for tasks using Dial.
	DialTaskBackoffIntervals = [...]time.Duration{
		100 * time.Millisecond,
		time.Second,
		10 * time.Second,
	}
	// DialTaskBackoffIntervalFunc is the default TaskBackoffIntervalFunc for tasks using Dial.
	DialTaskBackoffIntervalFunc = MakeTaskBackoffIntervalFunc(false, DefaultTaskBackoffResetDuration, DialTaskBackoffIntervals[:]...)
	// DialTaskBackoffConfig is the default task backoff config for tasks using Dial.
	DialTaskBackoffConfig = &TaskBackoffConfig{
		Jitter:       DefaultTaskBackoffJitter,
		IntervalFunc: DialTaskBackoffIntervalFunc,
	}
)

// TaskConfig represents task configuration.
type TaskConfig struct {
	Context context.Context
	ID      string
	Func    TaskFunc
	Done    func()
	Restart TaskRestart
	Backoff *TaskBackoffConfig
}

// RegisterTask registers a task, optionally with restart policy and backoff, to be started after the component started.
func (c *Component) RegisterTask(conf *TaskConfig) {
	c.taskConfigs = append(c.taskConfigs, conf)
}

// TaskStarter starts tasks with a TaskConfig.
type TaskStarter interface {
	// StartTask starts the specified task function, optionally with restart policy and backoff.
	StartTask(*TaskConfig)
}

// StartTaskFunc is a function that implements the TaskStarter interface.
type StartTaskFunc func(*TaskConfig)

// StartTask implements the TaskStarter interface.
func (f StartTaskFunc) StartTask(conf *TaskConfig) {
	f(conf)
}

var errTaskRecovered = errors.DefineInternal("task_recovered", "task recovered")

// DefaultStartTask is the default TaskStarter.
func DefaultStartTask(conf *TaskConfig) {
	logger := log.FromContext(conf.Context).WithField("task_id", conf.ID)
	go func() {
		defer func() {
			if done := conf.Done; done != nil {
				done()
			}
		}()
		for invocation := uint(1); ; invocation++ {
			if invocation == 0 {
				logger.Warn("Invocation count rollover detected")
				invocation = 1
			}
			logger := logger.WithField("invocation", invocation)
			startTime := time.Now()
			err := conf.Func.Execute(conf.Context, logger)
			executionDuration := time.Since(startTime)
			if err != nil && err != context.Canceled {
				logger.WithError(err).Warn("Task failed")
			}
			switch conf.Restart {
			case TaskRestartNever:
				return
			case TaskRestartAlways:
			case TaskRestartOnFailure:
				if err == nil {
					return
				}
			default:
				panic("Invalid TaskConfig.Restart value")
			}
			select {
			case <-conf.Context.Done():
				return
			default:
			}
			if conf.Backoff == nil {
				continue
			}
			s := conf.Backoff.IntervalFunc(conf.Context, executionDuration, invocation, err)
			if s == 0 {
				continue
			}
			if conf.Backoff.Jitter != 0 {
				s = random.Jitter(s, conf.Backoff.Jitter)
			}
			select {
			case <-conf.Context.Done():
				return
			case <-time.After(s):
			}
		}
	}()
}

// StartTask implements the TaskStarter interface.
func (c *Component) StartTask(conf *TaskConfig) {
	c.taskStarter.StartTask(conf)
}

func (c *Component) startTasks() {
	for _, conf := range c.taskConfigs {
		c.taskStarter.StartTask(conf)
	}
}
