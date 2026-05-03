package metrics

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

type DatabaseCollector struct {
	pool *pgxpool.Pool

	maxConns     *prometheus.Desc
	totalConns   *prometheus.Desc
	idleConns    *prometheus.Desc
	activeConns  *prometheus.Desc
	acquireCount *prometheus.Desc
	acquireWait  *prometheus.Desc
}

func NewDatabaseCollector(pool *pgxpool.Pool) *DatabaseCollector {
	return &DatabaseCollector{
		pool: pool,
		maxConns: prometheus.NewDesc(
			"db_pool_max_conns",
			"Maximum number of connections in the pool",
			nil, nil,
		),
		totalConns: prometheus.NewDesc(
			"db_pool_total_conns",
			"Total number of connections in the pool",
			nil, nil,
		),
		idleConns: prometheus.NewDesc(
			"db_pool_idle_conns",
			"Number of idle connections in the pool",
			nil, nil,
		),
		activeConns: prometheus.NewDesc(
			"db_pool_active_conns",
			"Number of active connections in the pool",
			nil, nil,
		),
		acquireCount: prometheus.NewDesc(
			"db_pool_acquire_count_total",
			"Total number of times a connection was acquired",
			nil, nil,
		),
		acquireWait: prometheus.NewDesc(
			"db_pool_acquire_wait_duration_seconds_total",
			"Total duration of all connection acquisitions",
			nil, nil,
		),
	}
}

func (c *DatabaseCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.maxConns
	ch <- c.totalConns
	ch <- c.idleConns
	ch <- c.activeConns
	ch <- c.acquireCount
	ch <- c.acquireWait
}

func (c *DatabaseCollector) Collect(ch chan<- prometheus.Metric) {
	stats := c.pool.Stat()
	ch <- prometheus.MustNewConstMetric(c.maxConns, prometheus.GaugeValue, float64(stats.MaxConns()))
	ch <- prometheus.MustNewConstMetric(c.totalConns, prometheus.GaugeValue, float64(stats.TotalConns()))
	ch <- prometheus.MustNewConstMetric(c.idleConns, prometheus.GaugeValue, float64(stats.IdleConns()))
	ch <- prometheus.MustNewConstMetric(c.activeConns, prometheus.GaugeValue, float64(stats.AcquiredConns()))
	ch <- prometheus.MustNewConstMetric(c.acquireCount, prometheus.CounterValue, float64(stats.AcquireCount()))
	ch <- prometheus.MustNewConstMetric(c.acquireWait, prometheus.CounterValue, stats.AcquireDuration().Seconds())
}

type RedisCollector struct {
	client *redis.Client

	hits       *prometheus.Desc
	misses     *prometheus.Desc
	timeouts   *prometheus.Desc
	totalConns *prometheus.Desc
	idleConns  *prometheus.Desc
}

func NewRedisCollector(client *redis.Client) *RedisCollector {
	return &RedisCollector{
		client: client,
		hits: prometheus.NewDesc(
			"redis_pool_hits_total",
			"Number of times a connection was found in the idle pool",
			nil, nil,
		),
		misses: prometheus.NewDesc(
			"redis_pool_misses_total",
			"Number of times a connection was not found in the idle pool",
			nil, nil,
		),
		timeouts: prometheus.NewDesc(
			"redis_pool_timeouts_total",
			"Number of times a connection pool timeout occurred",
			nil, nil,
		),
		totalConns: prometheus.NewDesc(
			"redis_pool_total_conns",
			"Total number of connections in the pool",
			nil, nil,
		),
		idleConns: prometheus.NewDesc(
			"redis_pool_idle_conns",
			"Number of idle connections in the pool",
			nil, nil,
		),
	}
}

func (c *RedisCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.hits
	ch <- c.misses
	ch <- c.timeouts
	ch <- c.totalConns
	ch <- c.idleConns
}

func (c *RedisCollector) Collect(ch chan<- prometheus.Metric) {
	stats := c.client.PoolStats()
	ch <- prometheus.MustNewConstMetric(c.hits, prometheus.CounterValue, float64(stats.Hits))
	ch <- prometheus.MustNewConstMetric(c.misses, prometheus.CounterValue, float64(stats.Misses))
	ch <- prometheus.MustNewConstMetric(c.timeouts, prometheus.CounterValue, float64(stats.Timeouts))
	ch <- prometheus.MustNewConstMetric(c.totalConns, prometheus.GaugeValue, float64(stats.TotalConns))
	ch <- prometheus.MustNewConstMetric(c.idleConns, prometheus.GaugeValue, float64(stats.IdleConns))
}
