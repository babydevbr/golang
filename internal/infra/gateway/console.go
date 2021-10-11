// Package gateway ...
package gateway

import "log"

// Console ....
type Console struct{}

// Send ....
func (c *Console) Send(recipient, msg string) error {
	log.Println(recipient, msg)

	return nil
}
