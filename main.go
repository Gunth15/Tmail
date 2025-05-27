package main

import (
	"fmt"
	"log"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
)

// NOTE: May need to explore possible custom chracterset readers which may make parsing more effiucient for future customizations

var token string

func main() {
	client, err := imapclient.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	if err = client.Login(testUsername, testPass).Wait(); err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("Successfully logged in")
	// ReturnStatus requires server support for IMAP4rev2 or LIST-STATUS
	listCmd := client.List("", "%", &imap.ListOptions{
		ReturnStatus: &imap.StatusOptions{
			NumMessages: true,
			NumUnseen:   true,
		},
	})
	for {
		mbox := listCmd.Next()
		if mbox == nil || mbox.Status == nil {
			break
		}
		fmt.Printf("Mailbox %q contains %v messages (%v unseen)\n", mbox.Mailbox, *mbox.Status.NumMessages, *mbox.Status.NumUnseen)
	}
	if err := listCmd.Close(); err != nil {
		log.Fatalf("LIST command failed: %v", err)
	}

	if err = client.Logout().Wait(); err != nil {
		log.Fatalf("Failed to loggout: %v", err)
	}
	fmt.Println("Successfully logged out")
}
