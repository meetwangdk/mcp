package data_access

import (
	"time"

	"github.com/olivere/elastic"

	"code.alipay.com/sigma/eavesdropping-client/go-client/model"
)

type AuditInterface interface {
	QueryAuditWithAuditId(data interface{}, auditid string) error
	QueryAuditlogWithAuditId(data interface{}, Uid string, opts ...model.OptionFunc) error
	QueryEventPodsWithPodUid(data interface{}, uid string) error
	QueryEventNodeWithPodUid(data interface{}, uid string) error
	QueryEventWithTimeRange(data interface{}, from, to time.Time) error
	QueryWorkloadEventWithUid(data interface{}, Uid string, opts ...model.OptionFunc) error
	QueryWorkloadEventWithName(data interface{}, Uid string, opts ...model.OptionFunc) error
}

type PodInterface interface {
	QueryPodYamlTimelineWithPodUID(uid string, opts ...model.OptionFunc) (timeline []string, err error)
	QueryPodYamlsV2WithPodUID(data interface{}, index, uid string, opts ...model.OptionFunc) error
	QueryPodYamlsWithPodUID(data interface{}, uid string, opts ...model.OptionFunc) error
	QueryPodYamlsWithPodName(data interface{}, name string, opts ...model.OptionFunc) error
	QueryPodYamlsWithHostName(data interface{}, hostname string, opts ...model.OptionFunc) error
	QueryPodYamlsWithPodIp(data interface{}, ip string, opts ...model.OptionFunc) error
	QueryPodYamlsWithNodeIP(data interface{}, ip string) error
	QueryPodListWithNodeip(data interface{}, nodeIp string, isDeleted bool) error
	QueryPodUIDListByHostname(data interface{}, hostname string) error
	QueryPodUIDListByPodIP(data interface{}, ip string) error
	QueryPodUIDListByPodName(data interface{}, podname string) error
	QueryPodYamlWithParams(data interface{}, opts *model.PodParams) error
}

type NodeInterface interface {
	QueryNodeYamlTimelineWithNodeUid(uid string, opts ...model.OptionFunc) (timeline []string, err error)
	QueryNodeYamlTimelineWithNodeName(nodeName string, opts ...model.OptionFunc) (timeline []string, err error)
	QueryNodeYamlsV2WithNodeName(data interface{}, index, name string, opts ...model.OptionFunc) error
	QueryNodeYamlsWithNodeUid(data interface{}, uid string) error
	QueryNodeYamlsWithNodeName(data interface{}, name string) error
	QueryNodeYamlsWithNodeIP(data interface{}, ip string) error
	QueryNodeYamlWithParams(data interface{}, opts *model.NodeParams) error
	QueryNodeUIDListWithNodeIp(data interface{}, ip string) error
}

type NodePhaseInterface interface {
	QueryNodephaseWithNodeName(data interface{}, name string) error
	QueryNodephaseWithNodeUID(data interface{}, uid string) error
}

type SloTraceDataInterface interface {
	QuerySloTraceDataWithPodUID(data interface{}, uid string) error
	QuerySloTraceDataWithPodRequestParams(data interface{}, requestParamsKey, requestParamsValue string) error
	QueryDeleteSloWithResult(data interface{}, opts *model.SloOptions) error
	QueryUpgradeSloWithResult(data interface{}, opts *model.SloOptions) error
	QueryCreateSloWithResult(data interface{}, opts *model.SloOptions) error
	QuerySloTraceDataWithOwnerId(data interface{}, ownerid string, opts ...model.OptionFunc) error
	QueryDeliveryWithDuration(data interface{}, params *model.SloOptions, opts ...model.OptionFunc) error
}

type PodLifePhaseInterface interface {
	QueryLifePhaseWithPodUid(data interface{}, uid string, opts ...model.OptionFunc) error
	QueryPodLifePhaseByID(data interface{}, uid string, opts ...model.OptionFunc) error
	QueryLifePhaseWithDatasourceId(data interface{}, uid string, opts ...model.OptionFunc) error
}

type NodeLifePhaseInterface interface {
	QueryLifePhaseWithNodeUid(data interface{}, uid string, opts ...model.OptionFunc) error
	QueryLifePhaseWithNodeName(data interface{}, name string, opts ...model.OptionFunc) error
	QueryNodeLifePhaseByID(data interface{}, uid string, opts ...model.OptionFunc) error
}

type SpanInterface interface {
	QuerySpanWithPodUid(data interface{}, uid string, opts ...model.OptionFunc) error
}

type PodInfoInterface interface {
	QueryPodInfoWithPodUid(data interface{}, uid string) error
}

type CommonFilterInterface interface {
	QueryObjectsFromTerms(result interface{}, index, doctype, filters, lte, gte string, size int) error
	QueryObjectsFromNativeQuery(result interface{}, query, gte, lte, timeAttr string, size int) error
	QueryObjectsFromNativeMultiQuery(result interface{}, item interface{}, query []string, gte, lte, timeAttr string, size int) error
}

type EavesdroppingMetaInterface interface {
	QueryEavesdroppingMeta(data interface{}, opts ...model.OptionFunc) error
	QueryEavesdroppingLatency(data interface{}, opts ...model.OptionFunc) error
	QueryLivePodLatency(data interface{}, opts ...model.OptionFunc) error
}

type GatherFeedbackInterface interface {
	InsertGatherFeedback(data interface{}) error
	QueryFeedback(feedbackType []string, opts ...model.OptionFunc) (*elastic.SearchHits, error)
	UpdateGatherFeedback(id string, data interface{}) error
	GetGatherFeedback(id string, data interface{}) error
}

type StorageInterface interface {
	AuditInterface
	PodInterface
	NodeInterface
	NodeLifePhaseInterface
	NodePhaseInterface
	SloTraceDataInterface
	PodLifePhaseInterface
	SpanInterface
	PodInfoInterface
	CommonFilterInterface
	EavesdroppingMetaInterface
	GatherFeedbackInterface
}
