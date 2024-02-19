package entity

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

// func (i *UpdateListInput) Validate() error {
// 	if i.Title == nil && i.Description == nil {
// 		return errors.New("update structure has no values")
// 	}

// 	return nil
// }

// type UpdateItemInput struct {
// 	Title       string    `json:"title"`
// 	Description string    `json:"description"`
// 	Start       time.Time `json:"start"`
// 	End         time.Time `json:"end"`
// 	Done        bool      `json:"done"`
// }

func TestValidate(t *testing.T) {
	input := UpdateItemInput{}
	json.Unmarshal([]byte(""), &input)
	tests := []struct {
		input  UpdateItemInput
		result error
	}{
		{input, errors.New("update structure has no values")},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.input)
		t.Run(testname, func(t *testing.T) {
			err := tt.input.Validate()
			if errors.Is(err, tt.result) {
				t.Errorf("got %d, want %d", err, tt.result)
			}
		})
	}
}
