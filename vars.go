package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

var counters = map[string]*prometheus.Desc{
	"cron_job_launch_failures": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "cron_job_launch_failures"),
		"Scheduled job failures total.",
		nil, nil,
	),
	"cron_jobs_loaded": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "cron_jobs_loaded"),
		"Cron scheduler loaded.",
		nil, nil,
	),
	"framework_registered": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "framework_registered"),
		"Framework registered total.",
		nil, nil,
	),
	"gc_executor_tasks_lost": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "gc_executor_tasks_lost"),
		"Lost garbage collection task total.",
		nil, nil,
	),
	"http_500_responses_events": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "http_500_responses_events"),
		"Number of HTTP 500 status responses sent by the scheduler total.",
		nil, nil,
	),
	"job_update_delete_errors": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "job_update_delete_errors"),
		"Failed delete job total.",
		nil, nil,
	),
	"job_update_recovery_errors": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "job_update_recovery_errors"),
		"Failed resume job updates total.",
		nil, nil,
	),
	"job_update_state_change_errors": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "job_update_state_change_errors"),
		"State change errors total.",
		nil, nil,
	),
	"jvm_gc_PS_MarkSweep_collection_count": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "jvm_gc_ps_marksweep_collection_count"),
		"Parallel mark and sweep collection run total.",
		nil, nil,
	),
	"jvm_gc_PS_Scavenge_collection_count": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "jvm_gc_ps_scavenge_collection_count"),
		"Parallel scavenge collector runs total.",
		nil, nil,
	),
	"jvm_gc_collection_count": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "jvm_gc_collection_count"),
		"Garbage collection total.",
		nil, nil,
	),
	"jvm_memory_heap_mb_max": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "jvm_memory_heap_mb_max"),
		"Maximum heap memory.",
		nil, nil,
	),
	"jvm_memory_max_mb": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "jvm_memory_max_mb"),
		"Maximum amount of memory that the Java virtual machine will attempt to use.",
		nil, nil,
	),
	"jvm_memory_mb_total": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "jvm_memory_mb_total"),
		"Total amount of memory in the Java virtual machine.",
		nil, nil,
	),
	"jvm_threads_peak": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "jvm_threads_peak"),
		"Peak thread count.",
		nil, nil,
	),
	"jvm_threads_started": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "jvm_threads_started"),
		"Total threads started.",
		nil, nil,
	),
	"jvm_uptime_secs": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "jvm_uptime_secs"),
		"Number of seconds the JVM process has been running.",
		nil, nil,
	),
	"offer_accept_races": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "offer_accept_races"),
		"Accepted offer no longer exists in offer queue.",
		nil, nil,
	),
	"process_max_fd_count": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "process_max_fd_count"),
		"Max open file descriptors.",
		nil, nil,
	),
	"process_open_fd_count": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "process_open_fd_count"),
		"Open file descriptors in use.",
		nil, nil,
	),
	"schedule_attempts_failed": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "schedule_attempts_failed"),
		"Number of failed attempts to schedule tasks.",
		nil, nil,
	),
	"schedule_attempts_fired": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "schedule_attempts_fired"),
		"Number of attempts to schedule tasks.",
		nil, nil,
	),
	"schedule_attempts_no_match": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "schedule_attempts_no_match"),
		"Number of task which could not be scheduled.",
		nil, nil,
	),
	"schedule_queue_size": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "schedule_queue_size"),
		"Task scheduler queue size.",
		nil, nil,
	),
	"scheduler_backup_failed": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "scheduler_backup_failed"),
		"Number of failed storage backup.",
		nil, nil,
	),
	"scheduler_backup_success": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "scheduler_backup_success"),
		"Number successful storage backup.",
		nil, nil,
	),
	"scheduler_gc_insufficient_offers": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "scheduler_gc_insufficient_offers"),
		"Number off resource offer that was too small for a garbage collection task.",
		nil, nil,
	),
	"scheduler_gc_offers_consumed": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "scheduler_gc_offers_consumed"),
		"Number of resource offers consumed for garbage collection tasks.",
		nil, nil,
	),
	"scheduler_gc_tasks_created": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "scheduler_gc_tasks_created"),
		"Number of garbage collection tasks created.",
		nil, nil,
	),
	"scheduler_log_native_append_nanos_total": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "scheduler_log_native_append_nanos_total"),
		"Timed append operations total.",
		nil, nil,
	),
	"scheduler_log_native_append_events": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "scheduler_log_native_append_events"),
		"Number of append operations total.",
		nil, nil,
	),
	"scheduler_resource_offers": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "scheduler_resource_offers"),
		"Number of resource offers that the scheduler has received.",
		nil, nil,
	),
	"task_kill_retries": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "task_kill_retries"),
		"Number of times the scheduler has retried to kill a Task.",
		nil, nil,
	),
	"timed_out_tasks": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "timed_out_tasks"),
		"Number of times the scheduler has given up waiting to hear back about a task in a transient state.",
		nil, nil,
	),
	"uncaught_exceptions": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "uncaught_exceptions"),
		"Uncaught java exceptions.",
		nil, nil,
	),
}

var gauges = map[string]*prometheus.Desc{
	"jvm_gc_PS_MarkSweep_collection_time_ms": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "jvm_gc_ps_marksweep_collection_time_ms"),
		"Parallel mark and sweep collection time.",
		nil, nil,
	),
	"jvm_gc_PS_Scavenge_collection_time_ms": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "jvm_gc_ps_scavenge_collection_time_ms"),
		"Parallel scavenge collector time.",
		nil, nil,
	),
	"jvm_gc_collection_time_ms": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "jvm_gc_collection_time_ms"),
		"Garbage collection time.",
		nil, nil,
	),
	"jvm_memory_free_mb": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "jvm_memory_free_mb"),
		"Amount of free memory in the Java Virtual Machine.",
		nil, nil,
	),
	"jvm_memory_heap_mb_committed": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "jvm_memory_heap_mb_committed"),
		"Commited heap memory.",
		nil, nil,
	),
	"jvm_memory_heap_mb_used": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "jvm_memory_heap_mb_used"),
		"Current memory usage of the heap.",
		nil, nil,
	),
	"jvm_memory_non_heap_mb_committed": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "jvm_memory_non_heap_mb_committed"),
		"Commited non heap memory used.",
		nil, nil,
	),
	"jvm_memory_non_heap_mb_max": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "jvm_memory_non_heap_mb_max"),
		"Max non heap memory used in MB.",
		nil, nil,
	),
	"jvm_memory_non_heap_mb_used": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "jvm_memory_non_heap_mb_used"),
		"Non heap memory used in MB.",
		nil, nil,
	),
	"jvm_threads_active": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "jvm_threads_active"),
		"Current number of live threads both daemon and non-daemon threads.",
		nil, nil,
	),
	"jvm_threads_daemon": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "jvm_threads_daemon"),
		"Current number of live daemon threads.",
		nil, nil,
	),
	"outstanding_offers": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "outstanding_offers"),
		"Outstanding offers waiting to be returned.",
		nil, nil,
	),
	"process_cpu_cores_utilized": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "process_cpu_cores_utilized"),
		"CPU time used by the process.",
		nil, nil,
	),
	"system_free_physical_memory_mb": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "system_free_physical_memory_mb"),
		"Free physical memory in MB.",
		nil, nil,
	),
	"system_load_avg": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "system_load_avg"),
		"1 minute load average.",
		nil, nil,
	),
	"sla_cluster_platform_uptime_percent": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "sla_cluster_platform_uptime_percent"),
		"Aggregate amount of time a job spends in a non-runnable state.",
		nil, nil,
	),
	"sla_cluster_mtta_ms": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "sla_cluster_mtta_ms"),
		"Median time to assigned.",
		nil, nil,
	),
	"sla_cluster_mttr_ms": prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "sla_cluster_mttr_ms"),
		"Median time to running.",
		nil, nil,
	),
}

var (
	slaLabels = []string{"role", "env", "job"}
	slaRe     = map[*regexp.Regexp]string{
		regexp.MustCompile("sla_(?P<role>.*)/(?P<env>.*)/(?P<job>.*)_mtta_ms"):                 "Median time to assigned.",
		regexp.MustCompile("sla_(?P<role>.*)/(?P<env>.*)/(?P<job>.*)_mttr_ms"):                 "Median time to running.",
		regexp.MustCompile("sla_(?P<role>.*)/(?P<env>.*)/(?P<job>.*)_mtta_nonprod_ms"):         "Median time to assigned nonprod.",
		regexp.MustCompile("sla_(?P<role>.*)/(?P<env>.*)/(?P<job>.*)_mttr_nonprod_ms"):         "Median time to running nonprod.",
		regexp.MustCompile("sla_(?P<role>.*)/(?P<env>.*)/(?P<job>.*)_platform_uptime_percent"): "Aggregate platform uptime.",
	}
)

func slaMetric(name string, value float64) (metric prometheus.Metric) {
	for r, desc := range slaRe {
		match := r.FindStringSubmatch(name)
		if len(match) == 4 {
			role, env, job := match[1], match[2], match[3]

			jobKey := fmt.Sprintf("_%s/%s/%s", role, env, job)
			metricName := strings.Replace(name, jobKey, "", 1)

			metric = prometheus.MustNewConstMetric(
				prometheus.NewDesc(
					prometheus.BuildFQName(namespace, "", metricName),
					desc, slaLabels, nil,
				),
				prometheus.GaugeValue,
				value, role, env, job,
			)
			break
		}
	}
	return metric
}

var (
	stateLabel  = []string{"state"}
	taskStoreRe = regexp.MustCompile("task_store_(?P<state>[A-Z]+)")
)

func taskStoreMetric(name string, value float64) (metric prometheus.Metric) {
	match := taskStoreRe.FindStringSubmatch(name)
	if len(match) == 2 {
		state := match[1]

		metric = prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "", "task_store"),
				"Task store state.",
				stateLabel, nil,
			),
			prometheus.CounterValue,
			value, state,
		)
	}
	return metric
}

var (
	tasksLabels = []string{"state", "role", "env", "job"}
	tasksRe     = regexp.MustCompile("tasks_(?P<state>.*)_(?P<role>.*)/(?P<env>.*)/(?P<job>.*)")
)

func tasksMetric(name string, value float64) (metric prometheus.Metric) {
	match := tasksRe.FindStringSubmatch(name)
	if len(match) == 5 {
		state, role, env, job := match[1], match[2], match[3], match[4]

		metric = prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "", "tasks"),
				"Task state per job.",
				tasksLabels, nil,
			),
			prometheus.CounterValue,
			value, state, role, env, job,
		)
	}
	return metric
}

var (
	rackLabels  = []string{"rack"}
	tasksRackRe = regexp.MustCompile("tasks_lost_rack_(?P<rack>.*)")
)

func tasksRackMetric(name string, value float64) (metric prometheus.Metric) {
	match := tasksRackRe.FindStringSubmatch(name)
	if len(match) == 2 {
		rack := match[1]

		metric = prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "", "tasks_lost_rack"),
				"Task lost per rack total.",
				rackLabels, nil,
			),
			prometheus.CounterValue,
			value, rack,
		)
	}
	return metric
}

var updateRe = regexp.MustCompile("update_transition_(?P<state>[A-Z]+)")

func updateMetric(name string, value float64) (metric prometheus.Metric) {
	match := updateRe.FindStringSubmatch(name)
	if len(match) == 2 {
		state := match[1]

		metric = prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "", "update_transition"),
				"Update transition.",
				stateLabel, nil,
			),
			prometheus.CounterValue,
			value, state,
		)
	}
	return metric
}

var prefixFunc = map[string]func(string, float64) prometheus.Metric{
	"tasks_":             tasksMetric,
	"tasks_lost_rack_":   tasksRackMetric,
	"task_store_":        taskStoreMetric,
	"sla_":               slaMetric,
	"update_transition_": updateMetric,
}

func labelVars(ch chan<- prometheus.Metric, name string, value float64) {
	var metric prometheus.Metric

	for prefix, f := range prefixFunc {
		if strings.HasPrefix(name, prefix) {
			metric = f(name, value)
		}
	}

	if metric != nil {
		ch <- metric
	}
}
