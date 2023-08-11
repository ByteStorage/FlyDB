package structure

import (
	"errors"
	"github.com/ByteStorage/FlyDB/config"
	"github.com/ByteStorage/FlyDB/engine"
	_const "github.com/ByteStorage/FlyDB/lib/const"
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"sort"
	"testing"
)

func initTestSetDb() (*SetStructure, *config.Options) {
	opts := config.DefaultOptions
	dir, _ := os.MkdirTemp("", "TestSetStructure")
	opts.DirPath = dir
	str, _ := NewSetStructure(opts)
	return str, &opts
}

func TestSAdd(t *testing.T) {
	// Creating a list of cases for testing.
	tests := []struct {
		name        string
		setup       func(*SetStructure)
		key         string
		membersAdd  []string
		membersWant []string
		wantErr     error
	}{
		{
			name:        "test when key empty",
			setup:       func(s *SetStructure) {},
			key:         "",
			membersAdd:  []string{"key1", "key2"},
			membersWant: []string{},
			wantErr:     _const.ErrKeyIsEmpty,
		},
		{
			name:        "test add two members",
			setup:       func(s *SetStructure) {},
			key:         "destination",
			membersAdd:  []string{"key1", "key2"},
			membersWant: []string{"key1", "key2"},
			wantErr:     nil,
		},
		{
			name: "test add three members one duplicate",
			setup: func(s *SetStructure) {
				_ = s.SAdds("destination", "key2")
			},
			key:         "destination",
			membersAdd:  []string{"key1", "key2"},
			membersWant: []string{"key1", "key2"},
			wantErr:     nil,
		},
		{
			name: "test add db not init",
			setup: func(s *SetStructure) {
				s.db = nil
			},
			key:         "destination",
			membersAdd:  []string{"key1", "key2"},
			membersWant: []string{"key1", "key2"},
			wantErr:     ErrSetNotInitialized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := initTestSetDb()
			tt.setup(s)
			// Call the method with test case parameters.
			for _, s2 := range tt.membersAdd {
				err := s.SAdd(tt.key, s2)
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("SRems() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if tt.wantErr == nil {
				// positive verify
				assert.True(t, s.exists(tt.key, tt.membersWant...))
			}
		})
	}

}

func TestSAdds(t *testing.T) {
	// Creating a list of cases for testing.
	tests := []struct {
		name        string
		setup       func(*SetStructure)
		key         string
		membersAdd  []string
		membersWant []string
		wantErr     error
	}{
		{
			name:        "test when key empty",
			setup:       func(s *SetStructure) {},
			key:         "",
			membersAdd:  []string{"key1", "key2"},
			membersWant: []string{},
			wantErr:     _const.ErrKeyIsEmpty,
		},
		{
			name:        "test add two members",
			setup:       func(s *SetStructure) {},
			key:         "destination",
			membersAdd:  []string{"key1", "key2"},
			membersWant: []string{"key1", "key2"},
			wantErr:     nil,
		},
		{
			name: "test add three members one duplicate",
			setup: func(s *SetStructure) {
				_ = s.SAdds("destination", "key2")
			},
			key:         "destination",
			membersAdd:  []string{"key1", "key2"},
			membersWant: []string{"key1", "key2"},
			wantErr:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := initTestSetDb()
			tt.setup(s)
			// Call the method with test case parameters.
			err := s.SAdds(tt.key, tt.membersAdd...)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("SRems() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == nil {
				// positive verify
				assert.True(t, s.exists(tt.key, tt.membersWant...))
			}
		})
	}
}

func TestSRems(t *testing.T) {
	// Creating a list of cases for testing.
	tests := []struct {
		name        string
		setup       func(*SetStructure)
		key         string
		membersRem  []string
		membersWant []string
		wantErr     error
	}{
		{
			name:        "test empty key ",
			setup:       func(s *SetStructure) {},
			key:         "",
			membersRem:  []string{"key1", "key2"},
			membersWant: []string{},
			wantErr:     _const.ErrKeyIsEmpty,
		},
		{
			name:        "test when key not found",
			setup:       func(s *SetStructure) {},
			key:         "destination",
			membersRem:  []string{"key1", "key2"},
			membersWant: []string{},
			wantErr:     _const.ErrKeyNotFound,
		},
		{
			name: "test when remove one key",
			setup: func(s *SetStructure) {
				_ = s.SAdds("destination", "key1", "key2")
			},
			key:         "destination",
			membersRem:  []string{"key1"},
			membersWant: []string{"key2"},
			wantErr:     nil,
		},
		{
			name: "test remove all keys",
			setup: func(s *SetStructure) {
				_ = s.SAdds("destination", "key1", "key2")
			},
			key:         "destination",
			membersRem:  []string{"key1", "key1"},
			membersWant: []string{},
			wantErr:     nil,
		},
		{
			name: "test remove FSets that don't exist",
			setup: func(s *SetStructure) {
				_ = s.SAdds("destination", "key1", "key2")
			},
			key:         "destination",
			membersRem:  []string{"key10", "key11"},
			membersWant: []string{},
			wantErr:     ErrMemberNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := initTestSetDb()
			tt.setup(s)
			// Call the method with test case parameters.
			err := s.SRems(tt.key, tt.membersRem...)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("SRems() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr == nil {
				// positive verify
				assert.True(t, s.exists(tt.key, tt.membersWant...))
				// verify removal
				assert.False(t, s.exists(tt.key, tt.membersRem...))
			}
		})
	}
}

func TestSRem(t *testing.T) {
	// Creating a list of cases for testing.
	tests := []struct {
		name        string
		setup       func(*SetStructure)
		key         string
		membersRem  []string
		membersWant []string
		wantErr     error
	}{
		{
			name:        "test when key not found",
			setup:       func(s *SetStructure) {},
			key:         "destination",
			membersRem:  []string{"key1", "key2"},
			membersWant: []string{},
			wantErr:     _const.ErrKeyNotFound,
		},
		{
			name: "test when remove one key",
			setup: func(s *SetStructure) {
				_ = s.SAdds("destination", "key1", "key2")
			},
			key:         "destination",
			membersRem:  []string{"key1"},
			membersWant: []string{"key2"},
			wantErr:     nil,
		},
		{
			name: "test remove all keys",
			setup: func(s *SetStructure) {
				_ = s.SAdds("destination", "key3", "key4")
			},
			key:         "destination",
			membersRem:  []string{"key3", "key4"},
			membersWant: []string{},
			wantErr:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := initTestSetDb()
			tt.setup(s)
			for _, s2 := range tt.membersRem {
				// Call the method with test case parameters.
				err := s.SRem(tt.key, s2)
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("SRems() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			if tt.wantErr == nil {
				// positive verify
				assert.True(t, s.exists(tt.key, tt.membersWant...))
				// verify removal
				assert.False(t, s.exists(tt.key, tt.membersRem...))
			}
		})
	}
}

func TestSCard(t *testing.T) {
	// Creating a list of cases for testing.
	tests := []struct {
		name    string
		setup   func(*SetStructure)
		key     string
		want    int
		wantErr error
	}{
		{
			name:    "test when key empty",
			setup:   func(s *SetStructure) {},
			key:     "",
			want:    -1,
			wantErr: _const.ErrKeyIsEmpty,
		},
		{
			name: "test two members",
			setup: func(s *SetStructure) {
				_ = s.SAdds("destination", "mem1", "mem2")
			},
			key:     "destination",
			want:    2,
			wantErr: nil,
		},
		{
			name: "test  three members one duplicate",
			setup: func(s *SetStructure) {
				_ = s.SAdds("destination", "mem1", "mem2", "mem1")
			},
			key:     "destination",
			want:    2,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := initTestSetDb()
			tt.setup(s)
			// Call the method with test case parameters.
			got, err := s.SCard(tt.key)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSMembers(t *testing.T) {
	// Creating a list of cases for testing.
	tests := []struct {
		name    string
		setup   func(*SetStructure)
		key     string
		want    []string
		wantErr error
	}{
		{
			name:    "test when key empty",
			setup:   func(s *SetStructure) {},
			key:     "",
			want:    nil,
			wantErr: _const.ErrKeyIsEmpty,
		},
		{
			name: "test two members",
			setup: func(s *SetStructure) {
				_ = s.SAdds("destination", "mem1", "mem2")
			},
			key:     "destination",
			want:    []string{"mem1", "mem2"},
			wantErr: nil,
		},
		{
			name: "test  three members one duplicate",
			setup: func(s *SetStructure) {
				_ = s.SAdds("destination", "mem1", "mem2", "mem1")
			},
			key:     "destination",
			want:    []string{"mem1", "mem2"},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := initTestSetDb()
			tt.setup(s)
			// Call the method with test case parameters.
			got, err := s.SMembers(tt.key)
			sort.Strings(got)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}

}

func TestSIsMember(t *testing.T) {
	// Creating a list of cases for testing.
	tests := []struct {
		name    string
		setup   func(*SetStructure)
		key     string
		member  string
		want    bool
		wantErr error
	}{
		{
			name:    "test when key empty",
			setup:   func(s *SetStructure) {},
			key:     "",
			want:    false,
			wantErr: _const.ErrKeyIsEmpty,
		},
		{
			name: "test two members, is member",
			setup: func(s *SetStructure) {
				_ = s.SAdds("destination", "mem1", "mem2")
			},
			key:     "destination",
			member:  "mem2",
			want:    true,
			wantErr: nil,
		},
		{
			name: "test  three members, not member",
			setup: func(s *SetStructure) {
				_ = s.SAdds("destination", "mem1", "mem2", "mem1")
			},
			key:     "destination",
			member:  "mem3",
			want:    false,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := initTestSetDb()
			tt.setup(s)
			// Call the method with test case parameters.
			got, err := s.SIsMember(tt.key, tt.member)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSUnion(t *testing.T) {
	// Creating a list of cases for testing.
	tests := []struct {
		name    string
		setup   func(*SetStructure)
		keys    []string
		want    []string
		wantErr error
	}{
		{
			name:    "test when key empty",
			setup:   func(s *SetStructure) {},
			keys:    []string{""},
			want:    nil,
			wantErr: _const.ErrKeyIsEmpty,
		},
		{
			name:    "test when no keys",
			setup:   func(s *SetStructure) {},
			keys:    []string{},
			want:    nil,
			wantErr: nil,
		},
		{
			name: "test two members",
			setup: func(s *SetStructure) {
				_ = s.SAdds("key1", "mem1", "mem2")
				_ = s.SAdds("key2", "mem1", "mem2")
			},
			keys:    []string{"key1", "key2"},
			want:    []string{"mem1", "mem2"},
			wantErr: nil,
		},
		{
			name: "test  three members ",
			setup: func(s *SetStructure) {
				_ = s.SAdds("key1", "mem2", "mem1")
				_ = s.SAdds("key2", "mem3")
			},
			keys:    []string{"key1", "key2"},
			want:    []string{"mem1", "mem2", "mem3"},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := initTestSetDb()
			tt.setup(s)
			got, err := s.SUnion(tt.keys...)
			assert.Equal(t, tt.wantErr, err)
			sort.Strings(got)
			assert.True(t, reflect.DeepEqual(tt.want, got))

		})
	}
}

func TestSInter(t *testing.T) {
	// Creating a list of cases for testing.
	tests := []struct {
		name    string
		setup   func(*SetStructure)
		keys    []string
		want    []string
		wantErr error
	}{
		{
			name:    "test when key empty",
			setup:   func(s *SetStructure) {},
			keys:    []string{""},
			want:    nil,
			wantErr: _const.ErrKeyIsEmpty,
		},
		{
			name:    "test when no keys",
			setup:   func(s *SetStructure) {},
			keys:    []string{},
			want:    nil,
			wantErr: nil,
		},
		{
			name:    "test when keys not found",
			setup:   func(s *SetStructure) {},
			keys:    []string{"notfound", "notfound2"},
			want:    nil,
			wantErr: _const.ErrKeyNotFound,
		},
		{
			name: "test when first keys found but second not found",
			setup: func(s *SetStructure) {
				_ = s.SAdds("key1", "mem1")
			},
			keys:    []string{"key1", "notfound"},
			want:    nil,
			wantErr: _const.ErrKeyNotFound,
		},
		{
			name: "test two members",
			setup: func(s *SetStructure) {
				_ = s.SAdds("key1", "mem1", "mem2")
				_ = s.SAdds("key2", "mem1", "mem2")
			},
			keys:    []string{"key1", "key2"},
			want:    []string{"mem1", "mem2"},
			wantErr: nil,
		},
		{
			name: "test three members with no intersect",
			setup: func(s *SetStructure) {
				_ = s.SAdds("key1", "mem2", "mem1")
				_ = s.SAdds("key2", "mem3")
			},
			keys:    []string{"key1", "key2"},
			want:    nil,
			wantErr: nil,
		},
		{
			name: "test three keys  with one intersect",
			setup: func(s *SetStructure) {
				_ = s.SAdds("key1", "mem2", "mem1", "mem3", "mem4")
				_ = s.SAdds("key2", "mem3")
				_ = s.SAdds("key3", "mem5", "mem3", "mem6")
			},
			keys:    []string{"key1", "key2", "key3"},
			want:    []string{"mem3"},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := initTestSetDb()
			tt.setup(s)
			got, err := s.SInter(tt.keys...)
			assert.Equal(t, tt.wantErr, err)
			sort.Strings(got)
			assert.True(t, reflect.DeepEqual(tt.want, got))

		})
	}
}

func TestSDiff(t *testing.T) {
	// Creating a list of cases for testing.
	tests := []struct {
		name    string
		setup   func(*SetStructure)
		keys    []string
		want    []string
		wantErr error
	}{
		{
			name:    "test when key empty",
			setup:   func(s *SetStructure) {},
			keys:    []string{""},
			want:    nil,
			wantErr: _const.ErrKeyIsEmpty,
		},
		{
			name:    "test when no keys",
			setup:   func(s *SetStructure) {},
			keys:    nil,
			want:    nil,
			wantErr: nil,
		},
		{
			name:    "test when keys not found",
			setup:   func(s *SetStructure) {},
			keys:    []string{"notfound", "notfound2"},
			want:    nil,
			wantErr: _const.ErrKeyNotFound,
		},
		{
			name: "test when first keys found but second not found",
			setup: func(s *SetStructure) {
				_ = s.SAdds("key1", "mem1")
			},
			keys:    []string{"key1", "notfound"},
			want:    nil,
			wantErr: _const.ErrKeyNotFound,
		},
		{
			name: "test two members no diff",
			setup: func(s *SetStructure) {
				_ = s.SAdds("key1", "mem1", "mem2")
				_ = s.SAdds("key2", "mem1", "mem2")
			},
			keys:    []string{"key1", "key2"},
			want:    nil,
			wantErr: nil,
		},
		{
			name: "test three members with two diffs",
			setup: func(s *SetStructure) {
				_ = s.SAdds("key1", "mem2", "mem1")
				_ = s.SAdds("key2", "mem3")
			},
			keys:    []string{"key1", "key2"},
			want:    []string{"mem2", "mem1"},
			wantErr: nil,
		},
		{
			name: "test three keys  with five three",
			setup: func(s *SetStructure) {
				_ = s.SAdds("key1", "mem2", "mem1", "mem3", "mem4")
				_ = s.SAdds("key2", "mem3")
				_ = s.SAdds("key3", "mem5", "mem3", "mem6")
			},
			keys:    []string{"key1", "key2", "key3"},
			want:    []string{"mem1", "mem2", "mem4"},
			wantErr: nil,
		},
		{
			name: "test three keys  with no diffs",
			setup: func(s *SetStructure) {
				_ = s.SAdds("key1", "mem2", "mem1", "mem3", "mem4")
				_ = s.SAdds("key2", "mem3")
				_ = s.SAdds("key3", "mem5", "mem3", "mem6")
			},
			keys:    []string{"key2", "key1", "key3"}, // the first key here will be `key2`
			want:    nil,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := initTestSetDb()
			tt.setup(s)
			got, err := s.SDiff(tt.keys...)
			assert.Equal(t, tt.wantErr, err)
			sort.Strings(got)
			sort.Strings(tt.want)
			assert.EqualValues(t, tt.want, got)

		})
	}
}

func TestSInterStore(t *testing.T) {
	// Creating a list of cases for testing.
	tests := []struct {
		name    string
		setup   func(*SetStructure)
		keys    []string
		want    []string
		wantErr error
	}{
		{
			name:    "test when key empty",
			setup:   func(s *SetStructure) {},
			keys:    []string{""},
			want:    nil,
			wantErr: _const.ErrKeyIsEmpty,
		},
		{
			name:    "test when no keys",
			setup:   func(s *SetStructure) {},
			keys:    nil,
			want:    nil,
			wantErr: nil,
		},
		{
			name:    "test when keys not found",
			setup:   func(s *SetStructure) {},
			keys:    []string{"notfound", "notfound2"},
			want:    nil,
			wantErr: _const.ErrKeyNotFound,
		},
		{
			name: "test when first keys found but second not found",
			setup: func(s *SetStructure) {
				_ = s.SAdds("key1", "mem1")
			},
			keys:    []string{"key1", "notfound"},
			want:    nil,
			wantErr: _const.ErrKeyNotFound,
		},
		{
			name: "test two members no two",
			setup: func(s *SetStructure) {
				_ = s.SAdds("key1", "mem1", "mem2")
				_ = s.SAdds("key2", "mem1", "mem2")
			},
			keys:    []string{"key1", "key2"},
			want:    []string{"mem1", "mem2"},
			wantErr: nil,
		},
		{
			name: "test three members with no intersect",
			setup: func(s *SetStructure) {
				_ = s.SAdds("key1", "mem2", "mem1")
				_ = s.SAdds("key2", "mem3")
			},
			keys:    []string{"key1", "key2"},
			want:    nil,
			wantErr: nil,
		},
		{
			name: "test three keys  with one intersect",
			setup: func(s *SetStructure) {
				_ = s.SAdds("key1", "mem2", "mem1", "mem3", "mem4")
				_ = s.SAdds("key2", "mem3")
				_ = s.SAdds("key3", "mem5", "mem3", "mem6")
			},
			keys:    []string{"key1", "key2", "key3"},
			want:    []string{"mem3"},
			wantErr: nil,
		},
		{
			name: "test three keys with two intersect",
			setup: func(s *SetStructure) {
				_ = s.SAdds("key1", "mem2", "mem5", "mem3", "mem4")
				_ = s.SAdds("key2", "mem3", "mem5")
				_ = s.SAdds("key3", "mem5", "mem3", "mem6")
			},
			keys:    []string{"key2", "key1", "key3"}, // the first key here will be `key2`
			want:    []string{"mem3", "mem5"},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := initTestSetDb()
			tt.setup(s)
			destKey := "dest"
			err := s.SInterStore(destKey, tt.keys...)
			assert.Equal(t, tt.wantErr, err)
			got, _ := s.SUnion(destKey) // get values stored in `destKey`
			sort.Strings(tt.want)
			sort.Strings(got)
			assert.EqualValues(t, tt.want, got)

		})
	}
}

func TestSUnionStore(t *testing.T) {
	// Creating a list of cases for testing.
	tests := []struct {
		name    string
		setup   func(*SetStructure)
		keys    []string
		want    []string
		wantErr error
	}{
		{
			name:    "test when key empty",
			setup:   func(s *SetStructure) {},
			keys:    []string{""},
			want:    nil,
			wantErr: _const.ErrKeyIsEmpty,
		},
		{
			name:    "test when no keys",
			setup:   func(s *SetStructure) {},
			keys:    nil,
			want:    nil,
			wantErr: nil,
		},
		{
			name:    "test when keys not found",
			setup:   func(s *SetStructure) {},
			keys:    []string{"notfound", "notfound2"},
			want:    nil,
			wantErr: _const.ErrKeyNotFound,
		},
		{
			name: "test when first keys found but second not found",
			setup: func(s *SetStructure) {
				_ = s.SAdds("key1", "mem1")
			},
			keys:    []string{"key1", "notfound"},
			want:    nil,
			wantErr: _const.ErrKeyNotFound,
		},
		{
			name: "test two members no diff",
			setup: func(s *SetStructure) {
				_ = s.SAdds("key1", "mem1", "mem2")
				_ = s.SAdds("key2", "mem1", "mem2")
			},
			keys:    []string{"key1", "key2"},
			want:    []string{"mem1", "mem2"},
			wantErr: nil,
		},
		{
			name: "test three members with two diffs",
			setup: func(s *SetStructure) {
				_ = s.SAdds("key1", "mem2", "mem1")
				_ = s.SAdds("key2", "mem3")
			},
			keys:    []string{"key1", "key2"},
			want:    []string{"mem2", "mem1", "mem3"},
			wantErr: nil,
		},
		{
			name: "test three keys  with five three",
			setup: func(s *SetStructure) {
				_ = s.SAdds("key1", "mem2")
				_ = s.SAdds("key2", "mem3")
				_ = s.SAdds("key3", "mem5", "mem3", "mem6")
			},
			keys:    []string{"key1", "key2", "key3"},
			want:    []string{"mem2", "mem3", "mem5", "mem6"},
			wantErr: nil,
		},
		{
			name: "test three keys  with no diffs",
			setup: func(s *SetStructure) {
				_ = s.SAdds("key1", "mem2", "mem1", "mem3", "mem4")
				_ = s.SAdds("key2", "mem3")
				_ = s.SAdds("key3", "mem5", "mem3", "mem6")
			},
			keys:    []string{"key2", "key1", "key3"}, // the first key here will be `key2`
			want:    []string{"mem1", "mem2", "mem3", "mem4", "mem5", "mem6"},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := initTestSetDb()
			tt.setup(s)
			destKey := "dest"
			err := s.SUnionStore(destKey, tt.keys...)
			assert.Equal(t, tt.wantErr, err)
			got, _ := s.SUnion(destKey) // get values stored in `destKey`
			sort.Strings(tt.want)
			sort.Strings(got)
			assert.EqualValues(t, tt.want, got)

		})
	}
}

func TestNewSetStructure(t *testing.T) {
	// expect error
	opts := config.DefaultOptions
	opts.DirPath = "" // the cause of error
	setDB, err := NewSetStructure(opts)
	assert.NotNil(t, err)
	assert.Nil(t, setDB)
	// expect no error
	opts = config.DefaultOptions
	dir, _ := os.MkdirTemp("", "TestSetStructure")
	opts.DirPath = dir
	setDB, _ = NewSetStructure(opts)
	assert.NotNil(t, setDB)
	assert.IsType(t, &engine.DB{}, setDB.db)
}

func TestSetStructure_exists(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*SetStructure)
		key     string
		members []string
		found   bool
	}{
		{
			name:    "test when empty Set",
			setup:   func(s *SetStructure) {},
			key:     "testKey",
			members: []string{"key1", "key2"},
			found:   false,
		},
		{
			name: "test Set with no members  found",
			setup: func(s *SetStructure) {
				_ = s.SAdds("testKey", "non1", "non2")
			},
			key:     "testKey",
			members: []string{"key1", "key2"},
			found:   false,
		},
		{
			name: "test Set with partial members found",
			setup: func(s *SetStructure) {
				_ = s.SAdds("testKey", "key1", "non2")
			},
			key:     "testKey",
			members: []string{"key1", "key2"},
			found:   false,
		},
		{
			name: "test Set with all members found",
			setup: func(s *SetStructure) {
				_ = s.SAdds("testKey", "key1", "non2", "key2")
			},
			key:     "testKey",
			members: []string{"key1", "key2"},
			found:   true,
		},
		{
			name: "test Set key empty",
			setup: func(s *SetStructure) {
				_ = s.SAdds("testKey", "key1", "non2", "key2")
			},
			key:     "",
			members: []string{"key1", "key2"},
			found:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := initTestSetDb()
			tt.setup(s)
			if tt.found {
				assert.True(t, s.exists(tt.key, tt.members...))
			} else {
				assert.False(t, s.exists(tt.key, tt.members...))
			}
		})
	}
}

func TestFSetsExists(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(sets *FSets)
		key     string
		members []string
		found   bool
	}{
		{
			name:    "test when empty Set",
			setup:   func(s *FSets) {},
			key:     "testKey",
			members: []string{"key1", "key2"},
			found:   false,
		},
		{
			name: "test Set with no members  found",
			setup: func(fs *FSets) {
				fs.add("non1", "non2")
			},
			key:     "testKey",
			members: []string{"key1", "key2"},
			found:   false,
		},
		{
			name: "test Set with partial members found",
			setup: func(fs *FSets) {
				fs.add("key1", "non2")
			},
			key:     "testKey",
			members: []string{"key1", "key2"},
			found:   false,
		},
		{
			name: "test Set with all members found",
			setup: func(fs *FSets) {
				fs.add("key1", "non2", "key2")
			},
			key:     "testKey",
			members: []string{"key1", "key2"},
			found:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFSets := &FSets{}
			tt.setup(testFSets)

			if tt.found {
				assert.True(t, testFSets.exists(tt.members...))
			} else {
				assert.False(t, testFSets.exists(tt.members...))
			}
		})
	}
}

func TestFSets_remove(t *testing.T) {
	tests := []struct {
		name       string
		setup      *FSets
		members    []string
		want       []string
		wantErrors error
	}{
		{
			name:       "test when empty Set",
			setup:      &FSets{},
			members:    []string{"key1", "key2"},
			want:       []string{},
			wantErrors: ErrMemberNotFound,
		},
		{
			name:       "test Set with wrong members ",
			setup:      &FSets{"non1": struct{}{}, "non2": struct{}{}},
			want:       []string{"non1", "non2"},
			members:    []string{"key1", "key2"},
			wantErrors: ErrMemberNotFound,
		},
		{
			name:       "test Set with partial members found", // this should not delete the found member neither
			setup:      &FSets{"key1": struct{}{}, "non2": struct{}{}},
			want:       []string{"key1", "non2"},
			members:    []string{"key1", "key2"},
			wantErrors: ErrMemberNotFound,
		},
		{
			name:    "test Set with all members found",
			setup:   &FSets{"key1": struct{}{}, "non2": struct{}{}, "key2": struct{}{}},
			want:    []string{"non2"},
			members: []string{"key1", "key2"},
		},
		{
			name:       "test Set not initialized",
			want:       []string{},
			members:    []string{"key1", "key2"},
			wantErrors: ErrSetNotInitialized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := tt.setup.remove(tt.members...)
			assert.Equal(t, tt.wantErrors, err)
			assert.True(t, tt.setup.exists(tt.want...))
		})
	}
}

func TestSetStructure_Keys(t *testing.T) {
	set, _ := initTestSetDb()
	defer set.db.Clean()

	err = set.SAdd("testKey11", "non1")
	assert.Nil(t, err)

	err = set.SAdd("testKey12", "non1")
	assert.Nil(t, err)

	err = set.SAdd("testKey3", "non1")
	assert.Nil(t, err)

	err = set.SAdd("testKey4", "non1")
	assert.Nil(t, err)

	keys, err := set.Keys("*")
	assert.Nil(t, err)
	assert.Equal(t, 4, len(keys))

	keys, err = set.Keys("testKey*")
	assert.Nil(t, err)
	assert.Equal(t, 4, len(keys))

	keys, err = set.Keys("testKey1*")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(keys))

	keys, err = set.Keys("testKey?*")
	assert.Nil(t, err)
	assert.Equal(t, 4, len(keys))

	keys, err = set.Keys("testKey111?*")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(keys))

}
