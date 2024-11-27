package middlewares

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

/*
	// 配置 MySQL 数据源名称（DSN）
	dsn := "root:88888888@tcp(localhost:33062)/afserver"

	// 创建 MySQL exporter
	exporter, err := middlewares.NewExporter(dsn)
	if err != nil {
		log.Fatalf("Error creating exporter: %v", err)
	}

	// 注册 Prometheus 指标
	exporter.RegisterPrometheusMetrics()
*/
// MySQLExporter 定义了一个 MySQL Exporter
type MySQLExporter struct {
	db *sql.DB
}

// NewExporter 创建一个新的 MySQL exporter
func NewExporter(dsn string) (*MySQLExporter, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening MySQL connection: %v", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging MySQL: %v", err)
	}

	fmt.Println("Successfully connected to MySQL")

	return &MySQLExporter{db: db}, nil
}

// CollectMySQLMetrics 从 MySQL 获取并返回指标
func (e *MySQLExporter) CollectMySQLMetrics() (map[string]float64, error) {
	metrics := make(map[string]float64)

	// 获取当前 MySQL 连接数
	rows, err := e.db.Query("SHOW STATUS LIKE 'Threads_connected'")
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var name string
	var value float64
	if rows.Next() {
		if err := rows.Scan(&name, &value); err != nil {
			return nil, fmt.Errorf("failed to scan result: %v", err)
		}
		metrics[name] = value
	}

	return metrics, nil
}

// RegisterPrometheusMetrics 将 MySQL 收集到的指标注册到 Prometheus
func (e *MySQLExporter) RegisterPrometheusMetrics() {
	// 创建 Prometheus 指标
	connections := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "mysql_threads_connected",
			Help: "Number of MySQL threads connected",
		},
	)

	// 注册 Prometheus 指标
	prometheus.MustRegister(connections)

	// 定时从 MySQL 拉取数据
	go func() {
		for {
			metrics, err := e.CollectMySQLMetrics()
			if err != nil {
				log.Printf("Error collecting metrics: %v", err)
				continue
			}

			// 将指标数据发送到 Prometheus
			for _, value := range metrics {
				connections.Set(value)
			}

			// 睡眠一段时间（例如，30秒）
			<-time.After(30 * time.Second)
		}
	}()
}
