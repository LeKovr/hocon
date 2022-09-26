<!--
MIT License

Source: https://gist.github.com/LeKovr/d6ee7d31c65a4b7e90d8d94295e4d535
Copyright (c) 2021 Aleksei Kovrizhkin (LeKovr)

Original version: https://github.com/pseudomuto/protoc-gen-doc/blob/v1.4.1/resources/markdown.tmpl
Copyright (c) 2017 David Muto (pseudomuto)

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
-->

<!-- use first file package name -->
# api.hocon.v1 API Documentation

<a name="top"></a>

## Table of Contents

- [proto/service.proto](#proto/service.proto)
  - Services
      - [HoconService](#api.hocon.v1.HoconService)
  
  - Messages
      - [LampStatus](#api.hocon.v1.LampStatus)
  
  - Enums
      - [LampScene](#api.hocon.v1.LampScene)
  
- [Scalar Value Types](#scalar-value-types)



<a name="proto/service.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## proto/service.proto




<a name="api.hocon.v1.HoconService"></a>

### HoconService

Home Control service

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| LampControl | [LampStatus](#api.hocon.v1.LampStatus) | [LampStatus](#api.hocon.v1.LampStatus) | управление лампочкой |

 <!-- end services -->


<a name="api.hocon.v1.LampStatus"></a>

### LampStatus

Lamp status attributes


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| scene | [LampScene](#api.hocon.v1.LampScene) |  |  |



          

 <!-- end messages -->


<a name="api.hocon.v1.LampScene"></a>

### LampScene

Lamp scenes

| Name | Number | Description |
| ---- | ------ | ----------- |
| UNKNOWN | 0 | A Standard tournament |
| OFF | 1 | Item is off |
| NIGHT | 2 | Item is in night mode |
| DAY | 3 | Item is in day mode |


 <!-- end enums -->

 <!-- end HasExtensions -->



## Scalar Value Types

| .proto Type | Notes | Go  | C++  | Java |
| ----------- | ----- | --- | ---- | ---- |
| <a name="double" /> double |  | float64 | double | double |
| <a name="float" /> float |  | float32 | float | float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int32 | int |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | int64 | long |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | uint32 | int |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | uint64 | long |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int32 | int |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | int64 | long |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | uint32 | int |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | uint64 | long |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int32 | int |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | int64 | long |
| <a name="bool" /> bool |  | bool | bool | boolean |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | string | String |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | []byte | string | ByteString |

