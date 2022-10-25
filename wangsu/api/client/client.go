// This file is auto-generated, don't edit it. Thanks.
package client

import (
	"github.com/alibabacloud-go/tea/tea"
)

type BandwidthPeakRankingRequest struct {
	// {"en":"cust_en_name of sub-client.
	// When a merged-account wants to  view the information of the subclient,the cust_en_name is required.", "zh_CN":"合并账号下的某个客户的英文名，当合并账号要查看子客户的信息时，必须填写子客户的英文名"}
	Cust *string `json:"cust,omitempty" xml:"cust,omitempty"`
	// {"en":"Specifies the query date:
	// 1)With format yyyy-mm-dd.
	// 2)If not specified,it means today as default.", "zh_CN":"查询的日期，日期格式为yyyy-mm-dd,不选或者为空时默认为当天；"}
	Date *string `json:"date,omitempty" xml:"date,omitempty"`
	// {"en":"1)Must work with 'enddate' and they  specify the query date scope.
	// 2)With format yyyy-mm-dd.
	// 3)If there is a 'date' parameter,this parameter will be omitted.", "zh_CN":"查询的起始日期，日期格式为yyyy-mm-dd；此参数需与enddate参数配合,若存在date参数,则该参数无效"}
	Startdate *string `json:"startdate,omitempty" xml:"startdate,omitempty"`
	// {"en":"1)Must work with 'startdate' and they  specify the query date scope.
	// 2)With format yyyy-mm-dd
	// 3)If there is a 'date' parameter,this parameter will be omitted.", "zh_CN":"查询的结束日期,日期格式为yyyy-mm-dd；此参数需与startdate参数配合,若存在date参数,则该参数无效。"}
	Enddate *string `json:"enddate,omitempty" xml:"enddate,omitempty"`
	// {"en":"domains that been queried:
	// 1)If there are multiple inputs,use  ';' as separator.
	// 2)If not specified, it means all the domains of the account .", "zh_CN":"查询的频道，多个频道值请用英文分号';'，不选或者为空时默认为所查询客户的所有频道"}
	Channel *string `json:"channel,omitempty" xml:"channel,omitempty"`
	// {"en":"1)If there are multiple inputs,use ';' as separator.For example,u can use 'region=cn;apac' to query data of cn and apac region.
	// 2)If not specified, it means all the regions.", "zh_CN":"查询的加速区域的缩写，多个区域请用英文分号';'分隔开，如查询大陆及亚太区域，参数填写为：'region=cn;apac'。不选或者为空时默认为全部区域。"}
	Region *string `json:"region,omitempty" xml:"region,omitempty"`
	// {"en":"1)If there area multiple inputs,use ';' as demimeter.
	// 2)optional values of isp: refers to the ISP-section of appendix.
	// 3) If not specified,means all the isp.", "zh_CN":"&nbsp;要查询的运营商的缩写，多个isp请用英文分号';'分隔开。运营商的缩写格式参考附录：具体运行商（ISP）信息的代号。备注：只有当地区只写了'cn'时，填写isp信息才有效。不选或者为空时默认为所有isp。"}
	Isp *string `json:"isp,omitempty" xml:"isp,omitempty"`
	// {"en":"acceleration type.
	// 1)If there are multiple inputs,use ';' as separator.
	// 2)If not specified or specified as 'all', it means all the accetypes.", "zh_CN":"加速类型参数，如accetype=web。多个请用英文分号';'分隔开，不填或值为all表示所有类型"}
	Accetype *string `json:"accetype,omitempty" xml:"accetype,omitempty"`
	// {"en":"The response format:
	// 1)optional values:xml, json.
	// 2)'xml' as default.", "zh_CN":"返回结果格式,支持格式为xml和json,默认为xml"}
	Dataformat *string `json:"dataformat,omitempty" xml:"dataformat,omitempty"`
	// {"en":"Specifies if  the 'channel' parameter should be exactly matched:
	// 1)'true' as default.
	// 2) If not 'true',it will query data of channels that ends with any item of input 'channel's.", "zh_CN":"&nbsp;频道是否完全匹配,为true时，必须填写完整的域名(此时会过滤用户输入的无效或重复频道,所有输入频道都无效时返403)。不为true时，显示以用户输入的频道为结尾的所有频道。默认为true"}
	IsExactMatch *string `json:"isExactMatch,omitempty" xml:"isExactMatch,omitempty"`
	// {"en":"Different data types.
	// 1)optional values:1,2,3
	// 2)'2' means bandwidth of http.'3' means bandwidth of https.'1' mean the total bandwidth.
	// 3)ISP parameter not supported.", "zh_CN":"datatype=1时，输出总带宽；datatype=2时输出http的带宽；datatype=3时，输出https的带宽。默认datatype=1。当datatype=2或者3时，不支持isp入参。"}
	Datatype *string `json:"datatype,omitempty" xml:"datatype,omitempty"`
}

func (s BandwidthPeakRankingRequest) String() string {
	return tea.Prettify(s)
}

func (s BandwidthPeakRankingRequest) GoString() string {
	return s.String()
}

func (s *BandwidthPeakRankingRequest) SetCust(v string) *BandwidthPeakRankingRequest {
	s.Cust = &v
	return s
}

func (s *BandwidthPeakRankingRequest) SetDate(v string) *BandwidthPeakRankingRequest {
	s.Date = &v
	return s
}

func (s *BandwidthPeakRankingRequest) SetStartdate(v string) *BandwidthPeakRankingRequest {
	s.Startdate = &v
	return s
}

func (s *BandwidthPeakRankingRequest) SetEnddate(v string) *BandwidthPeakRankingRequest {
	s.Enddate = &v
	return s
}

func (s *BandwidthPeakRankingRequest) SetChannel(v string) *BandwidthPeakRankingRequest {
	s.Channel = &v
	return s
}

func (s *BandwidthPeakRankingRequest) SetRegion(v string) *BandwidthPeakRankingRequest {
	s.Region = &v
	return s
}

func (s *BandwidthPeakRankingRequest) SetIsp(v string) *BandwidthPeakRankingRequest {
	s.Isp = &v
	return s
}

func (s *BandwidthPeakRankingRequest) SetAccetype(v string) *BandwidthPeakRankingRequest {
	s.Accetype = &v
	return s
}

func (s *BandwidthPeakRankingRequest) SetDataformat(v string) *BandwidthPeakRankingRequest {
	s.Dataformat = &v
	return s
}

func (s *BandwidthPeakRankingRequest) SetIsExactMatch(v string) *BandwidthPeakRankingRequest {
	s.IsExactMatch = &v
	return s
}

func (s *BandwidthPeakRankingRequest) SetDatatype(v string) *BandwidthPeakRankingRequest {
	s.Datatype = &v
	return s
}

type BandwidthPeakRankingResponse struct {
	// {'en':'provider', 'zh_CN':'结果'}
	Provider *BandwidthPeakRankingResponseProvider `json:"provider,omitempty" xml:"provider,omitempty" require:"true" type:"Struct"`
}

func (s BandwidthPeakRankingResponse) String() string {
	return tea.Prettify(s)
}

func (s BandwidthPeakRankingResponse) GoString() string {
	return s.String()
}

func (s *BandwidthPeakRankingResponse) SetProvider(v *BandwidthPeakRankingResponseProvider) *BandwidthPeakRankingResponse {
	s.Provider = v
	return s
}

type BandwidthPeakRankingResponseProvider struct {
	// {'en':'tenant', 'zh_CN':'租户'}
	Name *string `json:"name,omitempty" xml:"name,omitempty" require:"true"`
	// {'en':'type', 'zh_CN':'接口类型'}
	Type *string `json:"type,omitempty" xml:"type,omitempty" require:"true"`
	// {'en':'data', 'zh_CN':'数据'}
	Date *BandwidthPeakRankingResponseProviderDate `json:"date,omitempty" xml:"date,omitempty" require:"true" type:"Struct"`
}

func (s BandwidthPeakRankingResponseProvider) String() string {
	return tea.Prettify(s)
}

func (s BandwidthPeakRankingResponseProvider) GoString() string {
	return s.String()
}

func (s *BandwidthPeakRankingResponseProvider) SetName(v string) *BandwidthPeakRankingResponseProvider {
	s.Name = &v
	return s
}

func (s *BandwidthPeakRankingResponseProvider) SetType(v string) *BandwidthPeakRankingResponseProvider {
	s.Type = &v
	return s
}

func (s *BandwidthPeakRankingResponseProvider) SetDate(v *BandwidthPeakRankingResponseProviderDate) *BandwidthPeakRankingResponseProvider {
	s.Date = v
	return s
}

type BandwidthPeakRankingResponseProviderDate struct {
	// {'en':'startdate', 'zh_CN':'开始日期'}
	Startdate *string `json:"startdate,omitempty" xml:"startdate,omitempty" require:"true"`
	// {'en':'enddate', 'zh_CN':'结束日期'}
	Enddate *string `json:"enddate,omitempty" xml:"enddate,omitempty" require:"true"`
	// {'en':'channelPeak', 'zh_CN':'频道峰值数据'}
	ChannelPeak *BandwidthPeakRankingResponseProviderDateChannelPeak `json:"channelPeak,omitempty" xml:"channelPeak,omitempty" require:"true" type:"Struct"`
}

func (s BandwidthPeakRankingResponseProviderDate) String() string {
	return tea.Prettify(s)
}

func (s BandwidthPeakRankingResponseProviderDate) GoString() string {
	return s.String()
}

func (s *BandwidthPeakRankingResponseProviderDate) SetStartdate(v string) *BandwidthPeakRankingResponseProviderDate {
	s.Startdate = &v
	return s
}

func (s *BandwidthPeakRankingResponseProviderDate) SetEnddate(v string) *BandwidthPeakRankingResponseProviderDate {
	s.Enddate = &v
	return s
}

func (s *BandwidthPeakRankingResponseProviderDate) SetChannelPeak(v *BandwidthPeakRankingResponseProviderDateChannelPeak) *BandwidthPeakRankingResponseProviderDate {
	s.ChannelPeak = v
	return s
}

type BandwidthPeakRankingResponseProviderDateChannelPeak struct {
	// {'en':'channel', 'zh_CN':'频道'}
	Channel *string `json:"channel,omitempty" xml:"channel,omitempty" require:"true"`
	// {'en':'peakTime', 'zh_CN':'峰值时间'}
	PeakTime *string `json:"peakTime,omitempty" xml:"peakTime,omitempty" require:"true"`
	// {'en':'peakvalue（Mbps）', 'zh_CN':'带宽峰值，单位Mbps'}
	PeakValue *string `json:"peakValue,omitempty" xml:"peakValue,omitempty" require:"true"`
	// {'en':'total traffic,unit GB', 'zh_CN':'总流量，单位GB'}
	TotalFlow *string `json:"totalFlow,omitempty" xml:"totalFlow,omitempty" require:"true"`
}

func (s BandwidthPeakRankingResponseProviderDateChannelPeak) String() string {
	return tea.Prettify(s)
}

func (s BandwidthPeakRankingResponseProviderDateChannelPeak) GoString() string {
	return s.String()
}

func (s *BandwidthPeakRankingResponseProviderDateChannelPeak) SetChannel(v string) *BandwidthPeakRankingResponseProviderDateChannelPeak {
	s.Channel = &v
	return s
}

func (s *BandwidthPeakRankingResponseProviderDateChannelPeak) SetPeakTime(v string) *BandwidthPeakRankingResponseProviderDateChannelPeak {
	s.PeakTime = &v
	return s
}

func (s *BandwidthPeakRankingResponseProviderDateChannelPeak) SetPeakValue(v string) *BandwidthPeakRankingResponseProviderDateChannelPeak {
	s.PeakValue = &v
	return s
}

func (s *BandwidthPeakRankingResponseProviderDateChannelPeak) SetTotalFlow(v string) *BandwidthPeakRankingResponseProviderDateChannelPeak {
	s.TotalFlow = &v
	return s
}

type ReportBandwidthMultiDomainRealTimeEdgeServiceRequest struct {
	// {'en':'-', 'zh_CN':'开始时间：
	//
	// 时间格式为yyyy-MM-ddTHH:mm:ss+08:00，例如，2019-01-01T10:00:00+08:00（为北京时间2019年1月1日10点0分0秒）；
	// 不能大于当前时间
	// 最多可获取最近半年（183天）的数据。'}
	DateFrom *string `json:"dateFrom,omitempty" xml:"dateFrom,omitempty"`
	// {'en':'-', 'zh_CN':'结束时间：
	//
	// 时间格式为yyyy-MM-ddTHH:mm:ss+08:00
	// 结束时间需大于开始时间，结束时间如果大于当前时间，取当前时间。
	// dateFrom，dateTo二者都未传，默认查询过去的1小时；如仅有一个未传，抛异常
	// 允许查询最大时间间隔：1小时（可联系技术支持调整），即dateFrom和dateTo相差不能超过1小时。'}
	DateTo *string `json:"dateTo,omitempty" xml:"dateTo,omitempty"`
	// {'en':'-', 'zh_CN':'数据粒度：不传默认1m
	//
	// 支持1m（1分钟）、5m（5分钟）'}
	DataInterval *string `json:"dataInterval,omitempty" xml:"dataInterval,omitempty"`
	// {'en':'-', 'zh_CN':'域名：
	//
	// 可传递域名数量上限默认为20（可联系技术支持调整）；
	// 自动过滤掉非法域名（如传递非法域名，会被过滤掉，查询结果只返回合法域名的数据）'}
	Domain []*string `json:"domain,omitempty" xml:"domain,omitempty" require:"true" type:"Repeated"`
	// {'en':'-', 'zh_CN':'分组维度
	//
	// 可选值为domain；
	// 有传入则按照该维度展示明细数据；'}
	GroupBy *string `json:"groupBy,omitempty" xml:"groupBy,omitempty"`
}

func (s ReportBandwidthMultiDomainRealTimeEdgeServiceRequest) String() string {
	return tea.Prettify(s)
}

func (s ReportBandwidthMultiDomainRealTimeEdgeServiceRequest) GoString() string {
	return s.String()
}

func (s *ReportBandwidthMultiDomainRealTimeEdgeServiceRequest) SetDateFrom(v string) *ReportBandwidthMultiDomainRealTimeEdgeServiceRequest {
	s.DateFrom = &v
	return s
}

func (s *ReportBandwidthMultiDomainRealTimeEdgeServiceRequest) SetDateTo(v string) *ReportBandwidthMultiDomainRealTimeEdgeServiceRequest {
	s.DateTo = &v
	return s
}

func (s *ReportBandwidthMultiDomainRealTimeEdgeServiceRequest) SetDataInterval(v string) *ReportBandwidthMultiDomainRealTimeEdgeServiceRequest {
	s.DataInterval = &v
	return s
}

func (s *ReportBandwidthMultiDomainRealTimeEdgeServiceRequest) SetDomain(v []*string) *ReportBandwidthMultiDomainRealTimeEdgeServiceRequest {
	s.Domain = v
	return s
}

func (s *ReportBandwidthMultiDomainRealTimeEdgeServiceRequest) SetGroupBy(v string) *ReportBandwidthMultiDomainRealTimeEdgeServiceRequest {
	s.GroupBy = &v
	return s
}

type ReportBandwidthMultiDomainRealTimeEdgeServiceResponse struct {
	// {'en':'-', 'zh_CN':'请求结果状态码'}
	Code *string `json:"code,omitempty" xml:"code,omitempty" require:"true"`
	// {'en':'-', 'zh_CN':'请求结果信息'}
	Message *string `json:"message,omitempty" xml:"message,omitempty" require:"true"`
	// {'en':'-', 'zh_CN':'请求结果的详细数据'}
	Data []*ReportBandwidthMultiDomainRealTimeEdgeServiceResponseData `json:"data,omitempty" xml:"data,omitempty" require:"true" type:"Repeated"`
}

func (s ReportBandwidthMultiDomainRealTimeEdgeServiceResponse) String() string {
	return tea.Prettify(s)
}

func (s ReportBandwidthMultiDomainRealTimeEdgeServiceResponse) GoString() string {
	return s.String()
}

func (s *ReportBandwidthMultiDomainRealTimeEdgeServiceResponse) SetCode(v string) *ReportBandwidthMultiDomainRealTimeEdgeServiceResponse {
	s.Code = &v
	return s
}

func (s *ReportBandwidthMultiDomainRealTimeEdgeServiceResponse) SetMessage(v string) *ReportBandwidthMultiDomainRealTimeEdgeServiceResponse {
	s.Message = &v
	return s
}

func (s *ReportBandwidthMultiDomainRealTimeEdgeServiceResponse) SetData(v []*ReportBandwidthMultiDomainRealTimeEdgeServiceResponseData) *ReportBandwidthMultiDomainRealTimeEdgeServiceResponse {
	s.Data = v
	return s
}

type ReportBandwidthMultiDomainRealTimeEdgeServiceResponseData struct {
	// {'en':'-', 'zh_CN':'域名，如果不选择域名分组维度，该字段为所有域名以分号分隔的字符串'}
	Domain *string `json:"domain,omitempty" xml:"domain,omitempty" require:"true"`
	// {'en':'-', 'zh_CN':'峰值带宽，单位Mbps，示例 （931556.21)'}
	PeakValue *string `json:"peakValue,omitempty" xml:"peakValue,omitempty" require:"true"`
	// {'en':'-', 'zh_CN':'峰值时间，示例（2019-02-13 18:01）'}
	PeakTime *string `json:"peakTime,omitempty" xml:"peakTime,omitempty" require:"true"`
	// {'en':'-', 'zh_CN':'边缘总流量，单位MB，示例 ( 74099.92)'}
	Total         *string                                                                   `json:"total,omitempty" xml:"total,omitempty" require:"true"`
	BandwidthData []*ReportBandwidthMultiDomainRealTimeEdgeServiceResponseDataBandwidthData `json:"bandwidthData,omitempty" xml:"bandwidthData,omitempty" require:"true" type:"Repeated"`
}

func (s ReportBandwidthMultiDomainRealTimeEdgeServiceResponseData) String() string {
	return tea.Prettify(s)
}

func (s ReportBandwidthMultiDomainRealTimeEdgeServiceResponseData) GoString() string {
	return s.String()
}

func (s *ReportBandwidthMultiDomainRealTimeEdgeServiceResponseData) SetDomain(v string) *ReportBandwidthMultiDomainRealTimeEdgeServiceResponseData {
	s.Domain = &v
	return s
}

func (s *ReportBandwidthMultiDomainRealTimeEdgeServiceResponseData) SetPeakValue(v string) *ReportBandwidthMultiDomainRealTimeEdgeServiceResponseData {
	s.PeakValue = &v
	return s
}

func (s *ReportBandwidthMultiDomainRealTimeEdgeServiceResponseData) SetPeakTime(v string) *ReportBandwidthMultiDomainRealTimeEdgeServiceResponseData {
	s.PeakTime = &v
	return s
}

func (s *ReportBandwidthMultiDomainRealTimeEdgeServiceResponseData) SetTotal(v string) *ReportBandwidthMultiDomainRealTimeEdgeServiceResponseData {
	s.Total = &v
	return s
}

func (s *ReportBandwidthMultiDomainRealTimeEdgeServiceResponseData) SetBandwidthData(v []*ReportBandwidthMultiDomainRealTimeEdgeServiceResponseDataBandwidthData) *ReportBandwidthMultiDomainRealTimeEdgeServiceResponseData {
	s.BandwidthData = v
	return s
}

type ReportBandwidthMultiDomainRealTimeEdgeServiceResponseDataBandwidthData struct {
	// {'en':'-', 'zh_CN':'格式为yyyy-MM-dd HH:mm；每一个时间片数据值代表的是前一个时间粒度范围内的数据值。
	//
	// 一天开始的时间片是yyyy-MM-dd 00:01，最后一个时间片是第二天（yyyy-MM-dd） 00:00。
	//
	// 返回开始时间和结束时间包含的时间片'}
	Timestamp *string `json:"timestamp,omitempty" xml:"timestamp,omitempty" require:"true"`
	// {'en':'-', 'zh_CN':'带宽值，单位Mbps，保留2位小数。'}
	Value *string `json:"value,omitempty" xml:"value,omitempty" require:"true"`
}

func (s ReportBandwidthMultiDomainRealTimeEdgeServiceResponseDataBandwidthData) String() string {
	return tea.Prettify(s)
}

func (s ReportBandwidthMultiDomainRealTimeEdgeServiceResponseDataBandwidthData) GoString() string {
	return s.String()
}

func (s *ReportBandwidthMultiDomainRealTimeEdgeServiceResponseDataBandwidthData) SetTimestamp(v string) *ReportBandwidthMultiDomainRealTimeEdgeServiceResponseDataBandwidthData {
	s.Timestamp = &v
	return s
}

func (s *ReportBandwidthMultiDomainRealTimeEdgeServiceResponseDataBandwidthData) SetValue(v string) *ReportBandwidthMultiDomainRealTimeEdgeServiceResponseDataBandwidthData {
	s.Value = &v
	return s
}

type Paths struct {
}

func (s Paths) String() string {
	return tea.Prettify(s)
}

func (s Paths) GoString() string {
	return s.String()
}

type Parameters struct {
}

func (s Parameters) String() string {
	return tea.Prettify(s)
}

func (s Parameters) GoString() string {
	return s.String()
}

type RequestHeader struct {
}

func (s RequestHeader) String() string {
	return tea.Prettify(s)
}

func (s RequestHeader) GoString() string {
	return s.String()
}

type ResponseHeader struct {
}

func (s ResponseHeader) String() string {
	return tea.Prettify(s)
}

func (s ResponseHeader) GoString() string {
	return s.String()
}
