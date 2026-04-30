package storage

import (
	"context"
	sql2 "database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"zhacked.me/oxyl/shared/pkg/datasource"
	"zhacked.me/oxyl/shared/pkg/models"
)

type MonitoringStorage struct {
	conn *datasource.TimescaleConnection
}

func NewMonitoringStorage(conn *datasource.TimescaleConnection) *MonitoringStorage {
	return &MonitoringStorage{conn: conn}
}

// The bucket size is used to group metrics into time buckets for aggregation
// If the duration is less than 15 minutes, no bucketing is used (0)
// If the duration is less than 6 hours, bucketing is done in 1 minute intervals
// Otherwise, bucketing is done in 1 hour intervals
// This avoids sending too much data and straining the database and the device parsing everyting. We aren't grafana, we cannot handle that much data
func bucketSize(duration time.Duration) time.Duration {
	switch {
	case duration <= 15*time.Minute:
		return 0
	case duration <= 6*time.Hour:
		return time.Minute
	default:
		return time.Hour
	}
}

func (m *MonitoringStorage) InsertData(ctx context.Context, agentId string,
	general *models.AgentGeneralMetrics, mounts []*models.AgentMountPointMetrics,
	disks []*models.AgentPhysicalDiskMetrics, interfaces []*models.AgentNetworkMetrics) error {

	tx, err := m.conn.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	now := time.Now()

	sql := `INSERT INTO agent_general_metrics (timestamp, agent, cpu_usage, memory_usage, uptime) VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.Exec(ctx, sql, now, agentId, general.CPUUsage, general.MemoryUsage, general.Uptime)
	if err != nil {
		return fmt.Errorf("failed to insert agent general metrics: %w", err)
	}

	sql = `INSERT INTO agent_disk_metrics (timestamp, agent, mount_point, disk_usage) VALUES ($1, $2, $3, $4)`
	for _, mountPoint := range mounts {
		_, err = tx.Exec(ctx, sql, now, agentId, mountPoint.MountPoint, mountPoint.DiskUsage)
		if err != nil {
			return fmt.Errorf("failed to insert mount point metrics: %w", err)
		}
	}

	sql = `INSERT INTO agent_physical_disk_metrics (timestamp, agent, disk_path, health_left, media_errors_1, media_errors_2, error_rate, pending_sectors) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	for _, blockDevice := range disks {
		_, err = tx.Exec(ctx, sql, now, agentId,
			blockDevice.DiskPath, blockDevice.HealthUsed,
			blockDevice.MediaErrors1, blockDevice.MediaErrors2,
			blockDevice.ErrorRate, blockDevice.PendingSectors)
		if err != nil {
			return fmt.Errorf("failed to insert physical disk metrics: %w", err)
		}
	}

	sql = `INSERT INTO agent_network_metrics (timestamp, agent, interface_name, rx_bytes, tx_bytes, rx_packets, tx_packets, rx_rate, tx_rate, rx_packet_rate, tx_packet_rate) 
		   VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	for _, interfaceData := range interfaces {
		_, err = tx.Exec(ctx, sql, now, agentId, interfaceData.InterfaceName,
			interfaceData.RxBytes, interfaceData.TxBytes,
			interfaceData.RxPackets, interfaceData.TxPackets,
			interfaceData.RxRate, interfaceData.TxRate,
			interfaceData.RxPacketRate, interfaceData.TxPacketRate)
		if err != nil {
			return fmt.Errorf("failed to insert network metrics: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (m *MonitoringStorage) GetGeneralMetrics(ctx context.Context, agentId string, duration time.Duration) ([]*models.AgentGeneralMetrics, error) {
	metrics := make([]*models.AgentGeneralMetrics, 0)
	bucket := bucketSize(duration)

	var rows pgx.Rows
	var err error

	if bucket == 0 {
		rows, err = m.conn.Pool().Query(ctx,
			`SELECT timestamp, cpu_usage, memory_usage, uptime
			 FROM agent_general_metrics
			 WHERE agent = $1 AND timestamp > NOW() - make_interval(secs => $2)
			 ORDER BY timestamp ASC`,
			agentId, duration.Seconds())
	} else {
		rows, err = m.conn.Pool().Query(ctx,
			`SELECT time_bucket(make_interval(secs => $3), timestamp) AS timestamp,
					AVG(cpu_usage)::float8,
					AVG(memory_usage)::bigint,
					MAX(uptime)
			 FROM agent_general_metrics
			 WHERE agent = $1 AND timestamp > NOW() - make_interval(secs => $2)
			 GROUP BY 1 ORDER BY 1 ASC`,
			agentId, duration.Seconds(), bucket.Seconds())
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query general metrics: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		metric := new(models.AgentGeneralMetrics)
		if err := rows.Scan(&metric.When, &metric.CPUUsage, &metric.MemoryUsage, &metric.Uptime); err != nil {
			return nil, fmt.Errorf("failed to scan general metrics: %w", err)
		}
		metrics = append(metrics, metric)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get general metrics: %w", err)
	}

	return metrics, nil
}

func (m *MonitoringStorage) GetMountPointMetrics(ctx context.Context, agentId string, duration time.Duration) ([]*models.AgentMountPointMetrics, error) {
	metrics := make([]*models.AgentMountPointMetrics, 0)
	bucket := bucketSize(duration)

	var rows pgx.Rows
	var err error

	if bucket == 0 {
		rows, err = m.conn.Pool().Query(ctx,
			`SELECT timestamp, mount_point, disk_usage
			 FROM agent_disk_metrics
			 WHERE agent = $1 AND timestamp > NOW() - make_interval(secs => $2)
			 ORDER BY timestamp ASC`,
			agentId, duration.Seconds())
	} else {
		rows, err = m.conn.Pool().Query(ctx,
			`SELECT time_bucket(make_interval(secs => $3), timestamp) AS timestamp,
					mount_point,
					AVG(disk_usage)::bigint
			 FROM agent_disk_metrics
			 WHERE agent = $1 AND timestamp > NOW() - make_interval(secs => $2)
			 GROUP BY 1, 2 ORDER BY 1 ASC`,
			agentId, duration.Seconds(), bucket.Seconds())
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query mount point metrics: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		metric := new(models.AgentMountPointMetrics)
		if err := rows.Scan(&metric.When, &metric.MountPoint, &metric.DiskUsage); err != nil {
			return nil, fmt.Errorf("failed to scan mount point metrics: %w", err)
		}
		metrics = append(metrics, metric)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get mount point metrics: %w", err)
	}

	return metrics, nil
}

func (m *MonitoringStorage) GetPhysicalDiskMetrics(ctx context.Context, agentId string, duration time.Duration) ([]*models.AgentPhysicalDiskMetrics, error) {
	metrics := make([]*models.AgentPhysicalDiskMetrics, 0)
	bucket := bucketSize(duration)

	var rows pgx.Rows
	var err error

	if bucket == 0 {
		rows, err = m.conn.Pool().Query(ctx,
			`SELECT timestamp, disk_path, health_left, media_errors_1, media_errors_2, error_rate, pending_sectors
			 FROM agent_physical_disk_metrics
			 WHERE agent = $1 AND timestamp > NOW() - make_interval(secs => $2)
			 ORDER BY timestamp ASC`,
			agentId, duration.Seconds())
	} else {
		rows, err = m.conn.Pool().Query(ctx,
			`SELECT time_bucket(make_interval(secs => $3), timestamp) AS timestamp,
					disk_path,
					AVG(health_left)::int,
					AVG(media_errors_1)::bigint,
					AVG(media_errors_2)::bigint,
					AVG(error_rate)::bigint,
					AVG(pending_sectors)::int
			 FROM agent_physical_disk_metrics
			 WHERE agent = $1 AND timestamp > NOW() - make_interval(secs => $2)
			 GROUP BY 1, 2 ORDER BY 1 ASC`,
			agentId, duration.Seconds(), bucket.Seconds())
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query physical disk metrics: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		metric := new(models.AgentPhysicalDiskMetrics)

		var healthLeft sql2.NullInt64
		var mediaErrors1 sql2.NullInt64
		var mediaErrors2 sql2.NullInt64
		var errorRate sql2.NullInt64
		var pendingSectors sql2.NullInt64

		if err := rows.Scan(&metric.When, &metric.DiskPath, &healthLeft, &mediaErrors1, &mediaErrors2, &errorRate, &pendingSectors); err != nil {
			return nil, fmt.Errorf("failed to scan physical disk metrics: %w", err)
		}

		if healthLeft.Valid {
			metric.HealthUsed = uint32(healthLeft.Int64)
		}
		if mediaErrors1.Valid {
			metric.MediaErrors1 = uint64(mediaErrors1.Int64)
		}
		if mediaErrors2.Valid {
			metric.MediaErrors2 = uint64(mediaErrors2.Int64)
		}
		if errorRate.Valid {
			metric.ErrorRate = uint64(errorRate.Int64)
		}
		if pendingSectors.Valid {
			metric.PendingSectors = uint32(pendingSectors.Int64)
		}

		metrics = append(metrics, metric)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get physical disk metrics: %w", err)
	}

	return metrics, nil
}

func (m *MonitoringStorage) GetNetworkMetrics(ctx context.Context, agentId string, duration time.Duration) ([]*models.AgentNetworkMetrics, error) {
	metrics := make([]*models.AgentNetworkMetrics, 0)
	bucket := bucketSize(duration)

	var rows pgx.Rows
	var err error

	if bucket == 0 {
		rows, err = m.conn.Pool().Query(ctx,
			`SELECT timestamp, interface_name, rx_bytes, tx_bytes, rx_packets, tx_packets, rx_rate, tx_rate, rx_packet_rate, tx_packet_rate
			 FROM agent_network_metrics
			 WHERE agent = $1 AND timestamp > NOW() - make_interval(secs => $2)
			 ORDER BY timestamp ASC`,
			agentId, duration.Seconds())
	} else {
		rows, err = m.conn.Pool().Query(ctx,
			`SELECT time_bucket(make_interval(secs => $3), timestamp) AS timestamp,
					interface_name,
					AVG(rx_bytes)::bigint,
					AVG(tx_bytes)::bigint,
					AVG(rx_packets)::bigint,
					AVG(tx_packets)::bigint,
					AVG(rx_rate)::bigint,
					AVG(tx_rate)::bigint,
					AVG(rx_packet_rate)::bigint,
					AVG(tx_packet_rate)::bigint
			 FROM agent_network_metrics
			 WHERE agent = $1 AND timestamp > NOW() - make_interval(secs => $2)
			 GROUP BY 1, 2 ORDER BY 1 ASC`,
			agentId, duration.Seconds(), bucket.Seconds())
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query network metrics: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		metric := new(models.AgentNetworkMetrics)
		if err := rows.Scan(&metric.When, &metric.InterfaceName,
			&metric.RxBytes, &metric.TxBytes,
			&metric.RxPackets, &metric.TxPackets,
			&metric.RxRate, &metric.TxRate,
			&metric.RxPacketRate, &metric.TxPacketRate); err != nil {
			return nil, fmt.Errorf("failed to scan network metrics: %w", err)
		}
		metrics = append(metrics, metric)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get network metrics: %w", err)
	}

	return metrics, nil
}