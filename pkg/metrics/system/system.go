package system

import (
	"bytes"
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/procfs"
	"github.com/sankalp-r/metrics-agent/pkg/metrics"
	"github.com/shirou/gopsutil/v3/process"
	"go.uber.org/zap"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

const (
	uTime     = "algod_utime"
	sTime     = "algod_stime"
	startTime = "algod_starttime"
	inOctets  = "algod_inoctets"
	outOctets = "algod_outoctets"
)

// systemSource
type systemSource struct {
	processName    string
	pID            int32
	counterMetrics map[string]prometheus.Counter
	gaugeMetrics   map[string]prometheus.Gauge
	lock           sync.RWMutex
	metricRegistry *prometheus.Registry
}

// NewSystemSource return new system-source
func NewSystemSource(processName string) metrics.Source {
	s := &systemSource{
		processName:    processName,
		pID:            -1,
		counterMetrics: make(map[string]prometheus.Counter),
		gaugeMetrics:   make(map[string]prometheus.Gauge),
		lock:           sync.RWMutex{},
		metricRegistry: prometheus.NewRegistry(),
	}
	s.initMetrics()
	return s
}

// initMetrics creates and initializes system-metrics
func (h *systemSource) initMetrics() {

	h.gaugeMetrics[uTime] = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: uTime,
		Help: "Algod process utime",
	})
	h.metricRegistry.Register(h.gaugeMetrics[uTime])

	h.gaugeMetrics[sTime] = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: sTime,
		Help: "Algod process stime",
	})
	h.metricRegistry.Register(h.gaugeMetrics[sTime])

	h.gaugeMetrics[startTime] = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: startTime,
		Help: "Algod process start time",
	})
	h.metricRegistry.Register(h.gaugeMetrics[startTime])

	h.gaugeMetrics[inOctets] = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: inOctets,
		Help: "Algod inoctets (bytes)",
	})
	h.metricRegistry.Register(h.gaugeMetrics[inOctets])

	h.gaugeMetrics[outOctets] = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: outOctets,
		Help: "Algod outoctets (bytes)",
	})
	h.metricRegistry.Register(h.gaugeMetrics[outOctets])
}

// Collect system-metrics
func (h *systemSource) Collect() ([]byte, error) {
	if err := h.syncPID(); err != nil {
		return nil, err
	}
	var pid int
	h.lock.RLock()
	pid = int(h.pID)
	h.lock.RUnlock()

	stat, err := findProcessStat(pid)
	if err != nil {
		return nil, err
	}
	h.gaugeMetrics[uTime].Set(stat.Utime)
	h.gaugeMetrics[sTime].Set(stat.Stime)
	h.gaugeMetrics[startTime].Set(stat.StartTime)
	h.gaugeMetrics[inOctets].Set(stat.InOctets)
	h.gaugeMetrics[outOctets].Set(stat.OutOctets)

	var buffer bytes.Buffer
	encoder := expfmt.NewEncoder(&buffer, expfmt.FmtText)

	mf, err := h.metricRegistry.Gather()
	if err != nil {
		return nil, err
	}

	for _, metric := range mf {
		err = encoder.Encode(metric)
		if err != nil {
			zap.L().Error(err.Error())
		}
	}

	return buffer.Bytes(), nil
}

// syncPID of algod-node
func (h *systemSource) syncPID() error {
	h.lock.RLock()
	p := &process.Process{Pid: h.pID}
	h.lock.RUnlock()
	if exist, _ := p.IsRunning(); !exist {
		pid, err := findProcessID(h.processName)
		if err != nil {
			return errors.New("unable to find pid")
		}
		i, err := strconv.ParseInt(pid, 10, 64)
		if err != nil {
			return err
		}
		h.lock.Lock()
		h.pID = int32(i)
		h.lock.Unlock()
	}
	return nil
}

// findProcessID of a process
var findProcessID = func(processName string) (string, error) {
	var cmd *exec.Cmd
	cmd = exec.Command("pgrep", processName)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// findProcessStat of the specified process
var findProcessStat = func(pid int) (*Stat, error) {
	p, err := procfs.NewProc(pid)
	if err != nil {
		return nil, err
	}

	stat, err := p.Stat()
	if err != nil {
		return nil, err
	}
	pStat := Stat{
		Utime:     float64(stat.UTime),
		Stime:     float64(stat.STime),
		StartTime: float64(stat.Starttime),
	}
	netStat, err := p.Netstat()
	if err != nil {
		return nil, err
	}
	if netStat.InOctets != nil {
		pStat.InOctets = *netStat.InOctets
	}
	if netStat.OutOctets != nil {
		pStat.OutOctets = *netStat.OutOctets
	}
	return &pStat, nil
}
