package pkg

import (
	"reflect"
	"testing"
)

func TestGetNonEmptyFields(t *testing.T) {
	type MyStruct struct {
		Field1 string `json:"field1" dc:"测试1"`
		Field2 int    `json:"field2" dc:"测试2"`
	}
	type args struct {
		obj     interface{}
		tagName string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{name: "a struct", args: args{obj: MyStruct{Field1: "value1", Field2: 10}, tagName: "json"}, want: map[string]interface{}{"field1": "value1", "field2": 10}},
		{name: "pointer to a struct", args: args{obj: &MyStruct{Field1: "value1", Field2: 10}, tagName: "json"}, want: map[string]interface{}{"field1": "value1", "field2": 10}},
		{name: "pointer to a slice", args: args{obj: &[]MyStruct{{Field1: "value1", Field2: 10}}, tagName: "json"}, want: map[string]interface{}{"field1": "value1", "field2": 10}},
		{name: "a slice", args: args{obj: []MyStruct{{Field1: "value1", Field2: 10}}, tagName: "json"}, want: map[string]interface{}{"field1": "value1", "field2": 10}},
		{name: "a slice and pointer to a struct", args: args{obj: []*MyStruct{{Field1: "value1", Field2: 10}}, tagName: "json"}, want: map[string]interface{}{"field1": "value1", "field2": 10}},
		{name: "pointer to a slice and pointer to a struct", args: args{obj: &[]*MyStruct{{Field1: "value1", Field2: 10}}, tagName: "json"}, want: map[string]interface{}{"field1": "value1", "field2": 10}},
		{name: "pointer to a slice not struct", args: args{obj: &[]int{1, 2, 3}, tagName: "json"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNonEmptyFields(tt.args.obj, tt.args.tagName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNonEmptyFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNonEmptyFields() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTagsMap(t *testing.T) {
	type MyStruct struct {
		Field1 string `json:"field1" dc:"测试1"`
		Field2 int    `json:"field2" dc:"测试2"`
	}
	type args struct {
		obj      interface{}
		tagName  string
		vTagName string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantNil bool
	}{
		{name: "a struct", args: args{obj: MyStruct{Field1: "value1", Field2: 10}, tagName: "json", vTagName: "dc"}, want: map[string]string{"field1": "测试1", "field2": "测试2"}},
		{name: "pointer to a struct", args: args{obj: &MyStruct{Field1: "value1", Field2: 10}, tagName: "json", vTagName: "dc"}, want: map[string]string{"field1": "测试1", "field2": "测试2"}},
		{name: "pointer to a slice", args: args{obj: &[]MyStruct{{Field1: "value1", Field2: 10}}, tagName: "json", vTagName: "dc"}, want: map[string]string{"field1": "测试1", "field2": "测试2"}},
		{name: "a slice", args: args{obj: []MyStruct{{Field1: "value1", Field2: 10}}, tagName: "json", vTagName: "dc"}, want: map[string]string{"field1": "测试1", "field2": "测试2"}},
		{name: "a slice and pointer to a struct", args: args{obj: []*MyStruct{{Field1: "value1", Field2: 10}}, tagName: "json", vTagName: "dc"}, want: map[string]string{"field1": "测试1", "field2": "测试2"}},
		{name: "pointer to a slice and pointer to a struct", args: args{obj: &[]*MyStruct{{Field1: "value1", Field2: 10}}, tagName: "json", vTagName: "dc"}, want: map[string]string{"field1": "测试1", "field2": "测试2"}},
		{name: "pointer to a slice not struct", args: args{obj: &[]int{1, 2, 3}, tagName: "json", vTagName: "dc"}, want: nil, wantNil: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetTagsMap(tt.args.obj, tt.args.tagName, tt.args.vTagName)
			if (got == nil) != tt.wantNil {
				t.Errorf("GetNonEmptyFields(), wantNil %v", tt.wantNil)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTagsMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSlice(t *testing.T) {
	type MyStruct struct {
		Field1 string `json:"field1" dc:"测试1"`
		Field2 int    `json:"field2" dc:"测试2"`
	}
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "a struct", args: args{obj: MyStruct{Field1: "value1", Field2: 10}}, want: false},
		{name: "a slice", args: args{obj: []int{1, 2, 3}}, want: true},
		{name: "pointer to a slice", args: args{obj: &[]int{1, 2, 3}}, want: true},
		{name: "pointer to a slice", args: args{obj: &[]MyStruct{{Field1: "value1", Field2: 10}}}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSlice(tt.args.obj); got != tt.want {
				t.Errorf("IsSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
