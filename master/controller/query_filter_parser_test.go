package controller

import (
	"reflect"
	"testing"

	"github.com/valyala/fasthttp"
)

func TestParseIntSliceQuery_RepeatedAndCommaSeparatedValues(t *testing.T) {
	args := fasthttp.Args{}
	args.Add("distributor_id", "0")
	args.Add("distributor_id", "67,68")
	args.Add("distributor_id[]", "68,69")

	actual, err := parseIntSliceQuery(&args, "distributor_id", "distributor_id", "distributor_id[]")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []int{67, 68, 69}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func TestParseIntSliceQueryAllowZero_RepeatedAndCommaSeparatedDistributorValues(t *testing.T) {
	args := fasthttp.Args{}
	args.Add("distributor_id", "0")
	args.Add("distributor_id", "67,68")
	args.Add("distributor_id[]", "68,69")

	actual, err := parseIntSliceQueryAllowZero(&args, "distributor_id", "distributor_id", "distributor_id[]")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []int{0, 67, 68, 69}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func TestParseIntSliceQuery_InvalidValue(t *testing.T) {
	args := fasthttp.Args{}
	args.Add("ot_class_id", "91,abc")

	_, err := parseIntSliceQuery(&args, "ot_class_id", "ot_class_id")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestParseIntSliceQuery_NullLiteralMeansEmpty(t *testing.T) {
	args := fasthttp.Args{}
	args.Add("distributor_id", "null")

	actual, err := parseIntSliceQuery(&args, "distributor_id", "distributor_id", "distributor_id[]")
	if err != nil {
		t.Fatalf("expected no error for null, got %v", err)
	}
	if len(actual) != 0 {
		t.Fatalf("expected empty slice, got %v", actual)
	}
}

func TestParseIntSliceQuery_WorkingDayCalendarIDArrayForms(t *testing.T) {
	args := fasthttp.Args{}
	args.Add("working_day_calendar_id", "7")
	args.Add("working_day_calendar_id[]", "8,9")
	args.Add("working_day_calendar_id[3]", "9,10")

	actual, err := parseIntSliceQuery(&args, "working_day_calendar_id", "working_day_calendar_id", "working_day_calendar_id[]")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := []int{7, 8, 9, 10}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}
