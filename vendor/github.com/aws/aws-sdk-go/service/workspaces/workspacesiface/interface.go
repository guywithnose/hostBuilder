// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

// Package workspacesiface provides an interface to enable mocking the Amazon WorkSpaces service client
// for testing your code.
//
// It is important to note that this interface will have breaking changes
// when the service model is updated and adds new API operations, paginators,
// and waiters.
package workspacesiface

import (
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/workspaces"
)

// WorkSpacesAPI provides an interface to enable mocking the
// workspaces.WorkSpaces service client's API operation,
// paginators, and waiters. This make unit testing your code that calls out
// to the SDK's service client's calls easier.
//
// The best way to use this interface is so the SDK's service client's calls
// can be stubbed out for unit testing your code with the SDK without needing
// to inject custom request handlers into the the SDK's request pipeline.
//
//    // myFunc uses an SDK service client to make a request to
//    // Amazon WorkSpaces.
//    func myFunc(svc workspacesiface.WorkSpacesAPI) bool {
//        // Make svc.CreateTags request
//    }
//
//    func main() {
//        sess := session.New()
//        svc := workspaces.New(sess)
//
//        myFunc(svc)
//    }
//
// In your _test.go file:
//
//    // Define a mock struct to be used in your unit tests of myFunc.
//    type mockWorkSpacesClient struct {
//        workspacesiface.WorkSpacesAPI
//    }
//    func (m *mockWorkSpacesClient) CreateTags(input *workspaces.CreateTagsInput) (*workspaces.CreateTagsOutput, error) {
//        // mock response/functionality
//    }
//
//    TestMyFunc(t *testing.T) {
//        // Setup Test
//        mockSvc := &mockWorkSpacesClient{}
//
//        myfunc(mockSvc)
//
//        // Verify myFunc's functionality
//    }
//
// It is important to note that this interface will have breaking changes
// when the service model is updated and adds new API operations, paginators,
// and waiters. Its suggested to use the pattern above for testing, or using
// tooling to generate mocks to satisfy the interfaces.
type WorkSpacesAPI interface {
	CreateTagsRequest(*workspaces.CreateTagsInput) (*request.Request, *workspaces.CreateTagsOutput)

	CreateTags(*workspaces.CreateTagsInput) (*workspaces.CreateTagsOutput, error)

	CreateWorkspacesRequest(*workspaces.CreateWorkspacesInput) (*request.Request, *workspaces.CreateWorkspacesOutput)

	CreateWorkspaces(*workspaces.CreateWorkspacesInput) (*workspaces.CreateWorkspacesOutput, error)

	DeleteTagsRequest(*workspaces.DeleteTagsInput) (*request.Request, *workspaces.DeleteTagsOutput)

	DeleteTags(*workspaces.DeleteTagsInput) (*workspaces.DeleteTagsOutput, error)

	DescribeTagsRequest(*workspaces.DescribeTagsInput) (*request.Request, *workspaces.DescribeTagsOutput)

	DescribeTags(*workspaces.DescribeTagsInput) (*workspaces.DescribeTagsOutput, error)

	DescribeWorkspaceBundlesRequest(*workspaces.DescribeWorkspaceBundlesInput) (*request.Request, *workspaces.DescribeWorkspaceBundlesOutput)

	DescribeWorkspaceBundles(*workspaces.DescribeWorkspaceBundlesInput) (*workspaces.DescribeWorkspaceBundlesOutput, error)

	DescribeWorkspaceBundlesPages(*workspaces.DescribeWorkspaceBundlesInput, func(*workspaces.DescribeWorkspaceBundlesOutput, bool) bool) error

	DescribeWorkspaceDirectoriesRequest(*workspaces.DescribeWorkspaceDirectoriesInput) (*request.Request, *workspaces.DescribeWorkspaceDirectoriesOutput)

	DescribeWorkspaceDirectories(*workspaces.DescribeWorkspaceDirectoriesInput) (*workspaces.DescribeWorkspaceDirectoriesOutput, error)

	DescribeWorkspaceDirectoriesPages(*workspaces.DescribeWorkspaceDirectoriesInput, func(*workspaces.DescribeWorkspaceDirectoriesOutput, bool) bool) error

	DescribeWorkspacesRequest(*workspaces.DescribeWorkspacesInput) (*request.Request, *workspaces.DescribeWorkspacesOutput)

	DescribeWorkspaces(*workspaces.DescribeWorkspacesInput) (*workspaces.DescribeWorkspacesOutput, error)

	DescribeWorkspacesPages(*workspaces.DescribeWorkspacesInput, func(*workspaces.DescribeWorkspacesOutput, bool) bool) error

	DescribeWorkspacesConnectionStatusRequest(*workspaces.DescribeWorkspacesConnectionStatusInput) (*request.Request, *workspaces.DescribeWorkspacesConnectionStatusOutput)

	DescribeWorkspacesConnectionStatus(*workspaces.DescribeWorkspacesConnectionStatusInput) (*workspaces.DescribeWorkspacesConnectionStatusOutput, error)

	ModifyWorkspacePropertiesRequest(*workspaces.ModifyWorkspacePropertiesInput) (*request.Request, *workspaces.ModifyWorkspacePropertiesOutput)

	ModifyWorkspaceProperties(*workspaces.ModifyWorkspacePropertiesInput) (*workspaces.ModifyWorkspacePropertiesOutput, error)

	RebootWorkspacesRequest(*workspaces.RebootWorkspacesInput) (*request.Request, *workspaces.RebootWorkspacesOutput)

	RebootWorkspaces(*workspaces.RebootWorkspacesInput) (*workspaces.RebootWorkspacesOutput, error)

	RebuildWorkspacesRequest(*workspaces.RebuildWorkspacesInput) (*request.Request, *workspaces.RebuildWorkspacesOutput)

	RebuildWorkspaces(*workspaces.RebuildWorkspacesInput) (*workspaces.RebuildWorkspacesOutput, error)

	StartWorkspacesRequest(*workspaces.StartWorkspacesInput) (*request.Request, *workspaces.StartWorkspacesOutput)

	StartWorkspaces(*workspaces.StartWorkspacesInput) (*workspaces.StartWorkspacesOutput, error)

	StopWorkspacesRequest(*workspaces.StopWorkspacesInput) (*request.Request, *workspaces.StopWorkspacesOutput)

	StopWorkspaces(*workspaces.StopWorkspacesInput) (*workspaces.StopWorkspacesOutput, error)

	TerminateWorkspacesRequest(*workspaces.TerminateWorkspacesInput) (*request.Request, *workspaces.TerminateWorkspacesOutput)

	TerminateWorkspaces(*workspaces.TerminateWorkspacesInput) (*workspaces.TerminateWorkspacesOutput, error)
}

var _ WorkSpacesAPI = (*workspaces.WorkSpaces)(nil)
