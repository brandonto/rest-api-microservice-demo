# {{classname}}

All URIs are relative to */*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateMessage**](MessagesApi.md#CreateMessage) | **Post** /messages | Create a message
[**DeleteMessageById**](MessagesApi.md#DeleteMessageById) | **Delete** /messages/{messageId} | Delete a specific message by ID
[**GetMessageById**](MessagesApi.md#GetMessageById) | **Get** /messages/{messageId} | Get a specific message by ID
[**ListMessages**](MessagesApi.md#ListMessages) | **Get** /messages | List all messages
[**UpdateMessageById**](MessagesApi.md#UpdateMessageById) | **Put** /messages/{messageId} | Update a specific message by ID

# **CreateMessage**
> CreateMessage(ctx, body)
Create a message

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Message**](Message.md)|  | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteMessageById**
> DeleteMessageById(ctx, messageId)
Delete a specific message by ID

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **messageId** | **int32**| The ID of the message | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetMessageById**
> InlineResponse2001 GetMessageById(ctx, messageId, optional)
Get a specific message by ID

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **messageId** | **int32**| The ID of the message | 
 **optional** | ***MessagesApiGetMessageByIdOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a MessagesApiGetMessageByIdOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **detailed** | **optional.Bool**| Include metadata | [default to false]

### Return type

[**InlineResponse2001**](inline_response_200_1.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListMessages**
> InlineResponse200 ListMessages(ctx, optional)
List all messages

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***MessagesApiListMessagesOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a MessagesApiListMessagesOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **limit** | **optional.Int32**| How many messages to return at one time (max 100) | [default to 20]
 **afterId** | **optional.Int32**| Show messages after a specified ID | [default to 0]
 **detailed** | **optional.Bool**| Include metadata | [default to false]

### Return type

[**InlineResponse200**](inline_response_200.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateMessageById**
> UpdateMessageById(ctx, body, messageId)
Update a specific message by ID

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**Message**](Message.md)|  | 
  **messageId** | **int32**| The ID of the message | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

