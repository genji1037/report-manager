package service

import (
	"fmt"
	"report-manager/collector"
	"report-manager/config"
	"report-manager/report"
	"testing"
)

func TestRadarOTCNotice(t *testing.T) {
	tests := []struct {
		name       string
		collectors []collector.Collector
		expectErr  error
	}{
		{
			name: "both waiting real names and failed transfers",
			collectors: []collector.Collector{
				&HasRealNames{},
				&HasTrans{},
			},
			expectErr: nil,
		},
		{
			name: "only waiting real names no failed transfers",
			collectors: []collector.Collector{
				&HasRealNames{},
				&NoTrans{},
			},
			expectErr: nil,
		},
		{
			name: "only failed transfers no waiting real names",
			collectors: []collector.Collector{
				&NoRealNames{},
				&HasTrans{},
			},
			expectErr: nil,
		},
		{
			name: "neither failed transfers nor waiting real names",
			collectors: []collector.Collector{
				&NoRealNames{},
				&NoTrans{},
			},
			expectErr: report.DoNotReport,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reportMaker := report.NewMaker("radar notice", config.GetServer().Template.RadarOTCNotice.Content, tt.collectors)

			content, err := reportMaker.Make()
			if err != tt.expectErr {
				t.Fatalf("RadarOTCNotice() error = %v, wantErr %v", err, tt.expectErr)
			}
			fmt.Println(content)
		})
	}
}

type HasRealNames struct {
	collector.RadarWaitingRealNames
}

func (c *HasRealNames) Collect() error {
	c.Num = 100
	return nil
}

type NoRealNames struct {
	collector.RadarWaitingRealNames
}

func (c *NoRealNames) Collect() error {
	return nil
}

type HasTrans struct {
	collector.RadarFailedTransfer
}

func (c *HasTrans) Collect() error {
	c.Num = 100
	return nil
}

type NoTrans struct {
	collector.RadarFailedTransfer
}

func (c *NoTrans) Collect() error {
	return nil
}
