// This file was generated by counterfeiter
package pipelinesfakes

import (
	"sync"

	"github.com/concourse/atc/dbng"
	"github.com/concourse/atc/pipelines"
	"github.com/concourse/atc/radar"
	"github.com/concourse/atc/scheduler"
)

type FakeRadarSchedulerFactory struct {
	BuildScanRunnerFactoryStub        func(dbPipeline dbng.Pipeline, externalURL string) radar.ScanRunnerFactory
	buildScanRunnerFactoryMutex       sync.RWMutex
	buildScanRunnerFactoryArgsForCall []struct {
		dbPipeline  dbng.Pipeline
		externalURL string
	}
	buildScanRunnerFactoryReturns struct {
		result1 radar.ScanRunnerFactory
	}
	buildScanRunnerFactoryReturnsOnCall map[int]struct {
		result1 radar.ScanRunnerFactory
	}
	BuildSchedulerStub        func(pipeline dbng.Pipeline, externalURL string) scheduler.BuildScheduler
	buildSchedulerMutex       sync.RWMutex
	buildSchedulerArgsForCall []struct {
		pipeline    dbng.Pipeline
		externalURL string
	}
	buildSchedulerReturns struct {
		result1 scheduler.BuildScheduler
	}
	buildSchedulerReturnsOnCall map[int]struct {
		result1 scheduler.BuildScheduler
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeRadarSchedulerFactory) BuildScanRunnerFactory(dbPipeline dbng.Pipeline, externalURL string) radar.ScanRunnerFactory {
	fake.buildScanRunnerFactoryMutex.Lock()
	ret, specificReturn := fake.buildScanRunnerFactoryReturnsOnCall[len(fake.buildScanRunnerFactoryArgsForCall)]
	fake.buildScanRunnerFactoryArgsForCall = append(fake.buildScanRunnerFactoryArgsForCall, struct {
		dbPipeline  dbng.Pipeline
		externalURL string
	}{dbPipeline, externalURL})
	fake.recordInvocation("BuildScanRunnerFactory", []interface{}{dbPipeline, externalURL})
	fake.buildScanRunnerFactoryMutex.Unlock()
	if fake.BuildScanRunnerFactoryStub != nil {
		return fake.BuildScanRunnerFactoryStub(dbPipeline, externalURL)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.buildScanRunnerFactoryReturns.result1
}

func (fake *FakeRadarSchedulerFactory) BuildScanRunnerFactoryCallCount() int {
	fake.buildScanRunnerFactoryMutex.RLock()
	defer fake.buildScanRunnerFactoryMutex.RUnlock()
	return len(fake.buildScanRunnerFactoryArgsForCall)
}

func (fake *FakeRadarSchedulerFactory) BuildScanRunnerFactoryArgsForCall(i int) (dbng.Pipeline, string) {
	fake.buildScanRunnerFactoryMutex.RLock()
	defer fake.buildScanRunnerFactoryMutex.RUnlock()
	return fake.buildScanRunnerFactoryArgsForCall[i].dbPipeline, fake.buildScanRunnerFactoryArgsForCall[i].externalURL
}

func (fake *FakeRadarSchedulerFactory) BuildScanRunnerFactoryReturns(result1 radar.ScanRunnerFactory) {
	fake.BuildScanRunnerFactoryStub = nil
	fake.buildScanRunnerFactoryReturns = struct {
		result1 radar.ScanRunnerFactory
	}{result1}
}

func (fake *FakeRadarSchedulerFactory) BuildScanRunnerFactoryReturnsOnCall(i int, result1 radar.ScanRunnerFactory) {
	fake.BuildScanRunnerFactoryStub = nil
	if fake.buildScanRunnerFactoryReturnsOnCall == nil {
		fake.buildScanRunnerFactoryReturnsOnCall = make(map[int]struct {
			result1 radar.ScanRunnerFactory
		})
	}
	fake.buildScanRunnerFactoryReturnsOnCall[i] = struct {
		result1 radar.ScanRunnerFactory
	}{result1}
}

func (fake *FakeRadarSchedulerFactory) BuildScheduler(pipeline dbng.Pipeline, externalURL string) scheduler.BuildScheduler {
	fake.buildSchedulerMutex.Lock()
	ret, specificReturn := fake.buildSchedulerReturnsOnCall[len(fake.buildSchedulerArgsForCall)]
	fake.buildSchedulerArgsForCall = append(fake.buildSchedulerArgsForCall, struct {
		pipeline    dbng.Pipeline
		externalURL string
	}{pipeline, externalURL})
	fake.recordInvocation("BuildScheduler", []interface{}{pipeline, externalURL})
	fake.buildSchedulerMutex.Unlock()
	if fake.BuildSchedulerStub != nil {
		return fake.BuildSchedulerStub(pipeline, externalURL)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.buildSchedulerReturns.result1
}

func (fake *FakeRadarSchedulerFactory) BuildSchedulerCallCount() int {
	fake.buildSchedulerMutex.RLock()
	defer fake.buildSchedulerMutex.RUnlock()
	return len(fake.buildSchedulerArgsForCall)
}

func (fake *FakeRadarSchedulerFactory) BuildSchedulerArgsForCall(i int) (dbng.Pipeline, string) {
	fake.buildSchedulerMutex.RLock()
	defer fake.buildSchedulerMutex.RUnlock()
	return fake.buildSchedulerArgsForCall[i].pipeline, fake.buildSchedulerArgsForCall[i].externalURL
}

func (fake *FakeRadarSchedulerFactory) BuildSchedulerReturns(result1 scheduler.BuildScheduler) {
	fake.BuildSchedulerStub = nil
	fake.buildSchedulerReturns = struct {
		result1 scheduler.BuildScheduler
	}{result1}
}

func (fake *FakeRadarSchedulerFactory) BuildSchedulerReturnsOnCall(i int, result1 scheduler.BuildScheduler) {
	fake.BuildSchedulerStub = nil
	if fake.buildSchedulerReturnsOnCall == nil {
		fake.buildSchedulerReturnsOnCall = make(map[int]struct {
			result1 scheduler.BuildScheduler
		})
	}
	fake.buildSchedulerReturnsOnCall[i] = struct {
		result1 scheduler.BuildScheduler
	}{result1}
}

func (fake *FakeRadarSchedulerFactory) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.buildScanRunnerFactoryMutex.RLock()
	defer fake.buildScanRunnerFactoryMutex.RUnlock()
	fake.buildSchedulerMutex.RLock()
	defer fake.buildSchedulerMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeRadarSchedulerFactory) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ pipelines.RadarSchedulerFactory = new(FakeRadarSchedulerFactory)
