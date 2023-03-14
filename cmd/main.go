package main

import (
	"context"
	"github.com/golang/glog"
	"github.com/google/UUID"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		glog.Errorf("Get in cluster config error: %v", err)
	}
	run := func(ctx context.Context) {
		if err := rootCmd.Execute(); err != nil {
			glog.Errorf("execute err:%v", err)
			os.Exit(1)
		}
	}
	client, err := kubernetes.NewForConfig(config)
	ctx, cancelFunc := context.WithCancel(context.Background())
	id := uuid.New().String()
	defer cancelFunc()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		glog.Infof("Received termination, signaling shutdown")
		cancelFunc()
	}()
	lock := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      "test-lock",
			Namespace: "default",
		},
		Client: client.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: id,
		},
	}
	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   60 * time.Second,
		RenewDeadline:   30 * time.Second,
		RetryPeriod:     4 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				run(ctx)
			},
			OnStoppedLeading: func() {
				glog.Infof("leader lost, id: %v", id)
				os.Exit(0)
			},
			OnNewLeader: func(identity string) {
				// we're notified when new leader elected
				if identity == id {
					// I just got the lock
					return
				}
				glog.Infof("new leader elected: %s", identity)
			},
		},
	})
}
