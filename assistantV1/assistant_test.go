package assistantV1

import (
	"fmt"
	"strings"
	"testing"
	"io/ioutil"
	watson "golang-sdk"
)

func TestWorkspace(t *testing.T) {
	readCreds, readCredsErr := ioutil.ReadFile("credentials.txt")

	if readCredsErr != nil {
		fmt.Println(readCredsErr)
		t.Fail()
	}

	creds := strings.Split(string(readCreds), "\n")

	assistant, assistantErr := NewAssistantV1(watson.Credentials{
		ServiceURL: creds[0],
		Version: "2018-02-16",
		Username: creds[1],
		Password: creds[2],
	})

	if assistantErr != nil {
		fmt.Println(assistantErr)
		return
	}

	list := testList(assistant)

	testGet(assistant, list)

	testDelete(assistant, list)

	testList(assistant)

	testCreate(assistant)

	list = testList(assistant)

	testUpdate(assistant, list)
}

func testList(assistant *AssistantV1) *ListWorkspacesResponse {
	listReq := ListWorkspacesRequest{
		IncludeCount: true,
	}

	list, listErr := assistant.ListWorkspaces(&listReq)

	if listErr != nil {
		fmt.Println(listErr)
		return nil
	}

	result := list.ListWorkspacesResponse()

	if result != nil {
		fmt.Printf("FOUND %v WORKSPACES\n", len(result.Workspaces))
		return result
	}

	return nil
}

func testGet(assistant *AssistantV1, list *ListWorkspacesResponse) {
	for i, workspace := range list.Workspaces {
		workspaceID := workspace.WorkspaceID

		get, getErr := assistant.GetWorkspace(workspaceID, new(GetWorkspaceRequest))

		if getErr != nil {
			fmt.Println(getErr)
			return
		}

		result := get.GetWorkspaceResponse()

		if result != nil {
			fmt.Printf("WORKSPACE %v: %+v\n", i, result)
		}
	}
}

func testDelete(assistant *AssistantV1, list *ListWorkspacesResponse) {
	for i, workspace := range list.Workspaces {
		workspaceID := workspace.WorkspaceID

		_, delErr := assistant.DeleteWorkspace(workspaceID)

		if delErr != nil {
			fmt.Println(delErr)
			return
		}

		fmt.Printf("DELETED WORKSPACE %v: %v\n", i, workspaceID)
	}
}

func testCreate(assistant *AssistantV1) {
	for i := 0; i < 3; i++ {
		createReq := CreateWorkspace{
			Name: fmt.Sprintf("create%v", i),
		}

		create, createErr := assistant.CreateWorkspace(&createReq)

		if createErr != nil {
			fmt.Println(createErr)
			return
		}

		result := create.CreateWorkspaceResponse()

		if result != nil {
			fmt.Printf("CREATED WORKSPACE %v: %v\n", i, result.Name)
		}
	}
}

func testUpdate(assistant *AssistantV1, list *ListWorkspacesResponse) {
	for i, workspace := range list.Workspaces {
		workspaceID := workspace.WorkspaceID

		updateReq := CreateWorkspace{
			Name: fmt.Sprintf("create%v", i * 10),
		}

		update, updateErr := assistant.UpdateWorkspace(workspaceID, &updateReq, new(UpdateWorkspaceRequest))

		if updateErr != nil {
			fmt.Println(updateErr)
			return
		}

		result := update.UpdateWorkspaceResponse()

		if result != nil {
			fmt.Printf("UPDATED WORKSPACE %v: %v\n", i, result.Name)
		}
	}
}
