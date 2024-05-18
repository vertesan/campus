package master

import (
	"io"
	"vertesan/campus/proto/mastertag"
)

func UnmarshalPlain(reader io.Reader) (*mastertag.MasterGetResponse, error) {
  masterGetResp := &mastertag.MasterGetResponse{}
  return masterGetResp, nil
}
