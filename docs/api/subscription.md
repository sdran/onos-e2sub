# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api/e2/subscription/v1beta1/subscription.proto](#api/e2/subscription/v1beta1/subscription.proto)
    - [AddSubscriptionRequest](#subscription.v1beta1.AddSubscriptionRequest)
    - [AddSubscriptionResponse](#subscription.v1beta1.AddSubscriptionResponse)
    - [Event](#subscription.v1beta1.Event)
    - [GetSubscriptionRequest](#subscription.v1beta1.GetSubscriptionRequest)
    - [GetSubscriptionResponse](#subscription.v1beta1.GetSubscriptionResponse)
    - [Lifecycle](#subscription.v1beta1.Lifecycle)
    - [ListSubscriptionsRequest](#subscription.v1beta1.ListSubscriptionsRequest)
    - [ListSubscriptionsResponse](#subscription.v1beta1.ListSubscriptionsResponse)
    - [Payload](#subscription.v1beta1.Payload)
    - [RemoveSubscriptionRequest](#subscription.v1beta1.RemoveSubscriptionRequest)
    - [RemoveSubscriptionResponse](#subscription.v1beta1.RemoveSubscriptionResponse)
    - [ServiceModel](#subscription.v1beta1.ServiceModel)
    - [Subscription](#subscription.v1beta1.Subscription)
    - [WatchSubscriptionsRequest](#subscription.v1beta1.WatchSubscriptionsRequest)
    - [WatchSubscriptionsResponse](#subscription.v1beta1.WatchSubscriptionsResponse)
  
    - [Encoding](#subscription.v1beta1.Encoding)
    - [EventType](#subscription.v1beta1.EventType)
    - [Status](#subscription.v1beta1.Status)
  
    - [E2SubscriptionService](#subscription.v1beta1.E2SubscriptionService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="api/e2/subscription/v1beta1/subscription.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/e2/subscription/v1beta1/subscription.proto



<a name="subscription.v1beta1.AddSubscriptionRequest"></a>

### AddSubscriptionRequest
AddSubscriptionRequest a subscription request


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| subscription | [Subscription](#subscription.v1beta1.Subscription) |  |  |






<a name="subscription.v1beta1.AddSubscriptionResponse"></a>

### AddSubscriptionResponse
AddSubscriptionResponse a subscription response


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| subscription | [Subscription](#subscription.v1beta1.Subscription) |  |  |






<a name="subscription.v1beta1.Event"></a>

### Event
Event is a subscription event


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [EventType](#subscription.v1beta1.EventType) |  |  |
| subscription | [Subscription](#subscription.v1beta1.Subscription) |  |  |






<a name="subscription.v1beta1.GetSubscriptionRequest"></a>

### GetSubscriptionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="subscription.v1beta1.GetSubscriptionResponse"></a>

### GetSubscriptionResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| subscription | [Subscription](#subscription.v1beta1.Subscription) |  |  |






<a name="subscription.v1beta1.Lifecycle"></a>

### Lifecycle
Lifecycle is the subscription lifecycle


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [Status](#subscription.v1beta1.Status) |  |  |






<a name="subscription.v1beta1.ListSubscriptionsRequest"></a>

### ListSubscriptionsRequest







<a name="subscription.v1beta1.ListSubscriptionsResponse"></a>

### ListSubscriptionsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| subscriptions | [Subscription](#subscription.v1beta1.Subscription) | repeated |  |






<a name="subscription.v1beta1.Payload"></a>

### Payload
Payload is a subscription payload


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| encoding | [Encoding](#subscription.v1beta1.Encoding) |  |  |
| bytes | [bytes](#bytes) |  |  |






<a name="subscription.v1beta1.RemoveSubscriptionRequest"></a>

### RemoveSubscriptionRequest
RemoveSubscriptionRequest a subscription delete request


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="subscription.v1beta1.RemoveSubscriptionResponse"></a>

### RemoveSubscriptionResponse
RemoveSubscriptionResponse a subscription delete response






<a name="subscription.v1beta1.ServiceModel"></a>

### ServiceModel
ServiceModel is a service model definition


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="subscription.v1beta1.Subscription"></a>

### Subscription
Subscription is a subscription state


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| revision | [uint64](#uint64) |  |  |
| app_id | [string](#string) |  |  |
| e2_node_id | [string](#string) |  |  |
| service_model | [ServiceModel](#subscription.v1beta1.ServiceModel) |  |  |
| payload | [Payload](#subscription.v1beta1.Payload) |  |  |
| lifecycle | [Lifecycle](#subscription.v1beta1.Lifecycle) |  |  |






<a name="subscription.v1beta1.WatchSubscriptionsRequest"></a>

### WatchSubscriptionsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| noreplay | [bool](#bool) |  |  |






<a name="subscription.v1beta1.WatchSubscriptionsResponse"></a>

### WatchSubscriptionsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| event | [Event](#subscription.v1beta1.Event) |  |  |





 


<a name="subscription.v1beta1.Encoding"></a>

### Encoding
Encoding indicates a payload encoding

| Name | Number | Description |
| ---- | ------ | ----------- |
| ENCODING_ASN1 | 0 |  |
| ENCODING_PROTO | 1 |  |



<a name="subscription.v1beta1.EventType"></a>

### EventType
EventType is a subscription event type

| Name | Number | Description |
| ---- | ------ | ----------- |
| NONE | 0 |  |
| ADDED | 1 |  |
| UPDATED | 2 |  |
| REMOVED | 3 |  |



<a name="subscription.v1beta1.Status"></a>

### Status
Status is a subscription status

| Name | Number | Description |
| ---- | ------ | ----------- |
| ACTIVE | 0 |  |
| PENDING_DELETE | 1 |  |


 

 


<a name="subscription.v1beta1.E2SubscriptionService"></a>

### E2SubscriptionService
SubscriptionService manages subscription and subscription delete requests

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| AddSubscription | [AddSubscriptionRequest](#subscription.v1beta1.AddSubscriptionRequest) | [AddSubscriptionResponse](#subscription.v1beta1.AddSubscriptionResponse) | AddSubscription establishes E2 subscriptions on E2 Node. |
| RemoveSubscription | [RemoveSubscriptionRequest](#subscription.v1beta1.RemoveSubscriptionRequest) | [RemoveSubscriptionResponse](#subscription.v1beta1.RemoveSubscriptionResponse) | RemoveSubscription removes E2 subscriptions on E2 Node. |
| GetSubscription | [GetSubscriptionRequest](#subscription.v1beta1.GetSubscriptionRequest) | [GetSubscriptionResponse](#subscription.v1beta1.GetSubscriptionResponse) | GetSubscription retrieves information about a specific subscription in the list of existing subscriptions |
| ListSubscriptions | [ListSubscriptionsRequest](#subscription.v1beta1.ListSubscriptionsRequest) | [ListSubscriptionsResponse](#subscription.v1beta1.ListSubscriptionsResponse) | ListSubscriptions returns the list of current existing subscriptions |
| WatchSubscriptions | [WatchSubscriptionsRequest](#subscription.v1beta1.WatchSubscriptionsRequest) | [WatchSubscriptionsResponse](#subscription.v1beta1.WatchSubscriptionsResponse) stream | WatchSubscriptions returns a stream of subscription changes |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers ??? if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers ??? if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

