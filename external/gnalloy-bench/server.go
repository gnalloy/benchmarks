package main

import (
	"context"
	"fmt"
	"time"

	"gnalloy.org/gnalloy/bootstrap"
	"gnalloy.org/gnalloy/buffer"
	"gnalloy.org/gnalloy/channel"
	"gnalloy.org/gnalloy/transport"
	"gnalloy.org/transport-tcp"
)

const shutdownTimeout = 5 * time.Second

type echoServer struct {
	addr         string
	server       bootstrap.Server
	boss         *transport.EventLoopGroup
	workers      *transport.EventLoopGroup
	echoExecutor *tcpEchoExecutorGroup
}

func startEchoServer(ctx context.Context, cfg config) (*echoServer, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	boss, workers, err := newGroups(cfg)
	if err != nil {
		return nil, err
	}
	echoExecutor, err := newTCPEchoExecutorForConfig(cfg)
	if err != nil {
		shutdownGroups(boss, workers)
		return nil, err
	}
	server, err := bindEchoServer(ctx, cfg, boss, workers, echoExecutor)
	if err != nil {
		shutdownTCPEchoExecutor(echoExecutor)
		shutdownGroups(boss, workers)
		return nil, err
	}
	return &echoServer{addr: server.Addr(), server: server, boss: boss, workers: workers, echoExecutor: echoExecutor}, nil
}

func newGroups(cfg config) (*transport.EventLoopGroup, *transport.EventLoopGroup, error) {
	pollerConfig := transport.Config{
		Backend:         cfg.Backend,
		MultishotAccept: cfg.IOUringMultishotAccept,
		SQPoll:          cfg.IOUringSQPoll,
	}
	boss, err := transport.NewEventLoopGroup(transport.EventLoopGroupConfig{
		Size:           cfg.Boss,
		PollerConfig:   pollerConfig,
		EventBatchSize: cfg.EventBatchSize,
		CPUAffinity:    cfg.BossCPUSet,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("create boss group: %w", err)
	}
	workers, err := transport.NewEventLoopGroup(transport.EventLoopGroupConfig{
		Size:           cfg.Workers,
		PollerConfig:   pollerConfig,
		EventBatchSize: cfg.EventBatchSize,
		CPUAffinity:    cfg.WorkerCPUSet,
	})
	if err != nil {
		_ = boss.Close()
		return nil, nil, fmt.Errorf("create worker group: %w", err)
	}
	return boss, workers, nil
}

func benchmarkChildOptions(cfg config) []channel.ChannelOptionAssignment {
	return []channel.ChannelOptionAssignment{
		channel.OptionMaxMessagesPerRead.Assignment(cfg.MaxMessagesPerRead),
		channel.OptionFlushStrategy.Assignment(cfg.FlushStrategy),
	}
}

func bindEchoServer(ctx context.Context, cfg config, boss *transport.EventLoopGroup, workers *transport.EventLoopGroup, echoExecutor *tcpEchoExecutorGroup) (bootstrap.Server, error) {
	tcpConfig := tcp.DefaultConfig()
	tcpConfig.ReadBufferSize = cfg.ReadBufferSize
	tcpConfig.ReusePort = cfg.ReusePort
	tcpConfig.IOUringFixedBuffers = cfg.IOUringFixedBuffers
	if cfg.Mmap {
		tcpConfig.AllocatorFactory = tcp.NewMmapAllocatorFactory(buffer.MmapAllocatorConfig{
			BlockSize: cfg.MmapBlockSize,
			Blocks:    cfg.MmapBlocks,
		}, false)
	}
	return bootstrap.NewServerBootstrap().
		Group(boss, workers).
		Transport(tcp.NewTransport(tcpConfig)).
		ChildOption(benchmarkChildOptions(cfg)...).
		ChildInitializer(func(ch channel.Channel) error {
			handler, err := newTCPEchoHandler(cfg, echoExecutor)
			if err != nil {
				return err
			}
			return ch.Pipeline().AddLast("echo", handler)
		}).
		BindContext(ctx, cfg.Addr)
}

func (s *echoServer) stop() {
	if s == nil {
		return
	}
	if s.server != nil {
		_ = s.server.Close()
	}
	shutdownTCPEchoExecutor(s.echoExecutor)
	shutdownGroups(s.boss, s.workers)
}

func shutdownGroups(groups ...*transport.EventLoopGroup) {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	for _, group := range groups {
		_ = group.Shutdown(ctx)
	}
}

func newTCPEchoExecutorForConfig(cfg config) (*tcpEchoExecutorGroup, error) {
	if !tcpEchoModeUsesExecutor(cfg.TCPEchoMode) {
		return nil, nil
	}
	return newTCPEchoExecutorGroup(cfg.TCPEchoExecutorWorkers, cfg.TCPEchoExecutorQueueSize)
}

func shutdownTCPEchoExecutor(group *tcpEchoExecutorGroup) {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	_ = group.shutdown(ctx)
}
