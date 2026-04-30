package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"zhacked.me/oxyl/service/thresholds/internal/models"
	"zhacked.me/oxyl/shared/pkg/datasource"
	comm "zhacked.me/oxyl/shared/pkg/models"
)

type ThresholdStorage struct {
	conn *datasource.TimescaleConnection
}

func NewThresholdStorage(conn *datasource.TimescaleConnection) *ThresholdStorage {
	return &ThresholdStorage{
		conn: conn,
	}
}

func (s *ThresholdStorage) GetAllThresholds(ctx context.Context) (map[string]*models.CompanyThresholds, error) {
	sql := `SELECT holder, notification_type, value FROM company_notification_thresholds`
	rows, err := s.conn.Pool().Query(ctx, sql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	thresholds := make(map[string]*models.CompanyThresholds)

	for rows.Next() {
		var holder string
		var notificationType comm.NotificationType
		var value float64

		err := rows.Scan(&holder, &notificationType, &value)
		if err != nil {
			return nil, err
		}

		activeThresholds, found := thresholds[holder]
		if !found {
			activeThresholds = new(models.CompanyThresholds)
			thresholds[holder] = activeThresholds
		}

		switch notificationType {
		case comm.NotificationTypeAgentCpuUsageThreshold:
			activeThresholds.SetCPU(value)
		case comm.NotificationTypeAgentDiskUsageThreshold:
			activeThresholds.SetMount(value)
		case comm.NotificationTypeAgentDiskHealthThreshold:
			activeThresholds.SetDisk(value)
		case comm.NotificationTypeAgentMemoryUsageThreshold:
			activeThresholds.SetMemory(value)
		case comm.NotificationTypeAgentNetworkUsageThreshold:
			activeThresholds.SetNetworkRX(value)
			activeThresholds.SetNetworkTX(value)
		default:
			continue
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return thresholds, nil
}

func (s *ThresholdStorage) GetThresholds(ctx context.Context, holder string) (*models.CompanyThresholds, error) {
	sql := `SELECT notification_type, value FROM company_notification_thresholds WHERE holder = $1`
	rows, err := s.conn.Pool().Query(ctx, sql, holder)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	defer rows.Close()

	thresholds := new(models.CompanyThresholds)

	for rows.Next() {
		var notificationType comm.NotificationType
		var value float64

		err := rows.Scan(&notificationType, &value)
		if err != nil {
			return nil, err
		}

		switch notificationType {
		case comm.NotificationTypeAgentCpuUsageThreshold:
			thresholds.SetCPU(value)
		case comm.NotificationTypeAgentDiskUsageThreshold:
			thresholds.SetMount(value)
		case comm.NotificationTypeAgentDiskHealthThreshold:
			thresholds.SetDisk(value)
		case comm.NotificationTypeAgentMemoryUsageThreshold:
			thresholds.SetMemory(value)
		case comm.NotificationTypeAgentNetworkUsageThreshold:
			thresholds.SetNetworkRX(value)
			thresholds.SetNetworkTX(value)
		default:
			continue
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return thresholds, nil
}

func (s *ThresholdStorage) GetGeneralAvg(ctx context.Context, agentID string, windowSecs int) (*models.GeneralMetricsAvg, error) {
	sql := `
        SELECT COALESCE(AVG(cpu_usage), 0), COALESCE(AVG(memory_usage), 0)
        FROM agent_general_metrics
        WHERE agent = $1
          AND timestamp > NOW() - make_interval(secs => $2)`

	avg := new(models.GeneralMetricsAvg)

	err := s.conn.Pool().QueryRow(ctx, sql, agentID, windowSecs).Scan(
		&avg.AvgCPU,
		&avg.AvgMemory,
	)
	if err != nil {
		return nil, err
	}

	return avg, nil
}

func (s *ThresholdStorage) GetDiskUsageAvg(ctx context.Context, agentID string, windowSecs int) ([]*models.DiskUsageAvg, error) {
	sql := `
        SELECT mount_point, COALESCE(AVG(disk_usage), 0)
        FROM agent_disk_metrics
        WHERE agent = $1
          AND timestamp > NOW() - make_interval(secs => $2)
        GROUP BY mount_point`

	rows, err := s.conn.Pool().Query(ctx, sql, agentID, windowSecs)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	defer rows.Close()

	var results []*models.DiskUsageAvg

	for rows.Next() {
		avg := new(models.DiskUsageAvg)

		err := rows.Scan(&avg.MountPoint, &avg.AvgDiskUsage)
		if err != nil {
			return nil, err
		}

		results = append(results, avg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *ThresholdStorage) GetDiskHealthAvg(ctx context.Context, agentID string, windowSecs int) ([]*models.DiskHealthAvg, error) {
	sql := `
        SELECT disk_path, COALESCE(AVG(health_left), 0), COALESCE(AVG(error_rate), 0), COALESCE(AVG(pending_sectors), 0)
        FROM agent_physical_disk_metrics
        WHERE agent = $1
          AND timestamp > NOW() - make_interval(secs => $2)
        GROUP BY disk_path`

	rows, err := s.conn.Pool().Query(ctx, sql, agentID, windowSecs)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	defer rows.Close()

	var results []*models.DiskHealthAvg

	for rows.Next() {
		avg := new(models.DiskHealthAvg)

		err := rows.Scan(
			&avg.DiskPath,
			&avg.AvgHealthLeft,
			&avg.AvgErrorRate,
			&avg.AvgPendingSectors,
		)
		if err != nil {
			return nil, err
		}

		results = append(results, avg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *ThresholdStorage) GetNetworkAvg(ctx context.Context, agentID string, windowSecs int) ([]*models.NetworkAvg, error) {
	sql := `
        SELECT interface_name, COALESCE(AVG(rx_rate), 0), COALESCE(AVG(tx_rate), 0)
        FROM agent_network_metrics
        WHERE agent = $1
          AND timestamp > NOW() - make_interval(secs => $2)
        GROUP BY interface_name`

	rows, err := s.conn.Pool().Query(ctx, sql, agentID, windowSecs)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	defer rows.Close()

	var results []*models.NetworkAvg

	for rows.Next() {
		avg := new(models.NetworkAvg)

		err := rows.Scan(&avg.InterfaceName, &avg.AvgRXRate, &avg.AvgTXRate)
		if err != nil {
			return nil, err
		}

		results = append(results, avg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
