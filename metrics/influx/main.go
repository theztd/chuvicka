package influx

var Url string
var Token string

type MetricResult struct {
	Name       string  // _field:DNSLookup
	StatusCode int     // StatusCode:200
	Url        string  // url:http://troll.fejk.net/v1/v1/pomala_url]
	Result     string  // result:mean
	Time       string  // _time:2023-02-01 19:49:09.014048 +0000 UTC
	Value      float64 // _value:4.1759832e+07
}

type Tags struct {
	Key   string
	Value string
}

type Fields struct {
	Key   string
	Value float32
}

type Metric struct {
	Name   string
	Tags   []Tags
	Fields []Fields
}

func getBucket() {

}

/*

- createBucket()
- deleteBucket()
- createBucketAccess()
- deleteBucketAccess()

*/
