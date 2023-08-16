// Code generated by mockery v2.32.4. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	slack "github.com/slack-go/slack"

	slackbot "github.com/clambin/go-common/slackbot"
)

// SlackBot is an autogenerated mock type for the SlackBot type
type SlackBot struct {
	mock.Mock
}

type SlackBot_Expecter struct {
	mock *mock.Mock
}

func (_m *SlackBot) EXPECT() *SlackBot_Expecter {
	return &SlackBot_Expecter{mock: &_m.Mock}
}

// Register provides a mock function with given fields: name, command
func (_m *SlackBot) Register(name string, command slackbot.CommandFunc) {
	_m.Called(name, command)
}

// SlackBot_Register_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Register'
type SlackBot_Register_Call struct {
	*mock.Call
}

// Register is a helper method to define mock.On call
//   - name string
//   - command slackbot.CommandFunc
func (_e *SlackBot_Expecter) Register(name interface{}, command interface{}) *SlackBot_Register_Call {
	return &SlackBot_Register_Call{Call: _e.mock.On("Register", name, command)}
}

func (_c *SlackBot_Register_Call) Run(run func(name string, command slackbot.CommandFunc)) *SlackBot_Register_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(slackbot.CommandFunc))
	})
	return _c
}

func (_c *SlackBot_Register_Call) Return() *SlackBot_Register_Call {
	_c.Call.Return()
	return _c
}

func (_c *SlackBot_Register_Call) RunAndReturn(run func(string, slackbot.CommandFunc)) *SlackBot_Register_Call {
	_c.Call.Return(run)
	return _c
}

// Run provides a mock function with given fields: ctx
func (_m *SlackBot) Run(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SlackBot_Run_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Run'
type SlackBot_Run_Call struct {
	*mock.Call
}

// Run is a helper method to define mock.On call
//   - ctx context.Context
func (_e *SlackBot_Expecter) Run(ctx interface{}) *SlackBot_Run_Call {
	return &SlackBot_Run_Call{Call: _e.mock.On("Run", ctx)}
}

func (_c *SlackBot_Run_Call) Run(run func(ctx context.Context)) *SlackBot_Run_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *SlackBot_Run_Call) Return(_a0 error) *SlackBot_Run_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SlackBot_Run_Call) RunAndReturn(run func(context.Context) error) *SlackBot_Run_Call {
	_c.Call.Return(run)
	return _c
}

// Send provides a mock function with given fields: channel, attachments
func (_m *SlackBot) Send(channel string, attachments []slack.Attachment) error {
	ret := _m.Called(channel, attachments)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []slack.Attachment) error); ok {
		r0 = rf(channel, attachments)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SlackBot_Send_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Send'
type SlackBot_Send_Call struct {
	*mock.Call
}

// Send is a helper method to define mock.On call
//   - channel string
//   - attachments []slack.Attachment
func (_e *SlackBot_Expecter) Send(channel interface{}, attachments interface{}) *SlackBot_Send_Call {
	return &SlackBot_Send_Call{Call: _e.mock.On("Send", channel, attachments)}
}

func (_c *SlackBot_Send_Call) Run(run func(channel string, attachments []slack.Attachment)) *SlackBot_Send_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].([]slack.Attachment))
	})
	return _c
}

func (_c *SlackBot_Send_Call) Return(_a0 error) *SlackBot_Send_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *SlackBot_Send_Call) RunAndReturn(run func(string, []slack.Attachment) error) *SlackBot_Send_Call {
	_c.Call.Return(run)
	return _c
}

// NewSlackBot creates a new instance of SlackBot. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSlackBot(t interface {
	mock.TestingT
	Cleanup(func())
}) *SlackBot {
	mock := &SlackBot{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
