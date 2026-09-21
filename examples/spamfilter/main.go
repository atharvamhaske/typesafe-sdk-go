// Command spamfilter scores a message for spam, following the same pattern
// as the many Jev-powered moderation and antispam bots in the community:
// one Noul question per message, thresholded in code.
//
// See https://github.com/hellogumbo/awesome-jev for real examples of this
// pattern (Jev-Moderation-Bot, scam-shield, jev_antispam_bot, and others).
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

	messages := []string{
		"Reminder: your standup notes are due by 10am.",
		"CONGRATULATIONS! You've won a $1000 gift card, click here NOW to claim: bit.ly/claim-now",
	}

	for _, msg := range messages {
		resp, err := client.SystemOne(context.Background(), msg,
			map[string]typesafe.Question{
				"is_spam": typesafe.Noul{
					Instructions: "Is this message spam, a scam, or unsolicited advertising?",
					Criteria:     map[string]string{"true": "Spam or scam", "false": "Legitimate message"},
				},
			})
		if err != nil {
			log.Fatal(err)
		}

		score := resp.Nouls()["is_spam"].Noul
		verdict := "clean"
		if score > 0.5 {
			verdict = "SPAM"
		}
		fmt.Printf("%-6s (%.2f)  %s\n", verdict, score, msg)
	}
}
