package module

import (
	"github.com/name5566/leaf/chanrpc"
	"github.com/name5566/leaf/console"
	"github.com/name5566/leaf/go"
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/timer"
	clog "casino_common/common/log"
	"time"
)

type Skeleton struct {
	GoLen              int
	TimerDispatcherLen int
	ChanRPCServer      *chanrpc.Server
	g                  *g.Go
	dispatcher         *timer.Dispatcher
	server             *chanrpc.Server
	commandServer      *chanrpc.Server
}

func (s *Skeleton) Init() {
	log.Debug("package module SkeLeton.Init()")
	if s.GoLen <= 0 {
		s.GoLen = 0
	}
	if s.TimerDispatcherLen <= 0 {
		s.TimerDispatcherLen = 0
	}

	s.g = g.New(s.GoLen)
	s.dispatcher = timer.NewDispatcher(s.TimerDispatcherLen)
	s.server = s.ChanRPCServer

	if s.server == nil {
		s.server = chanrpc.NewServer(0)
	}
	s.commandServer = chanrpc.NewServer(0)
}

func (s *Skeleton) Run(closeSig chan bool) {
	for {

		select {
		case <-closeSig:
			log.Debug("Skeleton.Run(): 取出<-closeSig")  //TODO: 临时调试log
			s.commandServer.Close()
			s.server.Close()
			s.g.Close()
			return
		case ci := <-s.server.ChanCall:
			log.Debug("Skeleton.Run(): 取出<-s.server.ChanCall去Exec")  //TODO: 临时调试log
			err := s.server.Exec(ci)
			if err != nil {
				log.Error("%v", err)
				clog.E("%v",err)
			}
		case ci := <-s.commandServer.ChanCall:
			log.Debug("Skeleton.Run(): 取出<-s.commandServer.ChanCall")  //TODO: 临时调试log
			err := s.commandServer.Exec(ci)
			if err != nil {
				log.Error("%v", err)
			}
		case cb := <-s.g.ChanCb:
			log.Debug("Skeleton.Run(): 取出s.g.ChanCb")  //TODO: 临时调试log
			s.g.Cb(cb)
		case t := <-s.dispatcher.ChanTimer:
			log.Debug("Skeleton.Run(): 取出s.dispatcher.ChanTimer执行")  //TODO: 临时调试log
			t.Cb()
		}
	}
}

func (s *Skeleton) AfterFunc(d time.Duration, cb func()) *timer.Timer {
	if s.TimerDispatcherLen == 0 {
		panic("invalid TimerDispatcherLen")
	}

	return s.dispatcher.AfterFunc(d, cb)
}

func (s *Skeleton) CronFunc(cronExpr *timer.CronExpr, cb func()) *timer.Cron {
	if s.TimerDispatcherLen == 0 {
		panic("invalid TimerDispatcherLen")
	}

	return s.dispatcher.CronFunc(cronExpr, cb)
}

func (s *Skeleton) Go(f func(), cb func()) {
	if s.GoLen == 0 {
		panic("invalid GoLen")
	}

	s.g.Go(f, cb)
}

func (s *Skeleton) NewLinearContext() *g.LinearContext {
	if s.GoLen == 0 {
		panic("invalid GoLen")
	}

	return s.g.NewLinearContext()
}

func (s *Skeleton) RegisterChanRPC(id interface{}, f interface{}) {
	if s.ChanRPCServer == nil {
		panic("invalid ChanRPCServer")
	}

	s.server.Register(id, f)
}

func (s *Skeleton) RegisterCommand(name string, help string, f interface{}) {
	console.Register(name, help, f, s.commandServer)
}
