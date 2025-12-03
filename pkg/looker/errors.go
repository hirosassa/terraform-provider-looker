package looker

import "fmt"

// wrapSDKError wraps an SDK error with additional context about the operation and resource.
// This makes error messages more informative by showing which API call failed and on what resource.
//
// The format and args parameters work like fmt.Sprintf to build the resource identifier.
// Common patterns:
//   - Single identifier: wrapSDKError(err, "CreateConnection", "connection", "%s", "my-db-connection")
//   - Name and ID: wrapSDKError(err, "UpdateConnection", "connection", "name=%s, id=%s", name, id)
//   - Composite key: wrapSDKError(err, "AddGroupUser", "group_membership", "%s:%s", groupID, userID)
//
// Example outputs:
//
//	CreateConnection failed for connection "my-db-connection": response error. status=422. error={"message":"already_exists"}
//	UpdateConnection failed for connection "name=my-db-connection, id=123": response error. status=404. error={"message":"not_found"}
func wrapSDKError(err error, operation string, resourceType string, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	if format == "" {
		// For operations without a specific resource identifier (e.g., list operations)
		return fmt.Errorf("%s failed for %s: %w", operation, resourceType, err)
	}

	identifier := fmt.Sprintf(format, args...)
	return fmt.Errorf("%s failed for %s %q: %w", operation, resourceType, identifier, err)
}
