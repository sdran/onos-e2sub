# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api/e2/registry/v1beta1/registry.proto](#api/e2/registry/v1beta1/registry.proto)
    - [AddTerminationRequest](#registry.v1beta1.AddTerminationRequest)
    - [AddTerminationResponse](#registry.v1beta1.AddTerminationResponse)
    - [GetTerminationRequest](#registry.v1beta1.GetTerminationRequest)
    - [GetTerminationResponse](#registry.v1beta1.GetTerminationResponse)
    - [ListTerminationsRequest](#registry.v1beta1.ListTerminationsRequest)
    - [ListTerminationsResponse](#registry.v1beta1.ListTerminationsResponse)
    - [RemoveTerminationRequest](#registry.v1beta1.RemoveTerminationRequest)
    - [RemoveTerminationResponse](#registry.v1beta1.RemoveTerminationResponse)
    - [TerminationEndPoint](#registry.v1beta1.TerminationEndPoint)
    - [WatchTerminationsRequest](#registry.v1beta1.WatchTerminationsRequest)
    - [WatchTerminationsResponse](#registry.v1beta1.WatchTerminationsResponse)
  
    - [EventType](#registry.v1beta1.EventType)
  
    - [E2RegistryService](#registry.v1beta1.E2RegistryService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="api/e2/registry/v1beta1/registry.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/e2/registry/v1beta1/registry.proto



<a name="registry.v1beta1.AddTerminationRequest"></a>

### AddTerminationRequest
AddTerminationRequest is a request for adding a new termination point


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| end_point | [TerminationEndPoint](#registry.v1beta1.TerminationEndPoint) |  |  |






<a name="registry.v1beta1.AddTerminationResponse"></a>

### AddTerminationResponse
AddTerminationResponse is a response to adding a new termination point






<a name="registry.v1beta1.GetTerminationRequest"></a>

### GetTerminationRequest
GetTerminationRequest is a request for getting existing termination point


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="registry.v1beta1.GetTerminationResponse"></a>

### GetTerminationResponse
GetTerminationResponse is a response with invormation about a requested termination point


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| end_point | [TerminationEndPoint](#registry.v1beta1.TerminationEndPoint) |  |  |






<a name="registry.v1beta1.ListTerminationsRequest"></a>

### ListTerminationsRequest
ListTerminationsRequest is a request to list all available E2 terminations






<a name="registry.v1beta1.ListTerminationsResponse"></a>

### ListTerminationsResponse
ListTerminationsResponse is a response to list all available E2 terminations


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| end_points | [TerminationEndPoint](#registry.v1beta1.TerminationEndPoint) | repeated |  |






<a name="registry.v1beta1.RemoveTerminationRequest"></a>

### RemoveTerminationRequest
RemoveTerminationRequest is a request for removing termination point


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="registry.v1beta1.RemoveTerminationResponse"></a>

### RemoveTerminationResponse
RemoveTerminationResponse is a response to removing a termination point






<a name="registry.v1beta1.TerminationEndPoint"></a>

### TerminationEndPoint
Termination is a record identifying the IP address and TCP port coordinates where the E2 termination
service is available.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| ip | [string](#string) |  |  |
| port | [uint32](#uint32) |  |  |






<a name="registry.v1beta1.WatchTerminationsRequest"></a>

### WatchTerminationsRequest
WatchTerminationsRequest is a request to receive a stream of all E2 termination changes.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| noreplay | [bool](#bool) |  |  |






<a name="registry.v1beta1.WatchTerminationsResponse"></a>

### WatchTerminationsResponse
WatchTerminationsResponse is a response indicating a change in the available E2 termination end-points.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [EventType](#registry.v1beta1.EventType) |  |  |
| end_point | [TerminationEndPoint](#registry.v1beta1.TerminationEndPoint) |  |  |





 


<a name="registry.v1beta1.EventType"></a>

### EventType
Type of change

| Name | Number | Description |
| ---- | ------ | ----------- |
| NONE | 0 |  |
| ADDED | 1 |  |
| REMOVED | 3 |  |


 

 


<a name="registry.v1beta1.E2RegistryService"></a>

### E2RegistryService
E2RegistryService manages subscription and subscription delete requests

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| AddTermination | [AddTerminationRequest](#registry.v1beta1.AddTerminationRequest) | [AddTerminationResponse](#registry.v1beta1.AddTerminationResponse) | AddTermination registers new E2 termination end-point. |
| GetTermination | [GetTerminationRequest](#registry.v1beta1.GetTerminationRequest) | [GetTerminationResponse](#registry.v1beta1.GetTerminationResponse) | GetTermination retrieves information about a specific end-point |
| RemoveTermination | [RemoveTerminationRequest](#registry.v1beta1.RemoveTerminationRequest) | [RemoveTerminationResponse](#registry.v1beta1.RemoveTerminationResponse) | RemoveTermination removes the specified E2 termination end-point. |
| ListTerminations | [ListTerminationsRequest](#registry.v1beta1.ListTerminationsRequest) | [ListTerminationsResponse](#registry.v1beta1.ListTerminationsResponse) | ListTerminations returns the list of currently registered E2 terminations. |
| WatchTerminations | [WatchTerminationsRequest](#registry.v1beta1.WatchTerminationsRequest) | [WatchTerminationsResponse](#registry.v1beta1.WatchTerminationsResponse) stream | WatchTerminations returns a stream of changes in the set of available E2 terminations. |

 



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

