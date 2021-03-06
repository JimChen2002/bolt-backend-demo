package consts

import "time"

//TODO: (low priority)move some of these consts to config file
const PageSize = 30
const MsgPageSize = 10
const SystemLoadThreshold float64 = 5.0
const TokenExpireDays = 1
const WanderPageSize = 15
const MaxPage = 150
const SearchPageSize = 30
const SearchMaxPage = 100
const SearchMaxLength = 30
const PostMaxLength = 10000
const ImageMaxWidth = 10000
const ImageMaxHeight = 10000
const VoteOptionMaxCharacters = 15
const VoteMaxOptions = 4
const MaxDevicesPerUser = 6
const ReportMaxLength = 1000
const ImgMaxLength = 2000000
const Base64Rate = 1.33333333
const AesIv = "12345678901234567890123456789012"

const PushApiLogFile = "push.log"
const ServicesApiLogFile = "services-api.log"
const DetailLogFile = "detail-log.log"
const SecurityApiLogFile = "security-api.log"

const DatabaseReadFailedString = "Database Read Failed"
const DatabaseWriteFailedString = "Database Write Failed"
const DatabaseDamagedString = "Database Damaged"
const DatabaseEncryptFailedString = "Database Encrypt Failed"

const DzName = "Host"
const ExtraNamePrefix = "You Win "

var TimeLoc, _ = time.LoadLocation("America/New_York")

var Names0 = []string{
	"Angry",
	"Baby",
	"Crazy",
	"Diligent",
	"Excited",
	"Fat",
	"Greedy",
	"Hungry",
	"Interesting",
	"Jolly",
	"Kind",
	"Little",
	"Magic",
	"Naïve",
	"Old",
	"PKU",
	"Quiet",
	"Rich",
	"Superman",
	"Tough",
	"Undefined",
	"Valuable",
	"Wifeless",
	"Xiangbuchulai",
	"Young",
	"Zombie",
}

var Names1 = []string{
	"Alice",
	"Bob",
	"Carol",
	"Dave",
	"Eve",
	"Francis",
	"Grace",
	"Hans",
	"Isabella",
	"Jason",
	"Kate",
	"Louis",
	"Margaret",
	"Nathan",
	"Olivia",
	"Paul",
	"Queen",
	"Richard",
	"Susan",
	"Thomas",
	"Uma",
	"Vivian",
	"Winnie",
	"Xander",
	"Yasmine",
	"Zach",
}
