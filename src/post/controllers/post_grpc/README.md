# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [post.proto](#post.proto)
    - [CreateReq](#post_grpc.CreateReq)
    - [DeleteReq](#post_grpc.DeleteReq)
    - [ID](#post_grpc.ID)
    - [ListPost](#post_grpc.ListPost)
    - [ListReq](#post_grpc.ListReq)
    - [Post](#post_grpc.Post)
    - [UpdateReq](#post_grpc.UpdateReq)
  
  
  
    - [PostService](#post_grpc.PostService)
  

- [Scalar Value Types](#scalar-value-types)



<a name="post.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## post.proto



<a name="post_grpc.CreateReq"></a>

### CreateReq



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| title | [string](#string) |  | 1文字以上20文字以下 |
| content | [string](#string) |  |  |
| user_id | [int64](#int64) |  |  |






<a name="post_grpc.DeleteReq"></a>

### DeleteReq



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int64](#int64) |  |  |
| user_id | [int64](#int64) |  |  |






<a name="post_grpc.ID"></a>

### ID



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int64](#int64) |  |  |






<a name="post_grpc.ListPost"></a>

### ListPost



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| posts | [Post](#post_grpc.Post) | repeated |  |






<a name="post_grpc.ListReq"></a>

### ListReq



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| datetime | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |
| num | [int64](#int64) |  |  |






<a name="post_grpc.Post"></a>

### Post



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int64](#int64) |  |  |
| title | [string](#string) |  |  |
| content | [string](#string) |  |  |
| created_at | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |
| updated_at | [google.protobuf.Timestamp](#google.protobuf.Timestamp) |  |  |
| user_id | [int64](#int64) |  |  |






<a name="post_grpc.UpdateReq"></a>

### UpdateReq



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [int64](#int64) |  |  |
| title | [string](#string) |  |  |
| content | [string](#string) |  |  |
| user_id | [int64](#int64) |  |  |





 

 

 


<a name="post_grpc.PostService"></a>

### PostService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Create | [CreateReq](#post_grpc.CreateReq) | [Post](#post_grpc.Post) | 投稿を作成 |
| GetByID | [ID](#post_grpc.ID) | [Post](#post_grpc.Post) | post_idで投稿を取得 |
| GetList | [ListReq](#post_grpc.ListReq) | [ListPost](#post_grpc.ListPost) | 取得件数と日時を指定して投稿を複数取得 |
| Update | [UpdateReq](#post_grpc.UpdateReq) | [Post](#post_grpc.Post) | 投稿を更新 |
| Delete | [DeleteReq](#post_grpc.DeleteReq) | [.google.protobuf.BoolValue](#google.protobuf.BoolValue) | 投稿を削除 |

 



## Scalar Value Types

| .proto Type | Notes | C++ Type | Java Type | Python Type |
| ----------- | ----- | -------- | --------- | ----------- |
| <a name="double" /> double |  | double | double | float |
| <a name="float" /> float |  | float | float | float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long |
| <a name="bool" /> bool |  | bool | boolean | boolean |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str |

