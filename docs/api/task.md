# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api/e2/task/v1beta1/task.proto](#api/e2/task/v1beta1/task.proto)
    - [Event](#task.v1beta1.Event)
    - [GetSubscriptionTaskRequest](#task.v1beta1.GetSubscriptionTaskRequest)
    - [GetSubscriptionTaskResponse](#task.v1beta1.GetSubscriptionTaskResponse)
    - [Lifecycle](#task.v1beta1.Lifecycle)
    - [ListSubscriptionTasksRequest](#task.v1beta1.ListSubscriptionTasksRequest)
    - [ListSubscriptionTasksResponse](#task.v1beta1.ListSubscriptionTasksResponse)
    - [SubscriptionTask](#task.v1beta1.SubscriptionTask)
    - [UpdateSubscriptionTaskRequest](#task.v1beta1.UpdateSubscriptionTaskRequest)
    - [UpdateSubscriptionTaskResponse](#task.v1beta1.UpdateSubscriptionTaskResponse)
    - [WatchSubscriptionTasksRequest](#task.v1beta1.WatchSubscriptionTasksRequest)
    - [WatchSubscriptionTasksResponse](#task.v1beta1.WatchSubscriptionTasksResponse)
  
    - [EventType](#task.v1beta1.EventType)
    - [Phase](#task.v1beta1.Phase)
    - [Status](#task.v1beta1.Status)
  
    - [E2SubscriptionTaskService](#task.v1beta1.E2SubscriptionTaskService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="api/e2/task/v1beta1/task.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/e2/task/v1beta1/task.proto



<a name="task.v1beta1.Event"></a>

### Event
Event is a SubscriptionTask event


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [EventType](#task.v1beta1.EventType) |  |  |
| task | [SubscriptionTask](#task.v1beta1.SubscriptionTask) |  |  |






<a name="task.v1beta1.GetSubscriptionTaskRequest"></a>

### GetSubscriptionTaskRequest
GetSubscriptionTaskRequest is a request for getting existing SubscriptionTask


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="task.v1beta1.GetSubscriptionTaskResponse"></a>

### GetSubscriptionTaskResponse
GetSubscriptionTaskResponse is a response with invormation about a requested SubscriptionTask


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| task | [SubscriptionTask](#task.v1beta1.SubscriptionTask) |  |  |






<a name="task.v1beta1.Lifecycle"></a>

### Lifecycle
Lifecycle is a subscription task status


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| phase | [Phase](#task.v1beta1.Phase) |  |  |
| status | [Status](#task.v1beta1.Status) |  |  |






<a name="task.v1beta1.ListSubscriptionTasksRequest"></a>

### ListSubscriptionTasksRequest
ListSubscriptionTasksRequest is a request to list all available SubscriptionTasks


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| subscription_id | [string](#string) |  |  |
| endpoint_id | [string](#string) |  |  |






<a name="task.v1beta1.ListSubscriptionTasksResponse"></a>

### ListSubscriptionTasksResponse
ListSubscriptionTasksResponse is a response to list all available SubscriptionTasks


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| tasks | [SubscriptionTask](#task.v1beta1.SubscriptionTask) | repeated |  |






<a name="task.v1beta1.SubscriptionTask"></a>

### SubscriptionTask
SubscriptionTask is a task representing a subscription between an E2 termination and an E2 node


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| revision | [uint64](#uint64) |  |  |
| subscription_id | [string](#string) |  |  |
| endpoint_id | [string](#string) |  |  |
| lifecycle | [Lifecycle](#task.v1beta1.Lifecycle) |  |  |






<a name="task.v1beta1.UpdateSubscriptionTaskRequest"></a>

### UpdateSubscriptionTaskRequest
UpdateSubscriptionTaskRequest is a request for updating a SubscriptionTask status


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| task | [SubscriptionTask](#task.v1beta1.SubscriptionTask) |  |  |






<a name="task.v1beta1.UpdateSubscriptionTaskResponse"></a>

### UpdateSubscriptionTaskResponse
UpdateSubscriptionTaskResponse is a response to updating a SubscriptionTask status






<a name="task.v1beta1.WatchSubscriptionTasksRequest"></a>

### WatchSubscriptionTasksRequest
WatchSubscriptionTasksRequest is a request to receive a stream of all SubscriptionTask changes.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| noreplay | [bool](#bool) |  |  |
| subscription_id | [string](#string) |  |  |
| endpoint_id | [string](#string) |  |  |






<a name="task.v1beta1.WatchSubscriptionTasksResponse"></a>

### WatchSubscriptionTasksResponse
WatchSubscriptionTasksResponse is a response indicating a change in the available SubscriptionTasks.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| event | [Event](#task.v1beta1.Event) |  |  |





 


<a name="task.v1beta1.EventType"></a>

### EventType
Type of change

| Name | Number | Description |
| ---- | ------ | ----------- |
| NONE | 0 |  |
| CREATED | 1 |  |
| UPDATED | 2 |  |
| REMOVED | 3 |  |



<a name="task.v1beta1.Phase"></a>

### Phase
Phase is a subscription task phase

| Name | Number | Description |
| ---- | ------ | ----------- |
| OPEN | 0 | OPEN is a subscription task open phase |
| CLOSE | 1 | CLOSE is a subscription task close phase |



<a name="task.v1beta1.Status"></a>

### Status
Status is a subscription task status

| Name | Number | Description |
| ---- | ------ | ----------- |
| PENDING | 0 | PENDING indicates the subscription task phase is pending |
| COMPLETE | 1 | COMPLETE indicates the subscription task phase is complete |


 

 


<a name="task.v1beta1.E2SubscriptionTaskService"></a>

### E2SubscriptionTaskService
E2SubscriptionTaskService manages subscription tasks between E2 termination points and E2 nodes

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetSubscriptionTask | [GetSubscriptionTaskRequest](#task.v1beta1.GetSubscriptionTaskRequest) | [GetSubscriptionTaskResponse](#task.v1beta1.GetSubscriptionTaskResponse) | GetSubscriptionTask retrieves information about a specific task |
| ListSubscriptionTasks | [ListSubscriptionTasksRequest](#task.v1beta1.ListSubscriptionTasksRequest) | [ListSubscriptionTasksResponse](#task.v1beta1.ListSubscriptionTasksResponse) | ListSubscriptionTasks returns the list of currently registered E2 Tasks. |
| WatchSubscriptionTasks | [WatchSubscriptionTasksRequest](#task.v1beta1.WatchSubscriptionTasksRequest) | [WatchSubscriptionTasksResponse](#task.v1beta1.WatchSubscriptionTasksResponse) stream | WatchSubscriptionTasks returns a stream of changes in the set of available E2 Tasks. |
| UpdateSubscriptionTask | [UpdateSubscriptionTaskRequest](#task.v1beta1.UpdateSubscriptionTaskRequest) | [UpdateSubscriptionTaskResponse](#task.v1beta1.UpdateSubscriptionTaskResponse) | UpdateSubscriptionTask updates a task status |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
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

