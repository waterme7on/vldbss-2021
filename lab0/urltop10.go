package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

func URLTop10(nWorkers int) RoundsArgs {
	var args RoundsArgs
	// round 1: do url count
	args = append(args, RoundArgs{
		MapFunc:    URLCountMap,
		ReduceFunc: URLCountReduce,
		NReduce:    nWorkers,
	})
	// round 2: sort and get the 10 most frequent URLs
	args = append(args, RoundArgs{
		MapFunc:    URLTop10Map,
		ReduceFunc: URLTop10Reduce,
		NReduce:    1,
	})
	return args
}

func URLCountMap(filename string, contents string) []KeyValue {
	lines := strings.Split(contents, "\n")
	res := make([]KeyValue, 0, len(lines))
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if len(l) == 0 {
			continue
		}
		res = append(res, KeyValue{Key: l})
	}
	return res
}

func URLCountReduce(key string, values []string) string {
	return fmt.Sprintf("%s %s\n", key, strconv.Itoa(len(values)))
}

func URLTop10Map(filename string, contents string) []KeyValue {
	lines := strings.Split(contents, "\n")
	kvs := make([]KeyValue, 0, 10)
	cntMap := make(map[string]int, len(lines))
	for _, l := range lines {
		if len(l) == 0 {
			continue
		}
		words := strings.Split(l, " ")
		cntMap[words[0]], _ = strconv.Atoi(words[1])
	}
	us, cs := TopN(cntMap, 10)
	for i := range us {
		kvs = append(kvs, KeyValue{"", fmt.Sprintf("%s %d", us[i], cs[i])})
	}
	return kvs
}

func URLTop10Reduce(key string, values []string) string {
	cnts := make(map[string]int, len(values))
	for _, v := range values {
		v := strings.TrimSpace(v)
		if len(v) == 0 {
			continue
		}
		tmp := strings.Split(v, " ")
		n, err := strconv.Atoi(tmp[1])
		if err != nil {
			panic(err)
		}
		cnts[tmp[0]] = n
	}

	us, cs := TopN(cnts, 10)
	buf := new(bytes.Buffer)
	for i := range us {
		fmt.Fprintf(buf, "%s: %d\n", us[i], cs[i])
	}
	return buf.String()
}
