// Command prlabel labels a pull request diff by its conceptual scope,
// following the same pattern as jev-pr-labeler and commit-miner in the
// Jev community: a Choice question over the diff instead of line counts.
//
// See https://github.com/hellogumbo/awesome-jev for the real projects this
// pattern is drawn from.
package main

import (
	"context"
	"fmt"
	"log"

	typesafe "github.com/atharvamhaske/typesafe-sdk-go"
)

func main() {
	client, err := typesafe.NewClient() // reads TYPESAFE_API_KEY
	if err != nil {
		log.Fatal(err)
	}

	diff := `diff --git a/auth/session.go b/auth/session.go
-func ValidateToken(t string) bool {
-    return t != ""
-}
+func ValidateToken(t string) bool {
+    return jwt.Verify(t, publicKey) == nil
+}`

	resp, err := client.SystemOne(context.Background(), diff,
		map[string]typesafe.Question{
			"scope": typesafe.Choice{
				Instructions: "What is the conceptual scope of this diff?",
				Criteria: map[string]string{
					"security": "Fixes or hardens authentication, authorization, or input validation",
					"feature":  "Adds new user-facing functionality",
					"refactor": "Restructures existing code without changing behavior",
					"docs":     "Changes documentation or comments only",
				},
			},
		})
	if err != nil {
		log.Fatal(err)
	}

	label := resp.Choices()["scope"]
	fmt.Printf("label: %s (confidence %.2f)\n", label.Choice, label.Confidence)
}
