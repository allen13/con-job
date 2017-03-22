package scheduler

import (
	"fmt"
	"github.com/allen13/con-job/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SchedulerTestSuite struct {
	suite.Suite
	scheduler         *Scheduler
	mockKeyValueStore *mocks.KeyValueStore
}

func (s *SchedulerTestSuite) SetupTest() {
	s.scheduler = new(Scheduler)
	mockKeyValueStore := new(mocks.KeyValueStore)
	s.scheduler.kvStore = mockKeyValueStore
	s.mockKeyValueStore = mockKeyValueStore
}

func (s *SchedulerTestSuite) TestExample() {
	assert.Equal(s.T(), 6, 6)
}

func (s *SchedulerTestSuite) TestSchedule() {
	s.mockKeyValueStore.
		On("RunForLeader").
		Return(nil)

	s.mockKeyValueStore.
		On("Put", "/nodes/n1", "value").
		Return(nil)
	s.mockKeyValueStore.
		On("Put", "n1", "value").
		Return(nil)

	s.mockKeyValueStore.
		On("Watch").
		Return()

	s.scheduler.Start()
	s.mockKeyValueStore.Put("/nodes/n1", "value")
	fmt.Println(s.mockKeyValueStore.Calls)
}

func TestSchedulerTestSuite(t *testing.T) {
	suite.Run(t, new(SchedulerTestSuite))
}
