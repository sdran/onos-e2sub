# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api/e2/channel/v1beta1/channel.proto](#api/e2/channel/v1beta1/channel.proto)
    - [Channel](#channel.v1beta1.Channel)
    - [GetChannelRequest](#channel.v1beta1.GetChannelRequest)
    - [GetChannelResponse](#channel.v1beta1.GetChannelResponse)
    - [ListChannelsRequest](#channel.v1beta1.ListChannelsRequest)
    - [ListChannelsResponse](#channel.v1beta1.ListChannelsResponse)
    - [State](#channel.v1beta1.State)
    - [UpdateChannelStateRequest](#channel.v1beta1.UpdateChannelStateRequest)
    - [UpdateChannelStateResponse](#channel.v1beta1.UpdateChannelStateResponse)
    - [WatchChannelsRequest](#channel.v1beta1.WatchChannelsRequest)
    - [WatchChannelsResponse](#channel.v1beta1.WatchChannelsResponse)
  
    - [EventType](#channel.v1beta1.EventType)
    - [Status](#channel.v1beta1.Status)
  
    - [E2ChannelService](#channel.v1beta1.E2ChannelService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="api/e2/channel/v1beta1/channel.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/e2/channel/v1beta1/channel.proto



<a name="channel.v1beta1.Channel"></a>

### Channel
Channel is a record representing a subscription between an E2 termination and an E2 node


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| revision | [uint64](#uint64) |  |  |
| subscription_id | [string](#string) |  |  |
| termination_endpoint_id | [string](#string) |  |  |
| state | [State](#channel.v1beta1.State) |  |  |






<a name="channel.v1beta1.GetChannelRequest"></a>

### GetChannelRequest
GetChannelRequest is a request for getting existing Channel


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="channel.v1beta1.GetChannelResponse"></a>

### GetChannelResponse
GetChannelResponse is a response with invormation about a requested Channel


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| channel | [Channel](#channel.v1beta1.Channel) |  |  |






<a name="channel.v1beta1.ListChannelsRequest"></a>

### ListChannelsRequest
ListChannelsRequest is a request to list all available E2 Channels






<a name="channel.v1beta1.ListChannelsResponse"></a>

### ListChannelsResponse
ListChannelsResponse is a response to list all available E2 Channels


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| channel | [Channel](#channel.v1beta1.Channel) | repeated |  |






<a name="channel.v1beta1.State"></a>

### State
State is a channel state


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [Status](#channel.v1beta1.Status) |  |  |






<a name="channel.v1beta1.UpdateChannelStateRequest"></a>

### UpdateChannelStateRequest
UpdateChannelRequest is a request for updating a Channel state


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| channel | [Channel](#channel.v1beta1.Channel) |  |  |






<a name="channel.v1beta1.UpdateChannelStateResponse"></a>

### UpdateChannelStateResponse
UpdateChannelResponse is a response to updating a Channel state






<a name="channel.v1beta1.WatchChannelsRequest"></a>

### WatchChannelsRequest
WatchChannelsRequest is a request to receive a stream of all E2 Channel changes.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| noreplay | [bool](#bool) |  |  |






<a name="channel.v1beta1.WatchChannelsResponse"></a>

### WatchChannelsResponse
WatchChannelsResponse is a response indicating a change in the available E2 Channels.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [EventType](#channel.v1beta1.EventType) |  |  |
| channel | [Channel](#channel.v1beta1.Channel) |  |  |





 


<a name="channel.v1beta1.EventType"></a>

### EventType
Type of change

| Name | Number | Description |
| ---- | ------ | ----------- |
| NONE | 0 |  |
| CREATED | 1 |  |
| UPDATED | 2 |  |
| REMOVED | 3 |  |



<a name="channel.v1beta1.Status"></a>

### Status
Status is a channel status

| Name | Number | Description |
| ---- | ------ | ----------- |
| INACTIVE | 0 | INACTIVE indicates the channel is inactive |
| ACTIVE | 1 | ACTIVE indicates the channel is active |


 

 


<a name="channel.v1beta1.E2ChannelService"></a>

### E2ChannelService
E2ChannelService manages subscription channels between E2 termination points and E2 nodes

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetChannel | [GetChannelRequest](#channel.v1beta1.GetChannelRequest) | [GetChannelResponse](#channel.v1beta1.GetChannelResponse) | GetChannel retrieves information about a specific channel |
| ListChannels | [ListChannelsRequest](#channel.v1beta1.ListChannelsRequest) | [ListChannelsResponse](#channel.v1beta1.ListChannelsResponse) | ListChannels returns the list of currently registered E2 Channels. |
| WatchChannels | [WatchChannelsRequest](#channel.v1beta1.WatchChannelsRequest) | [WatchChannelsResponse](#channel.v1beta1.WatchChannelsResponse) stream | WatchChannels returns a stream of changes in the set of available E2 Channels. |
| UpdateChannelState | [UpdateChannelStateRequest](#channel.v1beta1.UpdateChannelStateRequest) | [UpdateChannelStateResponse](#channel.v1beta1.UpdateChannelStateResponse) | UpdateChannelState updates a channel state |

 



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

