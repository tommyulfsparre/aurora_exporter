package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

func newMetric(subsys string, name string, descr string) *prometheus.Desc {
	fqn := prometheus.BuildFQName(namespace, subsys, name)
	return prometheus.NewDesc(fqn, descr, nil, nil)
}

var counters = map[string]*prometheus.Desc{
	"async_tasks_completed": newMetric(
		"async_tasks", "completed",
		"Number of completed async tasks.",
	),
	"cron_job_launch_failures": newMetric(
		"cron", "job_launch_failures",
		"Scheduled job failures total.",
	),
	"cron_jobs_loaded": newMetric(
		"cron", "jobs_loaded",
		"Cron scheduler loaded.",
	),
	"dropped_update_events": newMetric(
		"", "dropped_update_events", "",
	),
	"framework_registered": newMetric(
		"", "framework_registered",
		"Framework registered total.",
	),
	"gc_executor_tasks_lost": newMetric(
		"", "gc_executor_tasks_lost",
		"Lost garbage collection task total.",
	),
	"http_200_responses_events": newMetric(
		"http_200", "responses_events", "",
	),
	"http_200_responses_nanos_total": newMetric(
		"http_200", "responses_nanos_total", "",
	),
	"http_500_responses_events": newMetric(
		"http_500", "responses_events",
		"Number of HTTP 500 status responses sent by the scheduler total.",
	),
	"job_update_delete_errors": newMetric(
		"job_update", "delete_errors",
		"Failed delete job total.",
	),
	"job_update_recovery_errors": newMetric(
		"job_update", "recovery_errors",
		"Failed resume job updates total.",
	),
	"job_update_state_change_errors": newMetric(
		"job_update", "state_change_errors",
		"State change errors total.",
	),
	"jvm_class_total_loaded_count": newMetric(
		"jvm", "class_total_loaded_count", "",
	),
	"jvm_gc_PS_MarkSweep_collection_count": newMetric(
		"jvm", "gc_ps_marksweep_collection_count",
		"Parallel mark and sweep collection run total.",
	),
	"jvm_gc_PS_Scavenge_collection_count": newMetric(
		"jvm", "gc_ps_scavenge_collection_count",
		"Parallel scavenge collector runs total.",
	),
	"jvm_gc_collection_count": newMetric(
		"jvm", "gc_collection_count",
		"Garbage collection total.",
	),
	"jvm_memory_heap_mb_max": newMetric(
		"jvm", "memory_heap_mb_max",
		"Maximum heap memory.",
	),
	"jvm_memory_max_mb": newMetric(
		"jvm", "memory_max_mb",
		"Maximum amount of memory that the Java virtual machine will attempt to use.",
	),
	"jvm_memory_mb_total": newMetric(
		"jvm", "memory_mb_total",
		"Total amount of memory in the Java virtual machine.",
	),
	"jvm_threads_peak": newMetric(
		"jvm", "threads_peak",
		"Peak thread count.",
	),
	"jvm_threads_started": newMetric(
		"jvm", "threads_started",
		"Total threads started.",
	),
	"jvm_uptime_secs": newMetric(
		"jvm", "uptime_secs",
		"Number of seconds the JVM process has been running.",
	),
	"log_storage_write_lock_wait_events": newMetric(
		"log_storage", "write_lock_wait_events", "",
	),
	"log_storage_write_lock_wait_ns_total": newMetric(
		"log_storage", "write_lock_wait_ns_total", "",
	),
	"offer_accept_races": newMetric(
		"", "offer_accept_races",
		"Accepted offer no longer exists in offer queue.",
	),
	"preemptor_missing_attributes": newMetric(
		"", "preemptor_missing_attributes", "",
	),
	"preemptor_slot_search_attempts_for_non_prod": newMetric(
		"preemptor_slot", "search_attempts_for_non_prod", "",
	),
	"preemptor_slot_search_attempts_for_prod": newMetric(
		"preemptor_slot", "search_attempts_for_prod", "",
	),
	"preemptor_slot_search_failed_for_non_prod": newMetric(
		"preemptor_slot", "search_failed_for_non_prod", "",
	),
	"preemptor_slot_search_failed_for_prod": newMetric(
		"preemptor_slot", "search_failed_for_prod", "",
	),
	"preemptor_slot_search_successful_for_non_prod": newMetric(
		"preemptor_slot", "search_successful_for_non_prod", "",
	),
	"preemptor_slot_search_successful_for_prod": newMetric(
		"preemptor_slot", "search_successful_for_prod", "",
	),
	"preemptor_slot_validation_failed": newMetric(
		"preemptor_slot", "validation_failed", "",
	),
	"preemptor_slot_validation_successful": newMetric(
		"preemptor_slot", "validation_successful", "",
	),
	"preemptor_task_processor_runs": newMetric(
		"preemptor", "task_processor_runs", "",
	),
	"preemptor_tasks_preempted_non_prod": newMetric(
		"preemptor_tasks", "preempted_non_prod", "",
	),
	"preemptor_tasks_preempted_prod": newMetric(
		"preemptor_tasks", "preempted_prod", "",
	),
	"process_max_fd_count": newMetric(
		"process", "max_fd_count",
		"Max open file descriptors.",
	),
	"process_open_fd_count": newMetric(
		"process", "open_fd_count",
		"Open file descriptors in use.",
	),
	"pubsub_executor_queue_size": newMetric(
		"", "pubsub_executor_queue_size", "",
	),
	"schedule_attempts_failed": newMetric(
		"schedule", "attempts_failed",
		"Number of failed attempts to schedule tasks.",
	),
	"schedule_attempts_fired": newMetric(
		"schedule", "attempts_fired",
		"Number of attempts to schedule tasks.",
	),
	"schedule_attempts_no_match": newMetric(
		"schedule", "attempts_no_match",
		"Number of task which could not be scheduled.",
	),
	"schedule_queue_size": newMetric(
		"schedule", "queue_size",
		"Task scheduler queue size.",
	),
	"scheduled_task_penalty_events": newMetric(
		"scheduled_task", "penalty_events", "",
	),
	"scheduled_task_penalty_ms_total": newMetric(
		"scheduled_task", "penalty_ms_total", "",
	),
	"scheduler_backup_failed": newMetric(
		"scheduler", "backup_failed",
		"Number of failed storage backup.",
	),
	"scheduler_backup_success": newMetric(
		"scheduler", "backup_success",
		"Number successful storage backup.",
	),
	"scheduler_driver_kill_failures": newMetric(
		"scheduler", "driver_kill_failures", "",
	),
	"scheduler_gc_insufficient_offers": newMetric(
		"scheduler_gc", "insufficient_offers",
		"Number off resource offer that was too small for a garbage collection task.",
	),
	"scheduler_gc_offers_consumed": newMetric(
		"scheduler_gc", "offers_consumed",
		"Number of resource offers consumed for garbage collection tasks.",
	),
	"scheduler_gc_tasks_created": newMetric(
		"scheduler_gc", "tasks_created",
		"Number of garbage collection tasks created.",
	),
	"scheduler_log_bad_frames_read": newMetric(
		"scheduler_log", "bad_frames_read", "",
	),
	"scheduler_log_bytes_read": newMetric(
		"scheduler_log", "bytes_read", "",
	),
	"scheduler_log_deflated_entries_read": newMetric(
		"scheduler_log", "deflated_entries_read", "",
	),
	"scheduler_log_entries_read": newMetric(
		"scheduler_log", "entries_read", "",
	),
	"scheduler_log_entries_written": newMetric(
		"scheduler_log", "entries_written", "",
	),
	"scheduler_log_native_append_events": newMetric(
		"scheduler_log", "native_append_events",
		"Number of append operations total.",
	),
	"scheduler_log_native_append_failures": newMetric(
		"scheduler_log", "native_append_failures", "Number of append failures total.",
	),
	"scheduler_log_native_append_nanos_total": newMetric(
		"scheduler_log", "native_append_nanos_total",
		"Timed append operations total.",
	),
	"scheduler_log_native_append_nanos_total_per_sec": newMetric(
		"scheduler_log", "native_append_nanos_total_per_sec", "",
	),
	"scheduler_log_native_append_timeouts": newMetric(
		"scheduler_log", "native_append_timeouts", "",
	),
	"scheduler_log_native_native_entries_skipped": newMetric(
		"scheduler_log", "native_native_entries_skipped", "",
	),
	"scheduler_log_native_read_events": newMetric(
		"scheduler_log", "native_read_events", "",
	),
	"scheduler_log_native_read_failures": newMetric(
		"scheduler_log", "native_read_failures", "",
	),
	"scheduler_log_native_read_nanos_total": newMetric(
		"scheduler_log", "native_read_nanos_total", "",
	),
	"scheduler_log_native_read_timeouts": newMetric(
		"scheduler_log", "native_read_timeouts", "",
	),
	"scheduler_log_native_truncate_events": newMetric(
		"scheduler_log", "native_truncate_events", "",
	),
	"scheduler_log_native_truncate_failures": newMetric(
		"scheduler_log", "native_truncate_failures", "",
	),
	"scheduler_log_native_truncate_nanos_total": newMetric(
		"scheduler_log", "native_truncate_nanos_total", "",
	),
	"scheduler_log_native_truncate_timeouts": newMetric(
		"scheduler_log", "native_truncate_timeouts", "",
	),
	"scheduler_log_snapshots": newMetric(
		"scheduler_log", "snapshots", "",
	),
	"scheduler_log_un_snapshotted_transactions": newMetric(
		"scheduler_log", "un_snapshotted_transactions", "",
	),
	"scheduler_resource_offers": newMetric(
		"scheduler", "resource_offers",
		"Number of resource offers that the scheduler has received.",
	),
	"scheduler_thrift_getJobSummary_events": newMetric(
		"scheduler_thrift", "getJobSummary_events", "",
	),
	"scheduler_thrift_getJobSummary_nanos_total": newMetric(
		"scheduler_thrift", "getJobSummary_nanos_total", "",
	),
	"scheduler_thrift_getQuota_events": newMetric(
		"scheduler_thrift", "getQuota_events", "",
	),
	"scheduler_thrift_getQuota_nanos_per_event": newMetric(
		"scheduler_thrift", "getQuota_nanos_per_event", "",
	),
	"scheduler_thrift_getQuota_nanos_total": newMetric(
		"scheduler_thrift", "getQuota_nanos_total", "",
	),
	"task_kill_retries": newMetric(
		"task", "kill_retries",
		"Number of times the scheduler has retried to kill a Task.",
	),
	"task_queries_all": newMetric(
		"task", "queries_all", "",
	),
	"task_queries_by_host": newMetric(
		"task", "queries_by_host", "",
	),
	"task_queries_by_id": newMetric(
		"task", "queries_by_id", "",
	),
	"task_queries_by_job": newMetric(
		"task", "queries_by_job", "",
	),
	"task_throttle_events": newMetric(
		"task", "throttle_events", "",
	),
	"task_throttle_ms_total": newMetric(
		"task", "throttle_ms_total", "",
	),
	"timed_out_tasks": newMetric(
		"", "timed_out_tasks",
		"Number of times the scheduler has given up waiting to hear back about a task in a transient state.",
	),
	"uncaught_exceptions": newMetric(
		"", "uncaught_exceptions",
		"Uncaught java exceptions.",
	),
}

var gauges = map[string]*prometheus.Desc{
	"http_200_responses_events_per_sec": newMetric(
		"http_200", "responses_events_per_sec", "",
	),
	"http_200_responses_nanos_per_event": newMetric(
		"http_200", "responses_nanos_per_event", "",
	),
	"http_200_responses_nanos_total_per_sec": newMetric(
		"http_200", "responses_nanos_total_per_sec", "",
	),
	"jvm_class_loaded_count": newMetric(
		"jvm", "class_loaded_count", "",
	),
	"jvm_class_unloaded_count": newMetric(
		"jvm", "class_unloaded_count", "",
	),
	"jvm_gc_PS_MarkSweep_collection_time_ms": newMetric(
		"jvm", "gc_ps_marksweep_collection_time_ms",
		"Parallel mark and sweep collection time.",
	),
	"jvm_gc_PS_Scavenge_collection_time_ms": newMetric(
		"jvm", "gc_ps_scavenge_collection_time_ms",
		"Parallel scavenge collector time.",
	),
	"jvm_gc_collection_time_ms": newMetric(
		"jvm", "gc_collection_time_ms",
		"Garbage collection time.",
	),
	"jvm_memory_free_mb": newMetric(
		"jvm", "memory_free_mb",
		"Amount of free memory in the Java Virtual Machine.",
	),
	"jvm_memory_heap_mb_committed": newMetric(
		"jvm", "memory_heap_mb_committed",
		"Commited heap memory.",
	),
	"jvm_memory_heap_mb_used": newMetric(
		"jvm", "memory_heap_mb_used",
		"Current memory usage of the heap.",
	),
	"jvm_memory_non_heap_mb_committed": newMetric(
		"jvm", "memory_non_heap_mb_committed",
		"Commited non heap memory used.",
	),
	"jvm_memory_non_heap_mb_max": newMetric(
		"jvm", "memory_non_heap_mb_max",
		"Max non heap memory used in MB.",
	),
	"jvm_memory_non_heap_mb_used": newMetric(
		"jvm", "memory_non_heap_mb_used",
		"Non heap memory used in MB.",
	),
	"jvm_threads_active": newMetric(
		"jvm", "threads_active",
		"Current number of live threads both daemon and non-daemon threads.",
	),
	"jvm_threads_daemon": newMetric(
		"jvm", "threads_daemon",
		"Current number of live daemon threads.",
	),
	"log_storage_write_lock_wait_events_per_sec": newMetric(
		"log_storage", "write_lock_wait_events_per_sec", "",
	),
	"log_storage_write_lock_wait_ns_per_event": newMetric(
		"log_storage", "write_lock_wait_ns_per_event", "",
	),
	"log_storage_write_lock_wait_ns_total_per_sec": newMetric(
		"log_storage", "write_lock_wait_ns_total_per_sec", "",
	),
	"outstanding_offers": newMetric(
		"", "outstanding_offers",
		"Outstanding offers waiting to be returned.",
	),
	"process_cpu_cores_utilized": newMetric(
		"", "process_cpu_cores_utilized",
		"CPU time used by the process.",
	),
	"scheduled_task_penalty_events_per_sec": newMetric(
		"scheduled_task", "penalty_events_per_sec", "",
	),
	"scheduled_task_penalty_ms_per_event": newMetric(
		"scheduled_task", "penalty_ms_per_event", "",
	),
	"scheduled_task_penalty_ms_total_per_sec": newMetric(
		"scheduled_task", "penalty_ms_total_per_sec", "",
	),
	"scheduler_log_bytes_written": newMetric(
		"scheduler_log", "bytes_written", "",
	),
	"scheduler_log_native_append_events_per_sec": newMetric(
		"scheduler_log", "native_append_events_per_sec", "",
	),
	"scheduler_log_native_append_nanos_per_event": newMetric(
		"scheduler_log", "native_append_nanos_per_event", "",
	),
	"scheduler_log_native_read_events_per_sec": newMetric(
		"scheduler_log", "native_read_events_per_sec", "",
	),
	"scheduler_log_native_read_nanos_per_event": newMetric(
		"scheduler_log", "native_read_nanos_per_event", "",
	),
	"scheduler_log_native_read_nanos_total_per_sec": newMetric(
		"scheduler_log", "native_read_nanos_total_per_sec", "",
	),
	"scheduler_log_native_truncate_events_per_sec": newMetric(
		"scheduler_log", "native_truncate_events_per_sec", "",
	),
	"scheduler_log_native_truncate_nanos_per_event": newMetric(
		"scheduler_log", "native_truncate_nanos_per_event", "",
	),
	"scheduler_log_native_truncate_nanos_total_per_sec": newMetric(
		"scheduler_log", "native_truncate_nanos_total_per_sec", "",
	),
	"scheduler_thrift_getJobSummary_events_per_sec": newMetric(
		"scheduler_thrift", "getJobSummary_events_per_sec", "",
	),
	"scheduler_thrift_getJobSummary_nanos_per_event": newMetric(
		"scheduler_thrift", "getJobSummary_nanos_per_event", "",
	),
	"scheduler_thrift_getJobSummary_nanos_total_per_sec": newMetric(
		"scheduler_thrift", "getJobSummary_nanos_total_per_sec", "",
	),
	"scheduler_thrift_getQuota_events_per_sec": newMetric(
		"scheduler_thrift", "getQuota_events_per_sec", "",
	),
	"scheduler_thrift_getQuota_nanos_total_per_sec": newMetric(
		"scheduler_thrift", "getQuota_nanos_total_per_sec", "",
	),
	"sla_cluster_mtta_ms": newMetric(
		"sla", "cluster_mtta_ms",
		"Median time to assigned.",
	),
	"sla_cluster_mttr_ms": newMetric(
		"sla", "cluster_mttr_ms",
		"Median time to running.",
	),
	"sla_cluster_platform_uptime_percent": newMetric(
		"sla", "cluster_platform_uptime_percent",
		"Aggregate amount of time a job spends in a non-runnable state.",
	),
	"system_free_physical_memory_mb": newMetric(
		"system", "free_physical_memory_mb",
		"Free physical memory in MB.",
	),
	"system_load_avg": newMetric(
		"system", "load_avg",
		"1 minute load average.",
	),
	"task_throttle_events_per_sec": newMetric(
		"task", "throttle_events_per_sec", "",
	),
	"task_throttle_ms_total_per_sec": newMetric(
		"task", "throttle_ms_total_per_sec", "",
	),
	"timeout_queue_size": newMetric(
		"", "timeout_queue_size", "",
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

var schedulerLifecycleRe = regexp.MustCompile("scheduler_lifecycle_(?P<state>[A-Z]+)")

func schedulerLifecycleMetric(name string, value float64) (metric prometheus.Metric) {
	match := schedulerLifecycleRe.FindStringSubmatch(name)
	if len(match) == 2 {
		state := match[1]

		metric = prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, "", "scheduler_lifecycle"),
				"Scheduler lifecycle.",
				stateLabel, nil,
			),
			prometheus.GaugeValue,
			value, state,
		)
	}
	return metric
}

var prefixFunc = map[string]func(string, float64) prometheus.Metric{
	"tasks_":               tasksMetric,
	"tasks_lost_rack_":     tasksRackMetric,
	"task_store_":          taskStoreMetric,
	"sla_":                 slaMetric,
	"update_transition_":   updateMetric,
	"scheduler_lifecycle_": schedulerLifecycleMetric,
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
