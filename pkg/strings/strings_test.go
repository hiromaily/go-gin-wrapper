package strings

import (
	"testing"
)

func TestSearchIndex(t *testing.T) {
	type args struct {
		target  string
		sources []string
	}
	type want struct {
		idx int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path 1",
			args: args{
				target: "foo-bar",
				sources: []string{
					"aaaaaaa",
					"bbbbbbb",
					"foo-ber",
					"",
					"foo-bar",
				},
			},
			want: want{
				idx: 4,
			},
		},
		{
			name: "happy path 2",
			args: args{
				target: "foo-bar",
				sources: []string{
					"foo-bar",
					"bbbbbbb",
					"foo-ber",
					"",
					"foo-bar",
				},
			},
			want: want{
				idx: 0,
			},
		},
		{
			name: "no target",
			args: args{
				target: "foo-bar",
				sources: []string{
					"aaaaaaa",
					"bbbbbbb",
					"foo-ber",
					"",
				},
			},
			want: want{
				idx: -1,
			},
		},
		{
			name: "upper case and lower case",
			args: args{
				target: "foo-bar",
				sources: []string{
					"aaaaaaa",
					"bbbbbbb",
					"Foo-Bar",
				},
			},
			want: want{
				idx: -1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SearchIndex(tt.args.target, tt.args.sources)
			if got != tt.want.idx {
				t.Errorf("SearchIndex(): got = %d, want %d", got, tt.want.idx)
			}
		})
	}
}

func TestSearchIndexLower(t *testing.T) {
	type args struct {
		target  string
		sources []string
	}
	type want struct {
		idx int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path 1",
			args: args{
				target: "foo-bar",
				sources: []string{
					"aaaaaaa",
					"bbbbbbb",
					"foo-ber",
					"",
					"foo-bar",
				},
			},
			want: want{
				idx: 4,
			},
		},
		{
			name: "happy path 2",
			args: args{
				target: "foo-bar",
				sources: []string{
					"foo-bar",
					"bbbbbbb",
					"foo-ber",
					"",
					"foo-bar",
				},
			},
			want: want{
				idx: 0,
			},
		},
		{
			name: "no target",
			args: args{
				target: "foo-bar",
				sources: []string{
					"aaaaaaa",
					"bbbbbbb",
					"foo-ber",
					"",
				},
			},
			want: want{
				idx: -1,
			},
		},
		{
			name: "upper case and lower case",
			args: args{
				target: "foo-bar",
				sources: []string{
					"aaaaaaa",
					"bbbbbbb",
					"Foo-Bar",
				},
			},
			want: want{
				idx: 2,
			},
		},
		{
			name: "upper case and lower case 2",
			args: args{
				target: "FOO-bar",
				sources: []string{
					"aaaaaaa",
					"bbbbbbb",
					"foo-BAR",
				},
			},
			want: want{
				idx: 2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SearchIndexLower(tt.args.target, tt.args.sources)
			if got != tt.want.idx {
				t.Errorf("SearchIndexLower(): got = %d, want %d", got, tt.want.idx)
			}
		})
	}
}

func TestItos(t *testing.T) {
	type args struct {
		target interface{}
	}
	type want struct {
		expected string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path 1",
			args: args{
				target: "foo-bar",
			},
			want: want{
				expected: "foo-bar",
			},
		},
		{
			name: "int can not be converted",
			args: args{
				target: 100,
			},
			want: want{
				expected: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Itos(tt.args.target)
			if got != tt.want.expected {
				t.Errorf("Itos(): got = %s, want %s", got, tt.want.expected)
			}
		})
	}
}

func TestAtoi(t *testing.T) {
	type args struct {
		target string
	}
	type want struct {
		expected int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "happy path 1",
			args: args{
				target: "100",
			},
			want: want{
				expected: 100,
			},
		},
		{
			name: "wrong value",
			args: args{
				target: "foo-bar",
			},
			want: want{
				expected: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Atoi(tt.args.target)
			if got != tt.want.expected {
				t.Errorf("Itos(): got = %d, want %d", got, tt.want.expected)
			}
		})
	}
}
