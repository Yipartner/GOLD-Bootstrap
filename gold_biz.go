package gold

import (
	"github.com/MarkLux/GOLD/serving/common"
	"log"
)

func (s *GoldService) OnInit() {
	log.Println("inited")
}

// the biz function
func (s *GoldService) OnHandle(req *common.GoldRequest, rsp *common.GoldResponse) error {
	// get data from request
	userName := req.Data["name"].(string)
	log.Println("userName: " + userName)

	greeting := "hello, " + userName
	rsp.Data = make(map[string]interface{})
	rsp.Data["greeting"] = greeting

	return nil
}
// on error
func (s *GoldService) OnError(err error) bool {
	log.Println("error!", err)
	return false
}
