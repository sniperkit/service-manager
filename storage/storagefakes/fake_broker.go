// Code generated by counterfeiter. DO NOT EDIT.
package storagefakes

import (
	"sync"

	"github.com/Peripli/service-manager/storage"
	"github.com/Peripli/service-manager/types"
)

type FakeBroker struct {
	CreateStub        func(broker *types.Broker) error
	createMutex       sync.RWMutex
	createArgsForCall []struct {
		broker *types.Broker
	}
	createReturns struct {
		result1 error
	}
	createReturnsOnCall map[int]struct {
		result1 error
	}
	GetStub        func(id string) (*types.Broker, error)
	getMutex       sync.RWMutex
	getArgsForCall []struct {
		id string
	}
	getReturns struct {
		result1 *types.Broker
		result2 error
	}
	getReturnsOnCall map[int]struct {
		result1 *types.Broker
		result2 error
	}
	GetAllStub        func() ([]types.Broker, error)
	getAllMutex       sync.RWMutex
	getAllArgsForCall []struct{}
	getAllReturns     struct {
		result1 []types.Broker
		result2 error
	}
	getAllReturnsOnCall map[int]struct {
		result1 []types.Broker
		result2 error
	}
	DeleteStub        func(id string) error
	deleteMutex       sync.RWMutex
	deleteArgsForCall []struct {
		id string
	}
	deleteReturns struct {
		result1 error
	}
	deleteReturnsOnCall map[int]struct {
		result1 error
	}
	UpdateStub        func(broker *types.Broker) error
	updateMutex       sync.RWMutex
	updateArgsForCall []struct {
		broker *types.Broker
	}
	updateReturns struct {
		result1 error
	}
	updateReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeBroker) Create(broker *types.Broker) error {
	fake.createMutex.Lock()
	ret, specificReturn := fake.createReturnsOnCall[len(fake.createArgsForCall)]
	fake.createArgsForCall = append(fake.createArgsForCall, struct {
		broker *types.Broker
	}{broker})
	fake.recordInvocation("Create", []interface{}{broker})
	fake.createMutex.Unlock()
	if fake.CreateStub != nil {
		return fake.CreateStub(broker)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.createReturns.result1
}

func (fake *FakeBroker) CreateCallCount() int {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return len(fake.createArgsForCall)
}

func (fake *FakeBroker) CreateArgsForCall(i int) *types.Broker {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return fake.createArgsForCall[i].broker
}

func (fake *FakeBroker) CreateReturns(result1 error) {
	fake.CreateStub = nil
	fake.createReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeBroker) CreateReturnsOnCall(i int, result1 error) {
	fake.CreateStub = nil
	if fake.createReturnsOnCall == nil {
		fake.createReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.createReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeBroker) Get(id string) (*types.Broker, error) {
	fake.getMutex.Lock()
	ret, specificReturn := fake.getReturnsOnCall[len(fake.getArgsForCall)]
	fake.getArgsForCall = append(fake.getArgsForCall, struct {
		id string
	}{id})
	fake.recordInvocation("Get", []interface{}{id})
	fake.getMutex.Unlock()
	if fake.GetStub != nil {
		return fake.GetStub(id)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getReturns.result1, fake.getReturns.result2
}

func (fake *FakeBroker) GetCallCount() int {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return len(fake.getArgsForCall)
}

func (fake *FakeBroker) GetArgsForCall(i int) string {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return fake.getArgsForCall[i].id
}

func (fake *FakeBroker) GetReturns(result1 *types.Broker, result2 error) {
	fake.GetStub = nil
	fake.getReturns = struct {
		result1 *types.Broker
		result2 error
	}{result1, result2}
}

func (fake *FakeBroker) GetReturnsOnCall(i int, result1 *types.Broker, result2 error) {
	fake.GetStub = nil
	if fake.getReturnsOnCall == nil {
		fake.getReturnsOnCall = make(map[int]struct {
			result1 *types.Broker
			result2 error
		})
	}
	fake.getReturnsOnCall[i] = struct {
		result1 *types.Broker
		result2 error
	}{result1, result2}
}

func (fake *FakeBroker) GetAll() ([]types.Broker, error) {
	fake.getAllMutex.Lock()
	ret, specificReturn := fake.getAllReturnsOnCall[len(fake.getAllArgsForCall)]
	fake.getAllArgsForCall = append(fake.getAllArgsForCall, struct{}{})
	fake.recordInvocation("GetAll", []interface{}{})
	fake.getAllMutex.Unlock()
	if fake.GetAllStub != nil {
		return fake.GetAllStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getAllReturns.result1, fake.getAllReturns.result2
}

func (fake *FakeBroker) GetAllCallCount() int {
	fake.getAllMutex.RLock()
	defer fake.getAllMutex.RUnlock()
	return len(fake.getAllArgsForCall)
}

func (fake *FakeBroker) GetAllReturns(result1 []types.Broker, result2 error) {
	fake.GetAllStub = nil
	fake.getAllReturns = struct {
		result1 []types.Broker
		result2 error
	}{result1, result2}
}

func (fake *FakeBroker) GetAllReturnsOnCall(i int, result1 []types.Broker, result2 error) {
	fake.GetAllStub = nil
	if fake.getAllReturnsOnCall == nil {
		fake.getAllReturnsOnCall = make(map[int]struct {
			result1 []types.Broker
			result2 error
		})
	}
	fake.getAllReturnsOnCall[i] = struct {
		result1 []types.Broker
		result2 error
	}{result1, result2}
}

func (fake *FakeBroker) Delete(id string) error {
	fake.deleteMutex.Lock()
	ret, specificReturn := fake.deleteReturnsOnCall[len(fake.deleteArgsForCall)]
	fake.deleteArgsForCall = append(fake.deleteArgsForCall, struct {
		id string
	}{id})
	fake.recordInvocation("Delete", []interface{}{id})
	fake.deleteMutex.Unlock()
	if fake.DeleteStub != nil {
		return fake.DeleteStub(id)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.deleteReturns.result1
}

func (fake *FakeBroker) DeleteCallCount() int {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return len(fake.deleteArgsForCall)
}

func (fake *FakeBroker) DeleteArgsForCall(i int) string {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return fake.deleteArgsForCall[i].id
}

func (fake *FakeBroker) DeleteReturns(result1 error) {
	fake.DeleteStub = nil
	fake.deleteReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeBroker) DeleteReturnsOnCall(i int, result1 error) {
	fake.DeleteStub = nil
	if fake.deleteReturnsOnCall == nil {
		fake.deleteReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeBroker) Update(broker *types.Broker) error {
	fake.updateMutex.Lock()
	ret, specificReturn := fake.updateReturnsOnCall[len(fake.updateArgsForCall)]
	fake.updateArgsForCall = append(fake.updateArgsForCall, struct {
		broker *types.Broker
	}{broker})
	fake.recordInvocation("Update", []interface{}{broker})
	fake.updateMutex.Unlock()
	if fake.UpdateStub != nil {
		return fake.UpdateStub(broker)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.updateReturns.result1
}

func (fake *FakeBroker) UpdateCallCount() int {
	fake.updateMutex.RLock()
	defer fake.updateMutex.RUnlock()
	return len(fake.updateArgsForCall)
}

func (fake *FakeBroker) UpdateArgsForCall(i int) *types.Broker {
	fake.updateMutex.RLock()
	defer fake.updateMutex.RUnlock()
	return fake.updateArgsForCall[i].broker
}

func (fake *FakeBroker) UpdateReturns(result1 error) {
	fake.UpdateStub = nil
	fake.updateReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeBroker) UpdateReturnsOnCall(i int, result1 error) {
	fake.UpdateStub = nil
	if fake.updateReturnsOnCall == nil {
		fake.updateReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.updateReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeBroker) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	fake.getAllMutex.RLock()
	defer fake.getAllMutex.RUnlock()
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	fake.updateMutex.RLock()
	defer fake.updateMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeBroker) recordInvocation(key string, args []interface{}) {
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

var _ storage.Broker = new(FakeBroker)
