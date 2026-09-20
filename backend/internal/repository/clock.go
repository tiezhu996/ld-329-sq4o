package repository

import "time"

// nowFunc 集中时间来源，便于测试替换。
var nowFunc = time.Now
