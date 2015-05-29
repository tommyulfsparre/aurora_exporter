package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/samuel/go-zookeeper/zk"
)

const (
	zkPath         = "/aurora/scheduler"
	zkLeaderPrefix = "singleton_candidate_"
)

type finder interface {
	leaderURL() (string, error)
}

func newFinder(url string) (f finder, err error) {
	if strings.HasPrefix(url, "http://") {
		f = &httpFinder{url: url}
	} else if strings.HasPrefix(url, "zk://") {
		f = newZkFinder(url)
	} else {
		err = errors.New("finder: bad address")
	}

	return f, err
}

type httpFinder struct {
	url string
}

func (f *httpFinder) leaderURL() (string, error) {
	// This will redirect us to the elected Aurora master
	schedulerURL := fmt.Sprintf("%s/scheduler", f.url)
	rr, err := http.NewRequest("GET", schedulerURL, nil)
	if err != nil {
		return "", err
	}

	rresp, err := httpClient.Transport.RoundTrip(rr)
	if err != nil {
		return "", err
	}
	defer rresp.Body.Close()

	masterLoc := rresp.Header.Get("Location")
	if masterLoc == "" {
		glog.V(6).Info("missing Location header in request")
		masterLoc = schedulerURL
	}

	return strings.TrimRight(masterLoc, "/scheduler"), nil
}

func hostsFromURL(urls string) (hosts []string, err error) {
	for _, s := range strings.Split(urls, ",") {
		u, err := url.Parse(s)
		if err != nil {
			return hosts, err
		}

		hosts = append(hosts, u.Host)
	}

	return hosts, err
}

type zkFinder struct {
	sync.RWMutex
	conn     *zk.Conn
	leaderIP string
}

func newZkFinder(url string) *zkFinder {
	zkSrvs, err := hostsFromURL(url)
	if err != nil {
		panic(err)
	}

	conn, events, err := zk.Connect(zkSrvs, 20*time.Second)
	if err != nil {
		panic(err)
	}

	go func() {
		for ev := range events {
			glog.V(6).Infof("conn: %s server: %s", ev.State, ev.Server)
		}
	}()

	f := zkFinder{conn: conn}
	go f.watch()

	return &f
}

func (f *zkFinder) leaderzNode() (string, error) {
	children, stat, err := f.conn.Children(zkPath)
	if err != nil {
		return "", err
	} else if stat == nil {
		return "", errors.New("zkFinder: children returned nil stat")
	}

	var leaderSeq int
	var leader string
	for _, child := range children {
		path := strings.Split(child, zkLeaderPrefix)
		if len(path) > 1 {
			seq, err := strconv.Atoi(path[1])
			if err != nil {
				return "", err
			}

			if leader == "" {
				leader = child
			}

			if seq <= leaderSeq {
				leaderSeq = seq
				leader = child
			}
		}
	}

	if leader == "" {
		return leader, errors.New("zkFinder: zNode not found")
	}

	return fmt.Sprintf("%s/%s", zkPath, leader), nil
}

func (f *zkFinder) leaderURL() (string, error) {
	f.RLock()
	defer f.RUnlock()

	if f.leaderIP == "" {
		return "", errors.New("zkFinder: no leader found via ZooKeeper")
	}

	return fmt.Sprintf("http://%s:8081", f.leaderIP), nil
}

func (f *zkFinder) watch() {
Retry:
	for {
		zNode, err := f.leaderzNode()
		if err != nil {
			glog.Warning(err)
			break
		}

		glog.V(6).Info("leader zNode at: ", zNode)

		data, stat, events, err := f.conn.GetW(zNode)
		if err != nil {
			glog.Warning(err)
			break
		} else if stat == nil {
			glog.Warning("get returned nil stat")
			break
		}

		f.Lock()
		f.leaderIP = string(data)
		f.Unlock()

		for ev := range events {
			if ev.Err != nil {
				glog.Warning("watcher error %v", ev.Err)
			} else if ev.Type == zk.EventNodeDeleted {
				glog.Info("leader zNode deleted")
			} else {
				continue
			}
			break
		}
	}
	time.Sleep(1 * time.Second)
	goto Retry
}
