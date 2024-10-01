package repo

import "testing"

func TestDataType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		d    DataType
		want bool
	}{
		{
			name: "text",
			d:    TextType,
			want: true,
		},
		{
			name: "card",
			d:    CardType,
			want: true,
		},
		{
			name: "binary",
			d:    BinaryType,
			want: true,
		},
		{
			name: "random",
			d:    DataType("random"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
