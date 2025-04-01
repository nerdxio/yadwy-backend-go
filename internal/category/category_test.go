package category

import (
	"testing"
)

func TestNewCategory(t *testing.T) {
	type args struct {
		name        string
		description string
	}
	tests := []struct {
		name string
		args args
		want Category
	}{
		{
			name: "should create new category with valid inputs",
			args: args{
				name:        "Electronics",
				description: "Electronic devices and gadgets",
			},
			want: Category{
				Name:        "Electronics",
				Description: "Electronic devices and gadgets",
			},
		},
		{
			name: "should create new category with empty description",
			args: args{
				name:        "Books",
				description: "",
			},
			want: Category{
				Name:        "Books",
				Description: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCategory(tt.args.name, tt.args.description)
			if got.Name != tt.want.Name || got.Description != tt.want.Description {
				t.Errorf("NewCategory() = %v, want %v", got, tt.want)
			}
			if got.CreatedAt.IsZero() || got.UpdatedAt.IsZero() {
				t.Error("NewCategory() timestamps should not be zero")
			}
		})
	}
}
