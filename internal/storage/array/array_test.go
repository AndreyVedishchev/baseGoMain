package array

import (
	"reflect"
	"strconv"
	"testing"

	"base-go/internal/models"
)

func TestArrayStorage_Clear(t *testing.T) {
	var exp [maxSize]*models.Resume
	tests := []struct {
		name string
		as   *ArrayStorage
	}{
		{
			name: "empty",
			as:   &ArrayStorage{resumes: [maxSize]*models.Resume{}},
		},
		{
			name: "partially filled",
			as:   &ArrayStorage{resumes: [maxSize]*models.Resume{{UUID: "1"}, {UUID: "2"}, {UUID: "3"}}},
		},
		{
			name: "full",
			as:   makeFullStorage(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			as := tt.as

			as.Clear()

			if as.resumes != exp {
				t.Errorf("Clear() resumes = %v, want %v", as.resumes, exp)
			}
		})
	}
}

func TestArrayStorage_Delete(t *testing.T) {
	tests := []struct {
		name string
		as   *ArrayStorage
		uuid string
		exp  [maxSize]*models.Resume
	}{
		{
			name: "empty",
			as:   &ArrayStorage{resumes: [maxSize]*models.Resume{}},
			uuid: "123",
			exp:  [maxSize]*models.Resume{},
		},
		{
			name: "no such key",
			as:   &ArrayStorage{resumes: [maxSize]*models.Resume{{UUID: "1"}, {UUID: "2"}, {UUID: "3"}}},
			uuid: "123",
			exp:  [maxSize]*models.Resume{{UUID: "1"}, {UUID: "2"}, {UUID: "3"}},
		},
		{
			name: "ok",
			as:   &ArrayStorage{resumes: [maxSize]*models.Resume{{UUID: "1"}, {UUID: "2"}, {UUID: "3"}}},
			uuid: "1",
			exp:  [maxSize]*models.Resume{{UUID: "2"}, {UUID: "3"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			as := tt.as
			as.Delete(tt.uuid)
			if as.resumes != tt.exp {
				t.Errorf("Delete() resumes = %v, want %v", as.resumes, tt.exp)
			}
		})
	}
}

func TestArrayStorage_Get(t *testing.T) {
	tests := []struct {
		name string
		as   *ArrayStorage
		uuid string
		exp  *models.Resume
	}{
		{
			name: "empty",
			as:   &ArrayStorage{resumes: [maxSize]*models.Resume{}},
			uuid: "123",
			exp:  nil,
		},
		{
			name: "no such key",
			as:   &ArrayStorage{resumes: [maxSize]*models.Resume{{UUID: "1"}, {UUID: "2"}, {UUID: "3"}}},
			uuid: "123",
			exp:  nil,
		},
		{
			name: "ok",
			as:   &ArrayStorage{resumes: [maxSize]*models.Resume{{UUID: "1"}, {UUID: "2"}, {UUID: "3"}}},
			uuid: "1",
			exp:  &models.Resume{UUID: "1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			as := tt.as
			if got := as.Get(tt.uuid); !reflect.DeepEqual(got, tt.exp) {
				t.Errorf("Get() = %v, want %v", got, tt.exp)
			}
		})
	}
}

func TestArrayStorage_GetAll(t *testing.T) {
	tests := []struct {
		name string
		as   *ArrayStorage
		exp  []*models.Resume
	}{
		{
			name: "empty",
			as:   &ArrayStorage{resumes: [maxSize]*models.Resume{}},
			exp:  nil,
		},
		{
			name: "partially filled",
			as:   &ArrayStorage{resumes: [maxSize]*models.Resume{{UUID: "1"}, {UUID: "2"}, {UUID: "3"}}},
			exp:  []*models.Resume{{UUID: "1"}, {UUID: "2"}, {UUID: "3"}},
		},
		{
			name: "full",
			as:   makeFullStorage(),
			exp:  makeFullSlice(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			as := tt.as
			if got := as.GetAll(); !reflect.DeepEqual(got, tt.exp) {
				t.Errorf("GetAll() = %v, want %v", got, tt.exp)
			}
		})
	}
}

func TestArrayStorage_Save(t *testing.T) {
	tests := []struct {
		name string
		as   *ArrayStorage
		in   *models.Resume
		exp  [maxSize]*models.Resume
	}{
		{
			name: "empty",
			as:   &ArrayStorage{resumes: [maxSize]*models.Resume{}},
			in:   &models.Resume{UUID: "1"},
			exp:  [maxSize]*models.Resume{{UUID: "1"}},
		},
		{
			name: "partially filled",
			as:   &ArrayStorage{resumes: [maxSize]*models.Resume{{UUID: "1"}, {UUID: "2"}}},
			in:   &models.Resume{UUID: "3"},
			exp:  [maxSize]*models.Resume{{UUID: "1"}, {UUID: "2"}, {UUID: "3"}},
		},
		{
			name: "full",
			as:   makeFullStorage(),
			exp:  makeFullStorage().resumes,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			as := tt.as
			as.Save(tt.in)

			if as.resumes != tt.exp {
				t.Errorf("Save(r *models.Resume) resumes = %v, want %v", as.resumes, tt.exp)
			}
		})
	}
}

func TestArrayStorage_Size(t *testing.T) {
	tests := []struct {
		name string
		as   *ArrayStorage
		want int
	}{
		{
			name: "empty",
			as:   &ArrayStorage{resumes: [maxSize]*models.Resume{}},
			want: 0,
		},
		{
			name: "partially filled",
			as:   &ArrayStorage{resumes: [maxSize]*models.Resume{{UUID: "1"}, {UUID: "2"}, {UUID: "3"}}},
			want: 3,
		},
		{
			name: "full",
			as:   makeFullStorage(),
			want: maxSize,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			as := tt.as
			if got := as.Size(); got != tt.want {
				t.Errorf("Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

// makeFullStorage возвращает заполненное моковое хранилище
func makeFullStorage() *ArrayStorage {
	res := &ArrayStorage{
		resumes: [maxSize]*models.Resume{},
	}

	for i := 0; i < maxSize; i++ {
		res.resumes[i] = &models.Resume{UUID: strconv.Itoa(i + 1)}
	}

	return res
}

// makeFullSlice возвращает моковый заполненный слайс
func makeFullSlice() []*models.Resume {
	res := make([]*models.Resume, maxSize)
	for i := 0; i < maxSize; i++ {
		res[i] = &models.Resume{UUID: strconv.Itoa(i + 1)}
	}

	return res
}
