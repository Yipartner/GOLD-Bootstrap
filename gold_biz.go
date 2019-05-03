package gold

import (
	"github.com/MarkLux/GOLD/serving/rpc/goldrpc"
	"github.com/MarkLux/GOLD/serving/wrapper/db"
	"log"
)

/**
  * function service example
  * show usage of rpc, db & cache
 */

// the model of user
// use annotation `bson` to control the key saved in db(for mongo).
type UserModel struct {
	Name string `bson:"name"`
	Sex  string `bson:"sex"`
	Mail string `bson:"mail"`
}

// would be invoked when the service launch at first time.
func (s *GoldService) OnInit() {
	// do something here..
}

// the biz function
func (s *GoldService) OnHandle(req *goldrpc.GoldRequest, rsp *goldrpc.GoldResponse) error {
	// get data from request
	userName := req.Data["name"].(string)
	log.Println("userName: " + userName)

	// cache example
	cacheKey := "prefix_" + userName
	u, err := s.CacheClient.Get(cacheKey)
	if err != nil {
		log.Println("fail to get info from cache, ", err)
	}
	// build response
	rsp.Data = make(map[string]interface{})

	useCache := true
	// if got nothing from cache, then query the db.
	if u == nil {
		useCache = false
		// db session example
		dbSession, err := s.DbFactory.NewDataBaseSession("test", "user", "tst", "123")

		if err != nil {
			log.Println("create db session failed, ", err)
			return err
		}
		defer dbSession.Close()
		// db query example
		param := make(map[string]string)
		param["data.name"] = userName
		qUsers, err := dbSession.Query(db.GoldDBQuery{Param: param})
		if err != nil {
			log.Println("fail to query db, ", err)
			return err
		}
		log.Printf("qUsers: %v\n", qUsers)
		if len(qUsers) > 0 {
			u = qUsers[0]
			// reset the cache
			err = s.CacheClient.Set(cacheKey, u, 300 * 1000)
			if err != nil {
				log.Println("fail to reset cache, ", err)
			}
		} else {
			rsp.Data["err"] = "no data from database, insert one instead."
		}
		u = nil
	}

	if u != nil {
		rsp.Data["userModel"] = u
		rsp.Data["useCache"] = useCache
	}

	return nil
}

// when an error is throwed, this function would be called
// if you don't want to break the process, return true
// otherwise you should return false
func (s *GoldService) OnError(err error) bool {
	log.Println(err)
	return false
}