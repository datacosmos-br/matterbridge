package dify

// SummarizeGroup triggers a Dify workflow to summarize/triage a WhatsApp group.
func (c *Client) SummarizeGroup(groupID, transcript string, user string) (*WorkflowResponse, error) {
	return c.RunWorkflow(WorkflowRequest{
		Inputs: map[string]any{
			"group_id":   groupID,
			"transcript": transcript,
		},
		User: user,
		Mode: "blocking",
	}, "summarize-"+groupID)
}
