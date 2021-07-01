package code

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/brocaar/lorawan"

	"github.com/mxc-foundation/lpwan-app-server/internal/backend/networkserver"
	metricsmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/metrics"
	"github.com/mxc-foundation/lpwan-app-server/internal/storage/store"
)

// MigrateGatewayStats imports the gateway stats from the network-server.
func MigrateGatewayStats(handler *store.Handler) error {
	ctx := context.Background()
	ids, err := handler.GetAllGatewayIDs(ctx)
	if err != nil {
		return errors.Wrap(err, "select gateway ids error")
	}

	for _, id := range ids {
		if err := migrateGatewayStatsForGatewayID(handler, id); err != nil {
			log.WithError(err).WithFields(log.Fields{
				"gateway_id": id,
			}).Error("migrate gateway stats error")
		}
	}

	return nil
}

func migrateGatewayStatsForGatewayID(handler *store.Handler, gatewayID lorawan.EUI64) error {
	gw, err := handler.GetGateway(context.Background(), gatewayID, true)
	if err != nil {
		return errors.Wrap(err, "get gateway error")
	}

	n, err := handler.GetNetworkServer(context.Background(), gw.NetworkServerID)
	if err != nil {
		return errors.Wrap(err, "get network-server error")
	}

	nsClient, err := networkserver.GetPool().Get(n.Server, []byte(n.CACert), []byte(n.TLSCert), []byte(n.TLSKey))
	if err != nil {
		return errors.Wrap(err, "get network-server client error")
	}

	nsGw, err := nsClient.GetGateway(context.Background(), &ns.GetGatewayRequest{
		Id: gatewayID[:],
	})
	if err != nil {
		return errors.Wrap(err, "get gateway from network-server error")
	}

	if nsGw.Gateway != nil && nsGw.Gateway.Location != nil {
		gw.Latitude = nsGw.Gateway.Location.Latitude
		gw.Longitude = nsGw.Gateway.Location.Longitude
		gw.Altitude = nsGw.Gateway.Location.Altitude
	}

	if err := handler.UpdateGateway(context.Background(), &gw); err != nil {
		return errors.Wrap(err, "update gateway error")
	}

	metricsStruct := metricsmod.GetMetricsSettings()
	if err := migrateGatewayStatsForGatewayIDInterval(nsClient, gatewayID, ns.AggregationInterval_MINUTE, time.Now().Add(-metricsStruct.Redis.MinuteAggregationTTL), time.Now()); err != nil {
		return err
	}

	if err := migrateGatewayStatsForGatewayIDInterval(nsClient, gatewayID, ns.AggregationInterval_HOUR, time.Now().Add(-metricsStruct.Redis.HourAggregationTTL), time.Now()); err != nil {
		return err
	}

	if err := migrateGatewayStatsForGatewayIDInterval(nsClient, gatewayID, ns.AggregationInterval_DAY, time.Now().Add(-metricsStruct.Redis.DayAggregationTTL), time.Now()); err != nil {
		return err
	}

	if err := migrateGatewayStatsForGatewayIDInterval(nsClient, gatewayID, ns.AggregationInterval_MONTH, time.Now().Add(-metricsStruct.Redis.MonthAggregationTTL), time.Now()); err != nil {
		return err
	}

	return nil
}

func migrateGatewayStatsForGatewayIDInterval(nsClient ns.NetworkServerServiceClient, gatewayID lorawan.EUI64, interval ns.AggregationInterval, start, end time.Time) error {
	startPB := timestamppb.New(start)
	endPB := timestamppb.New(end)

	metrics, err := nsClient.GetGatewayStats(context.Background(), &ns.GetGatewayStatsRequest{
		GatewayId:      gatewayID[:],
		Interval:       interval,
		StartTimestamp: startPB,
		EndTimestamp:   endPB,
	})
	if err != nil {
		return errors.Wrap(err, "get gateway stats from network-server error")
	}

	for _, m := range metrics.Result {
		ts := m.Timestamp.AsTime()
		err = metricsmod.SaveMetricsForInterval(context.Background(), metricsmod.AggregationInterval(interval.String()), fmt.Sprintf("gw:%s", gatewayID), metricsmod.MetricsRecord{
			Time: ts,
			Metrics: map[string]float64{
				"rx_count":    float64(m.RxPacketsReceived),
				"rx_ok_count": float64(m.RxPacketsReceivedOk),
				"tx_count":    float64(m.TxPacketsReceived),
				"tx_ok_count": float64(m.TxPacketsEmitted),
			},
		})
		if err != nil {
			return errors.Wrap(err, "save metrics for interval error")
		}
	}

	return nil
}
