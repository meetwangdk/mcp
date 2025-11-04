package data_access

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"code.alipay.com/sigma/eavesdropping-client/go-client/common"
	"code.alipay.com/sigma/eavesdropping-client/go-client/metrics"
	"code.alipay.com/sigma/eavesdropping-client/go-client/model"
	clientUtils "code.alipay.com/sigma/eavesdropping-client/go-client/pkg/utils"
	utilsh "code.alipay.com/sigma/eavesdropping-client/go-client/pkg/utils"
	"code.alipay.com/sigma/eavesdropping-client/go-client/utils"
	"github.com/olivere/elastic"
	"k8s.io/klog"
)

// 2. 定义一个 StorageSqlImpl struct, 该struct 包含了存储client
type StorageEsImpl struct {
	DB *elastic.Client
}

const (
	EavesdroppingLatencyTag = "LastReadTime"
	LatencyTag              = "Latency"
	LogicPool               = "mandatory.k8s.alipay.com/app-logic-pool.keyword"
	Keyword                 = ".keyword"
)

const (
	DELIVERY_ENV_TEST   = "test"
	DELIVERY_ENV_PROD   = "prod"
	DELIVERY_ENV_DRYRUN = "dryrun"
	PodResource         = "Pod"
	NodeResource        = "Node"
	podYamlIndexName    = "pod_yaml"
	podYamlTypeName     = "_doc"
	nodeYamlIndexName   = "node_yaml"
	nodeYamlTypeName    = "_doc"
)

// 4. 提供一个 ProvideSqlStorate 方法, 传入一个 MysqlOptions, 返回一个 StorageInterface 和 error
func ProvideEsStorage(conf *common.ESOptions) (StorageInterface, error) {
	var err error

	esClient, err := elastic.NewClient(
		elastic.SetURL(conf.EndPoint), elastic.SetBasicAuth(conf.Username, conf.Password), elastic.SetSniff(false))
	if err != nil {
		// panic(err)
		klog.Errorf("init es client error %s", err)
		return nil, err
	}
	return &StorageEsImpl{
		DB: esClient,
	}, nil
}

func (s *StorageEsImpl) QuerySpanWithPodUid(data interface{}, uid string, opts ...model.OptionFunc) error {
	options := &model.Options{
		Limit:     80,
		Ascending: false,
		OrderBy:   "Elapsed",
	}
	for _, do := range opts {
		do(options)
	}
	if uid == "" {
		return fmt.Errorf("the params is error, uid is nil")
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}

	query := elastic.NewBoolQuery().Must(elastic.NewQueryStringQuery(fmt.Sprintf("OwnerRef.UID: \"%s\" AND OwnerRef.Resource: pods", uid)))
	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(options.Limit).
		Sort(options.OrderBy, options.Ascending).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error%v", err)
	}
	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}
	return nil

}
func (s *StorageEsImpl) QueryLifePhaseWithPodUid(data interface{}, uid string, opts ...model.OptionFunc) error {
	options := &model.Options{
		Limit:     200,
		Ascending: false,
		OrderBy:   "startTime",
	}
	for _, do := range opts {
		do(options)
	}
	if uid == "" {
		return fmt.Errorf("the params is error, uid is nil")
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}

	begin := time.Now()
	defer func() {
		cost := utils.TimeSinceInMilliSeconds(begin)
		metrics.QueryMethodDurationMilliSeconds.WithLabelValues("QueryLifePhase").Observe(cost)
	}()

	stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("podUID: \"%s\"", uid))
	query := elastic.NewBoolQuery().Must(stringQuery)

	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(options.Limit).
		Sort(options.OrderBy, options.Ascending).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error%v", err)
	}
	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *StorageEsImpl) QueryPodYamlsWithPodUID(data interface{}, uid string, opts ...model.OptionFunc) error {
	options := &model.Options{
		Limit:     5,
		Ascending: false,
		OrderBy:   "stageTimestamp",
	}
	for _, do := range opts {
		do(options)
	}

	if uid == "" {
		return fmt.Errorf("the params is error, uid is nil")
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}

	begin := time.Now()
	defer func() {
		cost := utils.TimeSinceInMilliSeconds(begin)
		metrics.QueryMethodDurationMilliSeconds.WithLabelValues("QueryPodYaml").Observe(cost)
	}()

	stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("podUID.keyword: \"%s\"", uid))
	query := elastic.NewBoolQuery().Must(stringQuery)
	if !options.From.IsZero() || !options.To.IsZero() {
		rangeQuery := elastic.NewRangeQuery("stageTimestamp").TimeZone("UTC")
		if !options.From.IsZero() {
			rangeQuery.From(options.From)
		}
		if !options.To.IsZero() {
			rangeQuery.To(options.To)
		}
		query.Must(rangeQuery)
	}
	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(options.Limit).
		Sort(options.OrderBy, options.Ascending).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error%v", err)
	}
	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *StorageEsImpl) QueryPodYamlsV2WithPodUID(data interface{}, index, uid string, opts ...model.OptionFunc) error {
	options := &model.Options{
		Limit:     5,
		Ascending: false,
		OrderBy:   "stageTimestamp",
	}
	for _, do := range opts {
		do(options)
	}

	if uid == "" {
		return fmt.Errorf("the params is error, uid is nil")
	}
	_, _, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}

	begin := time.Now()
	defer func() {
		cost := utils.TimeSinceInMilliSeconds(begin)
		metrics.QueryMethodDurationMilliSeconds.WithLabelValues("QueryPodYaml").Observe(cost)
	}()

	stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("podUID.keyword: \"%s\"", uid))
	query := elastic.NewBoolQuery().Must(stringQuery)
	if !options.From.IsZero() || !options.To.IsZero() {
		rangeQuery := elastic.NewRangeQuery("stageTimestamp").TimeZone("UTC")
		if !options.From.IsZero() {
			rangeQuery.From(options.From)
		}
		if !options.To.IsZero() {
			rangeQuery.To(options.To)
		}
		query.Must(rangeQuery)
	}
	searchResult, err := s.DB.Search().Index(index).Type(esType).Query(query).Size(options.Limit).
		Sort(options.OrderBy, options.Ascending).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error%v", err)
	}
	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *StorageEsImpl) QueryPodYamlTimelineWithPodUID(uid string, opts ...model.OptionFunc) (timeline []string, err error) {
	options := &model.Options{
		Limit:     500,
		Ascending: false,
		OrderBy:   "stageTimestamp",
	}
	for _, do := range opts {
		do(options)
	}
	if uid == "" {
		return timeline, fmt.Errorf("the params is error, uid is nil")
	}
	begin := time.Now()
	defer func() {
		cost := utils.TimeSinceInMilliSeconds(begin)
		metrics.QueryMethodDurationMilliSeconds.WithLabelValues("QueryPodYamlTimeline").Observe(cost)
	}()
	stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("podUID.keyword: \"%s\"", uid))
	query := elastic.NewBoolQuery().Must(stringQuery)
	if !options.From.IsZero() || !options.To.IsZero() {
		rangeQuery := elastic.NewRangeQuery("stageTimestamp").TimeZone("UTC")
		if !options.From.IsZero() {
			rangeQuery.From(options.From)
		}
		if !options.To.IsZero() {
			rangeQuery.To(options.To)
		}
		query.Must(rangeQuery)
	}
	searchResult, err := s.DB.Search().Index(podYamlIndexName).Type(podYamlTypeName).Query(query).Size(options.Limit).
		Sort(options.OrderBy, options.Ascending).Do(context.Background())
	if err != nil {
		return timeline, fmt.Errorf("error %v", err)
	}
	for _, hit := range searchResult.Hits.Hits {
		timeline = append(timeline, hit.Index)
	}
	return
}

func (s *StorageEsImpl) QueryPodYamlsWithPodName(data interface{}, podName string, opts ...model.OptionFunc) error {
	options := &model.Options{
		Limit:     5,
		Ascending: false,
		OrderBy:   "stageTimestamp",
	}
	for _, do := range opts {
		do(options)
	}

	if podName == "" {
		return fmt.Errorf("the params is error, uid is nil")
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}

	begin := time.Now()
	defer func() {
		cost := utils.TimeSinceInMilliSeconds(begin)
		metrics.QueryMethodDurationMilliSeconds.WithLabelValues("QueryLifePhase").Observe(cost)
	}()

	stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("podName.keyword: \"%s\"", podName))
	query := elastic.NewBoolQuery().Must(stringQuery)
	if !options.From.IsZero() || !options.To.IsZero() {
		rangeQuery := elastic.NewRangeQuery("stageTimestamp").TimeZone("UTC")
		if !options.From.IsZero() {
			rangeQuery.From(options.From)
		}
		if !options.To.IsZero() {
			rangeQuery.To(options.To)
		}
		query.Must(rangeQuery)
	}
	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(options.Limit).
		Sort(options.OrderBy, options.Ascending).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error%v", err)
	}
	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *StorageEsImpl) QueryPodYamlsWithHostName(data interface{}, hostName string, opts ...model.OptionFunc) error {
	options := &model.Options{
		Limit:     5,
		Ascending: false,
		OrderBy:   "stageTimestamp",
	}
	for _, do := range opts {
		do(options)
	}

	if hostName == "" {
		return fmt.Errorf("the params is error, uid is nil")
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}

	begin := time.Now()
	defer func() {
		cost := utils.TimeSinceInMilliSeconds(begin)
		metrics.QueryMethodDurationMilliSeconds.WithLabelValues("QueryPodYamlsWithHostName").Observe(cost)
	}()

	stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("hostname.keyword: \"%s\"", hostName))
	query := elastic.NewBoolQuery().Must(stringQuery)
	if !options.From.IsZero() || !options.To.IsZero() {
		rangeQuery := elastic.NewRangeQuery("stageTimestamp").TimeZone("UTC")
		if !options.From.IsZero() {
			rangeQuery.From(options.From)
		}
		if !options.To.IsZero() {
			rangeQuery.To(options.To)
		}
		query.Must(rangeQuery)
	}
	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(options.Limit).
		Sort(options.OrderBy, options.Ascending).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error%v", err)
	}
	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *StorageEsImpl) QueryPodYamlsWithPodIp(data interface{}, podIp string, opts ...model.OptionFunc) error {
	options := &model.Options{
		Limit:     5,
		Ascending: false,
		OrderBy:   "stageTimestamp",
	}
	for _, do := range opts {
		do(options)
	}

	if podIp == "" {
		return fmt.Errorf("the params is error, uid is nil")
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}

	begin := time.Now()
	defer func() {
		cost := utils.TimeSinceInMilliSeconds(begin)
		metrics.QueryMethodDurationMilliSeconds.WithLabelValues("QueryPodYamlsWithPodIp").Observe(cost)
	}()

	stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("podIP.keyword: \"%s\"", podIp))
	query := elastic.NewBoolQuery().Must(stringQuery)
	if !options.From.IsZero() || !options.To.IsZero() {
		rangeQuery := elastic.NewRangeQuery("stageTimestamp").TimeZone("UTC")
		if !options.From.IsZero() {
			rangeQuery.From(options.From)
		}
		if !options.To.IsZero() {
			rangeQuery.To(options.To)
		}
		query.Must(rangeQuery)
	}
	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(options.Limit).
		Sort(options.OrderBy, options.Ascending).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error%v", err)
	}
	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *StorageEsImpl) QueryPodListWithNodeip(data interface{}, nodeIp string, isDeleted bool) error {

	if nodeIp == "" {
		return fmt.Errorf("the params is error, uid is nil")
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}
	begin := time.Now()
	defer func() {
		cost := utils.TimeSinceInMilliSeconds(begin)
		metrics.QueryMethodDurationMilliSeconds.WithLabelValues("QueryLifePhase").Observe(cost)
	}()

	stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("hostIP.keyword: \"%s\"", nodeIp))
	deleteFalse := elastic.NewQueryStringQuery(fmt.Sprintf("isDeleted.keyword: \"%t\"", isDeleted))
	query := elastic.NewBoolQuery().Must(stringQuery).Must(deleteFalse)
	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(300).
		Sort("stageTimestamp", false).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error%v", err)
	}
	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}
	return nil

}
func (s *StorageEsImpl) QueryPodUIDListByHostname(data interface{}, hostName string) error {

	if hostName == "" {
		return fmt.Errorf("the params is error, uid is nil")
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}

	begin := time.Now()
	defer func() {
		cost := utils.TimeSinceInMilliSeconds(begin)
		metrics.QueryMethodDurationMilliSeconds.WithLabelValues("QueryLifePhase").Observe(cost)
	}()

	stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("hostname.keyword: \"%s\"", hostName))
	query := elastic.NewBoolQuery().Must(stringQuery)
	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(300).
		Sort("stageTimestamp", false).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error%v", err)
	}
	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *StorageEsImpl) QueryPodUIDListByPodIP(data interface{}, podIp string) error {

	if podIp == "" {
		return fmt.Errorf("the params is error, uid is nil")
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}

	begin := time.Now()
	defer func() {
		cost := utils.TimeSinceInMilliSeconds(begin)
		metrics.QueryMethodDurationMilliSeconds.WithLabelValues("QueryLifePhase").Observe(cost)
	}()

	stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("podIP.keyword: \"%s\"", podIp))
	query := elastic.NewBoolQuery().Must(stringQuery)
	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(300).
		Sort("stageTimestamp", false).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error%v", err)
	}
	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *StorageEsImpl) QueryPodUIDListByPodName(data interface{}, podName string) error {

	if podName == "" {
		return fmt.Errorf("the params is error, uid is nil")
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}

	begin := time.Now()
	defer func() {
		cost := utils.TimeSinceInMilliSeconds(begin)
		metrics.QueryMethodDurationMilliSeconds.WithLabelValues("QueryLifePhase").Observe(cost)
	}()

	stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("podName.keyword: \"%s\"", podName))
	query := elastic.NewBoolQuery().Must(stringQuery)
	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(300).
		Sort("stageTimestamp", false).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error%v", err)
	}
	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *StorageEsImpl) QueryNodeYamlsWithNodeUid(data interface{}, nodeUid string) error {

	if nodeUid == "" {
		return fmt.Errorf("the params is error, uid is nil")
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}

	begin := time.Now()
	defer func() {
		cost := utils.TimeSinceInMilliSeconds(begin)
		metrics.QueryMethodDurationMilliSeconds.WithLabelValues("QueryNodeYaml").Observe(cost)
	}()

	stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("uid: \"%s\"", nodeUid))
	query := elastic.NewBoolQuery().Must(stringQuery)
	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(1).
		Sort("stageTimestamp", false).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error%v", err)
	}
	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *StorageEsImpl) QueryNodeYamlTimelineWithNodeUid(nodeUid string, opts ...model.OptionFunc) (timeline []string, err error) {
	options := &model.Options{
		Limit:     500,
		Ascending: false,
		OrderBy:   "stageTimestamp",
	}
	for _, do := range opts {
		do(options)
	}
	if nodeUid == "" {
		return timeline, fmt.Errorf("the params is error, uid is nil")
	}

	begin := time.Now()
	defer func() {
		cost := utils.TimeSinceInMilliSeconds(begin)
		metrics.QueryMethodDurationMilliSeconds.WithLabelValues("QueryNodeYamlTimeline").Observe(cost)
	}()

	stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("uid: \"%s\"", nodeUid))
	query := elastic.NewBoolQuery().Must(stringQuery)
	searchResult, err := s.DB.Search().Index(nodeYamlIndexName).Type(nodeYamlTypeName).Query(query).Size(options.Limit).
		Sort(options.OrderBy, options.Ascending).Do(context.Background())
	if err != nil {
		return timeline, fmt.Errorf("error%v", err)
	}
	for _, hit := range searchResult.Hits.Hits {
		timeline = append(timeline, hit.Index)
	}
	return
}

func (s *StorageEsImpl) QueryNodeYamlsWithNodeName(data interface{}, nodeName string) error {

	if nodeName == "" {
		return fmt.Errorf("the params is error, uid is nil")
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}

	begin := time.Now()
	defer func() {
		cost := utils.TimeSinceInMilliSeconds(begin)
		metrics.QueryMethodDurationMilliSeconds.WithLabelValues("QueryNodeYaml").Observe(cost)
	}()

	stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("nodeName: \"%s\"", nodeName))
	query := elastic.NewBoolQuery().Must(stringQuery)
	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(1).
		Sort("stageTimestamp", false).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error%v", err)
	}
	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *StorageEsImpl) QueryNodeYamlsV2WithNodeName(data interface{}, index, nodeName string, opts ...model.OptionFunc) error {
	options := &model.Options{
		Limit:     500,
		Ascending: false,
		OrderBy:   "stageTimestamp",
	}
	for _, do := range opts {
		do(options)
	}
	if nodeName == "" {
		return fmt.Errorf("the params is error, uid is nil")
	}
	_, _, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}

	begin := time.Now()
	defer func() {
		cost := utils.TimeSinceInMilliSeconds(begin)
		metrics.QueryMethodDurationMilliSeconds.WithLabelValues("QueryNodeYaml").Observe(cost)
	}()

	stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("nodeName: \"%s\"", nodeName))
	query := elastic.NewBoolQuery().Must(stringQuery)
	searchResult, err := s.DB.Search().Index(index).Type(esType).Query(query).Size(1).
		Sort(options.OrderBy, options.Ascending).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error%v", err)
	}
	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *StorageEsImpl) QueryNodeYamlTimelineWithNodeName(nodeName string, opts ...model.OptionFunc) (timeline []string, err error) {
	options := &model.Options{
		Limit:     500,
		Ascending: false,
		OrderBy:   "stageTimestamp",
	}
	for _, do := range opts {
		do(options)
	}
	if nodeName == "" {
		return timeline, fmt.Errorf("the params is error, nodeName is nil")
	}

	begin := time.Now()
	defer func() {
		cost := utils.TimeSinceInMilliSeconds(begin)
		metrics.QueryMethodDurationMilliSeconds.WithLabelValues("QueryNodeYamlTimeline").Observe(cost)
	}()

	stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("nodeName: \"%s\"", nodeName))
	query := elastic.NewBoolQuery().Must(stringQuery)
	searchResult, err := s.DB.Search().Index(nodeYamlIndexName).Type(nodeYamlTypeName).Query(query).Size(options.Limit).
		Sort(options.OrderBy, options.Ascending).Do(context.Background())
	if err != nil {
		return timeline, fmt.Errorf("error%v", err)
	}
	for _, hit := range searchResult.Hits.Hits {
		timeline = append(timeline, hit.Index)
	}
	return
}

func (s *StorageEsImpl) QueryNodeYamlsWithNodeIP(data interface{}, nodeIp string) error {

	if nodeIp == "" {
		return fmt.Errorf("the params is error, uid is nil")
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}

	begin := time.Now()
	defer func() {
		cost := utils.TimeSinceInMilliSeconds(begin)
		metrics.QueryMethodDurationMilliSeconds.WithLabelValues("QueryLifePhase").Observe(cost)
	}()

	stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("nodeIp: \"%s\"", nodeIp))
	query := elastic.NewBoolQuery().Must(stringQuery)
	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(1).
		Sort("stageTimestamp", false).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error%v", err)
	}
	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *StorageEsImpl) QueryNodeUIDListWithNodeIp(data interface{}, nodeIp string) error {

	if nodeIp == "" {
		return fmt.Errorf("the params is error, nodeIp is nil")
	}

	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}

	begin := time.Now()
	defer func() {
		cost := utils.TimeSinceInMilliSeconds(begin)
		metrics.QueryMethodDurationMilliSeconds.WithLabelValues("Querynodeuid").Observe(cost)
	}()

	stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("nodeIp: \"%s\"", nodeIp))
	query := elastic.NewBoolQuery().Must(stringQuery)
	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(1).
		Sort("stageTimestamp", false).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error%v", err)
	}
	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *StorageEsImpl) QueryPodYamlsWithNodeIP(data interface{}, nodeIp string) error {

	returnResult, ok := data.(*[]*model.PodYaml)
	if !ok {
		return fmt.Errorf("parse error")
	}
	if nodeIp == "" {
		return fmt.Errorf("the params is error, uid is nil")
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}

	begin := time.Now()
	defer func() {
		cost := utils.TimeSinceInMilliSeconds(begin)
		metrics.QueryMethodDurationMilliSeconds.WithLabelValues("QueryPodYaml").Observe(cost)
	}()
	stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("hostIP.keyword: \"%s\"", nodeIp))
	deleteFalse := elastic.NewQueryStringQuery(fmt.Sprintf("isDeleted.keyword: \"%t\"", false))
	query := elastic.NewBoolQuery().Must(stringQuery).Must(deleteFalse)
	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(300).
		Sort("stageTimestamp", false).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error%v", err)
	}

	dedup := make(map[string]string)
	sloMap := make(map[int][]*model.SloTraceData)
	var mutex1 sync.Mutex
	var wg sync.WaitGroup
	for _, hit := range searchResult.Hits.Hits {
		pyaml := &model.PodYaml{}
		if er := json.Unmarshal(*hit.Source, pyaml); er == nil {
			if pyaml.Pod != nil {
				if _, ok := dedup[pyaml.PodUid]; !ok {
					podrestun := &model.PodYaml{
						AuditID:           pyaml.AuditID,
						ClusterName:       pyaml.ClusterName,
						HostIP:            pyaml.HostIP,
						PodIP:             pyaml.PodIP,
						Namespace:         pyaml.Namespace,
						PodUid:            pyaml.PodUid,
						CreationTimestamp: pyaml.Pod.CreationTimestamp.Time,
						DebugUrl:          "http://eavesdropping.sigma-eu95.svc.alipay.net:8080/api/v1/debugpod?name=" + pyaml.PodName,
						PodName:           pyaml.PodName,
						Status:            string(pyaml.Pod.Status.Phase),
					}
					key := len(*returnResult)
					wg.Add(1)
					go func() {
						defer wg.Done()
						sloTrace := make([]*model.SloTraceData, 0)
						s.QuerySloTraceDataWithPodUID(&sloTrace, pyaml.PodUid)
						mutex1.Lock()
						sloMap[key] = sloTrace
						mutex1.Unlock()
					}()
					*returnResult = append(*returnResult, podrestun)
					dedup[pyaml.PodUid] = "true"
				}
			}
		}
	}
	wg.Wait()
	for k, v := range *returnResult {
		v.SLOType = "OutOfDate"
		v.SLOResult = "OutOfDate"
		v.SLO = "OutOfDate"
		if len(sloMap[k]) > 0 {
			v.SLO = time.Duration(sloMap[k][0].PodSLO).String()
			for i := range sloMap[k] {
				if sloMap[k][i].StartUpResultFromCreate == "success" || sloMap[k][i].DeleteResult == "success" || sloMap[k][i].UpgradeResult == "success" {
					v.SLOResult = "success"
					v.SLOType = sloMap[k][i].Type
				} else {
					v.SLOResult = "fail"
				}
			}
		}
	}

	return nil
}

func (s *StorageEsImpl) QueryPodInfoWithPodUid(data interface{}, podUid string) error {

	if podUid == "" {
		return fmt.Errorf("the params is error, uid is nil")
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}

	begin := time.Now()
	defer func() {
		cost := utils.TimeSinceInMilliSeconds(begin)
		metrics.QueryMethodDurationMilliSeconds.WithLabelValues("QueryNodeYaml").Observe(cost)
	}()

	stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("podUID: \"%s\"", podUid))
	query := elastic.NewBoolQuery().Must(stringQuery)

	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(1).
		Sort("stageTimestamp", false).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error%v", err)
	}
	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *StorageEsImpl) QueryNodephaseWithNodeUID(data interface{}, nodeUid string) error {
	if nodeUid == "" {
		return fmt.Errorf("the params is error, uid is nil")
	}
	begin := time.Now()
	defer func() {
		metrics.ObserveQueryMethodDuration("QueryNodephase", begin)
	}()

	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}

	stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("uid: \"%s\"", nodeUid))
	query := elastic.NewBoolQuery().Must(stringQuery)

	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(200).
		Sort("startTime", false).Do(context.Background())

	if err != nil {
		return fmt.Errorf("error%v", err)
	}

	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *StorageEsImpl) QueryNodephaseWithNodeName(data interface{}, nodeName string) error {
	if nodeName == "" {
		return fmt.Errorf("the params is error, uid is nil")
	}
	begin := time.Now()
	defer func() {
		metrics.ObserveQueryMethodDuration("QueryNodephase", begin)
	}()

	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}

	stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("nodeName: \"%s\"", nodeName))
	query := elastic.NewBoolQuery().Must(stringQuery)

	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(200).
		Sort("startTime", false).Do(context.Background())

	if err != nil {
		return fmt.Errorf("error%v", err)
	}

	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *StorageEsImpl) QueryPodLifePhaseByID(data interface{}, uid string, opts ...model.OptionFunc) error {
	options := &model.Options{
		Limit:     200,
		Ascending: false,
		OrderBy:   "startTime",
	}
	for _, do := range opts {
		do(options)
	}
	if uid == "" {
		return fmt.Errorf("the params is error, uid is nil")
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}

	begin := time.Now()
	defer func() {
		cost := utils.TimeSinceInMilliSeconds(begin)
		metrics.QueryMethodDurationMilliSeconds.WithLabelValues("QueryLifePhase").Observe(cost)
	}()

	stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("_id: \"%s\"", uid))
	query := elastic.NewBoolQuery().Must(stringQuery)

	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(options.Limit).
		Sort(options.OrderBy, options.Ascending).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error%v", err)
	}
	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *StorageEsImpl) QuerySloTraceDataWithPodUID(data interface{}, podUid string) error {
	if podUid == "" {
		return fmt.Errorf("the params is error, uid is nil")
	}
	begin := time.Now()
	defer func() {
		metrics.ObserveQueryMethodDuration("QuerySloTraceData", begin)
	}()

	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}

	stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("PodUID: \"%s\"", podUid))
	query := elastic.NewBoolQuery().Must(stringQuery)
	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(10).
		Sort("CreatedTime", false).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error%v", err)
	}

	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *StorageEsImpl) QuerySloTraceDataWithPodRequestParams(data interface{}, requestParamsKey, requestParamsValue string) error {
	if requestParamsKey == "" {
		return fmt.Errorf("the params is error, requestParamsKey is nil")
	}
	if requestParamsValue == "" {
		return fmt.Errorf("the params is error, requestParamsValue is nil")
	}
	begin := time.Now()
	defer func() {
		metrics.ObserveQueryMethodDuration("QuerySloTraceData", begin)
	}()

	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}

	stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("%s: \"%s\"", requestParamsKey, requestParamsValue)) // e.g. requestParamsKey: "PodUID", "PodName"

	query := elastic.NewBoolQuery().Must(stringQuery)
	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(10).
		Sort("CreatedTime", false).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error%v", err)
	}

	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *StorageEsImpl) QueryCreateSloWithResult(data interface{}, requestParams *model.SloOptions) error {
	if requestParams == nil {
		return fmt.Errorf("the params is error")
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}

	begin := time.Now()
	defer func() {
		cost := utils.TimeSinceInMilliSeconds(begin)
		metrics.QueryMethodDurationMilliSeconds.WithLabelValues("QuerySlo").Observe(cost)
	}()

	env := DELIVERY_ENV_PROD
	if requestParams.Env == "test" {
		env = DELIVERY_ENV_TEST
	}
	if requestParams.Env == "dryrun" {
		env = DELIVERY_ENV_DRYRUN
	}

	query := elastic.NewBoolQuery()

	//处理用户错误：
	if requestParams.RemovingUserErrors {
		split := strings.Split(requestParams.UserError, ",")
		query = query.MustNot(
			elastic.NewBoolQuery().Must(
				elastic.NewTermsQuery("SLOViolationReason", ToAnySlice(split)...),
			),
			elastic.NewBoolQuery().Must(
				elastic.NewTermQuery("DeliveryStatus", "KILL"),
				elastic.NewBoolQuery().
					Should(elastic.NewRegexpQuery("SLOViolationReason", ".*TooMuchTime"),
						elastic.NewTermsQuery("SLOViolationReason", "ScheduleDelay", "beforeFinish")).
					MinimumNumberShouldMatch(1)),
		)
	}
	if requestParams.Result != "" {
		stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("SLOViolationReason: \"%s\"", requestParams.Result))
		query = query.Must(stringQuery)
	}
	if requestParams.NodeName != "" {
		stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("NodeName: \"%s\"", requestParams.NodeName))
		query = query.Must(stringQuery)
	}

	if requestParams.Cluster != "" {
		query = query.Must(elastic.NewTermQuery("Cluster.keyword", requestParams.Cluster))
	}
	if requestParams.DeliveryDuration != "" {
		query = query.Filter(elastic.NewRangeQuery("DeliveryDuration").Gte(requestParams.DeliveryDuration))
	}
	if requestParams.Namespace != "" {
		query = query.Must(elastic.NewTermQuery("Namespace.keyword", requestParams.Namespace))
	}
	if requestParams.TimeField == "" {
		requestParams.TimeField = "CreatedTime"
	}
	// add range query
	if !requestParams.From.IsZero() || !requestParams.To.IsZero() {
		rangeQuery := elastic.NewRangeQuery(requestParams.TimeField).TimeZone("UTC")
		if !requestParams.From.IsZero() {
			rangeQuery = rangeQuery.From(requestParams.From)
		}
		if !requestParams.To.IsZero() {
			rangeQuery = rangeQuery.To(requestParams.To)
		}
		query = query.Must(rangeQuery)
	} else {
		rangeQuery := elastic.NewRangeQuery(requestParams.TimeField).TimeZone("UTC").Gte("now-24h")
		query = query.Must(rangeQuery)
	}
	if requestParams.BizName != "" {
		stringQuery3 := elastic.NewQueryStringQuery(fmt.Sprintf("BizName: \"%s\"", requestParams.BizName))
		query = query.Must(stringQuery3)
	}

	if requestParams.DeliveryStatus != "" {

		stringQuery4 := elastic.NewQueryStringQuery(fmt.Sprintf("DeliveryStatusOrig: \"%s\"", requestParams.DeliveryStatus))
		if env == DELIVERY_ENV_TEST {
			stringQuery4 = elastic.NewQueryStringQuery(fmt.Sprintf("DeliveryStatusNew: \"%s\"", requestParams.DeliveryStatus))
		}
		if env == DELIVERY_ENV_DRYRUN {
			stringQuery4 = elastic.NewQueryStringQuery(fmt.Sprintf("DeliveryStatus: \"%s\"", requestParams.DeliveryStatus))
		}
		query = query.Must(stringQuery4)
	}

	if requestParams.SloTime != "" {
		sloduration, err := time.ParseDuration(requestParams.SloTime)

		if err == nil {
			stringQuery5 := elastic.NewQueryStringQuery(fmt.Sprintf("PodSLO: \"%d\"", int(sloduration)))
			if env == DELIVERY_ENV_TEST {
				stringQuery5 = elastic.NewQueryStringQuery(fmt.Sprintf("DeliverySLONew: \"%d\"", int(sloduration)))
			}
			if env == DELIVERY_ENV_DRYRUN {
				stringQuery5 = elastic.NewQueryStringQuery(fmt.Sprintf("DeliverySLO: \"%d\"", int(sloduration)))
			}
			query = query.Must(stringQuery5)
		} else {
			fmt.Printf("Error slotime format %s \n", requestParams.SloTime)
		}
	}

	querySize := 300
	if requestParams.Count != "" {
		count, err := strconv.Atoi(requestParams.Count)
		if err == nil {
			querySize = count
		}
	}
	if querySize > 500 {
		querySize = 500
	}

	//通过filter 构造查询能力
	timelionQb := clientUtils.NewQueryBuilder(nil)
	var timelionOb []*elastic.BoolQuery
	filters := requestParams.Filter
	if len(filters) > 0 {
		filtersAry := strings.Split(filters, "*")

		// must include
		for _, f := range filtersAry {
			fAry := strings.Split(f, ":")
			if len(fAry) < 2 {
				continue
			}
			if strings.Contains(f, "AND NOT") || strings.Contains(f, "AND") || strings.Contains(f, "OR") {
				// 使用新的解析器处理复杂表达式
				expression, err := timelionQb.BuildFromExpression(f, nil)
				if err != nil {
					continue
				}
				timelionOb = append(timelionOb, expression)
			} else {
				if len(fAry[1]) == 0 && fAry[0] != LogicPool {
					continue
				}
				if _, err := strconv.ParseBool(fAry[1]); err == nil {
					fAry[0] = strings.TrimSuffix(fAry[0], Keyword)
				}
				// add operation for query
				if strings.HasPrefix(fAry[1], "!") {
					fVal := strings.TrimPrefix(fAry[1], "!")
					if len(fVal) == 0 && fAry[0] != LogicPool {
						continue
					}
					if strings.Contains(fAry[0], "Cluster") {
						fVal = strings.TrimPrefix(fVal, "sigma-")
					}
					query = query.MustNot(elastic.NewTermsQuery(fAry[0], fVal))
				} else {
					if strings.Contains(fAry[0], "Cluster") {
						fAry[1] = strings.TrimPrefix(fAry[1], "sigma-")
					}
					query = query.Filter(elastic.NewTermsQuery(fAry[0], fAry[1]))
				}
			}
		}
	}
	query = timelionQb.MergeBoolQueries(append(timelionOb, query), clientUtils.CombineMust)

	source, _ := query.Source()
	fmt.Println(utils.Dumps(source))

	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(querySize).
		Sort(requestParams.TimeField, false).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error%v", err)
	}
	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}

	return nil
}

func ToAnySlice[T any](collection []T) []any {
	result := make([]any, len(collection))
	for i := range collection {
		result[i] = collection[i]
	}
	return result
}

func (s *StorageEsImpl) QueryUpgradeSloWithResult(data interface{}, requestParams *model.SloOptions) error {
	if requestParams == nil {
		return fmt.Errorf("the params is error")
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}
	query := elastic.NewBoolQuery()
	if requestParams.Result != "" {
		resultQuery := elastic.NewQueryStringQuery(fmt.Sprintf("UpgradeResult: \"%s\"", requestParams.Result))
		query = query.Must(resultQuery)
	}

	if requestParams.Cluster != "" {
		query = query.Must(elastic.NewTermQuery("Cluster.keyword", requestParams.Cluster))
	}

	if requestParams.Namespace != "" {
		query = query.Must(elastic.NewTermQuery("Namespace.keyword", requestParams.Namespace))
	}

	// Type: delete AND NOT DeleteResult: success AND NOT DeleteResult: "pod.beta1.sigma.ali/cni-allocated"
	if requestParams.Type != "" {
		typeQuery := elastic.NewQueryStringQuery(fmt.Sprintf("Type: \"%s\"", requestParams.Type))
		query = query.Must(typeQuery)
	}

	if requestParams.TimeField == "" {
		requestParams.TimeField = "Created"
	}
	// add range query
	if !requestParams.From.IsZero() || !requestParams.To.IsZero() {
		rangeQuery := elastic.NewRangeQuery(requestParams.TimeField).TimeZone("UTC")
		if !requestParams.From.IsZero() {
			rangeQuery = rangeQuery.From(requestParams.From)
		}
		if !requestParams.To.IsZero() {
			rangeQuery = rangeQuery.To(requestParams.To)
		}
		query = query.Must(rangeQuery)
	} else {
		rangeQuery := elastic.NewRangeQuery(requestParams.TimeField).TimeZone("UTC").Gte("now-24h")
		query = query.Must(rangeQuery)
	}

	querySize := 300
	if requestParams.Count != "" {
		count, err := strconv.Atoi(requestParams.Count)
		if err == nil {
			querySize = count
		}
	}
	if querySize > 500 {
		querySize = 500
	}

	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(querySize).
		Sort(requestParams.TimeField, false).Do(context.Background())
	if err != nil {
		klog.Error(err)
		return fmt.Errorf("the params is error")
	}

	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}

	return nil
}
func (s *StorageEsImpl) QueryDeleteSloWithResult(data interface{}, requestParams *model.SloOptions) error {
	if requestParams == nil {
		return fmt.Errorf("the params is error")
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}
	query := elastic.NewBoolQuery()
	if requestParams.Result != "" {
		stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("DeleteResult: \"%s\"", requestParams.Result))
		query = query.Must(stringQuery)
	}

	if requestParams.Cluster != "" {
		query = query.Must(elastic.NewTermQuery("Cluster.keyword", requestParams.Cluster))
	}
	if requestParams.BizName != "" {
		stringQuery3 := elastic.NewQueryStringQuery(fmt.Sprintf("BizName: \"%s\"", requestParams.BizName))
		query = query.Must(stringQuery3)
	}

	if requestParams.Namespace != "" {
		query = query.Must(elastic.NewTermQuery("Namespace.keyword", requestParams.Namespace))
	}

	// Type: delete AND NOT DeleteResult: success AND NOT DeleteResult: "pod.beta1.sigma.ali/cni-allocated"
	if requestParams.Type != "" {
		stringQuery4 := elastic.NewQueryStringQuery(fmt.Sprintf("Type: \"%s\"", requestParams.Type))
		query = query.Must(stringQuery4)
	}

	if requestParams.Result != "" {
		stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("DeleteResult: \"%s\"", requestParams.Result))
		query = query.Must(stringQuery)
	}
	if requestParams.DeliveryDuration != "" {
		query = query.Filter(elastic.NewRangeQuery("DeleteDuration").Gte(requestParams.DeliveryDuration))
	}

	if requestParams.TimeField == "" {
		requestParams.TimeField = "CreatedTime"
	}
	// add range query
	if !requestParams.From.IsZero() || !requestParams.To.IsZero() {
		rangeQuery := elastic.NewRangeQuery(requestParams.TimeField).TimeZone("UTC")
		if !requestParams.From.IsZero() {
			rangeQuery = rangeQuery.From(requestParams.From)
		}
		if !requestParams.To.IsZero() {
			rangeQuery = rangeQuery.To(requestParams.To)
		}
		query = query.Must(rangeQuery)
	} else {
		rangeQuery := elastic.NewRangeQuery(requestParams.TimeField).TimeZone("UTC").Gte("now-24h")
		query = query.Must(rangeQuery)
	}
	querySize := 300
	if requestParams.Count != "" {
		count, err := strconv.Atoi(requestParams.Count)
		if err == nil {
			querySize = count
		}
	}
	if querySize > 500 {
		querySize = 500
	}

	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(querySize).
		Sort(requestParams.TimeField, false).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error%v", err)
	}
	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}

	return nil
}
func (s *StorageEsImpl) QueryNodeYamlWithParams(data interface{}, debugparams *model.NodeParams) error {

	if debugparams == nil {
		return fmt.Errorf("the params is error")
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}

	begin := time.Now()
	defer func() {
		cost := utils.TimeSinceInMilliSeconds(begin)
		metrics.QueryMethodDurationMilliSeconds.WithLabelValues("QueryLifePhase").Observe(cost)
	}()

	var stringQuery *elastic.QueryStringQuery
	if debugparams.NodeName != "" {
		stringQuery = elastic.NewQueryStringQuery(fmt.Sprintf("nodeName: \"%s\"", debugparams.NodeName))
	} else if debugparams.NodeUid != "" {
		stringQuery = elastic.NewQueryStringQuery(fmt.Sprintf("uid: \"%s\"", debugparams.NodeUid))
	} else if debugparams.NodeIp != "" {
		stringQuery = elastic.NewQueryStringQuery(fmt.Sprintf("nodeIp : \"%s\"", debugparams.NodeIp))
	}
	query := elastic.NewBoolQuery().Must(stringQuery)
	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(1).
		Sort("stageTimestamp", false).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error%v", err)
	}
	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *StorageEsImpl) QueryAuditWithAuditId(data interface{}, auditid string) error {
	if auditid == "" {
		return fmt.Errorf("the params is error, uid is nil")
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}
	begin := time.Now()
	defer func() {
		metrics.ObserveQueryMethodDuration("QueryAudit", begin)
	}()

	query := elastic.NewBoolQuery().Must(elastic.NewQueryStringQuery(fmt.Sprintf("auditID: \"%s\"", auditid)))
	searchReulst, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(1).
		Sort("stageTimestamp", false).Do(context.Background())
	if err != nil {
		klog.Error(err)
		return fmt.Errorf("the query is error")
	}

	if err != nil {
		return fmt.Errorf("error%v", err)
	}

	res, ok := data.(*model.Audit)
	if !ok {
		return fmt.Errorf("the query is error")
	}
	for _, hit := range searchReulst.Hits.Hits {
		err = json.Unmarshal(*hit.Source, &res.AuditLog)
		if err != nil {
			return err
		}
	}

	return nil
}
func (s *StorageEsImpl) QueryEventPodsWithPodUid(data interface{}, PodUid string) error {
	if PodUid == "" {
		return fmt.Errorf("the params is error, uid is nil")
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}
	begin := time.Now()
	defer func() {
		metrics.ObserveQueryMethodDuration("QueryNodephase", begin)
	}()

	query := elastic.NewBoolQuery()
	query.Must(elastic.NewTermQuery("objectRef.uid.keyword", PodUid))
	query.Must(elastic.NewTermQuery("objectRef.resource", "pods"))

	searchReulst, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(10).
		Sort("stageTimestamp", false).Do(context.Background())
	if err != nil {
		klog.Error(err)
		return fmt.Errorf("the query is error")
	}

	if err != nil {
		return fmt.Errorf("error%v", err)
	}

	res, ok := data.(*model.Audit)
	if !ok {
		return fmt.Errorf("the query is error")
	}
	for _, hit := range searchReulst.Hits.Hits {
		err = json.Unmarshal(*hit.Source, &res)
		if err != nil {
			return err
		}
	}

	return nil
}
func (s *StorageEsImpl) QueryEventNodeWithPodUid(data interface{}, PodUid string) error {
	if PodUid == "" {
		return fmt.Errorf("the params is error, uid is nil")
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}
	begin := time.Now()
	defer func() {
		metrics.ObserveQueryMethodDuration("QueryNodephase", begin)
	}()

	query := elastic.NewBoolQuery()
	query.Must(elastic.NewTermQuery("objectRef.uid.keyword", PodUid))
	query.Must(elastic.NewTermQuery("objectRef.resource", "nodes"))

	searchReulst, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(10).
		Sort("stageTimestamp", false).Do(context.Background())
	if err != nil {
		klog.Error(err)
		return fmt.Errorf("the query is error")
	}

	if err != nil {
		return fmt.Errorf("error%v", err)
	}

	res, ok := data.(*model.Audit)
	if !ok {
		return fmt.Errorf("the query is error")
	}
	for _, hit := range searchReulst.Hits.Hits {
		err = json.Unmarshal(*hit.Source, &res)
		if err != nil {
			return err
		}
	}

	return nil
}
func (s *StorageEsImpl) QueryEventWithTimeRange(data interface{}, from, to time.Time) error {

	query := elastic.NewBoolQuery()
	if !from.IsZero() || !to.IsZero() {
		rangeQuery := elastic.NewRangeQuery("stageTimestamp").TimeZone("UTC")
		if !from.IsZero() {
			rangeQuery = rangeQuery.From(from)
		}
		if !to.IsZero() {
			rangeQuery = rangeQuery.To(to)
		}
		query = query.Must(rangeQuery)
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}
	begin := time.Now()
	defer func() {
		metrics.ObserveQueryMethodDuration("QueryNodephase", begin)
	}()

	searchReulst, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(10).
		Sort("stageTimestamp", false).Do(context.Background())
	if err != nil {
		klog.Error(err)
		return fmt.Errorf("the query is error")
	}

	res, ok := data.(*model.Audit)
	if !ok {
		return fmt.Errorf("the query is error")
	}
	for _, hit := range searchReulst.Hits.Hits {
		err = json.Unmarshal(*hit.Source, &res)
		if err != nil {
			return err
		}
	}

	return nil
}
func (s *StorageEsImpl) QueryPodYamlWithParams(data interface{}, params *model.PodParams) error {

	if params == nil {
		return fmt.Errorf("the params is error")
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}

	begin := time.Now()
	defer func() {
		cost := utils.TimeSinceInMilliSeconds(begin)
		metrics.QueryMethodDurationMilliSeconds.WithLabelValues("QueryPodYamlWithParams").Observe(cost)
	}()

	var stringQuery *elastic.QueryStringQuery
	if params.Podip != "" {
		stringQuery = elastic.NewQueryStringQuery(fmt.Sprintf("podIP: \"%s\"", params.Podip))
	} else if params.Name != "" {
		stringQuery = elastic.NewQueryStringQuery(fmt.Sprintf("podName.keyword: \"%s\"", params.Name))
	} else if params.Hostname != "" {
		stringQuery = elastic.NewQueryStringQuery(fmt.Sprintf("hostname: \"%s\"", params.Hostname))
	}
	query := elastic.NewBoolQuery().Must(stringQuery)
	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(10).
		Sort("stageTimestamp", false).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error%v", err)
	}
	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *StorageEsImpl) QueryResourceYamlWithUID(kind, uid string) (interface{}, error) {
	result := make([]interface{}, 0)
	if uid == "" {
		return result, nil
	}

	index := ""
	typ := ""
	var stringQuery *elastic.QueryStringQuery
	if kind == PodResource {
		index = podYamlIndexName
		typ = podYamlTypeName
		stringQuery = elastic.NewQueryStringQuery(fmt.Sprintf("podUID: \"%s\"", uid))
	} else if kind == NodeResource {
		index = nodeYamlIndexName
		typ = nodeYamlTypeName
		stringQuery = elastic.NewQueryStringQuery(fmt.Sprintf("uid: \"%s\"", uid))
	}

	if index == "" {
		return result, nil
	}

	begin := time.Now()
	defer func() {
		cost := utils.TimeSinceInMilliSeconds(begin)
		metrics.QueryMethodDurationMilliSeconds.WithLabelValues("QueryResourceYamlWithUID").Observe(cost)
	}()

	query := elastic.NewBoolQuery().Must(stringQuery)
	searchResult, err := s.DB.Search().Index(index).Type(typ).Query(query).Size(1).
		Sort("stageTimestamp", false).Do(context.Background())
	if err != nil {
		err = fmt.Errorf("failed to get yaml of %s[%s], error: %v", kind, uid, err)
		klog.Error(err)
		return result, err
	}

	for _, hit := range searchResult.Hits.Hits {
		if kind == PodResource {
			objYaml := &model.PodYaml{}
			if er := json.Unmarshal(*hit.Source, objYaml); er == nil {
				if objYaml.Pod != nil {
					result = append(result, objYaml)
					fmt.Println("fetch pod yaml")
				}
			}
		} else if kind == NodeResource {
			objYaml := &model.NodeYaml{}
			if er := json.Unmarshal(*hit.Source, objYaml); er == nil {
				if objYaml.Node != nil {
					result = append(result, objYaml)
				}
			}
		}
	}

	return result, nil
}

func (s *StorageEsImpl) QueryResourceYamlWithName(kind, name string) (interface{}, error) {
	result := make([]interface{}, 0)
	if name == "" {
		return result, nil
	}

	index := ""
	typ := ""
	var stringQuery *elastic.QueryStringQuery
	if kind == PodResource {
		index = podYamlIndexName
		typ = podYamlTypeName
		stringQuery = elastic.NewQueryStringQuery(fmt.Sprintf("podName.keyword: \"%s\"", name))
	} else if kind == NodeResource {
		index = nodeYamlIndexName
		typ = nodeYamlTypeName
		stringQuery = elastic.NewQueryStringQuery(fmt.Sprintf("nodeName.keyword: \"%s\"", name))
	}

	if index == "" {
		return result, nil
	}

	begin := time.Now()
	defer func() {
		cost := utils.TimeSinceInMilliSeconds(begin)
		metrics.QueryMethodDurationMilliSeconds.WithLabelValues("QueryResourceYamlWithName").Observe(cost)
	}()

	query := elastic.NewBoolQuery().Must(stringQuery)
	searchResult, err := s.DB.Search().Index(index).Type(typ).Query(query).Size(1).
		Sort("stageTimestamp", false).Do(context.Background())
	if err != nil {
		err = fmt.Errorf("failed to get yaml of %s[%s], error: %v", kind, name, err)
		klog.Error(err)
		return result, err
	}

	for _, hit := range searchResult.Hits.Hits {
		if kind == PodResource {
			objYaml := &model.PodYaml{}
			if er := json.Unmarshal(*hit.Source, objYaml); er == nil {
				if objYaml.Pod != nil {
					result = append(result, objYaml)
				}
			}
		} else if kind == NodeResource {
			objYaml := &model.NodeYaml{}
			if er := json.Unmarshal(*hit.Source, objYaml); er == nil {
				if objYaml.Node != nil {
					result = append(result, objYaml)
				}
			}
		}
	}

	return result, nil
}
func (s *StorageEsImpl) QuerySloTraceDataWithOwnerId(data interface{}, ownerid string, opts ...model.OptionFunc) error {
	options := &model.Options{
		Limit: 5,
	}
	for _, do := range opts {
		do(options)
	}
	if len(ownerid) == 0 {
		return fmt.Errorf("the params is error, uid is nil")
	}
	begin := time.Now()
	defer func() {
		metrics.ObserveQueryMethodDuration("QuerySloTraceData", begin)
	}()
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}
	stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("ExtraProperties.ownerref.uid.Value.keyword: \"%s\"", ownerid))
	query := elastic.NewBoolQuery().Must(stringQuery)
	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(options.Limit).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error%v", err)
	}
	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}
	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *StorageEsImpl) QueryDeliveryWithDuration(data interface{}, requestParams *model.SloOptions, opts ...model.OptionFunc) error {
	options := &model.Options{
		Limit: 5,
	}
	for _, do := range opts {
		do(options)
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}
	query := elastic.NewBoolQuery()
	if requestParams.DeliveryDuration != "" {
		query = query.Filter(elastic.NewRangeQuery("DeliveryDuration").Gte(requestParams.DeliveryDuration))
	}
	if requestParams.Cluster != "" {
		query = query.Must(elastic.NewTermQuery("Cluster.keyword", requestParams.Cluster))
	}
	if requestParams.BizName != "" {
		stringQuery3 := elastic.NewQueryStringQuery(fmt.Sprintf("BizName: \"%s\"", requestParams.BizName))
		query = query.Must(stringQuery3)
	}
	if requestParams.Result != "" {
		stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("SLOViolationReason: \"%s\"", requestParams.Result))
		query = query.Must(stringQuery)
	}
	if requestParams.Namespace != "" {
		query = query.Must(elastic.NewTermQuery("Namespace.keyword", requestParams.Namespace))
	}
	if !requestParams.From.IsZero() || !requestParams.To.IsZero() {
		rangeQuery := elastic.NewRangeQuery("CreatedTime").TimeZone("UTC")
		if !requestParams.From.IsZero() {
			rangeQuery = rangeQuery.From(requestParams.From)
		}
		if !requestParams.To.IsZero() {
			rangeQuery = rangeQuery.To(requestParams.To)
		}
		query = query.Must(rangeQuery)
	}
	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(options.Limit).
		Sort("DeliveryDuration", false).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error%v", err)
	}
	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}
	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *StorageEsImpl) QueryObjectsFromTerms(result interface{}, index, doctype, filters, lte, gte string, size int) error {
	boolSearch := elastic.NewBoolQuery()
	if len(filters) == 0 {
		return nil
	}
	filtersSlice := strings.Split(filters, "*")

	// must include
	for _, f := range filtersSlice {
		fAry := strings.Split(f, ":")
		termKey := fAry[0]
		termVal := fAry[1]

		if len(termVal) == 0 {
			continue
		}

		// add operation for query
		if strings.HasPrefix(termVal, "!") {
			tVal := strings.TrimPrefix(termVal, "!")
			if len(tVal) == 0 {
				continue
			}
			boolSearch = boolSearch.MustNot(elastic.NewTermsQuery(termKey, tVal))
		} else {
			boolSearch = boolSearch.Filter(elastic.NewTermsQuery(termKey, termVal))
		}
	}

	// range filter
	rangeQuery := elastic.NewRangeQuery("CreatedTime")
	if len(gte) > 0 {
		rangeQuery.Gte(gte)
	}
	if len(lte) > 0 {
		rangeQuery.Lte(lte)
	}

	if len(gte) > 0 || len(lte) > 0 {
		boolSearch = boolSearch.Filter(rangeQuery)
	}

	query := s.DB.Search().Index(index).Type(doctype).Query(boolSearch).Size(size)
	searchResult, err := query.Do(context.Background())
	if err != nil {
		return err
	}

	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, &result)
	if err != nil {
		return err
	}
	return nil
}

func (s *StorageEsImpl) QueryObjectsFromNativeQuery(result interface{}, query, gte, lte, timeAttr string, size int) error {
	boolSearch := elastic.NewBoolQuery()
	if len(query) == 0 {
		return nil
	}

	_, esTableName, esType, err := utilsh.GetMetaName(result)
	if err != nil {
		return err
	}

	filterStrQuery := elastic.NewQueryStringQuery(query)
	boolSearch = boolSearch.Must(filterStrQuery)

	// range filter
	rangeQuery := elastic.NewRangeQuery(timeAttr)
	if len(gte) > 0 {
		rangeQuery.Gte(gte)
	}
	if len(lte) > 0 {
		rangeQuery.Lte(lte)
	}

	if len(gte) > 0 || len(lte) > 0 {
		boolSearch = boolSearch.Filter(rangeQuery)
	}

	queryService := s.DB.Search().Index(esTableName).Type(esType).Query(boolSearch).Size(size)
	searchResult, err := queryService.Do(context.Background())
	if err != nil {
		return err
	}

	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, &result)
	if err != nil {
		return err
	}
	return nil
}

func (s *StorageEsImpl) QueryObjectsFromNativeMultiQuery(result interface{}, item interface{}, query []string, gte, lte, timeAttr string, size int) error {
	if len(query) == 0 {
		return nil
	}

	_, esTableName, esType, err := utilsh.GetMetaName(item)
	if err != nil {
		return err
	}

	// range filter
	rangeQuery := elastic.NewRangeQuery(timeAttr)
	if len(gte) > 0 {
		rangeQuery.Gte(gte)
	}
	if len(lte) > 0 {
		rangeQuery.Lte(lte)
	}

	queryService := s.DB.MultiSearch()
	for _, q := range query {
		boolSearch := elastic.NewBoolQuery()
		if len(gte) > 0 || len(lte) > 0 {
			boolSearch = boolSearch.Filter(rangeQuery)
		}

		filterStrQuery := elastic.NewQueryStringQuery(q)
		boolSearch = boolSearch.Must(filterStrQuery)
		searchReq := elastic.NewSearchRequest().Index(esTableName).Type(esType).Query(boolSearch).Size(size)
		queryService.Add(searchReq)
	}

	doResult, err := queryService.Do(context.Background())
	if err != nil {
		return err
	}

	hitsAll := make([][]*json.RawMessage, 0)
	for _, searchResult := range doResult.Responses {
		var hits []*json.RawMessage
		for _, hit := range searchResult.Hits.Hits {
			hits = append(hits, hit.Source)
		}
		hitsAll = append(hitsAll, hits)
	}

	var hitsStr []byte
	hitsStr, err = json.Marshal(hitsAll)
	if err != nil {
		return err
	}
	err = json.Unmarshal(hitsStr, &result)
	if err != nil {
		return err
	}
	return nil
}

func (s *StorageEsImpl) QueryEavesdroppingLatency(data interface{}, opts ...model.OptionFunc) error {
	options := &model.Options{
		Limit: 5,
	}
	for _, do := range opts {
		do(options)
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}
	query := elastic.NewBoolQuery().Must(elastic.NewExistsQuery(EavesdroppingLatencyTag))
	begin := time.Now()
	defer func() {
		metrics.ObserveQueryMethodDuration("QueryEavesdroppingMeta", begin)
	}()

	searchReulst, err := s.DB.Search().
		Index(esTableName).
		Type(esType).
		Size(options.Limit).
		Query(query).
		Do(context.Background())
	if err != nil {
		klog.Error(err)
		return fmt.Errorf("the EavesdroppingMeta error")
	}

	if err != nil {
		return fmt.Errorf("error%v", err)
	}

	res, ok := data.(*[]*model.EavesdroppingMeta)
	if !ok {
		return fmt.Errorf("unmarshl EavesdroppingMeta error")
	}
	for _, hit := range searchReulst.Hits.Hits {
		meta := &model.EavesdroppingMeta{}
		meta.ClusterName = hit.Id
		err = json.Unmarshal(*hit.Source, meta)
		if err != nil {
			return err
		}
		*res = append(*res, meta)
	}
	return nil
}

func (s *StorageEsImpl) QueryLivePodLatency(data interface{}, opts ...model.OptionFunc) error {
	options := &model.Options{
		Limit: 5,
	}
	for _, do := range opts {
		do(options)
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}
	query := elastic.NewBoolQuery().Must(elastic.NewExistsQuery(LatencyTag))
	stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("Type: \"%s\"", "live_pod"))
	query = query.Must(stringQuery)
	begin := time.Now()
	defer func() {
		metrics.ObserveQueryMethodDuration("QueryLivePodLatency", begin)
	}()

	searchReulst, err := s.DB.Search().
		Index(esTableName).
		Type(esType).
		Size(options.Limit).
		Query(query).
		Do(context.Background())
	if err != nil {
		klog.Error(err)
		return fmt.Errorf("the EavesdroppingMeta error")
	}

	if err != nil {
		return fmt.Errorf("error%v", err)
	}

	res, ok := data.(*[]*model.EavesdroppingMeta)
	if !ok {
		return fmt.Errorf("unmarshl EavesdroppingMeta error")
	}
	for _, hit := range searchReulst.Hits.Hits {
		meta := &model.EavesdroppingMeta{}
		meta.ClusterName = strings.TrimSuffix(hit.Id, "_livepod_latency")
		err = json.Unmarshal(*hit.Source, meta)
		if err != nil {
			return err
		}
		*res = append(*res, meta)
	}
	return nil
}

func (s *StorageEsImpl) QueryEavesdroppingMeta(data interface{}, opts ...model.OptionFunc) error {
	options := &model.Options{
		Limit: 5,
	}
	for _, do := range opts {
		do(options)
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}
	begin := time.Now()
	defer func() {
		metrics.ObserveQueryMethodDuration("QueryEavesdroppingMeta", begin)
	}()

	searchReulst, err := s.DB.Search().Index(esTableName).Type(esType).Size(options.Limit).Do(context.Background())
	if err != nil {
		klog.Error(err)
		return fmt.Errorf("the EavesdroppingMeta error")
	}

	if err != nil {
		return fmt.Errorf("error%v", err)
	}

	res, ok := data.(*[]*model.EavesdroppingMeta)
	if !ok {
		return fmt.Errorf("unmarshl EavesdroppingMeta error")
	}
	for _, hit := range searchReulst.Hits.Hits {
		meta := &model.EavesdroppingMeta{}
		meta.ClusterName = hit.Id
		err = json.Unmarshal(*hit.Source, meta)
		if err != nil {
			return err
		}
		*res = append(*res, meta)
	}
	return nil
}
func (s *StorageEsImpl) QueryLifePhaseWithDatasourceId(data interface{}, uid string, opts ...model.OptionFunc) error {
	options := &model.Options{
		Limit:     200,
		Ascending: false,
		OrderBy:   "startTime",
	}
	for _, do := range opts {
		do(options)
	}
	if uid == "" {
		return fmt.Errorf("the params is error, uid is nil")
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}

	begin := time.Now()
	defer func() {
		cost := utils.TimeSinceInMilliSeconds(begin)
		metrics.QueryMethodDurationMilliSeconds.WithLabelValues("QueryLifePhase").Observe(cost)
	}()

	stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("dataSourceId.keyword: \"%s\"", uid))
	query := elastic.NewBoolQuery().Must(stringQuery)

	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(options.Limit).
		Sort(options.OrderBy, options.Ascending).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error%v", err)
	}
	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *StorageEsImpl) QueryAuditlogWithAuditId(data interface{}, Uid string, opts ...model.OptionFunc) error {
	if Uid == "" {
		return fmt.Errorf("the params is error, uid is nil")
	}
	options := &model.Options{
		Limit:     1,
		Ascending: false,
		OrderBy:   "stageTimestamp",
	}
	for _, do := range opts {
		do(options)
	}

	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}
	begin := time.Now()
	defer func() {
		metrics.ObserveQueryMethodDuration("QueryAuditlogWithAuditId", begin)
	}()

	query := elastic.NewBoolQuery().Must(elastic.NewQueryStringQuery(fmt.Sprintf("auditID: \"%s\"", Uid)))

	searchReulst, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(options.Limit).
		Sort(options.OrderBy, options.Ascending).Do(context.Background())
	if err != nil {
		klog.Error(err)
		return fmt.Errorf("the query is error")
	}

	if err != nil {
		return fmt.Errorf("error%v", err)
	}

	var hits []*json.RawMessage
	for _, hit := range searchReulst.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	if len(hits) == 0 {
		return nil
	}
	hitsStr, err := json.Marshal(hits[0])
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		klog.Errorf("the query is error:%v", err)
		return err
	}
	return nil
}

func (s *StorageEsImpl) QueryWorkloadEventWithUid(data interface{}, Uid string, opts ...model.OptionFunc) error {
	if Uid == "" {
		return fmt.Errorf("the params is error, uid is nil")
	}
	options := &model.Options{
		Limit:     5000,
		Ascending: false,
		OrderBy:   "stageTimestamp",
	}
	for _, do := range opts {
		do(options)
	}

	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}
	begin := time.Now()
	defer func() {
		metrics.ObserveQueryMethodDuration("QueryWorkloadEventWithUid", begin)
	}()

	query := elastic.NewBoolQuery()
	query.Should(elastic.NewTermQuery("responseObject.metadata.uid.keyword", Uid),
		elastic.NewTermQuery("responseObject.involvedObject.uid", Uid))

	searchReulst, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(options.Limit).
		Sort(options.OrderBy, options.Ascending).Do(context.Background())
	if err != nil {
		klog.Error(err)
		return fmt.Errorf("the query is error")
	}

	if err != nil {
		return fmt.Errorf("error%v", err)
	}

	var hits []*json.RawMessage
	for _, hit := range searchReulst.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		klog.Errorf("the query is error:%v", err)
		return err
	}
	return nil
}

func (s *StorageEsImpl) QueryWorkloadEventWithName(data interface{}, Name string, opts ...model.OptionFunc) error {
	if Name == "" {
		return fmt.Errorf("the params is error, name is nil")
	}
	options := &model.Options{
		Limit:     2000,
		Ascending: false,
		OrderBy:   "stageTimestamp",
	}

	for _, do := range opts {
		do(options)
	}

	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}
	begin := time.Now()
	defer func() {
		metrics.ObserveQueryMethodDuration("QueryWorkloadEventWithName", begin)
	}()

	query := elastic.NewBoolQuery()
	query.Must(elastic.NewTermQuery("responseObject.metadata.name.keyword", Name))
	//query.Must(elastic.NewTermQuery("objectRef.resource", type))

	searchReulst, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(options.Limit).
		Sort(options.OrderBy, options.Ascending).Do(context.Background())
	if err != nil {
		klog.Error(err)
		return fmt.Errorf("the query is error")
	}

	if err != nil {
		return fmt.Errorf("error%v", err)
	}

	var hits []*json.RawMessage
	for _, hit := range searchReulst.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *StorageEsImpl) QueryNodeLifePhaseByID(data interface{}, uid string, opts ...model.OptionFunc) error {
	options := &model.Options{
		Limit:     200,
		Ascending: false,
		OrderBy:   "startTime",
	}
	for _, do := range opts {
		do(options)
	}
	if uid == "" {
		return fmt.Errorf("the params is error, uid is nil")
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}

	begin := time.Now()
	defer func() {
		cost := utils.TimeSinceInMilliSeconds(begin)
		metrics.QueryMethodDurationMilliSeconds.WithLabelValues("QueryLifePhase").Observe(cost)
	}()

	stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("_id: \"%s\"", uid))
	query := elastic.NewBoolQuery().Must(stringQuery)

	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(options.Limit).
		Sort(options.OrderBy, options.Ascending).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error%v", err)
	}
	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *StorageEsImpl) QueryLifePhaseWithNodeUid(data interface{}, uid string, opts ...model.OptionFunc) error {
	options := &model.Options{
		Limit:     200,
		Ascending: false,
		OrderBy:   "startTime",
	}
	for _, do := range opts {
		do(options)
	}
	if uid == "" {
		return fmt.Errorf("the params is error, uid is nil")
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}

	begin := time.Now()
	defer func() {
		cost := utils.TimeSinceInMilliSeconds(begin)
		metrics.QueryMethodDurationMilliSeconds.WithLabelValues("QueryNodeLifePhase").Observe(cost)
	}()

	stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("uid: \"%s\"", uid))
	query := elastic.NewBoolQuery().Must(stringQuery)

	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(options.Limit).
		Sort(options.OrderBy, options.Ascending).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error%v", err)
	}
	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *StorageEsImpl) QueryLifePhaseWithNodeName(data interface{}, name string, opts ...model.OptionFunc) error {
	options := &model.Options{
		Limit:     200,
		Ascending: false,
		OrderBy:   "startTime",
	}
	for _, do := range opts {
		do(options)
	}
	if name == "" {
		return fmt.Errorf("the params is error, uid is nil")
	}
	_, esTableName, esType, err := utilsh.GetMetaName(data)
	if err != nil {
		return err
	}

	begin := time.Now()
	defer func() {
		cost := utils.TimeSinceInMilliSeconds(begin)
		metrics.QueryMethodDurationMilliSeconds.WithLabelValues("QueryNodeLifePhase").Observe(cost)
	}()

	stringQuery := elastic.NewQueryStringQuery(fmt.Sprintf("nodeName: \"%s\"", name))
	query := elastic.NewBoolQuery().Must(stringQuery)

	searchResult, err := s.DB.Search().Index(esTableName).Type(esType).Query(query).Size(options.Limit).
		Sort(options.OrderBy, options.Ascending).Do(context.Background())
	if err != nil {
		return fmt.Errorf("error%v", err)
	}
	var hits []*json.RawMessage
	for _, hit := range searchResult.Hits.Hits {
		hits = append(hits, hit.Source)
	}
	hitsStr, err := json.Marshal(hits)
	if err != nil {
		return err
	}

	err = json.Unmarshal(hitsStr, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *StorageEsImpl) InsertGatherFeedback(data interface{}) error {
	defer utils.IgnorePanic(("SaveGatherFeedback"))
	if data == nil {
		return nil
	}
	err := utils.ReTry(func() error {
		bulkService := s.DB.Bulk()
		doc := elastic.NewBulkIndexRequest().Index("gather_feedback").Type("_doc").Id(genRandomUUID()).Doc(data).UseEasyJSON(true)
		bulkService = bulkService.Add(doc)
		_, err := bulkService.Do(context.Background())
		if err != nil {
			return err
		}
		return nil
	}, 1*time.Second, 10)
	if err != nil {
		return err
	}
	return nil
}
func genRandomUUID() string {
	// 创建一个 16 字节的缓冲区
	b := make([]byte, 16)
	// 从加密安全的随机源读取 16 字节的数据
	rand.Read(b)

	// 设置 UUID 版本 4 (随机生成)
	b[6] = (b[6] & 0x0f) | 0x40
	// 设置 UUID 变体为 RFC 4122 (10xxxxxx)
	b[8] = (b[8] & 0x3f) | 0x80

	// 使用 hex.EncodeToString 将字节转换为十六进制字符串
	// 然后手动插入破折号
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}
func (s *StorageEsImpl) QueryFeedback(feedbackType []string, opts ...model.OptionFunc) (*elastic.SearchHits, error) {
	options := &model.Options{
		Limit:     500,
		Ascending: false,
		OrderBy:   "createdTime",
	}
	for _, do := range opts {
		do(options)
	}
	if len(feedbackType) == 0 {
		return nil, fmt.Errorf("the params is error, feedbackType is null")
	}
	var interfaceValue []interface{}
	for _, v := range feedbackType {
		interfaceValue = append(interfaceValue, v)
	}

	begin := time.Now()
	defer func() {
		cost := utils.TimeSinceInMilliSeconds(begin)
		metrics.QueryMethodDurationMilliSeconds.WithLabelValues("QueryFeedback").Observe(cost)
	}()
	// 创建 BoolQuery 对象
	termsQuery := elastic.NewTermsQuery("feedbackType.keyword", interfaceValue...)
	query := elastic.NewBoolQuery().Must(termsQuery)
	if !options.From.IsZero() || !options.To.IsZero() {
		rangeQuery := elastic.NewRangeQuery("createdTime").TimeZone("UTC")
		if !options.From.IsZero() {
			rangeQuery.From(options.From)
		}
		if !options.To.IsZero() {
			rangeQuery.To(options.To)
		}
		query.Must(rangeQuery)
	}

	searchResult, err := s.DB.Search().Index("gather_feedback").Type("_doc").Query(query).Size(options.Limit).
		Sort(options.OrderBy, options.Ascending).Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("QueryFeedback error %v", err)
	}
	return searchResult.Hits, err
}

func (s *StorageEsImpl) UpdateGatherFeedback(id string, data interface{}) error {
	defer utils.IgnorePanic(("UpdateGatherFeedback"))
	if data == nil {
		return nil
	}
	err := utils.ReTry(func() error {
		updateService := s.DB.Update().
			Index("gather_feedback").
			Type("_doc").
			Id(id).
			Doc(data)
		_, err := updateService.Do(context.Background())
		if err != nil {
			return err
		}
		return nil
	}, 1*time.Second, 10)
	if err != nil {
		return err
	}
	return nil
}

// GetGatherFeedback 通过ID获取文档
func (s *StorageEsImpl) GetGatherFeedback(id string, data interface{}) error {
	defer utils.IgnorePanic("GetGatherFeedback")
	// 构建获取请求
	getService := s.DB.Get().Index("gather_feedback").Id(id)

	getResult, err := getService.Do(context.Background())
	if err != nil {
		if elastic.IsNotFound(err) {
			return fmt.Errorf("document with ID [ %s ] not found", id)
		}
		return fmt.Errorf("failed to get gather feedback [ %v ]", err)
	}
	marshalJSON, err := getResult.Source.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to unmarshal document source: [ %v ]", err)
	}
	err = json.Unmarshal(marshalJSON, &data)
	return err
}
