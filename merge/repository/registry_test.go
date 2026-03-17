package repository

import (
	"testing"
	"time"
)

func partitionValuesKey(values [][2]string) string {
	key := ""
	for _, value := range values {
		if key != "" {
			key += "/"
		}
		key += value[0] + "=" + value[1]
	}
	return key
}

func TestBuildPartitionDescsDay(t *testing.T) {
	tsData := []int64{
		time.Date(2026, time.March, 17, 10, 0, 0, 0, time.UTC).UnixNano(),
		time.Date(2026, time.March, 17, 11, 0, 0, 0, time.UTC).UnixNano(),
		time.Date(2026, time.March, 18, 1, 0, 0, 0, time.UTC).UnixNano(),
	}

	parts, err := buildPartitionDescs(tsData, "day")
	if err != nil {
		t.Fatal(err)
	}
	if len(parts) != 2 {
		t.Fatalf("expected 2 partitions, got %d", len(parts))
	}

	seen := map[string][]byte{}
	for _, part := range parts {
		seen[partitionValuesKey(part.Values)] = part.IndexMap
	}

	if _, ok := seen["date=2026-03-17"]; !ok {
		t.Fatalf("missing day partition for 2026-03-17: %+v", seen)
	}
	if _, ok := seen["date=2026-03-18"]; !ok {
		t.Fatalf("missing day partition for 2026-03-18: %+v", seen)
	}
	if seen["date=2026-03-17"][0] != 0b00000011 {
		t.Fatalf("unexpected index map for 2026-03-17: %08b", seen["date=2026-03-17"][0])
	}
	if seen["date=2026-03-18"][0] != 0b00000100 {
		t.Fatalf("unexpected index map for 2026-03-18: %08b", seen["date=2026-03-18"][0])
	}
}

func TestBuildPartitionDescsDayHour(t *testing.T) {
	tsData := []int64{
		time.Date(2026, time.March, 17, 10, 5, 0, 0, time.UTC).UnixNano(),
		time.Date(2026, time.March, 17, 10, 45, 0, 0, time.UTC).UnixNano(),
		time.Date(2026, time.March, 17, 11, 0, 0, 0, time.UTC).UnixNano(),
	}

	parts, err := buildPartitionDescs(tsData, "day,hour")
	if err != nil {
		t.Fatal(err)
	}
	if len(parts) != 2 {
		t.Fatalf("expected 2 partitions, got %d", len(parts))
	}

	seen := map[string][]byte{}
	for _, part := range parts {
		seen[partitionValuesKey(part.Values)] = part.IndexMap
	}

	if _, ok := seen["date=2026-03-17/hour=10"]; !ok {
		t.Fatalf("missing hour partition for 10: %+v", seen)
	}
	if _, ok := seen["date=2026-03-17/hour=11"]; !ok {
		t.Fatalf("missing hour partition for 11: %+v", seen)
	}
	if seen["date=2026-03-17/hour=10"][0] != 0b00000011 {
		t.Fatalf("unexpected index map for hour 10: %08b", seen["date=2026-03-17/hour=10"][0])
	}
	if seen["date=2026-03-17/hour=11"][0] != 0b00000100 {
		t.Fatalf("unexpected index map for hour 11: %08b", seen["date=2026-03-17/hour=11"][0])
	}
}

func TestBuildPartitionDescsRejectsInvalidPartition(t *testing.T) {
	_, err := buildPartitionDescs([]int64{time.Now().UnixNano()}, "day, hour")
	if err == nil {
		t.Fatal("expected invalid partition option to return an error")
	}
	if err.Error() != "unsupported partition option: \"day, hour\". Supported: \"day\", \"day,hour\"" {
		t.Fatalf("unexpected error: %v", err)
	}
}
