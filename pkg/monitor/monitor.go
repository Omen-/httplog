package monitor

import (
	"log"
	"time"

	"github.com/omen-/httplog/pkg/commonformat"
)

const (
	traficReportPeriod = 10 * time.Second
	alertMonitorPeriod = 2 * time.Minute

	traficReportBufferSize = 4096
	traficAlertBufferSize  = 4096
)

type Monitor struct {
	treshold      int64
	logReader     commonformat.Reader
	traficReports chan TraficReport
	traficAlerts  chan Alert
}

func New(treshold int64, logReader commonformat.Reader) *Monitor {
	monitor := Monitor{
		treshold:      treshold,
		logReader:     logReader,
		traficReports: make(chan TraficReport, traficReportBufferSize),
		traficAlerts:  make(chan Alert, traficAlertBufferSize),
	}

	go monitor.monitorLogs()

	return &monitor
}

func (monitor *Monitor) TraficReports() chan TraficReport {
	return monitor.traficReports
}

func (monitor *Monitor) Alerts() chan Alert {
	return monitor.traficAlerts
}

func (monitor *Monitor) onTraficReport(traficReport TraficReport) {
	select {
	case monitor.traficReports <- traficReport:
	default:
		log.Println("Trafic report buffer is full, discarding new report")
	}
}

func (monitor *Monitor) onTraficAlert(alert Alert) {
	select {
	case monitor.traficAlerts <- alert:
	default:
		log.Println("Trafic alerts buffer is full, discarding new alert")
	}
}

func (monitor *Monitor) monitorLogs() {
	logsChannel := monitor.logReader.Logs()

	alertMonitor := newAlertMonitor(alertMonitorPeriod, monitor.treshold)

	traficReportTicker := time.NewTicker(traficReportPeriod)
	traficReport := TraficReport{RequestsBySection: make(map[string]int64)}

	for {
		select {
		case logEntry := <-logsChannel:
			traficReport.updateTraficReport(logEntry)
			alert := alertMonitor.addLogEntry(logEntry)
			if alert != nil {
				monitor.onTraficAlert(alert)
			}
		case <-traficReportTicker.C:
			monitor.onTraficReport(traficReport)
			traficReport = TraficReport{RequestsBySection: make(map[string]int64)}
		}
	}
}