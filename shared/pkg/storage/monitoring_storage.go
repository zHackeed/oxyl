package storage

import (
	"context"
	sql2 "database/sql"
	"fmt"
	"time"

	"zhacked.me/oxyl/shared/pkg/datasource"
	"zhacked.me/oxyl/shared/pkg/models"
)

type MonitoringStorage struct {
	conn *datasource.TimescaleConnection
}

func NewMonitoringStorage(conn *datasource.TimescaleConnection) *MonitoringStorage {
	return &MonitoringStorage{conn: conn}
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

	sql = `INSERT INTO agent_network_metrics (timestamp, 
                                   agent, interface_name, rx_bytes, tx_bytes, 
                                   rx_packets, tx_packets, rx_rate, 
                                   tx_rate, rx_packet_rate, tx_packet_rate) VALUES (
									$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
                                   ) `

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

	sql := `SELECT timestamp, cpu_usage, memory_usage, uptime FROM agent_general_metrics WHERE agent = $1 AND timestamp > NOW() - make_interval(secs => $2)`
	rows, err := m.conn.Pool().Query(ctx, sql, agentId, duration.Seconds())
	if err != nil {
		return nil, fmt.Errorf("failed to query general metrics: %w", err)
	}

	for rows.Next() {
		metric := new(models.AgentGeneralMetrics)

		if err := rows.Scan(&metric.When, &metric.CPUUsage, &metric.MemoryUsage, &metric.Uptime); err != nil {
			return nil, fmt.Errorf("failed to get general metrics: %w", err)
		}

		metrics = append(metrics, metric)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get general metrics: %w", err)
	}

	if len(metrics) == 0 {
		return nil, fmt.Errorf("failed to get general metrics: no metrics found")
	}

	return metrics, nil
}

func (m *MonitoringStorage) GetMountPointMetrics(ctx context.Context, agentId string, duration time.Duration) ([]*models.AgentMountPointMetrics, error) {
	metrics := make([]*models.AgentMountPointMetrics, 0)

	sql := `SELECT timestamp, mount_point, disk_usage FROM agent_disk_metrics WHERE agent = $1 AND timestamp > NOW() - make_interval(secs => $2)`
	rows, err := m.conn.Pool().Query(ctx, sql, agentId, duration.Seconds())
	if err != nil {
		return nil, fmt.Errorf("failed to query mount point metrics: %w", err)
	}

	for rows.Next() {
		metric := new(models.AgentMountPointMetrics)

		if err := rows.Scan(&metric.When, &metric.MountPoint, &metric.DiskUsage); err != nil {
			return nil, fmt.Errorf("failed to parse data. %w", err)
		}

		metrics = append(metrics, metric)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get mount point metrics: %w", err)
	}

	if len(metrics) == 0 {
		return nil, fmt.Errorf("failed to get mount point metrics: no metrics found")
	}

	return metrics, nil
}

func (m *MonitoringStorage) GetPhysicalDiskMetrics(ctx context.Context, agentId string, duration time.Duration) ([]*models.AgentPhysicalDiskMetrics, error) {
	metrics := make([]*models.AgentPhysicalDiskMetrics, 0)

	sql := `SELECT timestamp, disk_path, 
    				health_left, 
    				media_errors_1, media_errors_2, 
    				error_rate, pending_sectors 
			FROM agent_physical_disk_metrics WHERE agent = $1 AND timestamp > NOW() - make_interval(secs => $2)`

	rows, err := m.conn.Pool().Query(ctx, sql, agentId, duration.Seconds())
	if err != nil {
		return nil, fmt.Errorf("failed to query physical disk metrics: %w", err)
	}

	for rows.Next() {
		metric := new(models.AgentPhysicalDiskMetrics)

		var healthLeft sql2.NullInt64
		var mediaErrors1 sql2.NullInt64
		var mediaErrors2 sql2.NullInt64
		var errorRate sql2.NullInt64
		var pendingSectors sql2.NullInt64

		if err := rows.Scan(&metric.When, &metric.DiskPath, &healthLeft, &mediaErrors1, &mediaErrors2, &errorRate, &pendingSectors); err != nil {
			return nil, fmt.Errorf("failed to parse data. %w", err)
		}

		if healthLeft.Valid {
			metric.HealthUsed = uint32(healthLeft.Int64)
			metric.MediaErrors1 = uint64(mediaErrors1.Int64)
			metric.MediaErrors2 = uint64(mediaErrors2.Int64)
		} else {
			metric.ErrorRate = uint64(errorRate.Int64)
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

	sql := `SELECT timestamp, interface_name, 
   					rx_bytes, tx_bytes, 
    				rx_packets, tx_packets, 
    				rx_rate, tx_rate, 
    				rx_packet_rate, tx_packet_rate 
			FROM agent_network_metrics WHERE agent = $1 AND timestamp > NOW() - make_interval(secs => $2)`
	rows, err := m.conn.Pool().Query(ctx, sql, agentId, duration.Seconds())
	if err != nil {
		return nil, fmt.Errorf("failed to query network metrics: %w", err)
	}

	for rows.Next() {
		metric := new(models.AgentNetworkMetrics)

		if err := rows.Scan(&metric.When, &metric.InterfaceName, &metric.RxBytes, &metric.TxBytes, &metric.RxPackets, &metric.TxPackets, &metric.RxRate, &metric.TxRate, &metric.RxPacketRate, &metric.TxPacketRate); err != nil {
			return nil, fmt.Errorf("failed to parse data. %w", err)
		}

		metrics = append(metrics, metric)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to get network metrics: %w", err)
	}

	return metrics, nil
}

// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! ------------------------------
